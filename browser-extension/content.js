// content.js — Gringottss Content Script (Fixed Popdown & Overlap)
// Manifest V3 Compliance (2026 Standard)

(function () {
  "use strict";

  // ─── State & Cache ────────────────────────────────────────────────────────
  let pageCredentials = [];
  const processedFields = new WeakSet();
  const fieldValueMap = new Map();

  // ─── Shared SVG Icons ──────────────────────────────────────────────────────
  const SVGS = {
    //     brand: `<svg viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
    //   <path d="M12 2L4 6v5c0 5.55 3.84 9.74 8 11 4.16-1.26 8-5.45 8-11V6l-8-4z" fill="currentColor" fill-opacity="0.2"/>
    //   <path d="M12 2L4 6v5c0 5.55 3.84 9.74 8 11 4.16-1.26 8-5.45 8-11V6l-8-4z" stroke="currentColor" stroke-width="1.5" stroke-linejoin="round"/>
    //   <path d="M9 12l2 2 4-4" stroke="currentColor" stroke-width="1.5" stroke-linecap="round" stroke-linejoin="round"/>
    // </svg>`,
    brand: `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 500 500">

  <defs>
    <linearGradient id="gold" x1="0%" y1="0%" x2="100%" y2="100%">
      <stop offset="0%" stop-color="#F8E08E"/>
      <stop offset="50%" stop-color="#D4AF37"/>
      <stop offset="100%" stop-color="#9C6B00"/>
    </linearGradient>
  </defs>

  <!-- Outer Coin -->
  <circle cx="256" cy="256" r="256" fill="#080808"/>

  <!-- Coin Edge -->
  <circle
    cx="256"
    cy="256"
    r="242"
    fill="none"
    stroke="url(#gold)"
    stroke-width="8"
    stroke-dasharray="2 6"/>

  <!-- Main Body -->
  <circle
    cx="256"
    cy="256"
    r="213"
    fill="#151515"
    stroke="#222"
    stroke-width="2"/>

  <!-- Outer Rune Markers -->
  <g fill="url(#gold)">
    <path d="M256 18 L264 34 L248 34 Z"/>
    <path d="M494 256 L478 264 L478 248 Z"/>
    <path d="M256 494 L248 478 L264 478 Z"/>
    <path d="M18 256 L34 248 L34 264 Z"/>
  </g>

  <!-- Diagonal Studs -->
  <g fill="url(#gold)">
    <circle cx="90" cy="90" r="10"/>
    <circle cx="422" cy="90" r="10"/>
    <circle cx="90" cy="422" r="10"/>
    <circle cx="422" cy="422" r="10"/>
  </g>

  <!-- Vault Wheel -->
  <g fill="url(#gold)">

    <rect x="245" y="96" width="22" height="320" rx="11"/>

    <rect x="96" y="245" width="320" height="22" rx="11"/>

    <rect
      x="245"
      y="96"
      width="22"
      height="320"
      rx="11"
      transform="rotate(45 256 256)"/>

    <rect
      x="245"
      y="96"
      width="22"
      height="320"
      rx="11"
      transform="rotate(-45 256 256)"/>
  </g>

  <!-- Forged Ring -->
  <circle
    cx="256"
    cy="256"
    r="165"
    fill="none"
    stroke="#7A5600"
    stroke-width="4"
    stroke-dasharray="3 8"/>

  <!-- Engraving Ring -->
  <circle
    cx="256"
    cy="256"
    r="115"
    fill="none"
    stroke="#7A5600"
    stroke-width="2"
    opacity="0.5"/>

  <!-- Rivets -->
  <g fill="url(#gold)">
    <circle cx="256" cy="141" r="8"/>
    <circle cx="371" cy="256" r="8"/>
    <circle cx="256" cy="371" r="8"/>
    <circle cx="141" cy="256" r="8"/>
  </g>

  <!-- Vault Core -->
  <circle
    cx="256"
    cy="256"
    r="86"
    fill="#101010"
    stroke="url(#gold)"
    stroke-width="8"/>

  <!-- Keyhole -->
  <circle cx="256" cy="225" r="16" fill="url(#gold)"/>

  <path
    d="
      M242 245
      H270
      V292
      Q270 318 256 330
      Q242 318 242 292
      Z"
    fill="url(#gold)"/>

  <circle cx="256" cy="255" r="6" fill="#111"/>

  <!-- Magic Sparks -->
  <g fill="url(#gold)" opacity="0.9">

    <path d="M145 145 L149 155 L159 159 L149 163 L145 173 L141 163 L131 159 L141 155 Z"/>

    <path d="M367 145 L371 155 L381 159 L371 163 L367 173 L363 163 L353 159 L363 155 Z"/>

    <path d="M145 367 L149 377 L159 381 L149 385 L145 395 L141 385 L131 381 L141 377 Z"/>

    <path d="M367 367 L371 377 L381 381 L371 385 L367 395 L363 385 L353 381 L363 377 Z"/>

  </g>

</svg>`,
    key: `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><circle cx="7.5" cy="15.5" r="5.5"/><path d="M21 2l-9.6 9.6"/><path d="M15.5 7.5l3 3L22 7l-3-3"/></svg>`,
    save: `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><path d="M19 21H5a2 2 0 0 1-2-2V5a2 2 0 0 1 2-2h11l5 5v11a2 2 0 0 1-2 2z"/><polyline points="17 21 17 13 7 13 7 21"/><polyline points="7 3 7 8 15 8"/></svg>`,
    dice: `<svg viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2"><rect x="2" y="2" width="20" height="20" rx="3"/><circle cx="8" cy="8" r="1.5" fill="currentColor"/><circle cx="16" cy="8" r="1.5" fill="currentColor"/><circle cx="8" cy="16" r="1.5" fill="currentColor"/><circle cx="16" cy="16" r="1.5" fill="currentColor"/><circle cx="12" cy="12" r="1.5" fill="currentColor"/></svg>`
  };

  // ─── Initialization ───────────────────────────────────────────────────────
  async function init() {
    try {
      const response = await chrome.runtime.sendMessage({
        type: "GET_CREDENTIALS",
        url: window.location.href
      });
      pageCredentials =
        response?.data && Array.isArray(response.data) ? response.data : [];
    } catch (err) {
      console.warn("Gringottss: Could not reach background worker.", err);
      pageCredentials = [];
    }

    scanAndAttach();
    observeDOMChanges();
  }

  // ─── DOM Observers & Processors ───────────────────────────────────────────
  function scanAndAttach() {
    const inputs = document.querySelectorAll(
      "input:not([type='submit']):not([type='button']):not([type='hidden']):not([type='radio']):not([type='checkbox']):not([type='file'])"
    );
    inputs.forEach(processInput);
  }

  function observeDOMChanges() {
    let debounceTimer;
    const observer = new MutationObserver(() => {
      clearTimeout(debounceTimer);
      debounceTimer = setTimeout(scanAndAttach, 150);
    });
    observer.observe(document.body, { childList: true, subtree: true });
  }

  function processInput(input) {
    // // ─── ID Presence Guard ───────────────────────────────────────────────────
    // // If the input field does not have an ID attribute, abort immediately.
    // if (
    //   (!input.id || input.id.trim() === "") &&
    //   (!input.name || input.name.trim() === "")
    // )
    //   return;

    if (processedFields.has(input)) return;

    processedFields.add(input);

    input.addEventListener("input", () =>
      fieldValueMap.set(input, input.value)
    );
    input.addEventListener("change", () =>
      fieldValueMap.set(input, input.value)
    );

    ensurePositionedParent(input);
    injectGringottssWidget(input);
  }

  function ensurePositionedParent(input) {
    const parent = input.parentElement;
    if (!parent) return;
    const pos = getComputedStyle(parent).position;
    // We explicitly avoid setting 'relative' if it disrupts flex or grid wrappers
    if (pos === "static") {
      parent.style.position = "relative";
    }
  }

  // ─── Computational Padding Offset Engine (With Original Value Cache) ───────
  function calculateRightOffset(input) {
    if (!input) return 10;

    // Cache the absolute original native padding before Gringottss ever alters it
    if (input.dataset.gringottssOriginalPadding === undefined) {
      const computedStyle = window.getComputedStyle(input);
      const rawPadding = parseFloat(computedStyle.paddingRight) || 0;
      input.dataset.gringottssOriginalPadding = rawPadding;
    }

    const originalPadding =
      parseFloat(input.dataset.gringottssOriginalPadding) || 0;
    let baseOffset = 10;

    // If the webpage's true native padding is large, a native icon is there
    if (originalPadding > 24) {
      baseOffset = originalPadding + 4;
    }

    return baseOffset;
  }

  function adjustInputPadding(input, rightOffset) {
    // Ensure we read the true original padding baseline
    const originalPadding =
      parseFloat(input.dataset.gringottssOriginalPadding) || 0;

    // We need room for our 28px wide button plus its offset placement space
    const minimumNeededPadding = rightOffset + 32;

    // Only increase the padding if the current style is less than what Gringottss needs
    if (originalPadding < minimumNeededPadding) {
      input.style.setProperty(
        "padding-right",
        `${minimumNeededPadding}px`,
        "important"
      );
    }
  }

  // ─── Core Widget Injection (Absolute Boundary Position Lock) ───────────────
  function injectGringottssWidget(input) {
    const parent = input.parentElement;
    if (!parent) return;

    const container = document.createElement("div");
    container.className = "gringottss-shadow-host";

    // Explicitly reset container layout leakages
    Object.assign(container.style, {
      position: "absolute",
      display: "inline-flex",
      alignItems: "center",
      justifyContent: "center",
      height: "28px",
      width: "28px",
      zIndex: "2147483640",
      background: "transparent none !important",
      border: "none !important",
      padding: "0 !important",
      margin: "0 !important",
      boxShadow: "none !important",
      pointerEvents: "auto"
    });

    // ─── Real-time Coordinate Lock ───
    function repositionWidget() {
      const inputRect = input.getBoundingClientRect();
      const parentRect = parent.getBoundingClientRect();

      // 1. Vertical Sync
      const verticalOffset =
        inputRect.top - parentRect.top + (inputRect.height - 28) / 2;
      container.style.top = `${verticalOffset}px`;

      // 2. Advanced Horizontal Positioning
      // We look for any sibling elements (like eye-icons, clear buttons)
      // and determine their left-most position relative to the input
      let obstructionLeft = parentRect.right;

      Array.from(parent.children).forEach((child) => {
        if (child === container || child === input) return;

        const childRect = child.getBoundingClientRect();

        // If child is to the right of the input center, it's a potential obstruction
        if (childRect.left > inputRect.left + inputRect.width / 2) {
          if (childRect.left < obstructionLeft) {
            obstructionLeft = childRect.left;
          }
        }
      });

      // Calculate how much space is available between the obstruction and the input right edge
      const spaceAvailable = obstructionLeft - inputRect.right;

      // If there is an obstruction (space is small/negative), shift left.
      // Otherwise, pin it to the right edge of the input.
      if (spaceAvailable < 35) {
        // Obstructed: position based on the obstruction's left edge
        container.style.right = `${parentRect.right - obstructionLeft + 4}px`;
      } else {
        // Clear: pin to the right edge of the input (with a small 6px buffer)
        container.style.right = `${parentRect.right - inputRect.right + 6}px`;
      }
    }

    // Run initial positioning sync
    repositionWidget();

    const shadow = container.attachShadow({ mode: "open" });

    const styleSheet = document.createElement("style");
    styleSheet.textContent = getEncapsulatedStyles();
    shadow.appendChild(styleSheet);

    const trigger = document.createElement("button");
    trigger.className = "gringottss-trigger-btn";
    trigger.setAttribute("type", "button");
    // Safely parse the SVG string into a real DOM node
    const parser = new DOMParser();
    const svgDoc = parser.parseFromString(SVGS.brand, "image/svg+xml");
    const svgElement = svgDoc.documentElement;

    // Clear any existing content and cleanly append the secure element
    trigger.textContent = "";
    trigger.appendChild(svgElement);
    shadow.appendChild(trigger);

    const matchingCred = findMatchingCredential(input);
    const isPassword = input.type === "password";

    trigger.addEventListener("click", (e) => {
      e.stopPropagation();
      e.preventDefault();
      toggleDropdown(shadow, input, matchingCred, isPassword, container);
    });

    parent.appendChild(container);

    // Dynamic Layout Padding Push
    const computedStyle = window.getComputedStyle(input);
    const currentPadding = parseFloat(computedStyle.paddingRight) || 0;
    if (currentPadding < 36) {
      input.style.setProperty(
        "padding-right",
        `${currentPadding + 32}px`,
        "important"
      );
      // Recalculate once padding is pushed to adjust placement perfectly
      repositionWidget();
    }

    // Layout tracking fallback listeners
    window.addEventListener("resize", repositionWidget);
    input.addEventListener("focus", repositionWidget);
  }

  // ─── Dropdown Renderer & Controller ────────────────────────────────────────
  function toggleDropdown(shadow, input, cred, isPassword, hostContainer) {
    // Check if a global dropdown is already open anywhere on the page
    const existingDropdown = document.getElementById(
      "gringottss-global-dropdown"
    );
    if (existingDropdown) {
      const parentInput = existingDropdown._targetInput;
      existingDropdown.remove();
      if (parentInput === input) return;
    }

    const dropdown = document.createElement("div");
    dropdown.id = "gringottss-global-dropdown";
    dropdown._targetInput = input;

    // ─── Firefox Autocomplete Shield ─────────────────────────────────────────
    // Save original attribute so we don't permanently break the website's configuration
    dropdown._originalAutocomplete = input.getAttribute("autocomplete");
    input.setAttribute("autocomplete", "off"); // Forcefully disables Firefox's native dropdown

    // Prevent interactions inside dropdown from closing the modal
    dropdown.addEventListener("click", (e) => e.stopPropagation());
    dropdown.addEventListener("mouseup", (e) => e.stopPropagation());

    // THE FOCUS LOCK: Let the mousedown happen naturally so buttons can be clicked,
    // but instantly pull focus straight back to the input field before the webpage can blur.
    dropdown.addEventListener("mousedown", (e) => {
      e.stopPropagation();
      setTimeout(() => {
        if (input && typeof input.focus === "function") {
          input.focus();
        }
      }, 0);
    });

    let sectionsHtml = "";

    if (cred) {
      sectionsHtml += `
        <div class="gringottss-menu-item" id="action-autofill">
          <span class="gringottss-icon">${SVGS.key}</span>
          <div class="gringottss-text-stack">
            <span class="gringottss-title">Fill Credential</span>
            <span class="gringottss-subtitle">Insert saved match for this field</span>
          </div>
        </div>`;
    }

    sectionsHtml += `
      <div class="gringottss-menu-item" id="action-save">
        <span class="gringottss-icon">${SVGS.save}</span>
        <div class="gringottss-text-stack">
          <span class="gringottss-title">Save Credential</span>
          <span class="gringottss-subtitle">Upsert in Gringottss vault</span>
        </div>
      </div>`;

    if (isPassword) {
      sectionsHtml += `
        <div class="gringottss-divider"></div>
        <div class="gringottss-pwgen-container">
          <div class="gringottss-pwgen-preview" id="pw-preview">--------</div>
          <div class="gringottss-range-row">
            <input type="range" id="pw-len" min="8" max="32" value="16">
            <span id="pw-len-val">16</span>
          </div>
          <div class="gringottss-btn-row">
            <button class="gringottss-btn gringottss-btn-sec" type="button" id="pw-regen">Generate</button>
            <button class="gringottss-btn gringottss-btn-pri" type="button" id="pw-use" disabled>Apply</button>
          </div>
        </div>`;
    }

    // 1. Parse your complete sectionsHtml string using the native secure DOMParser
    const htmlParser = new DOMParser();
    const parsedDoc = htmlParser.parseFromString(sectionsHtml, "text/html");

    // 2. Clear any placeholder layout markers out of the dropdown structure
    dropdown.textContent = "";

    // 3. Extract the clean, parsed DOM child nodes and safely append them
    // This moves the fully constructed elements directly into your dropdown container
    const nodes = Array.from(parsedDoc.body.childNodes);
    nodes.forEach((node) => {
      dropdown.appendChild(node);
    });

    Object.assign(dropdown.style, {
      position: "absolute",
      width: "260px",
      backgroundColor: "#0f0f0f",
      border: "1px solid #333333",
      borderRadius: "8px",
      boxShadow:
        "0 8px 24px rgba(0, 0, 0, 0.7), 0 0 0 1px rgba(214, 175, 55, 0.25)",
      padding: "10px 8px",
      display: "flex",
      flexDirection: "column",
      fontFamily:
        '-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif',
      boxSizing: "border-box",
      zIndex: "2147483647",
      pointerEvents: "auto",
      opacity: "0"
    });

    const menuStyles = document.createElement("style");
    menuStyles.textContent = `
      #gringottss-global-dropdown { z-index: 2147483647 !important; pointer-events: auto !important; }
      #gringottss-global-dropdown .gringottss-menu-item { display: flex; align-items: center; gap: 10px; padding: 8px; border-radius: 5px; cursor: pointer; transition: background 0.15s; text-align: left; }
      #gringottss-global-dropdown .gringottss-menu-item:hover { background: #1a1a1a !important; }
      #gringottss-global-dropdown .gringottss-icon { width: 14px; height: 14px; color: #d4af37; flex-shrink: 0; display: flex; }
      #gringottss-global-dropdown .gringottss-icon svg { width:100%; height:100%; }
      #gringottss-global-dropdown .gringottss-text-stack { display: flex; flex-direction: column; line-height: 1.2; }
      #gringottss-global-dropdown .gringottss-title { color: #e0e0e0; font-size: 13px; font-weight: 600; }
      #gringottss-global-dropdown .gringottss-subtitle { color: #808080; font-size: 11px; }
      #gringottss-global-dropdown .gringottss-divider { height: 1px; background: #333; margin: 6px 4px; }
      #gringottss-global-dropdown .gringottss-pwgen-container { padding: 6px 4px; display: flex; flex-direction: column; gap: 8px; }
      #gringottss-global-dropdown .gringottss-pwgen-preview { background: #0a0a0a; border: 1px solid #333; border-radius: 4px; padding: 6px; font-family: monospace; font-size: 11px; color: #666; text-align: center; word-break: break-all; user-select: all; }
      #gringottss-global-dropdown .gringottss-pwgen-preview.active { color: #d4af37; border-color: #d4af37; }
      #gringottss-global-dropdown .gringottss-range-row { display: flex; align-items: center; justify-content: space-between; gap: 8px; }
      #gringottss-global-dropdown .gringottss-range-row input { flex: 1; accent-color: #d4af37; height: 4px; cursor: pointer; }
      #gringottss-global-dropdown .gringottss-range-row span { font-size: 11px; color: #e0e0e0; min-width: 14px; }
      #gringottss-global-dropdown .gringottss-btn-row { display: flex; gap: 6px; }
      #gringottss-global-dropdown .gringottss-btn { flex: 1; border: none; padding: 5px 0; font-size: 11px; border-radius: 4px; cursor: pointer; font-family: inherit; transition: opacity 0.2s; }
      #gringottss-global-dropdown .gringottss-btn-pri { background: #d4af37; color: #000; font-weight: 700; }
      #gringottss-global-dropdown .gringottss-btn-sec { background: #333; color: #e0e0e0; }
      #gringottss-global-dropdown .gringottss-btn:disabled { opacity: 0.3; cursor: not-allowed; }
    `;
    dropdown.appendChild(menuStyles);

    document.body.appendChild(dropdown);

    const rect = hostContainer.getBoundingClientRect();
    const dropdownHeight = dropdown.offsetHeight || 215;
    const spaceBelow = window.innerHeight - rect.bottom;

    let finalTop = rect.bottom + window.scrollY + 4;
    if (spaceBelow < dropdownHeight && rect.top > dropdownHeight) {
      finalTop = rect.top + window.scrollY - dropdownHeight - 4;
    }

    dropdown.style.top = `${finalTop}px`;
    dropdown.style.left = `${rect.right + window.scrollX - 260}px`;
    dropdown.style.opacity = "1";

    // ─── STANDARD, UNINTERRUPTED EVENT BINDINGS ─────────────────────────────
    if (cred) {
      dropdown
        .querySelector("#action-autofill")
        .addEventListener("click", () => {
          applyValueToInput(input, cred.FormInputVal);
          dropdown.remove();
        });
    }

    dropdown
      .querySelector("#action-save")
      .addEventListener("click", async () => {
        dropdown.remove();
        await executeUpsert(input);
      });

    if (isPassword) {
      setupGeneratorEngine(dropdown, input);
    }

    // Global Dismiss Engine
    const handleOutsideClick = (e) => {
      const path = e.composedPath();
      if (!path.includes(hostContainer) && !path.includes(dropdown)) {
        dropdown.remove();
        document.removeEventListener("click", handleOutsideClick);
        document.removeEventListener("mouseup", handleOutsideClick);
      }
    };

    setTimeout(() => {
      document.addEventListener("click", handleOutsideClick);
      document.addEventListener("mouseup", handleOutsideClick);
    }, 20);
  }

  // ─── Password Generation Subsystem ─────────────────────────────────────────
  function setupGeneratorEngine(dropdown, input) {
    const preview = dropdown.querySelector("#pw-preview");
    const lenSlider = dropdown.querySelector("#pw-len");
    const lenVal = dropdown.querySelector("#pw-len-val");
    const regenBtn = dropdown.querySelector("#pw-regen");
    const useBtn = dropdown.querySelector("#pw-use");

    function runEngine(e) {
      if (e) e.stopPropagation();
      const len = parseInt(lenSlider.value);
      lenVal.textContent = len;

      const pool =
        "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+";
      const arr = new Uint32Array(len);
      crypto.getRandomValues(arr);
      const pw = Array.from(arr)
        .map((v) => pool[v % pool.length])
        .join("");

      preview.textContent = pw;
      preview.classList.add("active");
      useBtn.removeAttribute("disabled");
    }

    lenSlider.addEventListener("input", runEngine);
    regenBtn.addEventListener("click", runEngine);
    useBtn.addEventListener("click", (e) => {
      e.stopPropagation();
      applyValueToInput(input, preview.textContent);
      dropdown.remove();
    });
  }

  // ─── Form Automation Utilities ─────────────────────────────────────────────
  function applyValueToInput(input, val) {
    input.value = val;
    fieldValueMap.set(input, val);
    input.dispatchEvent(new Event("input", { bubbles: true }));
    input.dispatchEvent(new Event("change", { bubbles: true }));
    showToast("Applied payload ✓");
  }

  async function executeUpsert(input) {
    const val = fieldValueMap.get(input) ?? input.value;
    if (!val) {
      showToast("Cannot commit empty field.");
      return;
    }

    const result = await chrome.runtime.sendMessage({
      type: "UPSERT_CREDENTIALS",
      payload: {
        url: window.location.href,
        formInputId: input.id || "",
        formInputType: input.type || "text",
        formInputName: input.name || "",
        formInputVal: val,
        formInputXPath: getFullXPath(input)
      }
    });

    showToast(
      result?.success ? "Saved to Gringottss Vault" : "API processing failure"
    );
  }

  function findMatchingCredential(input) {
    if (!pageCredentials.length) return null;

    const id = (input.id || "").toLowerCase().trim();
    const name = (input.name || "").toLowerCase().trim();

    return (
      pageCredentials.find((cred) => {
        // 1. Map by Input ID
        if (cred.FormInputID && id === cred.FormInputID.toLowerCase().trim()) {
          return true;
        }

        // 2. Map by Input Name
        if (
          cred.FormInputName?.Valid &&
          cred.FormInputName.String &&
          name &&
          name === cred.FormInputName.String.toLowerCase().trim()
        ) {
          return true;
        }

        // 3. Map by Full XPath
        if (cred.FormInputXPath && cred.FormInputXPath.trim() !== "") {
          const resolvedElement = getElementByXPath(cred.FormInputXPath);
          if (resolvedElement === input) {
            return true;
          }
        }

        return false;
      }) || null
    );
  }

  function showToast(msg) {
    document
      .querySelectorAll(".gringottss-toast-node")
      .forEach((el) => el.remove());
    const toast = document.createElement("div");
    toast.className = "gringottss-toast-node";
    toast.textContent = msg;

    Object.assign(toast.style, {
      position: "fixed",
      bottom: "20px",
      left: "50%",
      transform: "translateX(-50%) translateY(10px)",
      zIndex: "2147483647",
      background: "#11111b",
      border: "1px solid #45475a",
      color: "#cdd6f4",
      padding: "8px 16px",
      borderRadius: "30px",
      fontSize: "12px",
      fontFamily: "sans-serif",
      boxShadow: "0 4px 16px rgba(0,0,0,0.4)",
      opacity: "0",
      transition: "all 0.2s ease"
    });

    document.body.appendChild(toast);
    setTimeout(() => {
      toast.style.opacity = "1";
      toast.style.transform = "translateX(-50%) translateY(0)";
    }, 20);
    setTimeout(() => {
      toast.style.opacity = "0";
      setTimeout(() => toast.remove(), 200);
    }, 2200);
  }

  // ─── Shadow DOM Encapsulated CSS String ────────────────────────────────────
  function getEncapsulatedStyles() {
    return `
      .gringottss-trigger-btn {
        background: transparent;
        border: none;
        outline: none;
        cursor: pointer;
        padding: 2px;
        width: 28px;
        height: 28px;
        display: flex;
        align-items: center;
        justify-content: center;
        color: #6c71c4;
        opacity: 0.9;
        transition: opacity 0.15s, color 0.15s;
      }
      .gringottss-trigger-btn:hover {
        opacity: 1;
        color: #89b4fa;
      }
      .gringottss-trigger-btn svg {
        width: 22px;
        height: 22px;
      }
      .gringottss-menu-item {
        display: flex;
        align-items: center;
        gap: 10px;
        padding: 8px;
        border-radius: 5px;
        cursor: pointer;
        transition: background 0.1s;
      }
      .gringottss-menu-item:hover {
        background: #313244;
      }
      .gringottss-icon {
        width: 14px;
        height: 14px;
        color: #cba6f7;
        flex-shrink: 0;
        display: flex;
      }
      .gringottss-icon svg { width:100%; height:100%; }
      .gringottss-text-stack {
        display: flex;
        flex-direction: column;
        line-height: 1.2;
      }
      .gringottss-title {
        color: #cdd6f4;
        font-size: 12px;
        font-weight: 500;
      }
      .gringottss-subtitle {
        color: #a6adc8;
        font-size: 10px;
      }
      .gringottss-divider {
        height: 1px;
        background: #45475a;
        margin: 6px 4px;
      }
      .gringottss-pwgen-container {
        padding: 6px 4px;
        display: flex;
        flex-direction: column;
        gap: 8px;
      }
      .gringottss-pwgen-preview {
        background: #11111b;
        border: 1px solid #313244;
        border-radius: 4px;
        padding: 6px;
        font-family: monospace;
        font-size: 11px;
        color: #585b70;
        text-align: center;
        word-break: break-all;
        user-select: all;
      }
      .gringottss-pwgen-preview.active {
        color: #a6e3a1;
      }
      .gringottss-range-row {
        display: flex;
        align-items: center;
        justify-content: space-between;
        gap: 8px;
      }
      .gringottss-range-row input {
        flex: 1;
        accent-color: #cba6f7;
        height: 4px;
        cursor: pointer;
      }
      .gringottss-range-row span {
        font-size: 11px;
        color: #cdd6f4;
        min-width: 14px;
      }
      .gringottss-btn-row {
        display: flex;
        gap: 6px;
      }
      .gringottss-btn {
        flex: 1;
        border: none;
        padding: 5px 0;
        font-size: 11px;
        border-radius: 4px;
        cursor: pointer;
        font-family: inherit;
      }
      .gringottss-btn-pri { background: #89b4fa; color: #11111b; font-weight: 500; }
      .gringottss-btn-sec { background: #313244; color: #cdd6f4; }
      .gringottss-btn:disabled { opacity: 0.4; cursor: not-allowed; }
    `;
  }

  // ─── Execution ────────────────────────────────────────────────────────────
  init();
})();

// ─── XPath Related Helpers ────────────────────────────────────────────────
function getFullXPath(element) {
  if (element.tagName.toLowerCase() === "html") return "/html[1]";
  if (element === document.body) return "/html[1]/body[1]";

  let ix = 0;
  const siblings = element.parentNode ? element.parentNode.childNodes : [];
  for (let i = 0; i < siblings.length; i++) {
    const sibling = siblings[i];
    if (sibling === element) {
      return `${getFullXPath(element.parentNode)}/${element.tagName.toLowerCase()}[${ix + 1}]`;
    }
    if (
      sibling.nodeType === Node.ELEMENT_NODE &&
      sibling.tagName === element.tagName
    ) {
      ix++;
    }
  }
  return "";
}

function getElementByXPath(xpath) {
  try {
    return document.evaluate(
      xpath,
      document,
      null,
      XPathResult.FIRST_ORDERED_NODE_TYPE,
      null
    ).singleNodeValue;
  } catch (e) {
    return null;
  }
}
