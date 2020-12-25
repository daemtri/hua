package api

import "embed"

//go:embed *.api.go
var Protocol embed.FS
