# Go HTTP Server

This is a simple HTTP server implemented in Go. It listens for HTTP POST requests and saves the request body into a JSON file. The filename is determined by the "name" field in the JSON body.

## Features

- Listens on `127.0.0.1:3000`
- Supports CORS for specific origins
- Saves POST request bodies to JSON files
- Handles CORS preflight requests

## Prerequisites

- [Go](https://golang.org/dl/) installed on your system

## Installation

1. Clone the repository or download the `server.go` file.

2. Navigate to the directory containing `server.go`.

3. Build the executable:

```bash
go build server.go
```
