# Chirpy

Chirpy is a simple Go web application that interacts with a PostgreSQL database, exposes JSON APIs for user and "chirp" management, and provides a minimal admin dashboard and health check. The project includes middleware, environment variable support, and platform-dependent admin functionality.

***

## Features

- **RESTful API** for managing users and chirps (short messages)
- **Middleware** to track and display file server hits
- **Admin dashboard** showing metrics and support for reset operations (dev only)
- **Health check endpoint**
- **Pluggable platform mode** via environment variable
- **Environment variable configuration** (uses `.env`)
- **Polka webhooks support**
- **PostgreSQL database** integration

***

## Getting Started

### Prerequisites

- Go 1.18 or newer
- PostgreSQL database
- (Optional) [direnv](https://github.com/direnv/direnv) or another tool for loading `.env` files

### Installation

1. **Clone the repository**:

    ```bash
    git clone https://github.com/Sheikh-Fahad-Ahmed/chirpy.git
    cd chirpy
    ```

2. **Install dependencies**:

    ```bash
    go mod tidy
    ```

3. **Create a `.env` file**:

    ```
    DB_URL=postgres://username:password@localhost:5432/chirpydb?sslmode=disable
    PLATFORM=dev
    SECRET_KEY=your_secret_key
    POLKA_KEY=your_polka_api_key
    ```

    Replace the values with your actual database credentials and secret keys.

4. **Set up the database schema**

    The database expects to use the code in `internal/database`. Ensure your schema is created as required by that package.

5. **Run the server**:

    ```bash
    go run main.go
    ```

    The server will start on port `8080`.

***

## API Endpoints

| Method | Path                             | Description                               |
|--------|----------------------------------|-------------------------------------------|
| GET    | `/admin/metrics`                 | View dashboard with file server hits      |
| POST   | `/admin/reset`                   | Reset DB/users and metrics (**dev only**) |
| GET    | `/api/healthz`                   | Health check ("OK")                      |
| GET    | `/api/chirps`                    | List all chirps                           |
| GET    | `/api/chirps/{chirpID}`          | Get a specific chirp                      |
| POST   | `/api/chirps`                    | Create a new chirp                        |
| DELETE | `/api/chirps/{chirpID}`          | Delete a chirp                            |
| POST   | `/api/users`                     | Create a new user                         |
| PUT    | `/api/users`                     | Update current user                       |
| POST   | `/api/login`                     | User login                                |
| POST   | `/api/refresh`                   | Refresh authentication token              |
| POST   | `/api/revoke`                    | Revoke authentication token               |
| POST   | `/api/polka/webhooks`            | Polka payment webhook handler             |
| GET    | `/app/*`                         | Static file server (middleware-tracked)   |

***

## Development

- Modify environment variables as needed in your `.env` file.
- Code entry point is `main.go`.
- Database query methods are defined under `internal/database`.

***

## Security

- Admin reset endpoint is restricted to the "dev" platform (set `PLATFORM=dev` in `.env` for local development).
- Secret keys should **never** be committed to source control.
- Use HTTPS and proper environment variable management in production.

***

## License

MIT or as specified in repository.

***

## Credits

- [Sheikh-Fahad-Ahmed/chirpy](https://github.com/Sheikh-Fahad-Ahmed/chirpy)
- Uses: [github.com/joho/godotenv], [github.com/lib/pq]

***

> For questions or contributions, open an issue or PR at the project repository.