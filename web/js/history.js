const MAX = 4;

let panel = null;

export function initHistory(panelEl) {
  panel = panelEl;
}

export function pushHistory(swatch, label) {
  if (!panel) return;

  const item = document.createElement("div");
  item.className = "history-item enter";

  const chip = document.createElement("span");
  chip.className = "history-chip";
  chip.style.background = swatch;

  const text = document.createElement("span");
  text.className = "history-label";
  text.textContent = label;

  item.append(chip, text);
  panel.prepend(item);

  while (panel.children.length > MAX) {
    panel.lastElementChild.remove();
  }

  requestAnimationFrame(() => item.classList.remove("enter"));

  panel.hidden = false;
}

export function clearHistory() {
  if (!panel) return;
  panel.replaceChildren();
  panel.hidden = true;
}
