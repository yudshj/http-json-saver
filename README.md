# Go HTTP Server

This is a simple HTTP server written in Go that listens on `127.0.0.1:3000`. It handles POST requests to the `/save` endpoint, saving the request body to a JSON file. The server also supports CORS (Cross-Origin Resource Sharing) with configurable origin validation.

## Features

- Handles HTTP POST requests to `/save`.
- Saves JSON request bodies to files.
- Supports CORS with optional origin validation.
- Configurable via command-line arguments.

## Usage

### Building the Server

To build the server, ensure you have Go installed and run:

```bash
go build
```

This will create an executable named `http-json-saver`.

### Running the Server

You can run the server with the following command:

```bash
./http-json-saver
```

By default, the server will have origin validation enabled.

### Command-Line Arguments

- `-enable-origin-check`: A boolean flag to enable or disable origin validation. Defaults to `true`.

Example to disable origin validation:

```bash
./http-json-saver -enable-origin-check=false
```

### Endpoints

- **POST /save**: Accepts a JSON payload and writes it to a file. The JSON must include a `"name"` field, which is used as the filename.

### Example JSON Payload

```json
{
  "name": "example",
  "data": "This is some example data."
}
```

### CORS Configuration

The server allows CORS requests from the following origins by default:

- `https://yudshj.synology.me`
- `http://127.0.0.1`

## Directory Structure

- **json_out/**: Directory where JSON files are saved.

## Error Handling

The server returns appropriate HTTP error codes for invalid requests, such as:

- `403 Forbidden` for invalid origins (if validation is enabled).
- `400 Bad Request` for invalid JSON or missing `name` field.
- `405 Method Not Allowed` for unsupported HTTP methods.

## License

This project is licensed under the MIT License.
