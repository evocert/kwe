package mqtt

import (
	"fmt"
	"sync"
)

type Topic interface {
	Topic() string
	TopicPath() string
}

type activeTopic struct {
	topic     string
	topicpath string
}

func (atvtpc *activeTopic) Topic() string {
	return atvtpc.topic
}

func (atvtpc *activeTopic) TopicPath() string {
	return atvtpc.topicpath
}

type MqttMessaging func(message Message)

func (atvpc *activeTopic) processMessage(mqttmsng MqttMessaging, message Message) (err error) {
	if mqttmsng != nil {
		mqttmsng(message)
	}
	return
}

type MQTTManager struct {
	lck           *sync.RWMutex
	cntns         map[string]*MQTTConnection
	activeTopics  map[string]*activeTopic
	MqttMessaging MqttMessaging
	lcktpcs       *sync.RWMutex
}

func NewMQTTManager(a ...interface{}) (mqttmngr *MQTTManager) {
	var mqttmsng MqttMessaging = nil
	if al := len(a); al > 0 {
		for al > 0 {
			d := a[0]
			if mqttmsng == nil {
				if mqttmsng, _ = d.(MqttMessaging); mqttmsng != nil {
					al--
					continue
				}
			}
			al--
		}
	}

	mqttmngr = &MQTTManager{lck: &sync.RWMutex{}, cntns: map[string]*MQTTConnection{},
		activeTopics: map[string]*activeTopic{}, lcktpcs: &sync.RWMutex{}, MqttMessaging: mqttmsng}
	return
}

func (mqttmngr *MQTTManager) ActiveTopics() (atvtpcs map[string]string) {
	if mqttmngr != nil {
		func() {
			mqttmngr.lcktpcs.RLock()
			defer mqttmngr.lcktpcs.RUnlock()
			for tpck, tpc := range mqttmngr.activeTopics {
				if tpc != nil {
					atvtpcs[tpck] = tpc.topicpath
				}
			}
		}()
	}
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

func (mqttmngr *MQTTManager) Connection(alias string) (mqttcn *MQTTConnection) {
	if mqttmngr != nil {
		if mqttmngr.ConnectionExist(alias) {
			if mqttmngr != nil && alias != "" {
				func() {
					mqttmngr.lck.RLock()
					defer mqttmngr.lck.RUnlock()
					mqttcn = mqttmngr.cntns[alias]
				}()
			}
		}
	}
	return
}

func (mqttmngr *MQTTManager) ConnectionInfo(alias string) (mqttcninfo string) {
	if mqttmngr != nil {
		if mqttmngr.ConnectionExist(alias) {
			if mqttmngr != nil && alias != "" {
				func() {
					mqttmngr.lck.RLock()
					defer mqttmngr.lck.RUnlock()
					mqttcninfo = mqttmngr.cntns[alias].String()
				}()
			}
		}
	}
	return
}

func (mqttmngr *MQTTManager) ConnectionExist(alias string) (exists bool) {
	if mqttmngr != nil && alias != "" {
		func() {
			mqttmngr.lck.RLock()
			defer mqttmngr.lck.RUnlock()
			if len(mqttmngr.cntns) > 0 {
				_, exists = mqttmngr.cntns[alias]
			}
		}()
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

func (mqttmngr *MQTTManager) messageReceived(mqttcn *MQTTConnection, alias string, msg *mqttMessage) {
	if mqttcn.autoack && msg != nil {
		msg.Ack()
	}
	if mqttmngr.MqttMessaging != nil && len(mqttmngr.activeTopics) > 0 {
		var atvtpc *activeTopic = nil
		func() {
			mqttmngr.lcktpcs.RLock()
			defer mqttmngr.lcktpcs.RUnlock()
			atvtpc = mqttmngr.activeTopics[msg.Topic()]
		}()
		//go func() {
		if atvtpc != nil {
			msg.tokenpath = atvtpc.topicpath
			atvtpc.processMessage(mqttmngr.MqttMessaging, msg)
		}
		//}()
	}
}

func (mqttmngr *MQTTManager) Connected(alias string) {
	//chnls.GLOBALCHNL().DefaultServePath("")
	fmt.Println("Connected:" + alias)
}

func (mqttmngr *MQTTManager) Disconnected(alias string, err error) {
	//chnls.GLOBALCHNL().DefaultServePath("")
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

func (mqttmngr *MQTTManager) IsSubscribed(alias string, topic string) (issbscrbed bool) {
	if alias != "" && topic != "" {
		if exsist, mqttnc := func() (exists bool, mqttcn *MQTTConnection) {
			mqttmngr.lck.RLock()
			defer mqttmngr.lck.RUnlock()
			mqttcn, exists = mqttmngr.cntns[alias]
			return
		}(); exsist {
			func() {
				issbscrbed = mqttnc.IsSubscribed(topic)
			}()
		}
	}
	return
}

func (mqttmngr *MQTTManager) Subscriptions(alias string) (subscrptns []*mqttsubscription) {
	if alias != "" {
		if exsist, mqttnc := func() (exists bool, mqttcn *MQTTConnection) {
			mqttmngr.lck.RLock()
			defer mqttmngr.lck.RUnlock()
			mqttcn, exists = mqttmngr.cntns[alias]
			return
		}(); exsist {
			func() {
				subscrptns = mqttnc.Subscriptions()
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

func (mqttmngr *MQTTManager) Unsubscribe(alias string, topic string) (err error) {
	if alias != "" {
		if exsist, mqttnc := func() (exists bool, mqttcn *MQTTConnection) {
			mqttmngr.lck.RLock()
			defer mqttmngr.lck.RUnlock()
			mqttcn, exists = mqttmngr.cntns[alias]
			return
		}(); exsist {
			func() {
				err = mqttnc.Unsubscribe(topic)
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

func (mqttmngr *MQTTManager) ActivateTopic(topic string, topicpath ...string) {
	if topic != "" {
		func() {
			var atvtpc *activeTopic = nil
			mqttmngr.lcktpcs.Lock()
			defer mqttmngr.lcktpcs.Unlock()
			if atvtpc = mqttmngr.activeTopics[topic]; atvtpc == nil {
				var topicpth = topic

				if len(topicpath) == 1 && topicpath[0] != "" {
					topicpth = topicpath[0]
				}
				atvtpc = &activeTopic{topic: topic, topicpath: topicpth}
				mqttmngr.activeTopics[topic] = atvtpc
			}
		}()
	}
}

func (mqttmngr *MQTTManager) DeactivateTopic(topic string) {
	if topic != "" {
		func() {
			mqttmngr.lcktpcs.Lock()
			defer mqttmngr.lcktpcs.Unlock()
			if atvtpc := mqttmngr.activeTopics[topic]; atvtpc != nil {
				mqttmngr.activeTopics[topic] = nil
				delete(mqttmngr.activeTopics, topic)
				atvtpc = nil
			}
		}()
	}
}

var gblmqttmngr *MQTTManager

func GLOBALMQTTMANAGER() *MQTTManager {
	return gblmqttmngr
}

func init() {
	gblmqttmngr = NewMQTTManager()
}
