package webapi

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/ik4rd/colorist/internal/colormap"
	"github.com/ik4rd/colorist/internal/imageio"
)

const maxUploadBytes = 32 << 20 // 32 MiB

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

		res, err := svc.Upload(data)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnsupportedMediaType)
			return
		}

		writeJSON(w, svc.log, res)
	}
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

		view := colormap.View{}
		if req.View != nil {
			view = *req.View
		}

		out, err := colormap.ProcessPixelsView(px, view, req.Opts)
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

func (svc *Service) regionsHandler() http.HandlerFunc {
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

		regions, err := colormap.BuildFrom(px, req.Opts)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		infos := make([]regionInfo, len(regions))
		for i, rg := range regions {
			infos[i] = regionInfo{
				X: rg.X, Y: rg.Y, W: rg.W, H: rg.H,
				Hex: rg.Hex, RGB: rg.RGB, CMYK: rg.CMYK, Name: rg.Name,
			}
		}

		w.Header().Set("Cache-Control", "no-store")
		writeJSON(w, svc.log, infos)
	}
}
