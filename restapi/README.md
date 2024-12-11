# `fabric`

...

## REST API

This REST API is part of the Fabric project, which provides various endpoints to interact with the Fabric system. The API is built using the Gin framework and includes endpoints for managing chat sessions, configurations, contexts, models, patterns, and sessions.

### Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [Endpoints](#endpoints)
  - [Chat](#chat)
  - [Configuration](#configuration)
  - [Contexts](#contexts)
  - [Models](#models)
  - [Patterns](#patterns)
  - [Sessions](#sessions)
- [Contributing](#contributing)
- [License](#license)

### Installation

To install and run the Fabric REST API, follow these steps:

1. Clone the repository:
```markdown
    ```sh
    git clone https://github.com/yourusername/fabric.git
    ```

2. Navigate to the project directory:
    ```sh
    cd fabric/restapi
    ```

3. Install the dependencies:
    ```sh
    go mod download
    ```

4. Run the application:
    ```sh
    go run main.go
    ```

### Usage

To use the Fabric REST API, you can send HTTP requests to the endpoints defined in the API. Below are examples of how to interact with the API using `curl`.

### Endpoints

#### Chat

- **GET /chat**
  - Description: Retrieve a list of chat sessions.
  - Example:
     ```sh
     curl -X GET http://localhost:8080/chat
     ```

#### Configuration

- **GET /configuration**
  - Description: Retrieve the current configuration.
  - Example:
     ```sh
     curl -X GET http://localhost:8080/configuration
     ```

#### Contexts

- **GET /contexts**
  - Description: Retrieve a list of contexts.
  - Example:
     ```sh
     curl -X GET http://localhost:8080/contexts
     ```

#### Models

- **GET /models**
  - Description: Retrieve a list of models.
  - Example:
     ```sh
     curl -X GET http://localhost:8080/models
     ```

#### Patterns

- **GET /patterns**
  - Description: Retrieve a list of patterns.
  - Example:
     ```sh
     curl -X GET http://localhost:8080/patterns
     ```

#### Sessions

- **GET /sessions**
  - Description: Retrieve a list of sessions.
  - Example:
     ```sh
     curl -X GET http://localhost:8080/sessions
     ```

### Contributing

We welcome contributions to the Fabric project. To contribute, follow these steps:

1. Fork the repository.
2. Create a new branch for your feature or bugfix.
3. Make your changes.
4. Submit a pull request.

### License

This project is licensed under the MIT License. See the [LICENSE](../LICENSE) file for details.
```
### Testing

To run the tests, execute:

```sh
go test ./...
```

### Troubleshooting

If you encounter issues, check the logs for error messages and ensure all dependencies are installed correctly. For further assistance, please open an issue on the repository.

### Contact

For any questions or support, please contact [your.email@example.com](mailto:your.email@example.com).