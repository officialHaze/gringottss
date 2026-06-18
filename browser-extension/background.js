// background.js — Gringottss Service Worker
// Handles API communication and reads optional .env config

const DEFAULT_BASE_URL = ""; // fallback if nothing else is configured

// ─── .env Support ────────────────────────────────────────────────────────────
// If a ".env" file is dropped into the extension's own folder (next to
// manifest.json) before/when it's loaded as an unpacked extension, we read
// it as a packaged extension resource (not arbitrary disk access — the
// service worker can only see files inside its own install directory).
// A BASE_URL=... line there takes priority over the popup-configured value.

let envCache; // undefined = not loaded yet, null = no .env found
let envLoadPromise = null;

function parseEnv(text) {
  const out = {};
  text.split(/\r?\n/).forEach((rawLine) => {
    const line = rawLine.trim();
    if (!line || line.startsWith("#")) return;
    const eq = line.indexOf("=");
    if (eq === -1) return;
    const key = line.slice(0, eq).trim();
    let val = line.slice(eq + 1).trim();
    if (
      (val.startsWith('"') && val.endsWith('"')) ||
      (val.startsWith("'") && val.endsWith("'"))
    ) {
      val = val.slice(1, -1);
    }
    out[key] = val;
  });
  return out;
}

async function loadEnv() {
  if (envCache !== undefined) return envCache;
  if (!envLoadPromise) {
    envLoadPromise = (async () => {
      try {
        const res = await fetch(chrome.runtime.getURL(".env"));
        if (!res.ok) return null;
        const text = await res.text();
        return parseEnv(text);
      } catch (e) {
        return null;
      }
    })();
  }
  envCache = await envLoadPromise;
  return envCache;
}

// ─── Helpers ────────────────────────────────────────────────────────────────

async function getBaseUrl() {
  const env = await loadEnv();
  const envUrl = env?.BASE_URL || env?.SECVEIL_BASE_URL;
  if (envUrl) return envUrl.trim().replace(/\/$/, "");

  return new Promise((resolve) => {
    chrome.storage.local.get(["baseUrl"], (result) => {
      resolve(result.baseUrl || DEFAULT_BASE_URL);
    });
  });
}

async function getBaseUrlInfo() {
  const env = await loadEnv();
  const envUrl = env?.BASE_URL || env?.SECVEIL_BASE_URL;
  if (envUrl) {
    return { source: "env", baseUrl: envUrl.trim().replace(/\/$/, "") };
  }
  return new Promise((resolve) => {
    chrome.storage.local.get(["baseUrl"], (result) => {
      resolve({
        source: result.baseUrl ? "manual" : "none",
        baseUrl: result.baseUrl || ""
      });
    });
  });
}

// ─── API Calls ───────────────────────────────────────────────────────────────

async function fetchCredentials(url) {
  const base = await getBaseUrl();
  if (!base) return { data: null, error: "Base URL not configured." };

  try {
    const res = await fetch(
      `${base}/v1/credentials?url=${encodeURIComponent(url)}&nopwdencryption`,
      {
        method: "GET",
        headers: { "Content-Type": "application/json" }
      }
    );
    if (!res.ok) throw new Error(`HTTP ${res.status}`);
    const json = await res.json();
    return { data: json.data || null };
  } catch (e) {
    console.error("[Gringottss] GET credentials failed:", e);
    return { data: null, error: e.message };
  }
}

async function upsertCredentials(payload) {
  const base = await getBaseUrl();
  if (!base) return { success: false, error: "Base URL not configured." };

  try {
    const res = await fetch(`${base}/v1/credentials`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify(payload)
    });
    if (!res.ok) throw new Error(`HTTP ${res.status}`);
    return { success: true };
  } catch (e) {
    console.error("[Gringottss] POST credentials failed:", e);
    return { success: false, error: e.message };
  }
}

// ─── Message Router ──────────────────────────────────────────────────────────

chrome.runtime.onMessage.addListener((msg, sender, sendResponse) => {
  if (msg.type === "GET_CREDENTIALS") {
    fetchCredentials(msg.url).then(sendResponse);
    return true; // async
  }

  if (msg.type === "UPSERT_CREDENTIALS") {
    upsertCredentials(msg.payload).then(sendResponse);
    return true;
  }

  if (msg.type === "GET_BASE_URL_INFO") {
    getBaseUrlInfo().then(sendResponse);
    return true;
  }
});
