module glens/tools/api

go 1.25

require (
	github.com/google/uuid v1.6.0
	github.com/rs/zerolog v1.35.1
	github.com/stretchr/testify v1.12.1
	glens/pkg/logging v0.0.0
)

require (
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/sys v0.36.0 // indirect
)

replace glens/pkg/logging => ../../pkg/logging
