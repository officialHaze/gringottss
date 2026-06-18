// popup.js — Gringottss Popup Logic

const baseUrlInput = document.getElementById("baseUrl");
const saveBtn = document.getElementById("saveBtn");
const statusEl = document.getElementById("status");
const envNoticeEl = document.getElementById("envNotice");
const currentUrlEl = document.getElementById("currentUrl");
const credDot = document.getElementById("credDot");
const credStatus = document.getElementById("credStatus");

let baseUrlSource = "none"; // "env" | "manual" | "none"

// ─── Load resolved base URL (env takes priority over manual) ──────────────

async function loadBaseUrlInfo() {
  const info = await chrome.runtime.sendMessage({ type: "GET_BASE_URL_INFO" });
  baseUrlSource = info?.source || "none";

  if (baseUrlSource === "env") {
    baseUrlInput.value = info.baseUrl;
    baseUrlInput.disabled = true;
    saveBtn.disabled = true;
    envNoticeEl.style.display = "block";
  } else {
    baseUrlInput.disabled = false;
    saveBtn.disabled = false;
    envNoticeEl.style.display = "none";
    if (info?.baseUrl) baseUrlInput.value = info.baseUrl;
  }
}

loadBaseUrlInfo();

// ─── Save base URL (only used when no .env is present) ────────────────────

saveBtn.addEventListener("click", async () => {
  const url = baseUrlInput.value.trim().replace(/\/$/, ""); // strip trailing slash

  if (!url) {
    showStatus("Please enter a valid URL.", "err");
    return;
  }

  try {
    new URL(url); // validate format
  } catch {
    showStatus("Invalid URL format.", "err");
    return;
  }

  await chrome.storage.local.set({ baseUrl: url });
  showStatus("Saved! Reload the page to apply.", "ok");
});

// Enter key saves
baseUrlInput.addEventListener("keydown", (e) => {
  if (e.key === "Enter") saveBtn.click();
});

function showStatus(msg, type) {
  statusEl.textContent = msg;
  statusEl.className = `status ${type}`;
  setTimeout(() => {
    statusEl.className = "status";
  }, 3000);
}

// ─── Current Page Info ─────────────────────────────────────────────────────

async function loadCurrentPageInfo() {
  const [tab] = await chrome.tabs.query({ active: true, currentWindow: true });
  if (!tab?.url) return;

  const url = tab.url;

  // Show URL (truncated)
  const display = url.length > 55 ? url.slice(0, 52) + "…" : url;
  currentUrlEl.textContent = display;
  currentUrlEl.title = url;

  // Skip non-http pages
  if (!url.startsWith("http")) {
    credDot.className = "dot";
    credStatus.textContent = "Not a web page";
    return;
  }

  // Check base URL is configured (env or manual), then query via background
  const info = await chrome.runtime.sendMessage({ type: "GET_BASE_URL_INFO" });
  if (!info?.baseUrl) {
    credDot.className = "dot";
    credStatus.textContent = "API URL not configured";
    return;
  }

  credStatus.textContent = "Checking credentials…";

  const result = await chrome.runtime.sendMessage({
    type: "GET_CREDENTIALS",
    url
  });

  if (result?.data && Array.isArray(result.data) && result.data.length > 0) {
    credDot.className = "dot active";
    credStatus.textContent = `${result.data.length} credential field(s) found`;
  } else if (result?.error) {
    credDot.className = "dot";
    credStatus.textContent = "Could not reach API";
  } else {
    credDot.className = "dot";
    credStatus.textContent = "No credentials for this page";
  }
}

loadCurrentPageInfo();
