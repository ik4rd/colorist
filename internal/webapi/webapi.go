package webapi

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"image"
	"io"
	"net/http"

	"github.com/ik4rd/colorist/internal/colormap"
	_ "github.com/ik4rd/colorist/internal/colormap/algorithms"
	"github.com/ik4rd/colorist/internal/imageio"
	"github.com/ik4rd/colorist/internal/logger"
)

const maxUploadBytes = 32 << 20 // 32 MiB

type Service struct {
	log   *logger.Logger
	store *store
}

func New(log *logger.Logger, maxImages int) *Service {
	return &Service{log: log, store: newStore(maxImages)}
}

func (svc *Service) Register(mux *http.ServeMux) {
	mux.HandleFunc("/upload", svc.uploadHandler())
	mux.HandleFunc("/render", svc.renderHandler())
}

func Register(mux *http.ServeMux, log *logger.Logger, maxImages int) {
	New(log, maxImages).Register(mux)
}

func (svc *Service) Upload(data []byte) (id string, width, height int, err error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return "", 0, 0, fmt.Errorf("decode: %w", err)
	}

	px, err := colormap.NewPixels(img)
	if err != nil {
		return "", 0, 0, err
	}

	b := img.Bounds()
	return svc.store.put(px), b.Dx(), b.Dy(), nil
}

func (svc *Service) uploadHandler() http.HandlerFunc {
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

		data, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, "read: "+err.Error(), http.StatusBadRequest)
			return
		}

		id, width, height, err := svc.Upload(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
			return
		}

		writeJSON(w, svc.log, map[string]any{"id": id, "width": width, "height": height})
	}
}

type renderRequest struct {
	ID   string           `json:"id"`
	Opts colormap.Options `json:"opts"`
}

func (svc *Service) renderHandler() http.HandlerFunc {
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

		px, ok := svc.store.get(req.ID)
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
