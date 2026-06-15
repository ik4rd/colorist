export async function uploadImage(file) {
  const upload = window.go?.main?.App?.UploadImage;
  if (upload) {
    try {
      return await upload(await fileToBase64(file));
    } catch (e) {
      throw new Error(typeof e === "string" ? e : e.message || String(e));
    }
  }

  const form = new FormData();
  form.append("image", file);

  const res = await fetch("/upload", { method: "POST", body: form });
  if (!res.ok) throw new Error(await res.text());

  return res.json();
}

function fileToBase64(file) {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onload = () => resolve(reader.result.split(",", 2)[1]);
    reader.onerror = reject;
    reader.readAsDataURL(file);
  });
}

export async function renderImage(id, opts, signal) {
  const res = await fetch("/render", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ id, opts }),
    signal,
  });
  if (!res.ok) throw new Error(await res.text());

  return res.blob();
}

export async function fetchRegions(id, opts, signal) {
  const res = await fetch("/regions", {
    method: "POST",
    headers: { "Content-Type": "application/json" },
    body: JSON.stringify({ id, opts }),
    signal,
  });
  if (!res.ok) throw new Error(await res.text());

  return res.json();
}
