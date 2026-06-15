import { RENDER_DEBOUNCE_MS } from "./config.js";
import { initControls, collectOpts } from "./controls.js";
import { themeFromImage } from "./theme.js";
import { uploadImage, renderImage } from "./api.js";

const els = {
  controls: document.getElementById("controls"),
  stage: document.getElementById("stage"),
  dropzone: document.getElementById("dropzone"),
  result: document.getElementById("result"),
  status: document.getElementById("status"),
  file: document.getElementById("file"),
  save: document.getElementById("save"),
};

let imageID = null;
let lastObjectURL = null;
let renderAbort = null;
let debounceTimer = null;

function setStatus(msg) {
  els.status.textContent = msg;
}

async function handleFile(file) {
  if (!file) return;

  setStatus("Uploading…");
  themeFromImage(file);

  try {
    const info = await uploadImage(file);
    imageID = info.id;
    els.stage.classList.remove("empty");
    setStatus("");
    render();
  } catch (err) {
    setStatus("Upload failed: " + err.message);
  }
}

function scheduleRender() {
  clearTimeout(debounceTimer);
  debounceTimer = setTimeout(render, RENDER_DEBOUNCE_MS);
}

async function render() {
  if (!imageID) return;

  if (renderAbort) renderAbort.abort();
  renderAbort = new AbortController();

  try {
    const blob = await renderImage(imageID, collectOpts(), renderAbort.signal);
    const url = URL.createObjectURL(blob);

    els.result.src = url;
    els.result.hidden = false;
    els.save.href = url;
    els.save.hidden = false;

    if (lastObjectURL) URL.revokeObjectURL(lastObjectURL);

    lastObjectURL = url;
    setStatus("");
  } catch (err) {
    if (err.name === "AbortError") return;
    setStatus("Render failed: " + err.message);
  }
}

els.file.addEventListener("change", (e) => handleFile(e.target.files[0]));

["dragover", "dragenter"].forEach((ev) =>
  els.stage.addEventListener(ev, (e) => {
    e.preventDefault();
    els.stage.classList.add("dragging");
  }),
);

["dragleave", "drop"].forEach((ev) =>
  els.stage.addEventListener(ev, (e) => {
    e.preventDefault();
    els.stage.classList.remove("dragging");
  }),
);

els.stage.addEventListener("drop", (e) => handleFile(e.dataTransfer.files[0]));

initControls(els.controls, scheduleRender);
