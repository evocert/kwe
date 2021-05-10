package mqtt

import (
	"fmt"
	"strconv"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type MQTTConnection struct {
	mqttmngr    *MQTTManager
	pahomqtt    mqtt.Client
	pahooptions *mqtt.ClientOptions
	ClientId    string
	broker      string
	port        int
	user        string
	password    string
}

func newMQTTOptions(clientid string, broker string, port int, user string, password string) (pahooptions *mqtt.ClientOptions) {
	pahooptions = mqtt.NewClientOptions()
	pahooptions.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	pahooptions.SetClientID(clientid)
	pahooptions.SetUsername(user)
	pahooptions.SetPassword(password)
	return
}

func NewMQTTConnections(clientid string, a ...interface{}) (mqttcn *MQTTConnection) {
	if clientid != "" {
		var broker string = ""
		var port int = 0
		var user string = ""
		var password string = ""
		for {
			if al := len(a); al > 0 {
				k := a[0]
				a = a[1:]
				if mp, mpok := k.(map[string]interface{}); mpok {
					for mk, mv := range mp {
						mk = strings.ToLower(mk)
						if s, sok := mv.(string); sok && s != "" {
							if mk == "broker" && broker == "" {
								broker = s
							} else if (mk == "user" || mk == "username") && user == "" {
								user = s
							} else if mk == "password" && password == "" {
								password = s
							}
						} else if i, iok := mv.(int64); iok && i > 0 {
							if mk == "port" && port == 0 {
								port = int(i)
							}
						}
					}
				} else if mp, mpok := k.(map[string]string); mpok {
					for mk, mv := range mp {
						mk = strings.ToLower(mk)
						if mk == "broker" && mv != "" && broker == "" {
							broker = mv
						} else if (mk == "user" || mk == "username") && mv != "" && user == "" {
							user = mv
						} else if mk == "password" && mv != "" && password == "" {
							password = mv
						} else if mk == "port" && mv != "" && port == 0 {
							if prsint, prsinterr := strconv.ParseInt(mv, 0, 64); prsinterr == nil {
								port = int(prsint)
							}
						}
					}
				} else if al > 1 {
					if mk, mkok := k.(string); mkok && mk != "" {
						mk = strings.ToLower(mk)
						if mv := a[0]; mv != nil {
							a = a[1:]
							if s, sok := mv.(string); sok && s != "" {
								if mk == "broker" && broker == "" {
									broker = s
								} else if mk == "password" && password == "" {
									password = s
								} else if (mk == "user" || mk == "username") && user == "" {
									user = s
								}
							} else if i, iok := mv.(int64); iok && i > 0 {
								if mk == "port" && port == 0 {
									port = int(i)
								}
							}
						} else {
							break
						}
					} else {
						break
					}
				}
			} else {
				break
			}
		}
		if pahooptions := newMQTTOptions(clientid, broker, port, user, password); pahooptions != nil {
			var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
				if mqttcn != nil && mqttcn.mqttmngr != nil {
					mqttcn.mqttmngr.Connected(clientid)
				}
			}

			var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
				if mqttcn != nil && mqttcn.mqttmngr != nil {
					mqttcn.mqttmngr.Disconnected(clientid, err)
				}
				//fmt.Printf("Connect lost:"+clientid+"%v", err)
			}

			var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
				if mqttcn != nil && mqttcn.mqttmngr != nil {
					var mqttmsg Message = msg
					mqttcn.mqttmngr.MessageReceived(clientid, mqttmsg)
					mqttmsg = nil
				}
			}
			pahooptions.SetDefaultPublishHandler(messagePubHandler)
			pahooptions.OnConnect = connectHandler
			pahooptions.OnConnectionLost = connectLostHandler
			pahomqtt := mqtt.NewClient(pahooptions)
			mqttcn = &MQTTConnection{mqttmngr: nil, pahomqtt: pahomqtt, broker: broker, port: port, user: user, password: password, ClientId: clientid}

		}
	}
	return
}

type Message interface {
	Duplicate() bool
	Qos() byte
	Retained() bool
	Topic() string
	MessageID() uint16
	Payload() []byte
	Ack()
}

func (mqttcn *MQTTConnection) IsConnected() (connected bool) {
	if mqttcn != nil && mqttcn.pahomqtt != nil {
		connected = mqttcn.pahomqtt.IsConnected()
	}
	return
}

func (mqttcn *MQTTConnection) Disconnect(quiesce uint) (err error) {
	if mqttcn != nil {
		if client := mqttcn.pahomqtt; client != nil {
			client.Disconnect(quiesce)
		}
	}
	return
}

func (mqttcn *MQTTConnection) Connect() (err error) {
	if mqttcn != nil {
		if client := mqttcn.pahomqtt; client != nil {
			if token := client.Connect(); token.Wait() && token.Error() != nil {
				err = token.Error()
			}
		}
	}
	return
}

func (mqttcn *MQTTConnection) Publish(topic string, qos byte, retained bool, message string) (err error) {
	if mqttcn != nil && mqttcn.pahomqtt != nil {
		tkn := mqttcn.pahomqtt.Publish(topic, qos, retained, message)
		tkn.Wait()
		err = tkn.Error()
	}
	return err
}

func (mqttcn *MQTTConnection) Subscribe(topic string, qos byte) (err error) {
	if mqttcn != nil && mqttcn.pahomqtt != nil {
		var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
			if mqttcn != nil && mqttcn.mqttmngr != nil {
				var mqttmsg Message = msg
				mqttcn.mqttmngr.MessageReceived(mqttcn.ClientId, mqttmsg)
				mqttmsg = nil
			}
		}
		tkn := mqttcn.pahomqtt.Subscribe(topic, qos, messagePubHandler)
		tkn.Wait()
		err = tkn.Error()
	}
	return err
}

func init() {

}
