const clamp = (v, lo, hi) => Math.min(hi, Math.max(lo, v));

function dist2(a, b) {
  const dr = a[0] - b[0],
    dg = a[1] - b[1],
    db = a[2] - b[2];
  return dr * dr + dg * dg + db * db;
}

function rgbToHsl(r, g, b) {
  r /= 255;
  g /= 255;
  b /= 255;

  const max = Math.max(r, g, b),
    min = Math.min(r, g, b);
  const l = (max + min) / 2;

  if (max === min) return [0, 0, l];

  const d = max - min;
  const s = l > 0.5 ? d / (2 - max - min) : d / (max + min);

  let h;

  switch (max) {
    case r:
      h = (g - b) / d + (g < b ? 6 : 0);
      break;
    case g:
      h = (b - r) / d + 2;
      break;
    default:
      h = (r - g) / d + 4;
  }

  return [h / 6, s, l];
}

function hslToRgb(h, s, l) {
  if (s === 0) {
    const v = Math.round(l * 255);
    return [v, v, v];
  }

  const hue = (p, q, t) => {
    if (t < 0) t += 1;
    if (t > 1) t -= 1;
    if (t < 1 / 6) return p + (q - p) * 6 * t;
    if (t < 1 / 2) return q;
    if (t < 2 / 3) return p + (q - p) * (2 / 3 - t) * 6;
    return p;
  };

  const q = l < 0.5 ? l * (1 + s) : l + s - l * s;
  const p = 2 * l - q;

  return [hue(p, q, h + 1 / 3), hue(p, q, h), hue(p, q, h - 1 / 3)].map((v) =>
    Math.round(v * 255),
  );
}

export function extractPalette(img, count = 4) {
  const MAX = 96;
  const scale = Math.min(1, MAX / Math.max(img.width, img.height));
  const w = Math.max(1, Math.round(img.width * scale));
  const h = Math.max(1, Math.round(img.height * scale));

  const cv = document.createElement("canvas");
  cv.width = w;
  cv.height = h;

  const ctx = cv.getContext("2d", { willReadFrequently: true });
  ctx.drawImage(img, 0, 0, w, h);

  const data = ctx.getImageData(0, 0, w, h).data;

  const buckets = new Map();

  for (let i = 0; i < data.length; i += 4) {
    if (data[i + 3] < 128) continue;

    const r = data[i],
      g = data[i + 1],
      b = data[i + 2];

    const key = ((r >> 3) << 10) | ((g >> 3) << 5) | (b >> 3);

    let e = buckets.get(key);

    if (!e) {
      e = { r: 0, g: 0, b: 0, n: 0 };
      buckets.set(key, e);
    }

    e.r += r;
    e.g += g;
    e.b += b;
    e.n++;
  }

  const cand = [...buckets.values()]
    .map((e) => ({
      rgb: [
        Math.round(e.r / e.n),
        Math.round(e.g / e.n),
        Math.round(e.b / e.n),
      ],
      n: e.n,
    }))
    .sort((a, b) => b.n - a.n)
    .slice(0, 48);

  if (cand.length === 0) return [];

  const chosen = [cand[0].rgb];

  while (chosen.length < count && chosen.length < cand.length) {
    let best = null,
      bestScore = -1;

    for (const c of cand) {
      let dmin = Infinity;
      for (const s of chosen) dmin = Math.min(dmin, dist2(c.rgb, s));
      if (dmin > bestScore) {
        bestScore = dmin;
        best = c.rgb;
      }
    }

    if (!best || bestScore === 0) break;

    chosen.push(best);
  }

  while (chosen.length < count) chosen.push(chosen[chosen.length - 1]);

  return chosen;
}

export function deriveAccent(palette) {
  let best = null,
    bestScore = 0;

  for (const c of palette) {
    const [h, s, l] = rgbToHsl(...c);
    const score = s * (1 - Math.abs(l - 0.5));
    if (score > bestScore) {
      bestScore = score;
      best = [h, s, l];
    }
  }

  if (!best || best[1] < 0.12) return null;

  const [h, s] = best;
  const sat = clamp(Math.max(s, 0.5), 0, 1);

  return {
    base: hslToRgb(h, sat, clamp(best[2], 0.34, 0.44)),
    deep: hslToRgb(h, sat, clamp(best[2] - 0.1, 0.22, 0.4)),
    bright: hslToRgb(h, Math.min(1, sat + 0.1), 0.58),
    soft: hslToRgb(h, Math.min(0.45, sat * 0.4 + 0.08), 0.92),
  };
}
