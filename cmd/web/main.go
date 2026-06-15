package main

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"flag"
	"image"
	"net/http"
	"os"

	"github.com/ik4rd/colorist/internal/colormap"
	_ "github.com/ik4rd/colorist/internal/colormap/algorithms"
	"github.com/ik4rd/colorist/internal/imageio"
	"github.com/ik4rd/colorist/internal/logger"
)

func main() {
	log := logger.New(os.Stderr)
	defer log.Recover()

	addr := flag.String("addr", ":8080", "listen address")
	webDir := flag.String("web", "web", "directory with static frontend files")
	maxImages := flag.Int("max-images", 16, "number of decoded images kept in memory")
	flag.Parse()

	store := newStore(*maxImages)

	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir(*webDir)))
	mux.HandleFunc("/upload", uploadHandler(log, store))
	mux.HandleFunc("/render", renderHandler(log, store))

	log.Infof("listening on %s (serving %s)", *addr, *webDir)
	if err := http.ListenAndServe(*addr, mux); err != nil {
		log.Fatal(err)
	}
}

const maxUploadBytes = 32 << 20 // 32 MiB

func uploadHandler(log *logger.Logger, store *store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "POST only", http.StatusMethodNotAllowed)
			return
		}

		r.Body = http.MaxBytesReader(w, r.Body, maxUploadBytes)
		file, _, err := r.FormFile("image")
		if err != nil {
			http.Error(w, "missing image field: "+err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()

		img, _, err := image.Decode(file)
		if err != nil {
			http.Error(w, "decode: "+err.Error(), http.StatusUnsupportedMediaType)
			return
		}

		px, err := colormap.NewPixels(img)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		id := store.put(px)
		b := img.Bounds()
		writeJSON(w, log, map[string]any{"id": id, "width": b.Dx(), "height": b.Dy()})
	}
}

type renderRequest struct {
	ID   string           `json:"id"`
	Opts colormap.Options `json:"opts"`
}

func renderHandler(_ *logger.Logger, store *store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "POST only", http.StatusMethodNotAllowed)
			return
		}

		req := renderRequest{Opts: colormap.DefaultOptions()}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "bad json: "+err.Error(), http.StatusBadRequest)
			return
		}

		px, ok := store.get(req.ID)
		if !ok {
			http.Error(w, "unknown image id (re-upload)", http.StatusNotFound)
			return
		}

		req.Opts.Gap = 0
		req.Opts.Labels = true
		out, err := colormap.ProcessPixels(px, req.Opts)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		var buf bytes.Buffer
		if err := imageio.Encode(&buf, imageio.FormatPNG, out); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "image/png")
		w.Header().Set("Cache-Control", "no-store")
		_, _ = w.Write(buf.Bytes())
	}
}

func writeJSON(w http.ResponseWriter, log *logger.Logger, v any) {
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Infof("write json: %v", err)
	}
}

func newID() string {
	var b [12]byte
	_, _ = rand.Read(b[:])
	return hex.EncodeToString(b[:])
}
