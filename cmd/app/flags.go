package main

import (
	"flag"
	"fmt"

	"github.com/ik4rd/colorist/internal/colormap"
)

type config struct {
	input  string
	output string
	opts   colormap.Options
}

func parseFlags(args []string) (config, error) {
	fs := flag.NewFlagSet("colorist", flag.ContinueOnError)

	cfg := config{opts: colormap.DefaultOptions()}

	fs.StringVar(&cfg.input, "input", "", "path to input image (.png/.jpg)")
	fs.StringVar(&cfg.output, "output", "", "path to output image (.png/.jpg)")
	fs.StringVar(&cfg.opts.Algorithm, "algorithm", cfg.opts.Algorithm,
		fmt.Sprintf("partition algorithm %v", colormap.Algorithms()))
	fs.Float64Var(&cfg.opts.Threshold, "threshold", cfg.opts.Threshold,
		"avg variance threshold (0..255), lower = more detail")
	fs.Float64Var(&cfg.opts.Detail, "detail", cfg.opts.Detail,
		"RGB distance above which a pixel is an outlier feature; 0 disables")
	fs.Float64Var(&cfg.opts.DetailFrac, "detail-frac", cfg.opts.DetailFrac,
		"min fraction of outlier pixels (0..1) to force a split")
	fs.IntVar(&cfg.opts.MinSize, "min-size", cfg.opts.MinSize, "minimum rectangle side in pixels")
	fs.IntVar(&cfg.opts.MaxDepth, "max-depth", cfg.opts.MaxDepth, "maximum recursion depth")
	fs.IntVar(&cfg.opts.Gap, "gap", cfg.opts.Gap, "gap between rectangles in pixels")
	fs.IntVar(&cfg.opts.HalvesPerAxis, "halves-per-axis", cfg.opts.HalvesPerAxis,
		"quadtree: splits per axis (2 = 2x2 quadrants)")
	fs.BoolVar(&cfg.opts.Labels, "labels", cfg.opts.Labels, "draw hex labels on large rectangles")
	fs.BoolVar(&cfg.opts.ColorNames, "color-names", cfg.opts.ColorNames,
		"also draw the color name on large rectangles")

	if err := fs.Parse(args); err != nil {
		return config{}, err
	}
	if cfg.input == "" || cfg.output == "" {
		return config{}, fmt.Errorf("both --input and --output are required")
	}

	return cfg, nil
}
