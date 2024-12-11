#!/bin/bash

CGO_ENABLED=0 GOOS=linux GOARCH=arm64 \
go build -tags lambda.norpc -o bootstrap ./cmd
zip function.zip bootstrap
rm bootstrap
