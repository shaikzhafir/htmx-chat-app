# htmx-chat-app

## Description

This is simple chat server application where users can join and converse with each other

## What you need

1. Go installed on your machine

## How to run

### Locally

1. Create a .env file and set DEV to be true (same as in sample.env). Run `source .env` to set the variable
2. Run `go mod tidy` to install the required go packages
3. Run `go run main.go`
4. Check `localhost:8080` in your browser

### In remote server (TODO)

Basically the idea is to compile the go program, run it as a systemd service in a unix environment and point a domain to the remote server so others can use it.

if its still up, you can visit
https://cloud.shaikzhafir.com to see a live version of this
