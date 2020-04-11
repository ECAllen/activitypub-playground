#!/bin/sh

# reflex -G '*.db' -G '*.db-journal' -s -- sh -c 'clear; go run ./cmd/web'
reflex -G '*.db' -s -- sh -c 'go run ./cmd/web/server/'
