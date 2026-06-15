const MAX_ZOOM = 24;
const ZOOM_SENS = 0.0045;
const MINIMAP_MAX = 140;
const ARROW_STEP = 0.14;
const DRAG_THRESHOLD = 3;

const clamp = (v, lo, hi) => Math.min(hi, Math.max(lo, v));

let img = null;
let stage = null;
let requestRender = () => {};

let dims = { w: 0, h: 0 };
let zoom = 1;
let view = null;
let renderedView = null;
let baseURL = null;

let minimap = null;
let mmImg = null;
let mmBox = null;
let mmW = 0;
let mmH = 0;
let badge = null;

let drag = null;
let didPan = false;

export function initZoom(imgEl, stageEl, onRequestRender) {
  img = imgEl;
  stage = stageEl;
  requestRender = onRequestRender;

  minimap = document.createElement("div");
  minimap.className = "minimap";
  minimap.hidden = true;
  mmImg = document.createElement("img");
  mmImg.alt = "";
  mmBox = document.createElement("div");
  mmBox.className = "mmbox";
  minimap.append(mmImg, mmBox);

  badge = document.createElement("div");
  badge.className = "zoom-badge";
  badge.hidden = true;

  stage.append(minimap, badge);

  img.addEventListener("wheel", onWheel, { passive: false });
  img.addEventListener("mousedown", onMouseDown);
  window.addEventListener("mousemove", onMouseMove);
  window.addEventListener("mouseup", onMouseUp);
  window.addEventListener("keydown", onKeyDown);

  stage.addEventListener(
    "click",
    (e) => {
      if (didPan) {
        e.stopPropagation();
        e.preventDefault();
        didPan = false;
      }
    },
    true,
  );
}

export function getView() {
  if (!view) return null;
  return {
    x: Math.round(view.x),
    y: Math.round(view.y),
    w: Math.round(view.w),
    h: Math.round(view.h),
  };
}

export function setImage(w, h) {
  dims = { w, h };
  zoom = 1;
  view = null;
  renderedView = { x: 0, y: 0, w, h };
  baseURL = null;

  if (w >= h) {
    mmW = MINIMAP_MAX;
    mmH = Math.max(1, Math.round((MINIMAP_MAX * h) / w));
  } else {
    mmH = MINIMAP_MAX;
    mmW = Math.max(1, Math.round((MINIMAP_MAX * w) / h));
  }
  minimap.style.width = mmW + "px";
  minimap.style.height = mmH + "px";

  img.style.transform = "";
  updateOverlays();
}

export function onRendered(url, isFull) {
  img.style.transform = "";
  renderedView = view ? { ...view } : { x: 0, y: 0, w: dims.w, h: dims.h };

  if (isFull && url) {
    baseURL = url;
    mmImg.src = url;
  }

  updateOverlays();
}

export function resetZoom() {
  zoom = 1;
  view = null;
  renderedView = { x: 0, y: 0, w: dims.w, h: dims.h };
  drag = null;
  img.style.transform = "";
  img.style.cursor = "";
  updateOverlays();
}

function onWheel(e) {
  if (!dims.w || !dims.h) return;
  e.preventDefault();

  if (view && !e.ctrlKey && Math.abs(e.deltaX) > Math.abs(e.deltaY)) {
    const [dx, dy] = pixelDelta(e);
    panView(dx * (view.w / img.clientWidth), dy * (view.h / img.clientHeight));
    return;
  }

  zoomWheel(e);
}

function zoomWheel(e) {
  const cw = img.clientWidth;
  const ch = img.clientHeight;
  if (!cw || !ch) return;

  const sr = stage.getBoundingClientRect();
  const fx = clamp((e.clientX - sr.left - img.offsetLeft) / cw, 0, 1);
  const fy = clamp((e.clientY - sr.top - img.offsetTop) / ch, 0, 1);

  const cur = view || { x: 0, y: 0, w: dims.w, h: dims.h };
  const pux = cur.x + fx * cur.w; // source point under cursor
  const puy = cur.y + fy * cur.h;

  const [, dy] = pixelDelta(e);
  const nz = clamp(zoom * Math.exp(-dy * ZOOM_SENS), 1, MAX_ZOOM);
  if (nz === zoom) return;
  zoom = nz;

  if (zoom <= 1.0001) {
    zoom = 1;
    view = null;
  } else {
    const minDim = Math.min(dims.w, dims.h);
    const s = Math.max(1, minDim / zoom);
    view = clampView({ x: pux - fx * s, y: puy - fy * s, w: s, h: s });
  }

  applyPreview();
  updateOverlays();
  requestRender();
}

function onMouseDown(e) {
  if (e.button !== 0) return;
  didPan = false;
  if (!view) return;
  drag = { sx: e.clientX, sy: e.clientY, view: { ...view }, moved: false };
}

function onMouseMove(e) {
  if (!drag) return;

  const dx = e.clientX - drag.sx;
  const dy = e.clientY - drag.sy;
  if (!drag.moved && Math.hypot(dx, dy) < DRAG_THRESHOLD) return;

  drag.moved = true;
  didPan = true;
  img.style.cursor = "grabbing";

  const cw = img.clientWidth;
  const ch = img.clientHeight;
  view = clampView({
    x: drag.view.x - dx * (drag.view.w / cw),
    y: drag.view.y - dy * (drag.view.h / ch),
    w: drag.view.w,
    h: drag.view.h,
  });

  applyPreview();
  updateOverlays();
  requestRender();
}

function onMouseUp() {
  if (drag) img.style.cursor = view ? "grab" : "";
  drag = null;
}

function onKeyDown(e) {
  if (!view) return;

  const tag = document.activeElement && document.activeElement.tagName;
  if (tag === "INPUT" || tag === "SELECT" || tag === "TEXTAREA") return;

  const step = view.w * ARROW_STEP;
  let dx = 0;
  let dy = 0;
  switch (e.key) {
    case "ArrowLeft":
      dx = -step;
      break;
    case "ArrowRight":
      dx = step;
      break;
    case "ArrowUp":
      dy = -step;
      break;
    case "ArrowDown":
      dy = step;
      break;
    default:
      return;
  }

  e.preventDefault();
  panView(dx, dy);
}

function panView(dxSrc, dySrc) {
  view = clampView({
    x: view.x + dxSrc,
    y: view.y + dySrc,
    w: view.w,
    h: view.h,
  });
  applyPreview();
  updateOverlays();
  requestRender();
}

function clampView(v) {
  const s = Math.min(v.w, dims.w, dims.h);
  return {
    x: clamp(v.x, 0, dims.w - s),
    y: clamp(v.y, 0, dims.h - s),
    w: s,
    h: s,
  };
}

function pixelDelta(e) {
  let { deltaX, deltaY } = e;
  if (e.deltaMode === 1) {
    deltaX *= 16;
    deltaY *= 16;
  } else if (e.deltaMode === 2) {
    deltaX *= img.clientWidth || 800;
    deltaY *= img.clientHeight || 800;
  }

  return [deltaX, deltaY];
}

function applyPreview() {
  const v = view || { x: 0, y: 0, w: dims.w, h: dims.h };
  const r = renderedView || { x: 0, y: 0, w: dims.w, h: dims.h };
  const cw = img.clientWidth;
  const ch = img.clientHeight;

  const bw = (v.w / r.w) * cw;
  const bh = (v.h / r.h) * ch;
  if (bw <= 0 || bh <= 0) {
    img.style.transform = "";
    return;
  }

  const bx = ((v.x - r.x) / r.w) * cw;
  const by = ((v.y - r.y) / r.h) * ch;

  const s = Math.max(cw / bw, ch / bh);
  const tx = -bx * s + (cw - bw * s) / 2;
  const ty = -by * s + (ch - bh * s) / 2;

  img.style.transformOrigin = "0 0";
  img.style.transform = `translate(${tx}px, ${ty}px) scale(${s})`;
}

function updateOverlays() {
  const zoomed = !!view;

  badge.hidden = !zoomed;
  minimap.hidden = !zoomed || !baseURL;
  img.style.cursor = zoomed && !drag ? "grab" : drag ? "grabbing" : "";

  if (!zoomed) return;

  badge.textContent = zoom.toFixed(1) + "×";

  mmBox.style.left = (view.x / dims.w) * mmW + "px";
  mmBox.style.top = (view.y / dims.h) * mmH + "px";
  mmBox.style.width = (view.w / dims.w) * mmW + "px";
  mmBox.style.height = (view.h / dims.h) * mmH + "px";
}
