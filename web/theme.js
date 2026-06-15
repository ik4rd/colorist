import { extractPalette, deriveAccent } from "./color.js";

function loadImage(file) {
  return new Promise((resolve, reject) => {
    const url = URL.createObjectURL(file);
    const img = new Image();

    img.onload = () => {
      URL.revokeObjectURL(url);
      resolve(img);
    };

    img.onerror = (e) => {
      URL.revokeObjectURL(url);
      reject(e);
    };

    img.src = url;
  });
}

function applyPalette(colors) {
  if (!colors.length) return;
  const root = document.documentElement;
  colors.forEach((c, i) =>
    root.style.setProperty(`--c${i + 1}`, `${c[0]}, ${c[1]}, ${c[2]}`),
  );
}

function applyAccent(palette) {
  const a = deriveAccent(palette);
  if (!a) return;

  const root = document.documentElement;
  const css = (c) => `rgb(${c[0]}, ${c[1]}, ${c[2]})`;

  root.style.setProperty("--lime", css(a.base));
  root.style.setProperty("--lime-deep", css(a.deep));
  root.style.setProperty("--lime-bright", css(a.bright));
  root.style.setProperty("--lime-soft", css(a.soft));
  root.style.setProperty(
    "--accent-rgb",
    `${a.base[0]}, ${a.base[1]}, ${a.base[2]}`,
  );
}

export function themeFromImage(file) {
  loadImage(file)
    .then((img) => {
      const palette = extractPalette(img, 4);
      applyPalette(palette);
      applyAccent(palette);
    })
    .catch(() => {});
}
