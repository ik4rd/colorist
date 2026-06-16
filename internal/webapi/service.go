package webapi

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"image"
	"net/http"

	"github.com/ik4rd/colorist/internal/colormap"
	_ "github.com/ik4rd/colorist/internal/colormap/algorithms"
	"github.com/ik4rd/colorist/internal/logger"
	"github.com/ik4rd/colorist/internal/palette"
)

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
	mux.HandleFunc("/regions", svc.regionsHandler())
}

func Register(mux *http.ServeMux, log *logger.Logger, maxImages int) {
	New(log, maxImages).Register(mux)
}

func (svc *Service) Upload(data []byte) (UploadResult, error) {
	img, _, err := image.Decode(bytes.NewReader(data))
	if err != nil {
		return UploadResult{}, fmt.Errorf("decode: %w", err)
	}

	px, err := colormap.NewPixels(img)
	if err != nil {
		return UploadResult{}, err
	}

	b := img.Bounds()

	return UploadResult{
		ID:     svc.store.put(px),
		Width:  b.Dx(),
		Height: b.Dy(),
		Theme:  palette.FromImage(img),
	}, nil
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
