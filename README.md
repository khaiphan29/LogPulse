# LogPulse
This is a real-time log monitoring tool that allows you to track and analyze logs from various sources.

## Set up
### Kafka
***Note***: You need to have Kafka installed and running on your machine. You can download it from the [Apache Kafka website](https://kafka.apache.org/downloads) or use brew.
```bash
# Shell set up with Kafka installed by brew
export PATH="$(brew --prefix kafka)/bin:$PATH"
```
Format broker
```bash
kafka-storage.sh format --config [broker.properties] --cluster-id $(kafka-storage random-uuid)
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
make run
```

### Kafka
You have to be in project root directory to run the following command.
Moreover, you can configure `kafka_server.properties` to your desired logs path and ports.

```bash
kafka-server-start config/kafka_server.properties
```

