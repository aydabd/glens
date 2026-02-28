module glens/tools/api

go 1.25

require (
	github.com/rs/zerolog v1.34.0
	glens/pkg/logging v0.0.0
)

require (
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	golang.org/x/sys v0.36.0 // indirect
)

replace glens/pkg/logging => ../../pkg/logging
