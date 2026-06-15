package webapi

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"image"
	"net/http"

	"github.com/ik4rd/colorist/internal/colormap"
	_ "github.com/ik4rd/colorist/internal/colormap/algorithms"
	"github.com/ik4rd/colorist/internal/imageio"
	"github.com/ik4rd/colorist/internal/logger"
)

const maxUploadBytes = 32 << 20 // 32 MiB

func Register(mux *http.ServeMux, log *logger.Logger, maxImages int) {
	s := newStore(maxImages)
	mux.HandleFunc("/upload", uploadHandler(log, s))
	mux.HandleFunc("/render", renderHandler(s))
}

func uploadHandler(log *logger.Logger, s *store) http.HandlerFunc {
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

		id := s.put(px)
		b := img.Bounds()
		writeJSON(w, log, map[string]any{"id": id, "width": b.Dx(), "height": b.Dy()})
	}
}

type renderRequest struct {
	ID   string           `json:"id"`
	Opts colormap.Options `json:"opts"`
}

func renderHandler(s *store) http.HandlerFunc {
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

		px, ok := s.get(req.ID)
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
