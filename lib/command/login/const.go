package login

import (
	_ "embed"
)

//go:embed auth.html
var AUTH_PAGE_HTML []byte

//go:embed code.html
var CODE_PAGE_HTML []byte
