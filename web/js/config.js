const POW2 = [1, 2, 4, 8, 16, 32, 64];

export const CONTROLS = [
  {
    type: "select",
    key: "Algorithm",
    label: "algorithm",
    value: "quadtree",
    options: ["quadtree", "bsp", "adaptive_bsp", "kdtree"],
  },
  {
    type: "range",
    key: "LabelDensity",
    label: "labels",
    values: [0, 0.25, 0.5, 0.75, 1],
    display: ["none", "25%", "50%", "75%", "all"],
    value: 1,
  },
  {
    type: "select",
    key: "LabelFormat",
    label: "format",
    value: "hex",
    options: ["hex", "rgb", "cmyk", "names"],
  },
  {
    type: "range",
    key: "HalvesPerAxis",
    label: "splits",
    values: [2, 3, 4, 5],
    value: 2,
    int: true,
    showFor: ["quadtree", "bsp"],
  },
  {
    type: "range",
    key: "Threshold",
    label: "threshold",
    values: POW2,
    value: 32,
    int: true,
  },
  {
    type: "range",
    key: "MinSize",
    label: "size",
    values: POW2,
    value: 8,
    int: true,
  },
  {
    type: "range",
    key: "MaxDepth",
    label: "depth",
    min: 1,
    max: 12,
    step: 1,
    value: 12,
    int: true,
  },
];

export const RENDER_DEBOUNCE_MS = 80;
