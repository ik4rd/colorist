import { CONTROLS } from "./config.js";

let container = null;
let onChange = () => {};

export function initControls(el, changeHandler) {
  container = el;
  onChange = changeHandler;
  buildControls();
  updateVisibility();
}

function buildControls() {
  for (const c of CONTROLS) {
    const row = document.createElement("div");
    row.className = "control";
    row.showFor = c.showFor || null;

    const head = document.createElement("div");
    head.className = "control-head";

    const label = document.createElement("label");
    label.textContent = c.label;
    label.htmlFor = "ctl-" + c.key;
    head.appendChild(label);

    const readout = document.createElement("span");
    readout.className = "readout";
    readout.textContent = c.value;
    head.appendChild(readout);

    row.appendChild(head);

    let input;

    if (c.type === "select") {
      input = document.createElement("select");

      for (const opt of c.options) {
        const o = document.createElement("option");
        o.value = o.textContent = opt;
        input.appendChild(o);
      }

      input.value = c.value;
      input.addEventListener("change", () => {
        readout.textContent = input.value;
        updateVisibility();
        onChange();
      });
    } else {
      input = document.createElement("input");
      input.type = "range";

      if (c.values) {
        input.valuesList = c.values;
        input.min = 0;
        input.max = c.values.length - 1;
        input.step = 1;
        input.value = c.values.indexOf(c.value);
      } else {
        input.min = c.min;
        input.max = c.max;
        input.step = c.step;
        input.value = c.value;
      }

      input.addEventListener("input", () => {
        readout.textContent = controlValue(input);
        onChange();
      });
    }

    input.id = "ctl-" + c.key;
    input.dataset.key = c.key;
    input.dataset.int = c.int ? "1" : "";
    row.appendChild(input);

    container.appendChild(row);
  }
}

function updateVisibility() {
  const algoEl = document.getElementById("ctl-Algorithm");
  const algo = algoEl ? algoEl.value : null;
  for (const row of container.children) {
    const show = !row.showFor || row.showFor.includes(algo);
    row.style.display = show ? "" : "none";
  }
}

function controlValue(input) {
  if (input.tagName === "SELECT") return input.value;
  const raw = input.valuesList
    ? input.valuesList[parseInt(input.value, 10)]
    : input.value;
  return input.dataset.int ? parseInt(raw, 10) : parseFloat(raw);
}

export function collectOpts() {
  const opts = {};
  for (const input of container.querySelectorAll("[data-key]")) {
    opts[input.dataset.key] = controlValue(input);
  }
  return opts;
}
