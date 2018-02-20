#!/usr/bin/env bash

# This will need to be re-run in a year when this cert expires.
openssl req -newkey rsa:4096 -nodes -keyout ./assets/ssl/key -x509 -days 365 -out ./assets/ssl/cert
openssl req -new -key ./assets/ssl/key -out ./assets/ssl/csr
