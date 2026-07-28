console.log("EDITOR LOADED");

let currentFile = "";
let zoom = 1;

window.onerror = function (msg, url, line) {
    console.log("ERROR:", msg, "Line:", line);
};

const editor = grapesjs.init({
    container: "#gjs",
    fromElement: false,
    height: "100%",
    storageManager: false,
    telemetry: false,
    dragMode: "absolute",

    selectorManager: {
        componentFirst: true
    }
});

/* ********* MULTI SELECT FOR ALIGNMENT ************ */
// ---------------------
// Multi Selection
// ---------------------

let suppressSelectHandler = false;

let selectedComps = [];

editor.on("component:selected", comp => {
    if (suppressSelectHandler) return;

    const ev = window.event;

    // Ctrl pressed ? add/remove from group
    if (ev && ev.ctrlKey) {
        const idx = selectedComps.indexOf(comp);

        if (idx >= 0) {
            selectedComps.splice(idx, 1);

            setSelectionOutline(comp, false);
        } else {
            selectedComps.push(comp);

            setSelectionOutline(comp, true);
        }
    } else {
        clearMultiSelection();

        selectedComps = [comp];

        setSelectionOutline(comp, true);
    }

    syncFields();

    // --- merged from the second component:selected listener ---
    const el = comp.getEl();
    if (!el) return;

    const style = comp.getStyle() || {};
    const computed = window.getComputedStyle(el);

    fieldX.value = parseInt(style.left || computed.left) || 0;
    fieldY.value = parseInt(style.top || computed.top) || 0;
    fieldW.value = parseInt(style.width || computed.width) || 0;
    fieldH.value = parseInt(style.height || computed.height) || 0;

    const wrapper = editor.DomComponents.getWrapper();

    wrapper.find("*").forEach(fieldComp => {
        const attrs = fieldComp.getAttributes();

        // Convert your field placeholders
        if (attrs["data-field-id"]) {
            fieldComp.set({
                editable: false, // manual dblclick-to-edit handles this uniformly instead
                selectable: true,
                hoverable: true,
                draggable: true,
                copyable: true,
                removable: true,
                badgable: true
            });
        }
    });
});

function clearMultiSelection() {
    selectedComps.forEach(c => {
        setSelectionOutline(c, false);
        c.set("draggable", true); // NEW: defensive reset
    });
    selectedComps = [];
}

function syncFields() {
    if (!selectedComps.length) return;

    const ref = selectedComps[0];

    const style = ref.getStyle();

    fieldX.value = parseInt(style.left) || 0;

    fieldY.value = parseInt(style.top) || 0;

    fieldW.value = parseInt(style.width) || 0;

    fieldH.value = parseInt(style.height) || 0;
}

function setSelectionOutline(comp, on) {
    const el = comp.getEl();
    if (!el) return;

    if (on) {
        el.style.setProperty("outline", "2px solid #F08080", "important");
    } else {
        el.style.removeProperty("outline");
    }
}

// Properties worth preserving on copy - extend this list if you
// use other CSS properties in your components (letter-spacing,
// box-shadow, text-transform, etc.)
const COPYABLE_PROPS = [
    "width",
    "height",
    "fontFamily",
    "fontSize",
    "fontWeight",
    "fontStyle",
    "color",
    "backgroundColor",
    "background",
    "border",
    "borderRadius",
    "padding",
    "paddingTop",
    "paddingRight",
    "paddingBottom",
    "paddingLeft",
    "margin",
    "textAlign",
    "lineHeight",
    "letterSpacing",
    "boxSizing",
    "opacity",
    "textDecoration",
    "textTransform"
];

function captureComputedStyle(el) {
    const computed = window.getComputedStyle(el);
    const captured = {};

    COPYABLE_PROPS.forEach(prop => {
        const value = computed[prop];
        if (value !== undefined && value !== "") captured[prop] = value;
    });

    return captured;
}

// ---------------------
// Copy / Paste
// ---------------------

let clipboard = [];
let pasteCount = 0; // NEW: tracks repeated pastes of the same copy
const PASTE_OFFSET = 20; // px, shift per paste

function copySelected() {
    if (!selectedComps.length) {
        alert("Select at least 1 component to copy");
        return;
    }

    clipboard = selectedComps.map(comp => {
        const el = comp.getEl();
        const rect = el.getBoundingClientRect();
        const parentRect = el.offsetParent
            ? el.offsetParent.getBoundingClientRect()
            : { left: 0, top: 0 };

        const left = Math.round(rect.left - parentRect.left);
        const top = Math.round(rect.top - parentRect.top);

        return {
            json: comp.toJSON(),
            originalLeft: left,
            originalTop: top,
            computedStyle: captureComputedStyle(el),
            // Remember the real parent (e.g. the field container that
            // has pointer-events:auto) so paste can insert back into
            // the SAME container, not always the root wrapper.
            parent: comp.parent() || editor.DomComponents.getWrapper()
        };
    });

    pasteCount = 0;
}

// ---------------------
// Copy / Paste - ID + text numbering
// ---------------------

function getBaseName(name) {
    // Strip any existing _copyN suffix so copying a copy
    // still increments from the ORIGINAL's highest number,
    // not from "name_copy1_copy1_copy1..."
    return name.replace(/_copy\d+$/i, "");
}

function findNextCopyNumber(baseId, usedIds) {
    const wrapper = editor.DomComponents.getWrapper();
    const escaped = baseId.replace(/[-/\\^$*+?.()|[\]{}]/g, "\\$&");
    const re = new RegExp("^" + escaped + "_copy(\\d+)$", "i");

    let maxN = 0;

    // Check what's already on the canvas
    wrapper.find("*").forEach(c => {
        const attrs = c.getAttributes();
        const id = attrs.id || attrs["data-field-id"] || "";
        const match = id.match(re);
        if (match) {
            const n = parseInt(match[1], 10);
            if (n > maxN) maxN = n;
        }
    });

    // Also check IDs already claimed earlier in THIS SAME paste batch,
    // so pasting multiple components at once never collides
    usedIds.forEach(id => {
        const match = id.match(re);
        if (match) {
            const n = parseInt(match[1], 10);
            if (n > maxN) maxN = n;
        }
    });

    return maxN + 1;
}

function replaceIdInContent(content, baseId, newId) {
    if (typeof content !== "string" || !content.includes(baseId)) return content;

    // Whole-word match so "age" doesn't accidentally match inside
    // something like "average" or "storage"
    const escapedBase = baseId.replace(/[-/\\^$*+?.()|[\]{}]/g, "\\$&");
    const re = new RegExp("\\b" + escapedBase + "\\b", "g");

    return content.replace(re, newId);
}

function assignFreshIds(json, usedIds) {
    // Only real, named field components should get renamed. Nested
    // formatting wrappers or textnode children (common in bold/colored
    // RTE content) never had an id to begin with - leave them exactly
    // as they are, or GrapesJS can fail to properly rebuild that part
    // of the tree on paste.
    const hasOwnId =
        json.attributes &&
        (json.attributes["id"] || json.attributes["data-field-id"]);

    if (!hasOwnId) {
        // Still recurse in case a real named field is nested deeper
        // inside a formatting wrapper, but don't touch this node itself.
        if (Array.isArray(json.components)) {
            json.components.forEach(child => assignFreshIds(child, usedIds));
        }
        return json;
    }

    const oldId = json.attributes["id"] || json.attributes["data-field-id"];

    const baseId = getBaseName(oldId);
    const n = findNextCopyNumber(baseId, usedIds);
    const newId = baseId + "_copy" + n;

    usedIds.add(newId);

    json.attributes = {
        ...json.attributes,
        id: newId,
        class: newId
    };

    if (json.attributes["data-field-id"]) {
        json.attributes["data-field-id"] = newId;
    }

    // Direct substring replace - works whether content is raw text,
    // "&lt;age&gt;" entity-encoded, or has surrounding label text,
    // with NO decode/re-encode round trip needed
    if (typeof json.content === "string") {
        json.content = replaceIdInContent(json.content, baseId, newId);
    }

    // Some GrapesJS text components store their actual text as
    // child "textnode" components rather than top-level `content`.
    // Walk children and patch those too - but only rename children
    // that themselves have their own id (handled recursively above).
    if (Array.isArray(json.components)) {
        json.components.forEach(child => {
            if (typeof child.content === "string") {
                child.content = replaceIdInContent(child.content, baseId, newId);
            }

            assignFreshIds(child, usedIds);
        });
    }

    return json;
}

function pasteClipboard() {
    if (!clipboard.length) {
        alert("Nothing to paste - copy a component first");
        return;
    }

    pasteCount++;

    const wrapper = editor.DomComponents.getWrapper();
    const pasted = [];
    const usedIds = new Set();

    clipboard.forEach(item => {
        const json = JSON.parse(JSON.stringify(item.json));

        assignFreshIds(json, usedIds);

        json.style = {
            ...item.computedStyle,
            ...json.style,
            position: "absolute",
            left: item.originalLeft + PASTE_OFFSET * pasteCount + "px",
            top: item.originalTop + PASTE_OFFSET * pasteCount + "px"
        };

        // Insert into the SAME parent the original lived in (e.g. the
        // field container with pointer-events:auto), not the root
        // wrapper.
        const targetParent = item.parent || wrapper;
        const added = targetParent.append(json);
        const comp = added[0];

        // NEW - explicitly restore interactivity, don't rely on
        // toJSON having carried these over, or on the
        // component:selected loop catching it afterward
        comp.set({
            editable: false, // manual dblclick-to-edit handles this uniformly instead
            selectable: true,
            hoverable: true,
            draggable: true,
            copyable: true,
            removable: true,
            badgable: true
        });

        // pointer-events isn't a property GrapesJS's Style Manager
        // recognizes, so comp.addStyle({pointerEvents:...}) gets
        // silently dropped and never reaches the generated CSS -
        // confirmed via editor.getCss(). Set it directly on the live
        // DOM element instead, the same proven approach already used
        // for the selection outline elsewhere in this file.
        const pastedEl = comp.getEl();
        if (pastedEl) {
            pastedEl.style.setProperty("pointer-events", "auto", "important");
        }

        pasted.push(comp);
    });

    clearMultiSelection();

    pasted.forEach(comp => {
        selectedComps.push(comp);
        setSelectionOutline(comp, true);
    });

    suppressSelectHandler = true;
    editor.select(pasted[0]);
    suppressSelectHandler = false;
    syncFields();
}

// ---------------------
// Align Functions
// ---------------------

function requireMultiSelection() {
    if (selectedComps.length < 2) {
        alert("Select at least 2 components (Ctrl+Click) first");
        return false;
    }
    return true;
}

// Measures ACTUAL rendered position/size, not the style model - a
// component positioned via a CSS class (not inline style) returns
// undefined from getStyle(), which silently broke alignRight/alignBottom
// math before (they'd fall back to 0 and look identical to
// alignLeft/alignTop).
function getRect(comp) {
    const el = comp.getEl();
    const rect = el.getBoundingClientRect();

    const parentRect = el.offsetParent
        ? el.offsetParent.getBoundingClientRect()
        : { left: 0, top: 0 };

    return {
        left: Math.round(rect.left - parentRect.left),
        top: Math.round(rect.top - parentRect.top),
        width: Math.round(rect.width),
        height: Math.round(rect.height)
    };
}

function alignLeft() {
    if (!requireMultiSelection()) return;

    // Align to the LEFTMOST edge across the whole selection
    const targetLeft = Math.min(...selectedComps.map(c => getRect(c).left));

    selectedComps.forEach(comp => {
        comp.addStyle({ position: "absolute", left: targetLeft + "px" });
    });
}

function alignRight() {
    if (!requireMultiSelection()) return;

    // Align to the RIGHTMOST edge across the whole selection
    const targetRight = Math.max(
        ...selectedComps.map(c => {
            const r = getRect(c);
            return r.left + r.width;
        })
    );

    selectedComps.forEach(comp => {
        const width = getRect(comp).width;
        comp.addStyle({ position: "absolute", left: targetRight - width + "px" });
    });
}

function alignTop() {
    if (!requireMultiSelection()) return;

    // Align to the TOPMOST edge across the whole selection
    const targetTop = Math.min(...selectedComps.map(c => getRect(c).top));

    selectedComps.forEach(comp => {
        comp.addStyle({ position: "absolute", top: targetTop + "px" });
    });
}

function alignBottom() {
    if (!requireMultiSelection()) return;

    // Align to the BOTTOMMOST edge across the whole selection
    const targetBottom = Math.max(
        ...selectedComps.map(c => {
            const r = getRect(c);
            return r.top + r.height;
        })
    );

    selectedComps.forEach(comp => {
        const height = getRect(comp).height;
        comp.addStyle({ position: "absolute", top: targetBottom - height + "px" });
    });
}

function sameWidth() {
    if (!requireMultiSelection()) return;

    const width = getRect(selectedComps[0]).width;

    selectedComps.forEach(comp => {
        comp.addStyle({ width: width + "px" });
    });
}

function sameHeight() {
    if (!requireMultiSelection()) return;

    const height = getRect(selectedComps[0]).height;

    selectedComps.forEach(comp => {
        comp.addStyle({ height: height + "px" });
    });
}

function sameSize() {
    sameWidth();
    sameHeight();
}

// ---------------------
// Wire up toolbar buttons
// ---------------------

document.getElementById("btnCopy").addEventListener("click", copySelected);

document.getElementById("btnPaste").addEventListener("click", pasteClipboard);

document.getElementById("btnAlignLeft")?.addEventListener("click", alignLeft);

document.getElementById("btnAlignRight")?.addEventListener("click", alignRight);

document.getElementById("btnAlignTop")?.addEventListener("click", alignTop);

document.getElementById("btnAlignBottom")?.addEventListener("click", alignBottom);

document.getElementById("btnSameWidth")?.addEventListener("click", sameWidth);

document.getElementById("btnSameHeight")?.addEventListener("click", sameHeight);

document.getElementById("btnSameSize")?.addEventListener("click", sameSize);

/* ********* MULTI SELECT FOR ALIGNMENT ENDS ************ */

/* ensure reliable history tracking */
editor.UndoManager.start();

const fieldX = document.getElementById("txtX");

const fieldY = document.getElementById("txtY");

const fieldW = document.getElementById("txtW");

const fieldH = document.getElementById("txtH");

// Shared rename-commit logic, called both by GrapesJS's own
// rte:disable event (works reliably for components GrapesJS detected
// as "text" type, like a freshly-created <NEW> field) AND by our own
// manual double-click-to-edit below (which works uniformly on ANY
// field regardless of type - existing/pasted fields parsed from raw
// HTML aren't guaranteed to get GrapesJS's "text" type, so its
// built-in double-click RTE trigger can silently never activate for
// them even with editable:true set).
function handleFieldTextCommit(comp) {
    const el = comp.getEl();

    if (!el) return;

    const text = el.innerText?.trim();

    if (!text) return;

    const match = text.match(/^<(.+)>$/);

    if (!match) return;

    // Only the FIRST whitespace-separated token becomes the actual
    // field name used for id/class/data-field-id. Anything after the
    // first space is just descriptive text that stays in the visible
    // content as-is - e.g. "<policyterm d0>" renames the field to
    // "policyterm" while the displayed text remains "<policyterm d0>".
    const rawInner = match[1].trim();
    const firstToken = rawInner.split(/\s+/)[0] || "";

    const fieldName = firstToken.replace(/[^\w-]/g, "");

    if (!fieldName) return;

    const attrs = comp.getAttributes();

    const oldName = attrs["data-field-id"];

    if (oldName === fieldName) return;

    //------------------------------------
    // Store pending rename only
    //------------------------------------

    comp.set("pendingFieldName", fieldName);

    comp.setName(fieldName);
}

editor.on("rte:disable", view => {
    const comp = view?.model;

    if (!comp) return;

    handleFieldTextCommit(comp);
});

// Map each field to the single CSS property it controls, so changing
// one field only touches that one property on every selected
// component - the others stay exactly as they were per-component.
const FIELD_PROP_MAP = [
    { field: fieldX, prop: "left" },
    { field: fieldY, prop: "top" },
    { field: fieldW, prop: "width" },
    { field: fieldH, prop: "height" }
];

FIELD_PROP_MAP.forEach(({ field, prop }) => {
    // Track whether the user has actually TYPED into this field since
    // it was last synced. Setting field.value programmatically (e.g.
    // syncFields() populating it after a selection) never fires an
    // "input" event, so this flag only goes true on genuine user
    // input - which is what blur should require before overwriting
    // every selected component's position. Without this, a field
    // left focused from an earlier edit would silently reapply its
    // stale value to the WHOLE new selection the next time the user
    // clicked anywhere else on the canvas (e.g. right after a
    // marquee select).
    let userEdited = false;

    field.addEventListener("input", () => {
        userEdited = true;
    });

    const apply = () => {
        if (!selectedComps.length) return;

        const value = parseInt(field.value || 0) + "px";

        selectedComps.forEach(comp => {
            comp.addStyle({ position: "absolute", [prop]: value });
        });
    };

    // Press Enter
    field.addEventListener("keydown", e => {
        if (e.key === "Enter") {
            apply();
            userEdited = false;
            field.blur();
        }
    });

    // Click outside - only apply if the user actually changed this
    // field themselves.
    field.addEventListener("blur", () => {
        if (userEdited) {
            apply();
        }
        userEdited = false;
    });
});

editor.on("component:deselected", () => {
    fieldX.value = 0;
    fieldY.value = 0;
    fieldW.value = 0;
    fieldH.value = 0;
});

// -------------------------
// Startup
// -------------------------

editor.on("load", () => {
    editor.select(null);

    // Delete button
    editor.Panels.addButton("options", {
        id: "delete-component",
        className: "fa fa-trash",
        attributes: { title: "Delete selected component" },
        command: "delete-component"
    });

    // New button
    editor.Panels.addButton("options", {
        id: "new-component",
        className: "fa fa-plus",
        attributes: { title: "Add New Text Component" },
        command: "new-component"
    });

    // Existing fields should behave exactly like NEW components
    editor.DomComponents.getWrapper()
        .find("[data-field-id]")
        .forEach(comp => {
            comp.set({ editable: false }); // manual dblclick-to-edit handles this uniformly instead
        });

    const body = editor.Canvas.getBody();

    let dragGroup = [];
    let startMouseX = 0;
    let startMouseY = 0;
    let isDraggingGroup = false;
    let wasDragging = false; // NEW

    // Tracks the element currently in active manual text-edit mode
    // (contentEditable, set by startManualTextEdit below). While a
    // field is being edited, clicks inside it must behave like normal
    // text editing (placing the cursor to insert/append) - NOT get
    // hijacked by Case 1's group-drag handling below, which is what
    // was blocking single-click cursor placement during an edit.
    let currentlyEditingEl = null;

    // -------------------------
    // Keyboard shortcuts for every toolbar button. Registered on BOTH
    // the parent page's document AND the canvas iframe's own
    // document, since keyboard events firing while focus is inside
    // the iframe never reach the parent document (same cross-frame
    // isolation as mouse events, discovered earlier) - a shortcut
    // registered on only one side would silently stop working
    // depending on where focus currently is.
    // -------------------------

    function handleGlobalKeydown(e) {
        const ctrl = e.ctrlKey || e.metaKey;

        // Undo/redo stay available everywhere EXCEPT while actively
        // editing a field's text (to avoid fighting the browser's own
        // contentEditable undo stack there).
        if (ctrl && !currentlyEditingEl) {
            if (e.key === "z" && !e.shiftKey) {
                e.preventDefault();
                editor.UndoManager.undo();
                return;
            }
            if (e.key === "y") {
                e.preventDefault();
                editor.UndoManager.redo();
                return;
            }
            if (e.key === "z" && e.shiftKey) {
                e.preventDefault();
                editor.UndoManager.redo();
                return;
            }
        }

        // Don't hijack any other shortcut while the user is actively
        // typing - inside a manually-edited field, or any normal
        // input/textarea (like the X/Y/W/H coordinate boxes).
        const active = (e.target && e.target.ownerDocument ? e.target.ownerDocument : document).activeElement;
        const isTyping =
            !!currentlyEditingEl ||
            (active && (active.tagName === "INPUT" || active.tagName === "TEXTAREA" || active.isContentEditable));

        if (isTyping) return;

        if (ctrl) {
            switch (e.key.toLowerCase()) {
                case "o":
                    e.preventDefault();
                    document.getElementById("btnOpen").click();
                    return;
                case "s":
                    e.preventDefault();
                    document.getElementById("btnSave").click();
                    return;
                case "c":
                    e.preventDefault();
                    document.getElementById("btnCopy").click();
                    return;
                case "v":
                    e.preventDefault();
                    document.getElementById("btnPaste").click();
                    return;
                case "n":
                    // Browsers commonly reserve Ctrl+N for opening a
                    // new window and won't let JS override it - this
                    // still wires it up in case the app's runtime
                    // (e.g. an Electron-style shell) allows it.
                    e.preventDefault();
                    document.getElementById("btnNew").click();
                    return;
                case "=":
                case "+":
                    e.preventDefault();
                    document.getElementById("btnZoomIn").click();
                    return;
                case "-":
                    e.preventDefault();
                    document.getElementById("btnZoomOut").click();
                    return;
                case "0":
                    e.preventDefault();
                    document.getElementById("btnZoom100").click();
                    return;
            }
            return;
        }

        if (e.altKey) {
            switch (e.key.toLowerCase()) {
                case "arrowleft":
                    e.preventDefault();
                    document.getElementById("btnAlignLeft")?.click();
                    return;
                case "arrowright":
                    e.preventDefault();
                    document.getElementById("btnAlignRight")?.click();
                    return;
                case "arrowup":
                    e.preventDefault();
                    document.getElementById("btnAlignTop")?.click();
                    return;
                case "arrowdown":
                    e.preventDefault();
                    document.getElementById("btnAlignBottom")?.click();
                    return;
                case "w":
                    e.preventDefault();
                    document.getElementById("btnSameWidth")?.click();
                    return;
                case "h":
                    e.preventDefault();
                    document.getElementById("btnSameHeight")?.click();
                    return;
                case "e":
                    e.preventDefault();
                    document.getElementById("btnSameSize")?.click();
                    return;
            }
            return;
        }

        if (e.key === "Delete" || e.key === "Backspace") {
            e.preventDefault();
            document.getElementById("btnDelete").click();
        }
    }

    document.addEventListener("keydown", handleGlobalKeydown);
    body.ownerDocument.addEventListener("keydown", handleGlobalKeydown);

    // ---- Marquee (rubber-band) selection state ----
    let isMarqueeSelecting = false;
    let marqueeStart = null;
    let marqueeEl = null;
    let marqueeExtend = false;

    // The canvas iframe's own document scrolls internally (see the
    // auto-scroll code below) - this returns how far it's currently
    // scrolled, so marquee coordinates can be tracked relative to the
    // DOCUMENT instead of the (fixed) viewport. Falls back across
    // documentElement/body since browsers differ on which one they
    // report scroll position on.
    function getIframeScroll() {
        const doc = body.ownerDocument;
        const win = doc.defaultView;
        return {
            x: win.scrollX || doc.documentElement.scrollLeft || doc.body.scrollLeft || 0,
            y: win.scrollY || doc.documentElement.scrollTop || doc.body.scrollTop || 0
        };
    }

    function startMarquee(e) {
        isMarqueeSelecting = true;

        // Store the start point in DOCUMENT-relative coordinates
        // (viewport position + current scroll offset). This is what
        // lets the logical selection stay anchored to the same
        // content even as auto-scroll moves the canvas underneath a
        // visually-stationary marquee box.
        const scroll = getIframeScroll();
        marqueeStart = { x: e.clientX + scroll.x, y: e.clientY + scroll.y };
        marqueeExtend = e.ctrlKey || e.metaKey;

        marqueeEl = body.ownerDocument.createElement("div");
        marqueeEl.style.position = "fixed";
        marqueeEl.style.margin = "0";
        marqueeEl.style.boxSizing = "border-box";
        marqueeEl.style.left = e.clientX + "px";
        marqueeEl.style.top = e.clientY + "px";
        marqueeEl.style.width = "0px";
        marqueeEl.style.height = "0px";
        marqueeEl.style.border = "1px dashed #2D7FF9";
        marqueeEl.style.background = "rgba(45,127,249,0.10)";
        marqueeEl.style.zIndex = "999999";
        marqueeEl.style.pointerEvents = "none";

        // Append to <html> rather than <body> - if GrapesJS's zoom
        // feature applies a CSS transform to body (even a no-op
        // scale(1) at 100% zoom), any transform on an ancestor
        // changes what position:fixed measures against, from the
        // true viewport to that ancestor's own box - which is exactly
        // the kind of small, consistent offset that would make the
        // box start a few pixels/lines away from the actual cursor
        // position. <html> is a much less likely target for such a
        // transform.
        (body.ownerDocument.documentElement || body).appendChild(marqueeEl);
    }

    function updateMarquee(e) {
        const scroll = getIframeScroll();

        // Current cursor position in DOCUMENT-relative terms, same
        // coordinate space as marqueeStart.
        const curDocX = e.clientX + scroll.x;
        const curDocY = e.clientY + scroll.y;

        // The mousedown point is the anchor corner, but the box can
        // grow in ANY direction from there (up/down/left/right) -
        // standard rubber-band selection. Whichever corner ends up
        // top-left/bottom-right is computed via min/max against the
        // current cursor position.
        const docX1 = Math.min(marqueeStart.x, curDocX);
        const docY1 = Math.min(marqueeStart.y, curDocY);
        const docX2 = Math.max(marqueeStart.x, curDocX);
        const docY2 = Math.max(marqueeStart.y, curDocY);
        const x1 = docX1 - scroll.x;
        const y1 = docY1 - scroll.y;
        const x2 = docX2 - scroll.x;
        const y2 = docY2 - scroll.y;

        marqueeEl.style.left = x1 + "px";
        marqueeEl.style.top = y1 + "px";
        marqueeEl.style.width = x2 - x1 + "px";
        marqueeEl.style.height = y2 - y1 + "px";

        lastMarqueeEvent = e;
        updateAutoScroll(e);
    }

    // -------------------------
    // Auto-scroll the canvas when the marquee (or group-drag) cursor
    // nears the edge of the visible viewport. Without this, dragging
    // toward a field outside the current scroll position just freezes
    // at the edge - mousemove on the iframe body stops firing once
    // the cursor can't move any further within the iframe, so fields
    // further out could never be reached or matched.
    // -------------------------

    let autoScrollTimer = null;
    let lastMarqueeEvent = null;

    function updateAutoScroll(e) {
        const EDGE = 40; // px from viewport edge that triggers scrolling
        const SPEED = 18; // px per tick

        const view = body.ownerDocument.defaultView;
        const vw = view.innerWidth;
        const vh = view.innerHeight;

        let dx = 0;
        let dy = 0;

        if (e.clientX < EDGE) dx = -SPEED;
        else if (e.clientX > vw - EDGE) dx = SPEED;

        if (e.clientY < EDGE) dy = -SPEED;
        else if (e.clientY > vh - EDGE) dy = SPEED;

        if (!dx && !dy) {
            stopAutoScroll();
            return;
        }

        if (autoScrollTimer) return; // already running

        // The REAL scroll container is the canvas iframe's OWN
        // document (its <html>/<body>) - GrapesJS renders the canvas
        // at a fixed device size and scrolls internally, it does NOT
        // scroll an outer wrapper. window.scrollBy() on the iframe's
        // own window is the standards-correct way to scroll it,
        // regardless of whether the browser tracks scroll position on
        // <html> or <body>.
        const iframeWindow = body.ownerDocument.defaultView;

        autoScrollTimer = setInterval(() => {
            const ev = lastMarqueeEvent;
            if (!ev) return;

            let sdx = 0;
            let sdy = 0;

            if (ev.clientX < EDGE) sdx = -SPEED;
            else if (ev.clientX > vw - EDGE) sdx = SPEED;

            if (ev.clientY < EDGE) sdy = -SPEED;
            else if (ev.clientY > vh - EDGE) sdy = SPEED;

            if (!sdx && !sdy) {
                stopAutoScroll();
                return;
            }

            const before = getIframeScroll();
            iframeWindow.scrollBy(sdx, sdy);
            const after = getIframeScroll();

            // Re-run the marquee box update against the new scroll
            // position so the rubber-band keeps growing correctly.
            if (isMarqueeSelecting && marqueeEl) {
                updateMarquee(ev);
            }
        }, 30);
    }

    function stopAutoScroll() {
        if (autoScrollTimer) {
            clearInterval(autoScrollTimer);
            autoScrollTimer = null;
        }
    }

    function finishMarquee(e) {
        isMarqueeSelecting = false;
        stopAutoScroll();

        const scroll = getIframeScroll();
        const curDocX = e.clientX + scroll.x;
        const curDocY = e.clientY + scroll.y;

        // Standard rubber-band rectangle, matching updateMarquee -
        // any direction from the mousedown anchor.
        const docX1 = Math.min(marqueeStart.x, curDocX);
        const docY1 = Math.min(marqueeStart.y, curDocY);
        const docX2 = Math.max(marqueeStart.x, curDocX);
        const docY2 = Math.max(marqueeStart.y, curDocY);

        // getBoundingClientRect() is always VIEWPORT-relative, so
        // convert our document-relative rectangle back using the
        // CURRENT (final) scroll offset before comparing against it.
        const x1 = docX1 - scroll.x;
        const y1 = docY1 - scroll.y;
        const x2 = docX2 - scroll.x;
        const y2 = docY2 - scroll.y;

        if (marqueeEl) {
            marqueeEl.remove();
            marqueeEl = null;
        }

        // Ignore accidental micro-drags (essentially just a click,
        // not a real rectangle) so a stray click on empty canvas
        // doesn't wipe the current selection.
        if (x2 - x1 < 4 && y2 - y1 < 4) return;

        const wrapper = editor.DomComponents.getWrapper();
        const matched = [];

        wrapper.find("[data-field-id]").forEach(comp => {
            const el = comp.getEl();
            if (!el) return;

            const r = el.getBoundingClientRect();

            // Standard rectangle-intersection test - any component
            // that overlaps the marquee at all gets selected.
            const intersects = r.left < x2 && r.right > x1 && r.top < y2 && r.bottom > y1;

            if (intersects) matched.push(comp);
        });

        if (!matched.length) return;

        if (marqueeExtend) {
            // Ctrl/Cmd + marquee: ADD newly matched components to the
            // existing selection instead of replacing it. Skip anything
            // already selected so it doesn't get pushed twice.
            matched.forEach(comp => {
                if (!selectedComps.includes(comp)) {
                    selectedComps.push(comp);
                    setSelectionOutline(comp, true);
                }
            });
        } else {
            clearMultiSelection();

            matched.forEach(comp => {
                selectedComps.push(comp);
                setSelectionOutline(comp, true);
            });
        }

        // Hand the whole selection off to GrapesJS's own NATIVE
        // multi-select (the same mechanism its built-in ctrl+click
        // uses internally), instead of select(null). This is what
        // makes the Style Manager open and target every selected
        // component - select(null) cleared GrapesJS's native
        // selection entirely, so the panel had nothing to show.
        suppressSelectHandler = true;

        editor.select(selectedComps[0]);

        for (let i = 1; i < selectedComps.length; i++) {
            editor.selectAdd(selectedComps[i]);
        }

        suppressSelectHandler = false;

        syncFields();

        // Reuse the existing click-suppression flag so the click that
        // follows this mouseup doesn't collapse the new selection.
        wasDragging = true;
    }

    function getComponentFromElement(el) {
        // Walk up the DOM until we find an element GrapesJS
        // has a component model attached to
        while (el && el !== body) {
            const fieldId = el.getAttribute && el.getAttribute("data-field-id");

            const comp =
                editor.Components?.getModel?.(el) ||
                (el.id && editor.getWrapper()?.find?.(`#${el.id}`)?.[0]) ||
                (fieldId &&
                    editor.getWrapper()?.find?.(`[data-field-id="${fieldId}"]`)?.[0]);

            if (comp) return comp;
            el = el.parentElement;
        }
        return null;
    }

    body.addEventListener("mousedown", e => {
        // While a field is in active manual text-edit mode, let
        // clicks inside it behave like completely normal text editing
        // (placing the cursor to insert/append at that point) - don't
        // let Case 1 below hijack it into a group-drag, which was
        // blocking single-click cursor placement during an edit.
        if (currentlyEditingEl && (e.target === currentlyEditingEl || currentlyEditingEl.contains(e.target))) {
            return;
        }

        const ctrlHeld = e.ctrlKey || e.metaKey;

        // Case 1: clicked on a component that's already part of the
        // current multi-selection - start the existing group drag.
        // Skip this when Ctrl/Cmd is held: that means the user is
        // trying to start a NEW marquee to append to the selection
        // (possibly starting from on top of an already-selected
        // component), not drag the group.
        if (selectedComps.length && !ctrlHeld) {
            const clickedComp = getComponentFromElement(e.target);

            if (clickedComp && selectedComps.includes(clickedComp)) {
                // Stop the browser's own default drag/text-selection
                // behavior AND, critically, stop this mousedown from
                // reaching GrapesJS's own native dragger. Without this,
                // GrapesJS's built-in absolute-mode dragging can start
                // on the SAME gesture as our custom group-drag below -
                // two drag systems fighting over the same component
                // leaves it in a broken state afterward (unselectable,
                // double-click-to-edit dead).
                e.preventDefault();
                e.stopPropagation();

                isDraggingGroup = true;
                wasDragging = false;

                startMouseX = e.clientX;
                startMouseY = e.clientY;

                dragGroup = selectedComps.map(comp => {
                    const el = comp.getEl();

                    // Read the starting position from the SAME place
                    // the drag writes back to - the component's own
                    // CSS left/top - rather than computing it from
                    // rendered screen geometry (getBoundingClientRect).
                    // Those two aren't guaranteed to match: any canvas
                    // zoom other than exactly 100%, or any padding/
                    // border on the field's positioning container,
                    // throws rendered-pixel geometry out of sync with
                    // the logical CSS value, causing the group to jump
                    // the instant a drag starts.
                    const style = comp.getStyle();
                    const computed = el ? getComputedStyle(el) : {};

                    const geomLeft = (() => {
                        if (!el) return null;
                        const r = el.getBoundingClientRect();
                        const pr = el.offsetParent ? el.offsetParent.getBoundingClientRect() : { left: 0, top: 0 };
                        return Math.round(r.left - pr.left);
                    })();
                    const geomTop = (() => {
                        if (!el) return null;
                        const r = el.getBoundingClientRect();
                        const pr = el.offsetParent ? el.offsetParent.getBoundingClientRect() : { left: 0, top: 0 };
                        return Math.round(r.top - pr.top);
                    })();

                    const finalLeft = parseInt(style.left || computed.left) || 0;
                    const finalTop = parseInt(style.top || computed.top) || 0;

                    return {
                        comp,
                        left: finalLeft,
                        top: finalTop
                    };
                });

                dragGroup.forEach(item => {
                    item.comp.set("draggable", false);
                    const style = item.comp.getStyle();

                    if (style.position !== "absolute") {
                        // Setting position:absolute alone, without
                        // ALSO re-asserting left/top in the same
                        // call, can cause a real visible jump: if
                        // this component's left/top only ever lived
                        // in an external CSS class rule (not this
                        // component's own tracked style), it would
                        // render with no explicit left/top until the
                        // very next addStyle() call - falling back to
                        // wherever unpositioned absolute elements
                        // land, typically near the top-left of the
                        // container. Setting all three together
                        // closes that gap entirely.
                        item.comp.addStyle({
                            position: "absolute",
                            left: item.left + "px",
                            top: item.top + "px"
                        });
                    }
                });

                return;
            }
        }

        // Case 2: clicked on empty canvas (no field component under
        // the cursor at all) - start a marquee/rubber-band selection.
        const clickedField = e.target.closest && e.target.closest("[data-field-id]");

        if (!clickedField) {
            startMarquee(e);
        }

        // Case 3: clicked on a not-yet-selected single component -
        // fall through and let GrapesJS's own click-to-select handle it.
    }, true); // capture phase - see Case 1 comment above for why

    function handleCanvasMouseMove(e) {
        if (isMarqueeSelecting) {
            updateMarquee(e);
            return;
        }

        if (!isDraggingGroup || !dragGroup.length) return;

        const rawDx = e.clientX - startMouseX;
        const rawDy = e.clientY - startMouseY;

        // Ignore tiny jitter so a double-click (which naturally has a
        // pixel or two of movement between its two clicks) doesn't
        // get misread as a drag. Once wasDragging flips true here,
        // the capture-phase click-suppressor below stops that click
        // from ever reaching GrapesJS - which is exactly what was
        // blocking double-click-to-edit on an already-selected field.
        if (!wasDragging && Math.abs(rawDx) < 4 && Math.abs(rawDy) < 4) return;

        wasDragging = true; // NEW - mouse actually moved during this drag

        const dx = rawDx / zoom;
        const dy = rawDy / zoom;

        // Move every selected component relative to ITS OWN
        // original position, using addStyle so we don't wipe
        // out position/width/height/border etc.
        dragGroup.forEach(item => {
            item.comp.addStyle({
                left: item.left + dx + "px",
                top: item.top + dy + "px"
            });
        });

        // Update toolbar based on the "leader" (first selected)
        const leader = dragGroup[0].comp;
        fieldX.value = parseInt(leader.getStyle().left) || 0;
        fieldY.value = parseInt(leader.getStyle().top) || 0;
    }

    function handleCanvasMouseUp(e) {
        if (isMarqueeSelecting) {
            finishMarquee(e);
            return;
        }

        if (!isDraggingGroup) return;

        // ---- NEW: restore native dragging once the custom drag ends
        dragGroup.forEach(item => item.comp.set("draggable", true));

        // Styles are already committed via addStyle on every move,
        // nothing further to merge here - just reset drag state.
        isDraggingGroup = false;
        dragGroup = [];
    }

    // body.addEventListener("mousemove", handleCanvasMouseMove);
    // body.addEventListener("mouseup", handleCanvasMouseUp);

    // Rescue listeners on the PARENT page's own document. The
    // diagnostic run proved that once a drag starts inside the
    // iframe, mouse events get captured at the PARENT document level
    // for the rest of that drag - REGARDLESS of where the cursor
    // visually is. Classifying events as "inside vs outside the
    // iframe" and only acting on the "outside" ones was silently
    // dropping real movement (stalling the marquee's growth) and,
    // critically, dropping the mouseup itself (leaving the drag
    // permanently stuck "in progress"). So: no position-based
    // filtering. mouseup is naturally safe to process twice (a second
    // call just sees isMarqueeSelecting/isDraggingGroup already
    // false and no-ops), and mousemove is deduplicated by comparing
    // against the last processed coordinates instead of guessing
    // based on position.
    function toIframeCoords(e) {
        const frameEl = editor.Canvas.getFrameEl();
        const rect = frameEl.getBoundingClientRect();
        return {
           // clientX: e.clientX - rect.left,
           // clientY: e.clientY - rect.top,
			clientX: e.clientX, 
            clientY: e.clientY,
            ctrlKey: e.ctrlKey,
            metaKey: e.metaKey
        };
    }

    let lastProcessedX = null;
    let lastProcessedY = null;

    document.addEventListener("mousemove", e => {
        if (!isMarqueeSelecting && !isDraggingGroup) return;

        const coords = toIframeCoords(e);

        if (coords.clientX === lastProcessedX && coords.clientY === lastProcessedY) {
            return; // same position already processed (likely by body's own listener)
        }

        lastProcessedX = coords.clientX;
        lastProcessedY = coords.clientY;

        handleCanvasMouseMove(coords);
    });

    document.addEventListener("mouseup", e => {
        if (!isMarqueeSelecting && !isDraggingGroup) return;
        lastProcessedX = null;
        lastProcessedY = null;
        handleCanvasMouseUp(toIframeCoords(e));
    });

    // The mouse leaving the iframe body is EXPECTED and normal now -
    // it's exactly what happens when dragging toward a field near/
    // past the edge, which is the whole scenario auto-scroll exists
    // for. Cancelling here would throw away the drag every time it's
    // needed most. Tracking continues via the parent-document-level
    // listeners below instead; this is now just a no-op placeholder
    // kept for clarity of intent.
    body.addEventListener("mouseleave", () => {});

    // -------------------------
    // Manual double-click-to-edit, working uniformly on ANY field.
    //
    // GrapesJS's own built-in double-click RTE trigger only activates
    // for components it detected as a "text"-compatible type. A
    // freshly-created <NEW> field (created with explicit type:"text")
    // gets that reliably, but existing fields parsed from raw HTML on
    // Open, and fields created via copy/paste, aren't guaranteed to
    // get that same type - so GrapesJS's native trigger can silently
    // never fire for them even with editable:true correctly set. This
    // bypasses that entirely: any element with data-field-id becomes
    // directly editable on double-click, regardless of how it was
    // created.
    // -------------------------

    body.addEventListener("dblclick", e => {
        const target = e.target.closest && e.target.closest("[data-field-id]");

        if (!target) return;

        const comp = getComponentFromElement(target);

        if (!comp) return;

        // Native double-click text-selection never actually happens
        // here: Case 1's mousedown handler calls e.preventDefault()
        // on the field (to stop GrapesJS's own drag from competing
        // with ours), and text selection is part of mousedown's
        // default action - so it gets suppressed along with it. We
        // select the full content ourselves instead, so double-click
        // reliably enters edit mode with everything selected, ready
        // for full replacement.
        startManualTextEdit(comp, target);
    });

    function selectInnerContent(el) {
        const doc = el.ownerDocument;

        // Find the field's text node (its content is typically a
        // single text node like "<policyterm d0>").
        const walker = doc.createTreeWalker(el, NodeFilter.SHOW_TEXT);
        const textNode = walker.nextNode();

        const range = doc.createRange();

        if (!textNode) {
            range.selectNodeContents(el);
            return range;
        }

        const text = textNode.textContent;
        const start = text.indexOf("<");
        const end = text.lastIndexOf(">");

        if (start === -1 || end === -1 || end <= start) {
            // No brackets found - fall back to selecting everything.
            range.selectNodeContents(el);
            return range;
        }

        // Select only what's BETWEEN the brackets, not the brackets
        // themselves.
        range.setStart(textNode, start + 1);
        range.setEnd(textNode, end);
        return range;
    }

    function startManualTextEdit(comp, el) {
        // Explicitly focus the iframe's OWN window first. Pasted
        // components get added to selectedComps programmatically (not
        // via an actual click inside the canvas), so the iframe may
        // never have received real window-level focus - and Case 1's
        // mousedown handler calling preventDefault() on an already-
        // selected field (to stop GrapesJS's native drag) also
        // suppresses the browser's normal "clicking inside an iframe
        // shifts window focus to it" behavior. Without true window
        // focus, the element can still receive a DOM selection, but
        // it renders in the inactive/grey color and keyboard input
        // doesn't route to it - exactly the pasted-field symptom.
        const win = el.ownerDocument.defaultView;

        if (win && win.focus) win.focus();

        el.contentEditable = "true";
        el.focus();
        currentlyEditingEl = el;

        // Select only the content BETWEEN the angle brackets, not the
        // brackets themselves, so double-click is ready for full
        // replacement of just the field name/text. A subsequent
        // SINGLE click (handled by the currentlyEditingEl check at
        // the top of the mousedown handler, which lets clicks through
        // untouched) places the cursor normally for inserting/
        // appending instead.
        const doc = el.ownerDocument;
        const range = selectInnerContent(el);
        const sel = doc.defaultView.getSelection();
        sel.removeAllRanges();
        sel.addRange(range);

        function finish() {
            el.contentEditable = "false";
            currentlyEditingEl = null;
            el.removeEventListener("blur", finish);
            el.removeEventListener("keydown", onKeydown);
            el.removeEventListener("beforeinput", onBeforeInput);

            const newHtml = el.innerHTML;

            // GrapesJS's component view renders from the components()
            // children collection - this is the single source of
            // truth for both display AND serialization/export.
            // (A previous version of this code ALSO called
            // comp.set("content", newHtml) afterward, believing it
            // was needed for export - that call was actually what
            // wiped the rendered content back to blank, confirmed via
            // diagnostic logging. components() alone is correct.)
            comp.components(newHtml);

            handleFieldTextCommit(comp);
        }

        function onKeydown(ke) {
            if (ke.key === "Enter" || ke.key === "Escape") {
                ke.preventDefault();
                el.blur();
            }
        }

        // keydown's preventDefault() alone isn't fully reliable for
        // blocking a newline in contentEditable across every browser -
        // beforeinput fires right before the DOM would actually be
        // mutated and is the robust way to block a specific insertion
        // action. Without this, Enter could still insert a line break
        // that wraps inside the field's single-line box, making the
        // wrapped part look like it "vanished" (clipped/overflowing).
        function onBeforeInput(ie) {
            if (ie.inputType === "insertParagraph" || ie.inputType === "insertLineBreak") {
                ie.preventDefault();
                el.blur();
            }
        }

        el.addEventListener("blur", finish);
        el.addEventListener("keydown", onKeydown);
        el.addEventListener("beforeinput", onBeforeInput);
    }

    // NEW - capture-phase click suppressor. Runs BEFORE GrapesJS's own
    // click-to-select handler, so a drag-release click never collapses
    // your multi-selection back down to one component.
    body.addEventListener(
        "click",
        e => {
            if (wasDragging) {
                e.stopPropagation();
                e.preventDefault();
                wasDragging = false;
            }
        },
        true
    );

    // Hide Body/root only
    const root = editor.getWrapper();
    if (root) {
        root.set({
            selectable: false,
            hoverable: false,
            highlightable: false,
            badgable: false
        });
    }

    setTimeout(() => {
        const frame = editor.Canvas.getFrameEl();
        const frameDoc = frame.contentDocument;

        const link = frameDoc.createElement("link");
        link.rel = "stylesheet";
        link.href = "editor-sb.css";
        frameDoc.head.appendChild(link);
    }, 1000);
});

editor.Commands.add("delete-component", {
    run(editor) {
        const selected = editor.getSelected();

        if (!selected) {
            alert("Select component first");
            return;
        }

        if (confirm("Delete selected component ?")) {
            selected.remove();
        }
    }
});

editor.Commands.add("new-component", {
    run(editor) {
        createNewComponent();
    }
});

// ---------------------
// Clean Up CSS
// ---------------------
function cleanupCss() {
    const wrapper = editor.getWrapper();

    const rulesToRemove = [];

    wrapper.find("[data-field-id]").forEach(comp => {
        const attrs = comp.getAttributes();

        //----------------------------------
        // Use modified text if present
        //----------------------------------

        let fieldName = comp.get("pendingFieldName") || attrs["data-field-id"];

        if (!fieldName) return;

        fieldName = fieldName
            .trim()
            .replace(/\s+/g, "_")
            .replace(/[^\w-]/g, "");

        //----------------------------------
        // Previous class name
        //----------------------------------

        const oldClass = attrs["data-field-id"];

        //----------------------------------
        // Merge all duplicate css rules
        //----------------------------------

        let classStyle = {};

        const matchingRules = [];

        editor.Css.getRules().forEach(rule => {
            const selectors = rule.getSelectors();

            if (!selectors) return;

            selectors.forEach(sel => {
                const name = sel.getFullName();

                if (
                    name === "." + oldClass ||
                    name === oldClass || // for body selector
                    name === "body"
                ) {
                    matchingRules.push(rule);
                }
            });
        });

        // Merge in order
        matchingRules.forEach(rule => {
            classStyle = {
                ...classStyle,
                ...rule.getStyle()
            };
        });

        //----------------------------------
        // Runtime incremental styles
        // (drag / resize changes)
        //----------------------------------

        const currentStyle = comp.getStyle() || {};

        //----------------------------------
        // Merge
        //----------------------------------

        const mergedStyle = {
            ...classStyle,

            ...currentStyle
        };

        //----------------------------------
        // Strip selection-outline artifacts - this is a purely
        // visual editing aid (see setSelectionOutline) and must
        // never be persisted into saved CSS, regardless of how
        // it ended up in an existing rule.
        //----------------------------------

        delete mergedStyle["outline"];
        delete mergedStyle["outline-color"];
        delete mergedStyle["outline-style"];
        delete mergedStyle["outline-width"];
        delete mergedStyle["outlineColor"];
        delete mergedStyle["outlineStyle"];
        delete mergedStyle["outlineWidth"];

        //----------------------------------
        // Remove old css rule
        //----------------------------------

        if (oldClass !== fieldName) {
            matchingRules.forEach(rule => {
                editor.Css.remove(rule);
            });
        }

        //----------------------------------
        // Create new css rule
        //----------------------------------

        editor.Css.setRule(
            "." + fieldName,

            mergedStyle
        );

        //----------------------------------
        // Sync component names
        //----------------------------------

        // IMPORTANT: comp.setClass() *resets* the component's entire
        // class collection, discarding any other classes (e.g. a
        // shared layout/base class) the element depends on for its
        // position/size. Only rename the field-specific class and
        // leave everything else untouched.
        if (oldClass && oldClass !== fieldName) {
            comp.removeClass(oldClass);
        }

        comp.addClass(fieldName);

        // attrs (captured above) still contains the OLD id/class/
        // data-field-id values from before this rename. Strip them
        // all out and rebuild the attributes object in the desired
        // order (data-field-id, id, then everything else) - using
        // comp.set("attributes", ...) directly for a full replace,
        // rather than setAttributes(), since setAttributes() may
        // internally merge with the existing attributes object and
        // preserve ITS original key order regardless of what order we
        // pass values in here.
        const {
            class: _staleClass,
            id: _staleId,
            "data-field-id": _staleFieldId,
            ...restAttrs
        } = attrs;

        comp.set("attributes", {
            "data-field-id": fieldName,
            id: fieldName,
            ...restAttrs
        });

        comp.setName(fieldName);

        //----------------------------------
        // Remove temp value
        //----------------------------------

        comp.unset("pendingFieldName");
    });

    //----------------------------------
    // Collect dynamic rules
    //----------------------------------

    editor.Css.getRules().forEach(rule => {
        const selectors = rule.getSelectors();

        if (!selectors) return;

        selectors.forEach(sel => {
            const name = sel.getFullName();

            if (name.startsWith("#")) {
                rulesToRemove.push(rule);
            }
        });
    });

    //----------------------------------
    // Remove afterwards
    //----------------------------------

    rulesToRemove.forEach(rule => {
        editor.Css.remove(rule);
    });

    //----------------------------------
    // Merge all duplicate CSS rules
    //----------------------------------

    const mergedRules = {};

    // Snapshot current rules first
    const allRules = [...editor.Css.getRules()];

    // Build merged styles
    allRules.forEach(rule => {
        const selectors = rule.getSelectors();

        if (!selectors) return;

        selectors.forEach(sel => {
            const selector = sel.getFullName();

            if (!selector) return;

            if (!mergedRules[selector]) {
                mergedRules[selector] = {};
            }

            // Last occurrence wins
            Object.assign(mergedRules[selector], rule.getStyle());
        });
    });

    //----------------------------------
    // Strip selection-outline artifacts here too - a stray duplicate
    // rule object sharing the same selector (e.g. left behind by
    // GrapesJS's own state/hover handling while a field was selected)
    // can carry the outline back in via the merge above, even though
    // the per-field loop already deleted it from its own rule.
    //----------------------------------

    Object.values(mergedRules).forEach(style => {
        delete style["outline"];
        delete style["outline-color"];
        delete style["outline-style"];
        delete style["outline-width"];
        delete style["outlineColor"];
        delete style["outlineStyle"];
        delete style["outlineWidth"];
    });

    //----------------------------------
    // Remove ALL existing rules first
    //----------------------------------

    allRules.forEach(rule => {
        editor.Css.remove(rule);
    });

    //----------------------------------
    // Recreate merged rules
    //----------------------------------

    Object.entries(mergedRules).forEach(([selector, style]) => {
        editor.Css.setRule(selector, style);
    });
}

// -------------------------
// Zoom
// -------------------------

function updateZoomDisplay() {
    const btn = document.getElementById("btnZoom100");

    if (btn) {
        btn.textContent = Math.round(zoom * 100) + "%";
    }
}

function applyZoom() {
    editor.Canvas.setZoom(zoom * 100);

    updateZoomDisplay();
}

// -------------------------
// Variable fields
// -------------------------

function makeEditableFields() {
    const comps = editor.getWrapper().find("*");

    comps.forEach(comp => {
        const el = comp.getEl();

        if (!el) return;

        // Background image
        if (el.tagName === "IMG" && el.classList.contains("bg-image")) {
            comp.set({
                selectable: false,
                draggable: false,
                editable: false,
                hoverable: false,
                highlightable: false,
                removable: false,
                copyable: false
            });

            return;
        }

        // Container ONLY:
        // Make it transparent to selection
        if (el.classList && el.classList.contains("template-container")) {
            comp.set({
                selectable: false,
                hoverable: false,
                highlightable: false,
                draggable: false,
                layerable: false,
                badgable: false
            });

            // Critical
            el.style.pointerEvents = "none";

            return;
        }

        // Actual variable fields
        if (el.hasAttribute("data-field-id")) {
            const field = el.getAttribute("data-field-id");

            console.log("ENABLE:", field);

            // Restore clickability
            el.style.pointerEvents = "auto";

            comp.set({
                selectable: true,
                draggable: true,
                editable: false, // manual dblclick-to-edit handles this uniformly instead
                droppable: false,
                hoverable: true,
                highlightable: true,
                removable: true,
                copyable: true,
                badgable: true
            });
        }
    });

    editor.select(null);
}

// -------------------------
// Position Field Coordinates
// -------------------------

function updatePositionFields(comp) {
    if (!comp) return;

    const el = comp.getEl();

    if (!el) return;

    const rect = el.getBoundingClientRect();

    const parentRect = el.offsetParent
        ? el.offsetParent.getBoundingClientRect()
        : { left: 0, top: 0 };

    const x = Math.round(rect.left - parentRect.left);

    const y = Math.round(rect.top - parentRect.top);

    const w = Math.round(rect.width);

    const h = Math.round(rect.height);

    document.getElementById("txtX").value = x;

    document.getElementById("txtY").value = y;

    document.getElementById("txtW").value = w;

    document.getElementById("txtH").value = h;
}

// -------------------------
// Apply Position Changed Coordinates
// -------------------------

// -------------------------
// CreateNewComponent function
// -------------------------

function createNewComponent() {
    const wrapper = editor.DomComponents.getWrapper();

    // New fields must live inside .template-container, same as every
    // other field, so they inherit the shared position/z-index/
    // pointer-events rule and are positioned relative to the
    // template instead of the page root.
    const container = wrapper.find(".template-container")[0] || wrapper;

    const id = "NEW_" + Date.now();

    const added = container.append({
        type: "text",

        editable: false, // manual dblclick-to-edit handles this uniformly instead
        selectable: true,
        hoverable: true,
        highlightable: true,
        draggable: true,
        droppable: false,

        content: "&lt;NEW&gt;",

        attributes: {
            "data-field-id": id,
            id: id,
            class: id
        },

        style: {
            position: "absolute",
            left: "100px",
            top: "100px",
            width: "100px",
            height: "30px",
            border: "1px dashed gray",
            padding: "2px",

            // Every OTHER field gets these three from either
            // makeEditableFields()'s explicit inline overrides (after
            // Open) or its own persisted CSS rule (after Save). A
            // brand-new field has been through neither, so it was
            // silently inheriting pointer-events:none from
            // .template-container with nothing to override it -
            // making it completely unclickable. Set them explicitly
            // instead of relying on the external editor.css rule,
            // which isn't reliably in effect inside the canvas iframe.
            "pointer-events": "auto",
            "z-index": "9999",
            display: "inline-block"
        }
    });

    // append() returns an ARRAY of added components, not a single
    // component - grab the first (only) one out of it.
    const comp = Array.isArray(added) ? added[0] : added;

    // Belt-and-suspenders: also force it as an inline style directly
    // on the live DOM element, same as makeEditableFields() does for
    // every other field, so it's not solely dependent on the model
    // style having been correctly translated to CSS.
    const newEl = comp && comp.getEl && comp.getEl();
    if (newEl) {
        newEl.style.pointerEvents = "auto";
        newEl.style.zIndex = "9999";
    }

    // editable stays false (set above) - manual dblclick-to-edit
    // handles text editing uniformly for this and every other field,
    // instead of relying on GrapesJS's own type-based RTE trigger.
}

// -------------------------
// Open
// -------------------------

document.getElementById("btnOpen").addEventListener("click", async () => {
    const result = await window.api.openHtml();
    if (!result) return;
    currentFile = result.path;

    // Build the absolute folder path as a file:// base
    const normalizedPath = result.path.replace(/\\/g, "/");
    const folder = normalizedPath.substring(0, normalizedPath.lastIndexOf("/") + 1);
    const baseHref = "file:///" + folder;

    const parser = new DOMParser();
    const doc = parser.parseFromString(result.content, "text/html");

    // Rewrite relative asset paths (img, video, source, etc.) to absolute file:// paths
    doc.querySelectorAll("[src]").forEach(el => {
        const src = el.getAttribute("src");
        if (src && !/^(https?:|data:|file:)/i.test(src)) {
            el.setAttribute("src", baseHref + src);
        }
    });

    // Also fix relative background-image urls inside inline style attributes
    doc.querySelectorAll("[style]").forEach(el => {
        const style = el.getAttribute("style");
        if (style && style.includes("url(")) {
            const fixed = style.replace(/url\(["']?(?!https?:|data:|file:)([^"')]+)["']?\)/gi,
                (match, p1) => `url("${baseHref}${p1}")`);
            el.setAttribute("style", fixed);
        }
    });

    let css = "";
    doc.querySelectorAll("style").forEach(style => {
        // Fix relative background-image urls inside <style> blocks too
        const fixedCss = style.textContent.replace(/url\(["']?(?!https?:|data:|file:)([^"')]+)["']?\)/gi,
            (match, p1) => `url("${baseHref}${p1}")`);
        css += fixedCss + "\n";
    });

    editor.select(null);
    editor.setComponents("");
    editor.setComponents(doc.body.innerHTML);
    editor.setStyle(css);

    setTimeout(() => {
        makeEditableFields();
        const root = editor.getWrapper();
        if (root) {
            root.set({
                selectable: false,
                hoverable: false,
                highlightable: false,
                badgable: false
            });
        }
        editor.select(null);
    }, 500);
});

// -------------------------
// Reorder fields to match visual reading order (top-to-bottom,
// left-to-right) before saving, so a newly-created field's <div>
// lands in the HTML wherever it visually belongs among the others,
// instead of always at the end (where it got appended on creation).
// -------------------------

function reorderFieldsByPosition() {
    const wrapper = editor.DomComponents.getWrapper();
    const container = wrapper.find(".template-container")[0];

    if (!container) return;

    const fields = container
        .components()
        .filter(c => c.getAttributes()["data-field-id"]);

    const withPos = fields.map(comp => {
        const style = comp.getStyle();
        const el = comp.getEl();
        const computed = el ? getComputedStyle(el) : {};

        return {
            comp,
            top: parseInt(style.top || computed.top) || 0,
            left: parseInt(style.left || computed.left) || 0
        };
    });

    // Reading order: top first, then left to break ties on the same row.
    withPos.sort((a, b) => a.top - b.top || a.left - b.left);

    withPos.forEach((item, i) => {
        item.comp.move(container, { at: i });
    });
}

// ---------------------------------
// Save without Prettier formatting
// ---------------------------------

document.getElementById("btnSave").addEventListener("click", async () => {
    if (!currentFile) return;

    cleanupCss();
    reorderFieldsByPosition();

    const html = `<!DOCTYPE html>
<html>
<head>
<style>
${editor.getCss()}
</style>
</head>
<body>
${editor.getHtml()}
</body>
</html>`;

    await window.api.saveHtml({
        path: currentFile,

        content: html
    });
});

// -------------------------
// Undo
// -------------------------

document.getElementById("btnUndo").addEventListener("click", () => {
    editor.UndoManager.undo();
});

// -------------------------
// Redo
// -------------------------

document.getElementById("btnRedo").addEventListener("click", () => {
    editor.UndoManager.redo();
});

// --------------------------------------------------
// keyboard shortcuts Ctrl+Z, Ctrl+Y, Ctrl+Shift+Z are now handled by
// handleGlobalKeydown (registered on both the parent document and
// the canvas iframe's document) - see inside editor.on("load", ...).
// --------------------------------------------------

// -------------------------
// New
// -------------------------

document.getElementById("btnNew").addEventListener("click", () => {
    createNewComponent();
});

// -------------------------
// Delete
// -------------------------

document.getElementById("btnDelete").addEventListener("click", () => {
    const selected = editor.getSelected();

    if (!selected) {
        alert("Please select a component first");

        return;
    }

    const content = selected.get("content") || "";

    const answer = confirm(`Delete ${content} ?`);

    if (!answer) return;

    selected.remove();

    editor.select(null);
});

// -------------------------
// Zoom buttons
// -------------------------

document.getElementById("btnZoomIn").onclick = () => {
    zoom += 0.1;

    applyZoom();
};

document.getElementById("btnZoomOut").onclick = () => {
    zoom -= 0.1;

    if (zoom < 0.2) zoom = 0.2;

    applyZoom();
};

document.getElementById("btnZoom100").onclick = () => {
    zoom = 1;

    applyZoom();
};

updateZoomDisplay();
