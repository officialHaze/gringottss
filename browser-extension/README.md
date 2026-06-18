# Gringottss Browser Extension

A Chrome (MV3) browser extension for the Gringottss credential manager.

## Setup

### 1. Configure API URL

Two ways to set this — `.env` takes priority if present:

**Option A — `.env` file (recommended, no popup needed)**

1. Copy `.env.example` to `.env` in this same folder (next to `manifest.json`).
2. Set `BASE_URL=https://your-api.example.com` inside it.
3. Load/reload the extension — the popup will show "Loaded from .env" and the field will be locked.

> Note: this is read as a file packaged inside the extension's own folder, not arbitrary disk access — it only works for `.env` files placed in the extension directory itself. On Windows, naming a file starting with a dot via Explorer's rename dialog can be finicky; easiest is to create it from a text editor's "Save As" dialog with the name `.env` (quotes included), or via a terminal (`type nul > .env` / `touch .env`). After editing `.env`, reload the extension from `chrome://extensions` for the change to take effect.

**Option B — Popup (used only if no `.env` is found)**
Open the extension popup and enter your API base URL, then click **Save**. Stored in `chrome.storage.local` under `baseUrl`.

### 2. Load in Chrome (Developer Mode)

1. Go to `chrome://extensions`
2. Enable **Developer mode** (top right toggle)
3. Click **Load unpacked**
4. Select this `browser-extension` folder

---

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
