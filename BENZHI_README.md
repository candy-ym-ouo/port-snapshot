# Port Snapshot Task

This is a cross-platform Go command-line tool that captures active listening
ports with their protocol, PID, process name, and owning user. It can filter
ports by range and export results as CSV.

## Commands

```bash
go build ./...
go test ./...
go vet ./...
go run . -min 8000 -max 9000
go run . -csv snapshot.csv
```

## Container Build

```bash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh port-snapshot linux/arm64
./build_benzhi_docker.sh port-snapshot linux/amd64
docker run -it port-snapshot:latest
```

The image retains the Go toolchain. It downloads module dependencies during
the build, then verifies `go build ./...` inside the image.
