# Gringottss Engine

The **Gringottss Engine** is the backbone of the Gringottss ecosystem.

It acts as the bridge between all other components and the local vault database, handling credential storage, retrieval, encryption, decryption and all other vault operations by exposing its API via REST server.

The browser extension communicates exclusively with the engine, which in turn manages the local SQLite database.

---

## Prerequisites

Before starting the Engine, ensure the following dependencies are installed:

- Go **1.25 or newer**
- Docker (recommended for running the engine in the background)

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

To start the Engine directly:

```bash
go run gringottss.go start_server
```

This command will:

- Create the SQLite database at:

```text
./data/gringottss.db
```

- Execute all required database migrations.
- Start the Engine.

---

## Running with Docker (Recommended)

Running through Docker allows the Engine to operate in the background and automatically restart when configured.

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

Start the engine:

```bash
docker compose up -d
```

This command will:

- Build the Docker image
- Create the container
- Start the Engine in the background

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

This SQLite database contains your credentials, settings, and other vault-related data.

---

## Upgrading Between Releases

Gringottss is designed with long-term portability in mind.

Your vault is not tied to a specific release, installation, operating system, or machine. By carrying forward your encryption keys and database, you can continue using your vault across future Gringottss releases without losing access to your stored credentials.

### Encryption Keys

Password encryption in Gringottss is driven by user-generated encryption keys.

For new installations, encryption keys should be generated using the **Gringottss CLI**.

Detailed instructions are available in the [CLI documentation](./cli/README.md).

Once generated, the resulting:

```text
encryption_keys.yml
```

file can be reused across future releases.

To preserve access to previously encrypted passwords:

1. Copy `encryption_keys.yml` from the existing installation.
2. Place it in the root directory of the new Gringottss Engine installation.
3. Start the Engine normally.

As long as the same encryption keys are preserved, previously stored passwords remain accessible.

### Vault Database

Credential data is stored in the Gringottss SQLite database.

The Gringottss CLI provides tooling for carrying vault data between releases and systems.

Detailed instructions are available in [CLI Documentation](./cli/README.md).

To migrate existing data:

1. Move the old database inside `<PROJECT_ROOT>/engine/migrate`.
2. Use [Gringottss CLI](./cli/README.md) to run the migration.
3. (Optional) Restart Gringottss Engine.

### What Should Be Preserved?

For a seamless upgrade experience, preserve:

```text
encryption_keys.yml
```

and

```text
gringottss.db
```

These two files together represent your vault.

- `gringottss.db` contains your stored credential data.
- `encryption_keys.yml` contains the keys required to decrypt protected password values.

Losing either may result in partial or complete loss of access to encrypted data.

### Recommended Upgrade Process

1. Stop the currently running Gringottss Engine.
2. Back up `gringottss.db`.
3. Back up `encryption_keys.yml`.
4. Install the new Gringottss release.
5. Copy both files into the new installation.
6. Start the new release.

Your vault should be immediately available with no additional configuration required.

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

## Accessing the Engine

Once the engine is running, it can be accessed locally at:

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
