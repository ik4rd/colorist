package main

import (
	"context"
	"encoding/base64"
	"os"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/ik4rd/colorist/internal/webapi"
)

type App struct {
	ctx context.Context
	svc *webapi.Service
}

func NewApp(svc *webapi.Service) *App {
	return &App{svc: svc}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) UploadImage(b64 string) (map[string]any, error) {
	data, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return nil, err
	}

	id, width, height, err := a.svc.Upload(data)
	if err != nil {
		return nil, err
	}

	return map[string]any{"id": id, "width": width, "height": height}, nil
}

func (a *App) SavePNG(b64 string) (string, error) {
	data, err := base64.StdEncoding.DecodeString(b64)
	if err != nil {
		return "", err
	}

	path, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:           "Save color-map",
		DefaultFilename: "colorist.png",
		Filters: []runtime.FileFilter{
			{DisplayName: "PNG image (*.png)", Pattern: "*.png"},
		},
	})
	if err != nil || path == "" {
		return "", err
	}

	if err := os.WriteFile(path, data, 0o644); err != nil {
		return "", err
	}

	return path, nil
}
