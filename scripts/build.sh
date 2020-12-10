#!/usr/bin/env bash

VERSION=0.1.0-pre

# Delete the old dir
echo "==> Removing old directory..."
rm -rf bin/*
mkdir -p bin/windows/
mkdir -p bin/linux/
mkdir -p bin/darwin/

echo "==> Creating new binaries for windows, darwin and linux amd64..."

env GOOS=linux GOARCH=amd64 go build -o bin/linux/terraform-provider-dynatrace_v${VERSION}
env GOOS=windows GOARCH=amd64 go build -o bin/windows/terraform-provider-dynatrace_v${VERSION}
env GOOS=darwin GOARCH=amd64 go build -o bin/darwin/terraform-provider-dynatrace_v${VERSION}

cd bin/darwin
zip terraform-provider-dynatrace_${VERSION}_darwin_amd64.zip terraform-provider-dynatrace_v${VERSION}
rm terraform-provider-dynatrace_v${VERSION}
cd ../

cd linux
zip terraform-provider-dynatrace_${VERSION}_linux_amd64.zip terraform-provider-dynatrace_v${VERSION}
rm terraform-provider-dynatrace_v${VERSION}
cd ../

cd windows
zip terraform-provider-dynatrace_${VERSION}_windows_amd64.zip terraform-provider-dynatrace_v${VERSION}
rm terraform-provider-dynatrace_v${VERSION}
