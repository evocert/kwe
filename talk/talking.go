package talk

import (
	"github.com/evocert/kwe/chnls"
)

//Talker - struct
type Talker struct {
	chnl *chnls.Channel
	dne  chan bool
}
