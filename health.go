package health

import _ "embed"

const (
	PROGRAM = "health"
)

var (
	//go:embed VERSION
	VERSION string

	//go:embed LICENSE
	LICENSE string
)
