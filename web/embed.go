package web

import "embed"

//go:embed index.html style.css *.js
var Assets embed.FS
