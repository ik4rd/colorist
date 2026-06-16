package webapi

import (
	"github.com/ik4rd/colorist/internal/colormap"
	"github.com/ik4rd/colorist/internal/palette"
)

type UploadResult struct {
	ID     string        `json:"id"`
	Width  int           `json:"width"`
	Height int           `json:"height"`
	Theme  palette.Theme `json:"theme"`
}

type renderRequest struct {
	ID   string           `json:"id"`
	Opts colormap.Options `json:"opts"`
	View *colormap.View   `json:"view,omitempty"`
}

type regionInfo struct {
	X    int    `json:"x"`
	Y    int    `json:"y"`
	W    int    `json:"w"`
	H    int    `json:"h"`
	Hex  string `json:"hex"`
	RGB  string `json:"rgb"`
	CMYK string `json:"cmyk"`
	Name string `json:"name"`
}
