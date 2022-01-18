// +build database,ora database,all

package main

import (
	//To use ora import use go 1.6+
	_ "github.com/evocert/kwe/database/ora"
)