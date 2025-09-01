# LogPulse
This is a low latency log monitoring tool that allows you to track and analyze logs from various sources.

## Set up
### Create Kafka topics and ES indexes

```bash
make start-containers
make provision
```

```
## How to run
1. Install air for code live-updating.
```bash
go install github.com/air-verse/air@latest
```
Config ~/.zshrc to use Go binaries
```bash
export PATH=$(go env GOPATH)/bin:$PATH
```
2. Run the service
***Note***: Change the ports on [Makefile] to your desired ports.
```bash
make dev
```

## Future Work
- [ ] Tokenized logs search
- [ ] More logs analysis features
