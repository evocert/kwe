package messaging

import (
	"github.com/evocert/kwe/iorw"
)

//MessageManager - struct
type MessageManager struct {
	messangers map[string]*Messanger
}

//Messanger - struct
type Messanger struct {
	guid    string
	name    string
	prntrdr iorw.PrinterReader
}

//Print print to
func (msgmnr *MessageManager) Print(pr iorw.PrinterReader, a ...interface{}) (err error) {

	return
}
