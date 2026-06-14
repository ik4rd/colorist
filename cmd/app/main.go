package main

import (
	"errors"
	"flag"
	"os"

	"github.com/ik4rd/colorist/internal/colormap"
	"github.com/ik4rd/colorist/internal/imageio"
	"github.com/ik4rd/colorist/internal/logger"
)

func main() {
	log := logger.New(os.Stderr)
	defer log.Recover()

	if err := run(log, os.Args[1:]); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			return
		}
		log.Fatal(err)
	}
}

func run(log *logger.Logger, args []string) error {
	cfg, err := parseFlags(args)
	if err != nil {
		return err
	}

	img, err := imageio.Load(cfg.input)
	if err != nil {
		return err
	}

	regions, err := colormap.Build(img, cfg.opts)
	if err != nil {
		return err
	}

	b := img.Bounds()
	out := colormap.Render(regions, b.Dx(), b.Dy(), cfg.opts)

	if err := imageio.Save(cfg.output, out); err != nil {
		return err
	}

	log.Infof("%d regions -> %s", len(regions), cfg.output)

	return nil
}
