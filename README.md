# LogPulse
This is a real-time log monitoring tool that allows you to track and analyze logs from various sources.

## Set up
### Kafka
You have to be in project root directory to run the following command.
Moreover, you can configure `kafka_server.properties` to your desired logs path and ports.

```bash
make setup-kafka-brokers
make create-kafka-topics
```

### ElasticSearch
```
make start-es
make setup-es-indexes
```
## How to run
1. Install air for live-updating the code.
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
make start-kafka-broker
make start-es
make run
```

## Future Work
- [ ] Tokenized logs search
- [ ] More logs analysis features
