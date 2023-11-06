compile:
	env GOOS=linux GOARCH=arm64 go build -o ./build/raspi-kvm-exporter ./cmd/server.go