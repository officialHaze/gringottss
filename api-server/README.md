# Gringottss API Server

The **Gringottss API Server** is the backbone of the Gringottss ecosystem.

It acts as the bridge between the browser extension and the local vault database, handling credential storage, retrieval, encryption, decryption and all other vault operations.

The browser extension communicates exclusively with this server, which in turn manages the local SQLite database.

---

## Prerequisites

Before starting the API server, ensure the following dependencies are installed:

- Go **1.25 or newer**
- Docker (recommended for running the server in the background)

Docker is recommended because it provides a consistent runtime environment regardless of the underlying operating system or architecture.

---

## Initial Setup

### 1. Configure Environment Variables

Copy the example environment file:

```bash
cp .env.example .env
```

Review and update the values in `.env` as required.

---

### 2. Configure Application Settings

Copy the example settings file:

```bash
cp ./settings/settings.conf.example.yaml ./settings/settings.conf.yaml
```

You may modify any configuration values that are **not explicitly marked as**:

```text
DO NOT EDIT
```

or

```text
DO NOT CHANGE
```

> **Important**
>
> The `env_file_name` field must not be modified.

---

## Running the Server

To start the API server directly:

```bash
go run gringottss.go start_server
```

This command will:

- Create the SQLite database at:

```text
./data/gringottss.db
```

- Execute all required database migrations.
- Start the API server.

---

## Running with Docker (Recommended)

Running through Docker allows the API server to operate in the background and automatically restart when configured.

### Step 1: Build the Binary

Docker uses the compiled binary to run the application.

Build the binary using:

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./bin/gringottss
```

---

### Step 2: Prepare the Database

The database file must exist before starting the Docker container.

#### Option A: Existing Database

If migrating from another Gringottss release, copy your existing database file in:

```text
./data/gringottss.db
```

into the current release.

#### Option B: New Installation

Create the database and execute migrations:

```bash
go run gringottss.go build_db
```

This command will:

- Create `./data/gringottss.db`
- Execute all required migrations

---

### Step 3: Start Docker

A ready-to-use Docker Compose file is included with the project.

Start the server:

```bash
docker compose up -d
```

This command will:

- Build the Docker image
- Create the container
- Start the API server in the background

---

## Port Configuration

If you change the server port in the configuration files, ensure the same port is updated in:

- `Dockerfile`
- `docker-compose.yml`

Failure to keep these values synchronized may prevent the container from exposing the correct service port.

---

## Persistent Data

All vault data is stored in:

```text
./data/gringottss.db
```

This SQLite database can be copied between Gringottss releases to retain credentials, settings, and stored vault data.

To migrate to a newer release:

1. Stop the current server.
2. Copy `gringottss.db` into the new release's `./data/` directory.
3. Start the new version.

No export or import process is required.

---

## Customizing Docker

The provided `docker-compose.yml` serves as a sensible default configuration.

You may customize:

- Port mappings
- Container names
- Restart policies
- Volume mappings
- Network configuration

to suit your deployment requirements.

---

## Accessing the API Server

Once the server is running, it can be accessed locally at:

```text
http://localhost:6230
```

If you have configured a different port, replace `6230` with your configured value:

```text
http://localhost:<PORT>
```

For example:

```text
http://localhost:8080
```

The browser extension communicates with this local endpoint to perform all vault operations.
