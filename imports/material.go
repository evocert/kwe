//go:build fonts && material || fonts && web
// +build fonts,material fonts,web

package imports

import (
	_ "github.com/evocert/kwe/fonts/material"
)