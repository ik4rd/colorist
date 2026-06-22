import { pushHistory } from "./history.js";

let img = null;
let highlight = null;
let tooltip = null;
let getFormat = () => "hex";
let getView = () => null;

let regions = [];
let dims = { w: 0, h: 0 };
let flashTimer = null;

export function initPicker(imgEl, stageEl, formatGetter, viewGetter) {
  img = imgEl;
  getFormat = formatGetter;
  getView = viewGetter || (() => null);

  highlight = document.createElement("div");
  highlight.className = "region-highlight";
  highlight.hidden = true;
  stageEl.appendChild(highlight);

  tooltip = document.createElement("div");
  tooltip.className = "region-tip";
  tooltip.hidden = true;
  stageEl.appendChild(tooltip);

  img.addEventListener("mousemove", onMove);
  img.addEventListener("mouseleave", hideHover);
  img.addEventListener("click", onClick);
}

export function setRegions(list, width, height) {
  regions = Array.isArray(list) ? list : [];
  dims = { w: width, h: height };
  hideHover();
}

function viewRect() {
  return getView() || { x: 0, y: 0, w: dims.w, h: dims.h };
}

function cursorSource(e) {
  const cw = img.clientWidth;
  const ch = img.clientHeight;
  if (!cw || !ch) return null;

  const sr = img.offsetParent.getBoundingClientRect();
  const fx = (e.clientX - sr.left - img.offsetLeft) / cw;
  const fy = (e.clientY - sr.top - img.offsetTop) / ch;
  if (fx < 0 || fx > 1 || fy < 0 || fy > 1) return null;

  const v = viewRect();

  return { sx: v.x + fx * v.w, sy: v.y + fy * v.h };
}

function regionAt(e) {
  if (!regions.length || !dims.w || !dims.h) return null;

  const p = cursorSource(e);
  if (!p) return null;

  for (const r of regions) {
    if (p.sx >= r.x && p.sx < r.x + r.w && p.sy >= r.y && p.sy < r.y + r.h) {
      return r;
    }
  }

  return null;
}

function valueFor(r) {
  switch (getFormat()) {
    case "rgb":
      return r.rgb;
    case "cmyk":
      return r.cmyk;
    case "names":
      return r.name;
    default:
      return r.hex;
  }
}

function onMove(e) {
  const r = regionAt(e);
  if (!r) {
    hideHover();
    return;
  }

  const v = viewRect();
  const kx = img.clientWidth / v.w;
  const ky = img.clientHeight / v.h;

  highlight.style.left = img.offsetLeft + (r.x - v.x) * kx + "px";
  highlight.style.top = img.offsetTop + (r.y - v.y) * ky + "px";
  highlight.style.width = r.w * kx + "px";
  highlight.style.height = r.h * ky + "px";
  highlight.hidden = false;

  if (!flashTimer) {
    const stage = img.offsetParent;
    const sr = stage.getBoundingClientRect();
    tooltip.textContent = valueFor(r);
    tooltip.classList.remove("copied");
    tooltip.style.left = e.clientX - sr.left + 14 + "px";
    tooltip.style.top = e.clientY - sr.top + 14 + "px";
    tooltip.hidden = false;
  }
}

function hideHover() {
  if (highlight) highlight.hidden = true;
  if (tooltip && !flashTimer) tooltip.hidden = true;
}

async function onClick(e) {
  const r = regionAt(e);
  if (!r) return;

  const text = valueFor(r);
  const ok = await copy(text);

  if (ok) pushHistory(r.hex, text);

  flash(ok ? "copied " + text : "copy failed", e);
}

function flash(message, e) {
  clearTimeout(flashTimer);

  const stage = img.offsetParent;
  const sr = stage.getBoundingClientRect();
  tooltip.textContent = message;
  tooltip.classList.add("copied");
  tooltip.style.left = e.clientX - sr.left + 14 + "px";
  tooltip.style.top = e.clientY - sr.top + 14 + "px";
  tooltip.hidden = false;

  flashTimer = setTimeout(() => {
    flashTimer = null;
    tooltip.classList.remove("copied");
    tooltip.hidden = true;
  }, 900);
}

async function copy(text) {
  try {
    if (navigator.clipboard?.writeText) {
      await navigator.clipboard.writeText(text);
      return true;
    }
  } catch {
    // :)
  }

  const ta = document.createElement("textarea");
  ta.value = text;
  ta.style.position = "fixed";
  ta.style.opacity = "0";
  document.body.appendChild(ta);
  ta.select();

  let ok = false;
  try {
    ok = document.execCommand("copy");
  } catch {
    ok = false;
  }

  document.body.removeChild(ta);

  return ok;
}
