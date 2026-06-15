export async function uploadImage(file) {
  const form = new FormData();
  form.append("image", file);

  const res = await fetch("/upload", { method: "POST", body: form });
  if (!res.ok) throw new Error(await res.text());

  return res.json();
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
