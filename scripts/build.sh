#!/bin/bash

cd cmd
env GOOS=linux GOARCH=arm64 go build -o ../build/raspi-kvm-exporter