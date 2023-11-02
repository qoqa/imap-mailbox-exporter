package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/joho/godotenv"
	"github.com/jop-software/imap-mailbox-exporter/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var cfg *config.Config

func countMailsInMailbox(server config.ConfigServer, account config.ConfigAcccount, mailbox string) (uint32, int, error) {
	c, err := client.DialTLS(server.HostPort(), nil)
	if err != nil {
		return 0, 0, err
	}

	defer c.Logout()

	// Login
	if err := c.Login(account.Username, account.Password); err != nil {
		return 0, 0, err
	}

	// Select INBOX
	mbox, err := c.Select(mailbox, true)
	if err != nil {
		return 0, 0, err
	}

	// Count unread emails
	criteria := imap.NewSearchCriteria()
	criteria.WithoutFlags = []string{imap.SeenFlag}
	ids, err := c.Search(criteria)
	if err != nil {
		log.Fatal(err)
	}

	return mbox.Messages, len(ids), nil
}

func main() {
	// We can ignore errors here, because the config might be set from environment variables
	_ = godotenv.Load()

	// Intialize Config
	conf, err := config.NewConfig("./config.yaml")
	if err != nil {
		log.Fatalf("Could not load configuration: %v", err)
	}

	cfg = conf

	http.HandleFunc("/-/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Healthy"))
	})

	http.HandleFunc("/probe", func(w http.ResponseWriter, r *http.Request) {
		// Maily use the target parameter, but support mailbox as a fallback.
		target := r.URL.Query().Get("target")
		if target == "" {
			target = r.URL.Query().Get("mailbox")
			if target == "" {
				http.Error(w, "Mailbox parameter is missing", http.StatusBadRequest)
				return
			}
		}

		mailbox := target

		hostname := r.URL.Query().Get("hostname")
		if hostname == "" {
			http.Error(w, "Hostname parameter is missing", http.StatusBadRequest)
			return
		}

		username := r.URL.Query().Get("username")
		if username == "" {
			http.Error(w, "Username parameter is missing", http.StatusBadRequest)
			return
		}

		probeCountGauge := prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "probe_mailbox_count",
			Help: "Displays the count of mails found in the mailbox",
		})

		probeUnreadGauge := prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "probe_mailbox_unread_count",
			Help: "Displays the count of unread mails found in the mailbox",
		})

		registry := prometheus.NewRegistry()
		registry.MustRegister(probeCountGauge)
		registry.MustRegister(probeUnreadGauge)

		server, account, err := cfg.FindAccountInServer(hostname, username)
		if err != nil {
			log.Fatalf("Error: %v", err)
		}

		// TODO: Proper error handling
		count, unread, err := countMailsInMailbox(*server, *account, mailbox)
		if err != nil {
			log.Printf("Cound not load mailbox data: %v", err)
			http.Error(w, fmt.Sprintf("Cound not load mailbox data: %v", err), http.StatusInternalServerError)
			return
		}

		log.Printf("Load mailbox count for %s of %s on %s: %d/%d", mailbox, username, hostname, count, unread)

		probeCountGauge.Set(float64(count))
		probeUnreadGauge.Set(float64(unread))

		h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
		h.ServeHTTP(w, r)
	})

	http.Handle("/metrics", promhttp.Handler())

	log.Fatal(http.ListenAndServe(":9101", nil))
}
