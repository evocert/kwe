//go:build (fonts && material) || (fonts && web) || ui
// +build fonts,material fonts,web ui

package imports

import (
	_ "github.com/evocert/kwe/fonts/material"
)
