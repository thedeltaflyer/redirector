# Redirector

The **Redirector** project is a URL redirection service designed to take long URLs and create shorter, more manageable links. It also includes support for generating QR codes and accessing data in different formats like JSON or plain text. The service is built using **Go (Golang)** and leverages the **Gin** framework for handling HTTP requests.

**Redirector** was built for use on the [lnk.now](https://lnk.now) URL shortener, so its usefulness in other contexts may be limited.

---

## Features

1. **URL Redirection**
    - Shorten long URLs with customizable keys.
    - Supports automatic key generation.

2. **Formats**
    - Access the URL data in multiple formats:
        - Plain text (`/text`)
        - JSON (`/json`)
        - QR Code (`/qr`)

3. **QR Code Generation**
    - Generate QR codes for URLs with configurable parameters:
        - Size
        - Error correction levels
        - Custom background and foreground colors
        - Border options.

4. **API Authentication**
    - Protect API endpoints with token-based authentication middleware.

5. **Health Monitoring**
    - Expose a simple health check endpoint.

6. **Database**
    - Uses a local **BoltDB** database to store URL mappings.

---

## Requirements

Ensure you have the following installed:

- **Go** version 1.23.4 or higher *(for local execution, not required for Docker)*
- **Docker** installed and running *(for container-based deployment)*
- A terminal or command-line client for running the project.

---

### Dependencies

The project relies on these Go packages (defined in `go.mod`):

- `github.com/gin-gonic/gin`: HTTP web framework.
- `go.etcd.io/bbolt`: Embedded key-value database.
- `github.com/skip2/go-qrcode`: QR code generation library.
- `github.com/spf13/pflag`: Command-line flag parsing.
- `github.com/sirupsen/logrus`: For structured application logging.
- `github.com/matoous/go-nanoid/v2`: For generating unique IDs for shortened URLs.

---

## Installation and Running the Application

### Running Locally

1. Clone the repository:

   ```bash
   git clone https://github.com/yourusername/redirector.git
   cd redirector
   ```

2. Install dependencies:

   ```bash
   go mod tidy
   ```

3. **Run in development mode:**
   ```bash
   go run main.go --debug
   ```

4. **Run in production mode:**
   ```bash
   go run main.go --bind=":8080" --db="./data/redirect.db"
   ```

Default values:
- Bind Address: `:8080`
- Database Path: `./db/db.bolt`

To view the full list of supported flags, use:
```bash
go run main.go --help
```

---

### Running with Docker

The application is containerized, and you can run it using Docker. A `Dockerfile` is included for your convenience.

#### Steps to Run with Docker:

1. Build the Docker image:

   ```bash
   docker build -t redirector .
   ```

2. Run the Docker container and mount the database directory:

   ```bash
   docker run -d \
     --name redirector \
     -p 8080:8080 \
     -v $(pwd)/db:/app/db \
     redirector
   ```

    - `-p 8080:8080`: Maps the container's exposed port (8080) to your local machine's port (8080).
    - `-v $(pwd)/db:/app/db`: Mounts the `db` directory from your local machine into the container. This ensures that the BoltDB database persists even when the container is stopped or removed.

3. Access the service:

    - Open your browser or make API requests to `http://localhost:8080`.

4. Stopping and removing the container:

   To stop the container:
   ```bash
   docker stop redirector
   ```

   To remove the container:
   ```bash
   docker rm redirector
   ```

#### Example Using Docker Compose

You can also use Docker Compose to manage the container.

Run the service:

```bash
docker compose up -d
```

Stop the service:

```bash
docker compose down
```

---

## Usage

The application exposes endpoints for health checks, redirection, and API interaction.

### Endpoints

1. **Health Check:**
   ```http
   GET /health
   ```

   Returns status information about the server.

2. **Redirection:**
   ```http
   GET /:key[/format]
   ```

    - `:key`: The key for the shortened URL.
    - `/json`, `/text`, `/qr`: Optional format specifiers.

   Example:
    - JSON: `/abc123/json` → `{ "key": "abc123", "url": "https://example.com" }`
    - QR Code: `/abc123/qr` (returns a QR PNG image).

3. **Shorten a URL (Requires Authentication):**
   ```http
   POST /
   POST /:key
   PUT /:key
   ```

    - Use `POST /` for a key-less request (key is auto-generated).
    - Use `POST /:key` to define a custom key.
    - Use `PUT /:key` to update an existing redirect. Note: this is a separate call to ensure that the replacement is intentional.

   **Request Body Example:**
   ```json
   {
     "url": "https://example.com/"
   }
   ```
   Returns:
   ```json
   {
     "status": "success",
     "redirect": {
       "key": "abc123",
       "url": "https://example.com/"
     }
   }
   ```
   
   Note: Authentication is provided via a `Bearer` Authentication token. This token must be added directly to the DB in the `api_keys` bucket.

---

## Custom QR Configurations

When using the `/qr` endpoint, you can configure the QR code by supplying query parameters:

- **size**: QR image size (default: 256).
- **level**: Error-correction level (L, M, H, or B; default: M).
- **bg_color**: Background color in hex (default: `#FFFFFFFF` (white)). 3, 4, 6, and 8 character RGB(A) values are supported.
- **fg_color**: Foreground color in hex (default: `#000000FF` (black)). 3, 4, 6, and 8 character RGB(A) values are supported.
- **border**: Boolean for enabling/disabling border (default: true).

### Example Request
```http
GET /abc123/qr?size=300&level=H&bg_color=#ffffff&fg_color=#000000&border=true
```

---

## File Structure

```plaintext
├── controllers/         # Application controllers (e.g., redirector logic)
├── database/            # Database initialization and configuration
├── helpers/             # Helper functions (e.g., QR code utilities)
├── logging/             # Logging utilities with logrus
├── models/              # Data models (e.g., URL mappings)
├── middleware/          # Authentication and other middleware
├── server/              # HTTP server setup
├── go.mod               # Dependency management
├── Dockerfile           # Docker container specification
├── main.go              # Entry point for the application
```

---

## Contributing

We welcome contributions to improve this project! Please follow these steps:

1. Fork the repository.
2. Create a new branch (`feature/your-feature`).
3. Commit your updates.
4. Open a pull request.

Please make sure your code passes all tests before submission.

---

## Testing

To test the application, write unit tests or run the HTTP endpoints with appropriate tools like **curl**, **Postman**, or automated testing suites.

Run all tests:
```bash
go test ./...
```

---

## License

This project is licensed under the MIT License. See the `LICENSE` file for details.

---

## Acknowledgments

- Framework: [Gin Web Framework](https://github.com/gin-gonic/gin).
- Database: [BoltDB](https://github.com/etcd-io/bbolt).
- QR Code Generation: [Skip2 QR Code](https://github.com/skip2/go-qrcode).
- Logging: [Logrus](https://github.com/sirupsen/logrus).

For any issues or feature requests, feel free to create an issue in this repository.

---