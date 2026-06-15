export function applyTheme(theme) {
  if (!theme) return;

  applyPalette(theme.palette);
  applyAccent(theme.accent);
}

function applyPalette(colors) {
  if (!Array.isArray(colors) || !colors.length) return;

  const root = document.documentElement;
  colors.forEach((c, i) =>
    root.style.setProperty(`--c${i + 1}`, `${c[0]}, ${c[1]}, ${c[2]}`),
  );
}

function applyAccent(a) {
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
