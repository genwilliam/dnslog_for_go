package web

import "embed"

//go:embed templates/* static/*
var EmbedFiles embed.FS
