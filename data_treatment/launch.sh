#!/bin/bash

if [ -z "$1" ]; then
  echo "Usage: $0 <connection_string>"
  exit 1
fi

cd code
go run main.go $1
cd ..