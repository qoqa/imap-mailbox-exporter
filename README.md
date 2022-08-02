# Imap Mailbox Exporter

> Export the amount of mails in a mailbox for use in prometheus.

## Usage

### Configuration

```dotenv
IMAP_SERVER=""
IMAP_USERNAME=""
IMAP_PASSWORD=""
```

### Probe

```txt
http://127.0.0.1:9101/probe?mailbox=INBOX
```

### Provided metrics

```txt
# HELP probe_mailbox_count Displays the count of mails found in the mailbox
# TYPE probe_mailbox_count gauge
probe_mailbox_count 0
```

## License

This project is licensed under the [MIT License](./LICENCE)

<div align="center">
    <span>&copy; 2022, jop-software Inh. Johannes Przymusinski</span>
</div>