package main

import (
	"net/http"
	"os"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	"github.com/ik4rd/colorist/internal/logger"
	"github.com/ik4rd/colorist/internal/webapi"
	"github.com/ik4rd/colorist/web"
)

const maxImages = 16

func main() {
	log := logger.New(os.Stderr)
	defer log.Recover()

	svc := webapi.New(log, maxImages)
	mux := http.NewServeMux()
	svc.Register(mux)

	app := NewApp(svc)

	err := wails.Run(&options.App{
		Title:     "colorist",
		Width:     1100,
		Height:    820,
		MinWidth:  640,
		MinHeight: 520,
		AssetServer: &assetserver.Options{
			Assets:  web.Assets,
			Handler: mux,
		},
		OnStartup: app.startup,
		Bind:      []any{app},
	})
	if err != nil {
		log.Fatal(err)
	}
}
