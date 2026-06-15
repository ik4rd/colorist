const ZOOM_KEYS = new Set(["+", "-", "=", "_", "0"]);

function onWheel(e) {
  if (e.ctrlKey) e.preventDefault();
}

function onKey(e) {
  if ((e.ctrlKey || e.metaKey) && ZOOM_KEYS.has(e.key)) e.preventDefault();
}

function onGesture(e) {
  e.preventDefault();
}

export function preventPageZoom() {
  window.addEventListener("wheel", onWheel, { passive: false });
  window.addEventListener("keydown", onKey);
  document.addEventListener("gesturestart", onGesture);
  document.addEventListener("gesturechange", onGesture);
  document.addEventListener("gestureend", onGesture);
}
