package docs

import (
	"embed"
	_ "embed"
)

//go:embed *
var StaticFiles embed.FS
