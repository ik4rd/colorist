package main

import (
	"flag"
	"net/http"
	"os"

	"github.com/ik4rd/colorist/internal/logger"
	"github.com/ik4rd/colorist/internal/webapi"
)

func main() {
	log := logger.New(os.Stderr)
	defer log.Recover()

	addr := flag.String("addr", ":8080", "listen address")
	webDir := flag.String("web", "web", "directory with static frontend files")
	maxImages := flag.Int("max-images", 16, "number of decoded images kept in memory")
	flag.Parse()

	mux := http.NewServeMux()
	webapi.Register(mux, log, *maxImages)
	mux.Handle("/", http.FileServer(http.Dir(*webDir)))

	log.Infof("listening on %s (serving %s)", *addr, *webDir)
	if err := http.ListenAndServe(*addr, mux); err != nil {
		log.Fatal(err)
	}
}
