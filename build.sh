#!/bin/bash

version="$1"
mkdir -p "binairies/""$version"
name="binairies/""$version""/linx-client-v""$version""_"

GOOS=darwin GOARCH=amd64 go build -o "$name"osx-amd64

GOOS=darwin GOARCH=386 go build -o "$name"osx-386

GOOS=freebsd GOARCH=amd64 go build -o "$name"freebsd-amd64

GOOS=freebsd GOARCH=386 go build -o "$name"freebsd-386

GOOS=linux GOARCH=arm go build -o "$name"linux-arm

GOOS=linux GOARCH=amd64 go build -o "$name"linux-amd64

GOOS=linux GOARCH=386 go build -o "$name"linux-386

GOOS=windows GOARCH=amd64 go build -o "$name"windows-amd64.exe

GOOS=windows GOARCH=386 go build -o "$name"windows-386.exe
