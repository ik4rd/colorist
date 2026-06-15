let img = null;
let highlight = null;
let tooltip = null;
let getFormat = () => "hex";

let regions = [];
let dims = { w: 0, h: 0 };
let flashTimer = null;

export function initPicker(imgEl, stageEl, formatGetter) {
  img = imgEl;
  getFormat = formatGetter;

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

function regionAt(e) {
  if (!regions.length || !dims.w || !dims.h) return null;

  const box = img.getBoundingClientRect();
  if (box.width === 0 || box.height === 0) return null;

  const sx = ((e.clientX - box.left) / box.width) * dims.w;
  const sy = ((e.clientY - box.top) / box.height) * dims.h;

  for (const r of regions) {
    if (sx >= r.x && sx < r.x + r.w && sy >= r.y && sy < r.y + r.h) return r;
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

  const kx = img.clientWidth / dims.w;
  const ky = img.clientHeight / dims.h;

  highlight.style.left = img.offsetLeft + r.x * kx + "px";
  highlight.style.top = img.offsetTop + r.y * ky + "px";
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
