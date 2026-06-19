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

Gringottss consists of two primary components:

### API Server

The API server is responsible for:

- Managing stored credentials
- Handling encryption and decryption operations
- Processing credential lookup requests
- Generating and storing credential data
- Providing a local API interface for extensions

### Browser Extension

The browser extension provides the user-facing experience and communicates directly with the locally running API server.

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

## How It Works

1. Start the Gringottss API Server. [API Server README](./api-server/README.md)
2. Install the Gringottss Browser Extension. [Browser Extension README](./browser-extension/README.md)
3. Configure the extension to communicate with the locally running API server as described in the extension documentation. [Browser Extension README](./browser-extension/README.md)
4. Browse normally.

Whenever a supported webpage contains credential input fields, Gringottss injects its widget directly into the form.

The widget allows users to:

- Save newly entered credentials
- Retrieve previously saved credentials
- Autofill usernames and passwords
- Generate strong passwords
- Save generated passwords
- Generate and instantly autofill passwords
- Manage credential interactions without leaving the current page

### Browser Support

Gringottss supports both Firefox and Chromium-based browsers, including:

- Google Chrome (Tested)
- Mozilla Firefox (Tested)
- Microsoft Edge (Untested)
- Brave Browser (Untested)

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
