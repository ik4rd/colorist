import { RENDER_DEBOUNCE_MS } from "./config.js";
import { initControls, collectOpts } from "./controls.js";
import { applyTheme, resetTheme } from "./theme.js";
import { uploadImage, renderImage, fetchRegions } from "./api.js";
import { initPicker, setRegions } from "./picker.js";
import { initHistory, clearHistory } from "./history.js";
import { initZoom, getView, setImage, onRendered, resetZoom } from "./zoom.js";
import { preventPageZoom } from "./noscale.js";

const els = {
  controls: document.getElementById("controls"),
  stage: document.getElementById("stage"),
  history: document.getElementById("history"),
  dropzone: document.getElementById("dropzone"),
  result: document.getElementById("result"),
  status: document.getElementById("status"),
  file: document.getElementById("file"),
  save: document.getElementById("save"),
  reset: document.getElementById("reset"),
};

let imageID = null;
let imageW = 0;
let imageH = 0;
let lastObjectURL = null;
let lastBlob = null;
let renderAbort = null;
let debounceTimer = null;

function setStatus(msg) {
  els.status.textContent = msg;
}

async function handleFile(file) {
  if (!file) return;

  setStatus("Uploading…");

  try {
    const info = await uploadImage(file);
    imageID = info.id;
    imageW = info.width;
    imageH = info.height;
    applyTheme(info.theme);
    setImage(imageW, imageH);
    els.stage.classList.remove("empty");
    els.reset.hidden = false;
    setStatus("");
    render();
  } catch (err) {
    setStatus("Upload failed: " + err.message);
  }
}

function resetImage() {
  if (renderAbort) renderAbort.abort();
  clearTimeout(debounceTimer);

  imageID = null;
  imageW = 0;
  imageH = 0;

  if (lastObjectURL) URL.revokeObjectURL(lastObjectURL);
  lastObjectURL = null;
  lastBlob = null;

  els.result.removeAttribute("src");
  els.result.hidden = true;
  els.save.hidden = true;
  els.save.removeAttribute("href");
  els.reset.hidden = true;
  els.file.value = "";
  els.stage.classList.add("empty");

  setRegions([], 0, 0);
  clearHistory();
  resetZoom();
  resetTheme();
  setStatus("");
}

function scheduleRender() {
  clearTimeout(debounceTimer);
  debounceTimer = setTimeout(render, RENDER_DEBOUNCE_MS);
}

async function render() {
  if (!imageID) return;

  if (renderAbort) renderAbort.abort();
  renderAbort = new AbortController();

  const opts = collectOpts();
  const view = getView();
  const isFull = !view;

  try {
    const [blob, regions] = await Promise.all([
      renderImage(imageID, opts, view, renderAbort.signal),
      fetchRegions(imageID, opts, renderAbort.signal),
    ]);
    const url = URL.createObjectURL(blob);

    els.result.onload = () => onRendered(url, isFull);
    els.result.src = url;
    els.result.hidden = false;
    els.save.href = url;
    els.save.hidden = false;

    if (lastObjectURL) URL.revokeObjectURL(lastObjectURL);

    lastObjectURL = url;
    lastBlob = blob;
    setRegions(regions, imageW, imageH);
    setStatus("");
  } catch (err) {
    if (err.name === "AbortError") return;
    setStatus("Render failed: " + err.message);
  }
}

function desktopSave() {
  return window.go?.main?.App?.SavePNG ?? null;
}

function blobToBase64(blob) {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onload = () => resolve(reader.result.split(",", 2)[1]);
    reader.onerror = reject;
    reader.readAsDataURL(blob);
  });
}

els.save.addEventListener("click", async (e) => {
  const save = desktopSave();
  if (!save) return;

  e.preventDefault();
  if (!lastBlob) return;

  try {
    await save(await blobToBase64(lastBlob));
  } catch (err) {
    setStatus("Save failed: " + err);
  }
});

els.reset.addEventListener("click", resetImage);

els.file.addEventListener("change", (e) => handleFile(e.target.files[0]));

["dragover", "dragenter"].forEach((ev) =>
  els.stage.addEventListener(ev, (e) => {
    e.preventDefault();
    if (imageID) return;
    els.stage.classList.add("dragging");
  }),
);

["dragleave", "drop"].forEach((ev) =>
  els.stage.addEventListener(ev, (e) => {
    e.preventDefault();
    els.stage.classList.remove("dragging");
  }),
);

els.stage.addEventListener("drop", (e) => {
  if (imageID) return;
  handleFile(e.dataTransfer.files[0]);
});

preventPageZoom();
initControls(els.controls, scheduleRender);
initHistory(els.history);
initPicker(els.result, els.stage, () => collectOpts().LabelFormat, getView);
initZoom(els.result, els.stage, scheduleRender);
