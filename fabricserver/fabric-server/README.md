# Fabric Server

A Go server that provides an API for interacting with the Fabric command-line tool.

## Prerequisites

- Go 1.23.1 or later
- Fabric command-line tool installed and accessible in your system PATH

## Installation

To install the Fabric server, run:

```
go install github.com/natlamir/fabric-server/cmd/fabric-server@latest
```

This will download the source, compile it, and install the binary in your `$GOPATH/bin` directory.

## Usage

After installation, you can run the server with:

```
fabric-server
```

The server will start and listen on `http://localhost:3001`.

Download and open the [FabricPatternsUI.html](https://raw.githubusercontent.com/natlamir/fabric-server/master/FabricPatternsUI.html) in a web browser.

## API Endpoints

- GET `/api/fabric/options`: Fetches available Fabric options
- POST `/api/fabric`: Executes a Fabric command

For more detailed API documentation, please refer to the source code.

## Development

To work on the Fabric server locally:

1. Clone the repository:
   ```
   git clone https://github.com/natlamir/fabric-server.git
   ```
2. Change to the project directory:
   ```
   cd fabric-server
   ```
3. Install dependencies:
   ```
   go mod download
   ```
4. Run the server:
   ```
   go run cmd/fabric-server/main.go
   ```

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
