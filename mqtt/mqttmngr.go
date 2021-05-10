package mqtt

import (
	"fmt"
	"sync"
)

type MQTTManager struct {
	lck   *sync.RWMutex
	cntns map[string]*MQTTConnection
}

func NewMQTTManager() (mqttmngr *MQTTManager) {
	mqttmngr = &MQTTManager{lck: &sync.RWMutex{}, cntns: map[string]*MQTTConnection{}}
	return
}

func (mqttmngr *MQTTManager) Connections() (aliases []string) {
	if mqttmngr != nil {
		if len(mqttmngr.cntns) > 0 {
			func() {
				mqttmngr.lck.RLock()
				defer mqttmngr.lck.RUnlock()
				aliases = make([]string, len(mqttmngr.cntns))
				aliasi := 0
				for alias := range mqttmngr.cntns {
					aliases[aliasi] = alias
					aliasi++
				}
			}()
		}
	}
	return
}

func (mqttmngr *MQTTManager) RegisterConnection(alias string, a ...interface{}) {
	if alias != "" {
		if !func() (exists bool) {
			mqttmngr.lck.RLock()
			defer mqttmngr.lck.RUnlock()
			_, exists = mqttmngr.cntns[alias]
			return
		}() {
			func() {
				mqttmngr.lck.Lock()
				defer mqttmngr.lck.Unlock()
				if mqttcn := NewMQTTConnections(alias, a...); mqttcn != nil {
					mqttmngr.cntns[alias] = mqttcn
					mqttcn.mqttmngr = mqttmngr
				}
			}()
		}
	}
}

func (mqttmngr *MQTTManager) MessageReceived(alias string, msg Message) {
	fmt.Printf("%s:Received message: %s from topic: %s\n", alias, msg.Payload(), msg.Topic())
}

func (mqttmngr *MQTTManager) Connected(alias string) {
	fmt.Println("Connected:" + alias)
}

func (mqttmngr *MQTTManager) Disconnected(alias string, err error) {
	if err != nil {
		fmt.Println("Disconnected:" + alias + "=> " + err.Error())
	} else {
		fmt.Println("Disconnected:" + alias)
	}
}

func (mqttmngr *MQTTManager) IsConnect(alias string) (connected bool) {
	if alias != "" {
		if exsist, mqttnc := func() (exists bool, mqttcn *MQTTConnection) {
			mqttmngr.lck.RLock()
			defer mqttmngr.lck.RUnlock()
			mqttcn, exists = mqttmngr.cntns[alias]
			return
		}(); exsist {
			func() {
				connected = mqttnc.IsConnected()
			}()
		}
	}
	return
}

func (mqttmngr *MQTTManager) Connect(alias string) (err error) {
	if alias != "" {
		if exsist, mqttnc := func() (exists bool, mqttcn *MQTTConnection) {
			mqttmngr.lck.RLock()
			defer mqttmngr.lck.RUnlock()
			mqttcn, exists = mqttmngr.cntns[alias]
			return
		}(); exsist {
			func() {
				err = mqttnc.Connect()
			}()
		}
	}
	return
}

func (mqttmngr *MQTTManager) Disconnect(alias string, quiesce uint) (err error) {
	if alias != "" {
		if exsist, mqttnc := func() (exists bool, mqttcn *MQTTConnection) {
			mqttmngr.lck.RLock()
			defer mqttmngr.lck.RUnlock()
			mqttcn, exists = mqttmngr.cntns[alias]
			return
		}(); exsist {
			func() {
				err = mqttnc.Disconnect(quiesce)
			}()
		}
	}
	return
}

func (mqttmngr *MQTTManager) Subscribe(alias string, topic string, qos byte) (err error) {
	if alias != "" {
		if exsist, mqttnc := func() (exists bool, mqttcn *MQTTConnection) {
			mqttmngr.lck.RLock()
			defer mqttmngr.lck.RUnlock()
			mqttcn, exists = mqttmngr.cntns[alias]
			return
		}(); exsist {
			func() {
				err = mqttnc.Subscribe(topic, qos)
			}()
		}
	}
	return
}

func (mqttmngr *MQTTManager) Publish(alias string, topic string, qos byte, retained bool, message string) (err error) {
	if alias != "" {
		if exsist, mqttnc := func() (exists bool, mqttcn *MQTTConnection) {
			mqttmngr.lck.RLock()
			defer mqttmngr.lck.RUnlock()
			mqttcn, exists = mqttmngr.cntns[alias]
			return
		}(); exsist {
			func() {
				err = mqttnc.Publish(topic, qos, retained, message)
			}()
		}
	}
	return
}

var gblmqttmngr *MQTTManager

func GLOBALMQTTMANAGER() *MQTTManager {
	return gblmqttmngr
}

func init() {
	gblmqttmngr = NewMQTTManager()
}
