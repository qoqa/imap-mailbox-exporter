package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/emersion/go-imap/client"
	"github.com/joho/godotenv"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Config holds the base IMAP credentialsf
type Config struct {
	ImapUsername string
	ImapPassword string
	ImapServer   string
}

func (c *Config) validate() bool {
	return c.ImapUsername != "" && c.ImapServer != "" && c.ImapPassword != ""
}

// NewConfig creates and validated a new Config struct from environment variables
func NewConfig() (*Config, error) {
	config := &Config{
		ImapServer:   os.Getenv("IMAP_SERVER"),
		ImapUsername: os.Getenv("IMAP_USERNAME"),
		ImapPassword: os.Getenv("IMAP_PASSWORD"),
	}

	if !config.validate() {
		return nil, errors.New("not all needed configuration flags could be found")
	}

	return config, nil
}

var config *Config

func countMailsInMailbox(mailbox string) (uint32, error) {
	c, err := client.DialTLS(config.ImapServer, nil)
	if err != nil {
		return 0, err
	}

	defer c.Logout()

	// Login
	if err := c.Login(config.ImapUsername, config.ImapPassword); err != nil {
		return 0, err
	}

	// Select INBOX
	mbox, err := c.Select(mailbox, true)
	if err != nil {
		return 0, err
	}

	return mbox.Messages, nil
}

func main() {
	// We can ignore errors here, because the config might be set from environment variables
	_ = godotenv.Load()

	// Intialize Config
	conf, err := NewConfig()
	if err != nil {
		log.Fatalf("Could not load configuration: %v", err)
	}

	config = conf

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

		probeCountGauge := prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "probe_mailbox_count",
			Help: "Displays the count of mails found in the mailbox",
		})

		registry := prometheus.NewRegistry()
		registry.MustRegister(probeCountGauge)

		// TODO: Proper error handling
		count, err := countMailsInMailbox(mailbox)
		if err != nil {
			log.Printf("Cound not load mailbox data: %v", err)
			http.Error(w, fmt.Sprintf("Cound not load mailbox data: %v", err), http.StatusInternalServerError)
			return
		}

		log.Printf("Loaded mail count for mailbox %s: %d", mailbox, count)

		probeCountGauge.Set(float64(count))

		h := promhttp.HandlerFor(registry, promhttp.HandlerOpts{})
		h.ServeHTTP(w, r)
	})

	http.Handle("/metrics", promhttp.Handler())

	log.Fatal(http.ListenAndServe(":9101", nil))
}
