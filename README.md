# Imap Mailbox Exporter

> Export the amount of mails in a mailbox for use in prometheus.

## Usage

### Probe

```txt
http://127.0.0.1:9101/probe?target=INBOX&hostname=imap.example.com&username=me@example.com
```

### Provided metrics

```txt
# HELP probe_mailbox_count Displays the count of mails found in the mailbox
# TYPE probe_mailbox_count gauge
probe_mailbox_count 5
# HELP probe_mailbox_unread_count Displays the count of unread mails found in the mailbox
# TYPE probe_mailbox_unread_count gauge
probe_mailbox_unread_count 2
```

### Configuration

The `imap-mailbox-exporter` can be configures with a `config.yaml` file and environment variables.

```yaml
server:
- hostname: 'imap.example.com'
  port: '993'
  accounts:
    - username: 'me@example.com'
      password: 'env:E_AT_MAIL_COM_PASSWORD'
```

You can use environment variables with the `env:VARIABLE_NAME` directive in YAML.

The configuration file is expected in `./config.yaml` relative to the `imap-mailbox-exporter` binary.

### Example Usage

You can find a example docker compose configuration.

Make sure to update `examples/imap-exporter.env` with your imap credentials.

**Start the example container**

```shell
pushd examples

docker compose pull
docker compose up -d
```

## Compilation

You can compile the source-code with the `go build` command.

```bash
go build -o imap-mailbox-exporter main.go
```

Alternativly you can use [`gnu make`](https://www.gnu.org/software/make/) with the `make build` command to execute the `go build` command.


## License

This project is licensed under the [MIT License](./LICENCE)

<div align="center">
    <span>&copy; 2022, jop-software Inh. Johannes Przymusinski</span>
</div>
