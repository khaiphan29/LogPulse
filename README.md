# LogPulse
This is a real-time log monitoring tool that allows you to track and analyze logs from various sources.

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


