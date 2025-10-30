package main

import "embed"

//go:embed all:static
var staticFS embed.FS
