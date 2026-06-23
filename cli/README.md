# Gringottss CLI

**Gringottss CLI** is the administrative and maintenance companion for the Gringottss ecosystem.

It communicates with the **Gringottss Engine** and provides utilities for vault initialization, encryption key management, database creation, and migration of data between Gringottss releases.

While most users interact with Gringottss through the browser extension, the CLI is responsible for the behind-the-scenes operations that keep your vault portable, secure, and upgradeable.

---

## Overview

The Gringottss CLI can be used to:

- Generate encryption keys
- Initialize vault databases
- Migrate data from older Gringottss releases
- Perform setup and maintenance operations
- Preserve vault compatibility across releases

---

## Prerequisites

Before using the CLI, ensure you have:

- Go **1.25 or newer**
- A valid Gringottss Engine installation

---

## Configuration

Before using the CLI, create an environment configuration file.

Copy:

```bash
cp .env.example .env
```

The `.env.example` file contains:

```env
GRINGOTTSS_ENGINE_HOST=localhost
GRINGOTTSS_ENGINE_PORT=6230
```

Update these values to match your Gringottss Engine configuration.

### Configuration Variables

| Variable                 | Description                                                    |
| ------------------------ | -------------------------------------------------------------- |
| `GRINGOTTSS_ENGINE_HOST` | Hostname or IP address where the Gringottss Engine is running. |
| `GRINGOTTSS_ENGINE_PORT` | Port exposed by the Gringottss Engine.                         |

### Example

If the Engine is running locally using the default configuration:

```env
GRINGOTTSS_ENGINE_HOST=localhost
GRINGOTTSS_ENGINE_PORT=6230
```

If the Engine is running on a different host or port, update the values accordingly.

> **Important**
>
> The CLI communicates directly with the Gringottss Engine. If the host or port values are incorrect, CLI operations may fail to connect to the Engine.

---

## Running the CLI

View all available commands:

```bash
go run cli.go --help
```

Example output:

```text
Gringottss CLI interacts with the Gringottss Engine to perform several operations.

Usage:
  gringottss-cli [command]

Available Commands:
  help        Help about any command
  keygen      Generate keys used by gringottss engine. Engine restart required.
  makedb      Call Gringottss Engine to build the DB once.
  migrate     Migrate old Gringottss DB into current.
```

---

## First-Time Setup

For a fresh Gringottss installation:

### 1. Generate Encryption Keys

```bash
go run cli.go keygen --type encryption
```

### 2. Create the Database

```bash
go run cli.go makedb
```

### 3. Start the Gringottss Engine

Refer to the Engine documentation:

```text
../gringottss-engine/README.md
```

### 4. Install the Browser Extension

Refer to the Browser Extension documentation:

```text
../browser-extension/README.md
```

---

## Commands

### keygen

Generate encryption keys used by the Gringottss Engine.

```bash
go run cli.go keygen --type encryption
```

#### Available Flags

| Flag         | Description                                                |
| ------------ | ---------------------------------------------------------- |
| `--type`     | Type of keys to generate. Currently supports `encryption`. |
| `-h, --help` | Display command help.                                      |

#### Generated Files

```text
encryption_keys.yml
```

The generated file contains the encryption keys used by the Gringottss Engine for password encryption and decryption.

#### Important Notes

- Generate keys before storing credentials.
- Store the generated file securely.
- The Gringottss Engine must be restarted after generating new keys.
- Losing the encryption keys may result in loss of access to encrypted password data.

View detailed help:

```bash
go run cli.go keygen --help
```

---

### makedb

Create and initialize the Gringottss database.

```bash
go run cli.go makedb
```

This command instructs the Gringottss Engine to:

- Create the SQLite database
- Build required tables
- Execute necessary migrations
- Prepare the vault for use

This operation is generally only required during first-time setup.

---

### migrate

Migrate data from an older Gringottss database into the current release.

#### Prepare the Source Database

Place the old database inside:

```text
<PROJECT_ROOT>/engine/migrate/
```

Example:

```text
<PROJECT_ROOT>/engine/migrate/old.db
```

#### Run Migration

```bash
go run cli.go migrate --dbname old.db
```

#### Available Flags

| Flag         | Description                                                       |
| ------------ | ----------------------------------------------------------------- |
| `--dbname`   | Name of the old database file located in the migration directory. |
| `-h, --help` | Display command help.                                             |

#### Example

```bash
go run cli.go migrate --dbname old.db
```

The migration process transfers all required Gringottss tables from the older database into the current release database.

This allows users to:

- Upgrade between releases
- Move vaults between systems
- Preserve stored credentials
- Maintain long-term vault continuity

View detailed help:

```bash
go run cli.go migrate --help
```

---

## Vault Portability

One of the core goals of Gringottss is portability.

Your vault is represented by two important components:

### Database

```text
gringottss.db
```

Contains:

- Stored credentials
- User data
- Vault metadata

### Encryption Keys

```text
encryption_keys.yml
```

Contains:

- Password encryption keys
- Password decryption keys

Both files should be preserved when upgrading releases or moving to another machine.

---

## Upgrading Between Releases

To preserve access to your vault across Gringottss releases:

### Step 1

Back up:

```text
gringottss.db
```

### Step 2

Back up:

```text
encryption_keys.yml
```

### Step 3

Install the newer Gringottss release.

### Step 4

Copy both files into the new installation.

### Step 5

Run migration commands if required.

### Step 6

Start the Gringottss Engine.

Your vault data should now be available in the new release.

---

## Best Practices

- Keep backups of both `gringottss.db` and `encryption_keys.yml`.
- Never share your encryption keys publicly.
- Always back up your vault before running migration operations.
- Restart the Engine after generating new encryption keys.
- Test migrations using backup copies whenever possible.

---

## Getting Help

View all commands:

```bash
go run cli.go --help
```

View help for a specific command:

```bash
go run cli.go <command> --help
```

Examples:

```bash
go run cli.go keygen --help
```

```bash
go run cli.go migrate --help
```

---

## Philosophy

Gringottss is built around a simple idea:

**Your vault belongs to you.**

No cloud lock-in.

No vendor dependency.

No subscriptions.

Just your data, your keys, and your control.
