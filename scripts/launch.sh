#!/bin/bash

set -e
go build ../cmd/serverht.go
./serverht -clientcert ../certs/client.crt -serverkey ../certs/server.key -servercert ../certs/server.crt
