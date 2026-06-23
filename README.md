<p align="center">
  <img src="./assets/icon512.png" width="180" alt="Gringottss Logo">
</p>

<h1 align="center">Gringottss</h1>

<p align="center">
  A minimalistic, lightweight, free and open source local credential vault.
</p>

<p align="center">
  Save your credentials once and use them across multiple browsers.
</p>

---

## Overview

**Gringottss** is a minimalistic, lightweight, easy to setup, easy to use, free and open source credential vault designed to run entirely on your local machine.

Store your credentials once and use them across multiple browsers without relying on cloud services, subscriptions, or third-party infrastructure.

Everything runs locally and all data remains under your control.

### Why Gringottss?

- 🔒 Local-first architecture
- 🪶 Lightweight and resource efficient
- 🌐 Cross-browser credential management
- 🆓 Free and open source
- 🛡️ Privacy focused by design

No cloud synchronization.

No external servers.

No subscriptions.

Just your vault, your machine, and your keys.

---

## Architecture

Gringottss consists of three primary components:

### Gringottss Engine

The Gringottss Engine is the core service of the platform.

It is responsible for:

- Managing stored credentials
- Handling encryption and decryption operations
- Processing credential lookup requests
- Generating and storing credential data
- Managing vault persistence
- Providing a local API interface for clients

The Engine serves as the bridge between the vault database and all Gringottss clients.

---

### Gringottss CLI

The Gringottss CLI provides a command-line interface for interacting with the Gringottss Engine.

It can be used to perform various vault-related operations directly from the terminal without requiring a browser.

Examples include:

- Managing credentials
- Querying stored entries
- Administrative operations
- Automation and scripting workflows
- Other Engine-supported actions

The CLI communicates directly with the locally running Gringottss Engine.

---

### Browser Extension

The browser extension provides the browser-integrated user experience and communicates directly with the locally running Gringottss Engine.

Its responsibilities include:

- Detecting login and registration forms
- Saving credentials
- Retrieving saved credentials
- Autofilling credentials
- Generating strong passwords
- Managing credential interactions within supported browsers

The extension never communicates with external Gringottss servers because none are required for normal operation.

---

## Security

Security is a core design principle of Gringottss.

### Password Encryption

Only password values are stored encrypted.

Other credential metadata such as usernames, URLs, labels, or descriptions may be stored in their original form to enable efficient searching and retrieval.

Passwords are encrypted using **PASETO (Platform-Agnostic Security Tokens)** with a **user-defined encryption key**.

### Storage Model

- Passwords are never stored in plaintext.
- Passwords are encrypted before being written to the database.
- Encryption keys remain under user control.
- Plaintext passwords are not persisted.

### API Responses

By default, password values are returned in encrypted form.

If a client explicitly requests plaintext values, the API server temporarily decrypts the password before returning the response.

The decrypted value is never stored back to the database.

---

## How It Works? (First time usage)

1. Download the latest Gringottss release.
2. Start the Gringottss Engine. [Engine README](./engine/README.md)
3. Install and configure the Gringottss Browser Extension. [Browser Extension README](./browser-extension/README.md)
4. Use the Gringottss CLI to generate `encryption_keys.yml`. This will be used by the gringottss engine to encrypt/decrypt sensitive data. [CLI README](./cli/README.md)
5. Restart gringottss engine.
6. Browse normally or interact with the vault through the CLI.

Whenever a supported webpage contains credential input fields, Gringottss injects its widget directly into the form.

The widget allows users to:

- Save newly entered credentials
- Retrieve previously saved credentials
- Autofill usernames and passwords
- Generate strong passwords
- Save generated passwords
- Generate and instantly autofill passwords
- Manage credential interactions without leaving the current page

Alternatively, users may interact with the vault directly through the Gringottss CLI for scripting, automation, and terminal-based workflows.

### Browser Support

Gringottss supports both Firefox and Chromium-based browsers, including:

- Google Chrome (✅ Tested)
- Mozilla Firefox (✅ Tested)
- Microsoft Edge (⌛ Untested)
- Brave Browser (⌛ Untested)

The installation process differs slightly between Firefox and Chromium-based browsers and is documented in the browser extension README.

---

## Data Persistence

Gringottss is designed to be lightweight, portable, and easy to maintain.

All application data is stored locally in a **SQLite database**, eliminating the need for external database servers or additional infrastructure.

### Benefits

- 📦 Single database file
- 🪶 Lightweight and resource efficient
- 🚀 No database server setup required
- 💾 Easy backup and migration
- 🔄 Consistent storage across Gringottss releases

### Version Compatibility

Your vault lives inside a single SQLite database file.

Moving to a newer Gringottss release is as simple as copying your existing database into the new installation directory.

No exports.
No imports.
No migration headaches.

Carry the same database file forward and keep using your vault across future Gringottss releases.

### Local Ownership

Because the database resides entirely on your machine:

- Your data remains under your control.
- Backups can be performed using standard file-copy operations.
- Credentials remain accessible even without an internet connection.
- No third-party services are required for storage or synchronization.

Gringottss keeps persistence simple: one application, one database, complete ownership.

---

## Upgrading Between Releases

Gringottss is designed to preserve vault data across releases.

To migrate an existing vault to a newer release:

- Carry forward your database (`gringottss.db`)
- Carry forward your encryption keys (`encryption_keys.yml`)

These files together form your vault and allow future releases to access previously stored credentials.

Detailed migration and key management instructions can be found in the:

- [Gringottss Engine README](./engine/README.md)
- [Gringottss CLI README](./cli/README.md)

---

## Contributing

Contributions are welcome and appreciated.

Whether you're fixing bugs, improving documentation, adding features, or suggesting improvements, we'd love your help in making Gringottss better.

### Getting Started

1. Fork the repository.
2. Clone your fork locally.

```bash
git clone <your-fork-url>
```

3. Create a new branch from `main`.

```bash
git checkout main
git pull
git checkout -b feature/my-feature
```

4. Make your changes and commit them.

```bash
git add .
git commit -m "feat: add awesome feature"
```

5. Push your branch to your fork.

```bash
git push origin feature/my-feature
```

6. Open a Pull Request against the `main` branch.

### Pull Request Process

- Ensure your changes are properly tested.
- Keep pull requests focused and reasonably scoped.
- Provide a clear description of the changes.
- Be responsive to review feedback.

Once your Pull Request has been reviewed and approved, it will be merged into the project.

### Reporting Issues

If you discover a bug or have a feature request, please open an issue and provide as much relevant information as possible.

Every contribution, no matter how small, helps improve Gringottss.
