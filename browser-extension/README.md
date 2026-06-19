# Gringottss Browser Extension

A Chrome (MV3) browser extension for the Gringottss credential manager.

## Setup

### 1. Configure API URL

Two ways to set this — `.env` takes priority if present:

**Option A — `.env` file (recommended, no popup needed)**

1. Copy `.env.example` to `.env` in this same folder (next to `manifest.json`).
2. Set `BASE_URL=http://localhost:6230/api` inside it.
3. Load/reload the extension — the popup will show "Loaded from .env" and the field will be locked.

> Note: this is read as a file packaged inside the extension's own folder, not arbitrary disk access — it only works for `.env` files placed in the extension directory itself. On Windows, naming a file starting with a dot via Explorer's rename dialog can be finicky; easiest is to create it from a text editor's "Save As" dialog with the name `.env` (quotes included), or via a terminal (`type nul > .env` / `touch .env`). After editing `.env`, reload the extension from `chrome://extensions` for the change to take effect.

**Option B — Popup (used only if no `.env` is found)**
Open the extension popup and enter `http://localhost:6230/api`, then click **Save**. Stored in `chrome.storage.local` under `baseUrl`.

## Installation

### A. Chromium-Based Browsers (Chrome, Edge, Brave, Chromium)

Because publishing browser extensions on the Chrome Web Store requires a developer registration fee, Gringottss is distributed as an unpacked extension. Since Gringottss is a free and open source project, the extension must be loaded manually.

1. Open:

```text
chrome://extensions
```

2. Enable **Developer Mode** using the toggle in the top-right corner.
3. Click **Load unpacked**.
4. Select the `browser-extension` directory.
5. The extension will be installed.
6. Pin the extension to your browser toolbar for quick access.

---

### B. Mozilla Firefox

Firefox allows local installation of packaged extensions without requiring a store listing.

1. Open:

```text
about:addons
```

2. Click the ⚙️ settings icon beside **Manage Your Extensions**.
3. Select **Install Add-on From File...**
4. Browse to the `browser-extension` directory and select:

```text
gringottss-extension.xpi
```

5. Grant the requested permissions:
   - Access data for all websites
   - (Optional) Allow the extension in Private Windows

6. The extension will be installed.

7. Pin the extension to the Firefox toolbar for easy access.

---

### C. Other Browsers

For browsers other than Chromium-based browsers and Firefox, refer to the browser's official documentation for instructions on installing extensions locally or loading unpacked extensions.

Since installation procedures vary between browsers, Gringottss does not provide browser-specific instructions for every platform.

If you encounter any issues during installation or usage, please open an issue in the project repository with details about your browser, version, and the steps you followed.

---

### Verifying Installation

After installation:

1. Ensure the Gringottss API Server is running.
2. Open any website containing a login or registration form.
3. Gringottss widgets should automatically appear alongside supported input fields.

You are now ready to save, retrieve, generate, and autofill credentials directly from your browser.

## API Contract

### GET `/v1/credentials`

**Query:** `?url=<encoded-url>`

**Response:**

```json
{
  "data": [
    {
      "ID": "...",
      "Url": "...",
      "FormInputID": { "String": "email", "Valid": true },
      "FormInputName": { "String": "email", "Valid": true },
      "FormInputType": "email",
      "FormInputVal": "user@example.com"
    }
  ],
  "message": "Credential found."
}
```

Or `"data": null` if nothing found.

### POST `/v1/credentials` (upsert)

```json
{
  "url": "https://example.com/login",
  "formInputId": "email",
  "formInputType": "email",
  "formInputName": "email",
  "formInputVal": "user@example.com"
}
```
