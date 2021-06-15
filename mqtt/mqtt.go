package mqtt

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"strconv"
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/evocert/kwe/iorw"
)

type MQTTConnection struct {
	mqttmngr *MQTTManager
	pahomqtt mqtt.Client
	ClientId string
	broker   string
	port     int
	user     string
	password string
	autoack  bool
}

func newMQTTOptions(clientid string, broker string, port int, user string, password string) (pahooptions *mqtt.ClientOptions) {
	pahooptions = mqtt.NewClientOptions()
	var schema = "tcp"
	if broker != "" && strings.HasPrefix(broker, "ws://") {
		schema = "ws"
		broker = broker[len("ws://"):]
	}
	pahooptions.AddBroker(fmt.Sprintf("%s://%s:%d", schema, broker, port))
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
		var autoack bool = false
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
						} else {
							if mk == "port" && port == 0 {
								if prsint, prsinterr := strconv.ParseInt(fmt.Sprint(mv), 0, 64); prsinterr == nil {
									port = int(prsint)
								} else if prsint, prsinterr := strconv.ParseInt(fmt.Sprint(mv), 0, 32); prsinterr == nil {
									port = int(prsint)
								} else if prsint, prsinterr := strconv.ParseInt(fmt.Sprint(mv), 0, 16); prsinterr == nil {
									port = int(prsint)
								}
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
							if prsint, prsinterr := strconv.ParseInt(fmt.Sprint(mv), 0, 64); prsinterr == nil {
								port = int(prsint)
							} else if prsint, prsinterr := strconv.ParseInt(fmt.Sprint(mv), 0, 32); prsinterr == nil {
								port = int(prsint)
							} else if prsint, prsinterr := strconv.ParseInt(fmt.Sprint(mv), 0, 16); prsinterr == nil {
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
							} else if b, bok := mv.(bool); bok && b {
								if mk == "autoack" && b {
									autoack = b
								}
							} else {
								if mk == "port" && port == 0 {
									if prsint, prsinterr := strconv.ParseInt(fmt.Sprint(mv), 0, 64); prsinterr == nil {
										port = int(prsint)
									} else if prsint, prsinterr := strconv.ParseInt(fmt.Sprint(mv), 0, 32); prsinterr == nil {
										port = int(prsint)
									} else if prsint, prsinterr := strconv.ParseInt(fmt.Sprint(mv), 0, 16); prsinterr == nil {
										port = int(prsint)
									}
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
			}

			var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
				if mqttcn != nil && mqttcn.mqttmngr != nil {
					func() {
						var mqttmsg *mqttMessage = &mqttMessage{msg: msg, mqttcn: mqttcn, mqttmmng: mqttcn.mqttmngr}
						defer func() {
							mqttmsg.mqttcn = nil
							mqttmsg.msg = nil
							mqttmsg.tokenpath = ""
							mqttmsg = nil
							mqttmsg.mqttmmng = nil
						}()
						mqttcn.mqttmngr.messageReceived(mqttcn, clientid, mqttmsg)
					}()
				}
			}
			pahooptions.SetDefaultPublishHandler(messagePubHandler)
			pahooptions.OnConnect = connectHandler
			pahooptions.OnConnectionLost = connectLostHandler
			pahomqtt := mqtt.NewClient(pahooptions)
			mqttcn = &MQTTConnection{mqttmngr: nil, pahomqtt: pahomqtt, broker: broker, port: port, user: user, password: password, ClientId: clientid, autoack: autoack}

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
	Connection() *MQTTConnection
	TopicPath() string
	Manager() *MQTTManager
	Ack()
}

type mqttMessage struct {
	mqttcn    *MQTTConnection
	mqttmmng  *MQTTManager
	msg       mqtt.Message
	tokenpath string
}

func (mqttmsg *mqttMessage) FPrint(w io.Writer) {
	if mqttmsg != nil && w != nil {
		enc := json.NewEncoder(w)
		enc.SetIndent("", "")
		iorw.Fprint(w, "{")
		enc.Encode("msgid")
		iorw.Fprint(w, ":")
		enc.Encode(mqttmsg.msg.MessageID())
		iorw.Fprint(w, ",")
		enc.Encode("clientid")
		iorw.Fprint(w, ":")
		enc.Encode(mqttmsg.mqttcn.ClientId)
		iorw.Fprint(w, ",")
		enc.Encode("duplicate")
		iorw.Fprint(w, ":")
		enc.Encode(mqttmsg.msg.Duplicate())
		iorw.Fprint(w, ",")
		payload := mqttmsg.msg.Payload()
		//if len(payload) > 0 {
		enc.Encode("payload")
		iorw.Fprint(w, ":")
		enc.Encode(string(payload))
		iorw.Fprint(w, ",")
		enc.Encode("bin-payload")
		iorw.Fprint(w, ":")
		arrpayload := make([]interface{}, len(payload))
		for pn, p := range payload {
			arrpayload[pn] = p
		}
		enc.Encode(arrpayload)
		arrpayload = nil
		//} else {
		//	enc.Encode("payload")
		//	iorw.Fprint(w, ":\"\"")
		//	iorw.Fprint(w, ",")
		//	enc.Encode("bin-payload")
		//	iorw.Fprint(w, ":[]")
		//}
		iorw.Fprint(w, ",")
		enc.Encode("topic")
		iorw.Fprint(w, ":")
		enc.Encode(mqttmsg.msg.Topic())
		iorw.Fprint(w, ",")
		enc.Encode("qos")
		iorw.Fprint(w, ":")
		enc.Encode(mqttmsg.msg.Qos())
		iorw.Fprint(w, ",")
		enc.Encode("retained")
		iorw.Fprint(w, ":")
		enc.Encode(mqttmsg.msg.Retained())
		iorw.Fprint(w, ",")
		enc.Encode("topicpath")
		iorw.Fprint(w, ":")
		enc.Encode(mqttmsg.tokenpath)
		iorw.Fprint(w, "}")
	}
}

func (mqttmsg *mqttMessage) String() (s string) {
	pr, pw := io.Pipe()
	defer pr.Close()
	ctx, ctxcancel := context.WithCancel(context.Background())
	go func() {
		defer pw.Close()
		ctxcancel()
		mqttmsg.FPrint(pw)
	}()
	<-ctx.Done()
	s, _ = iorw.ReaderToString(pr)
	s = strings.Replace(s, "\n", "", -1)
	return
}

func (mqttmsg *mqttMessage) TopicPath() (topicpath string) {
	if mqttmsg != nil {
		topicpath = mqttmsg.tokenpath
	}
	return
}

func (mqttmsg *mqttMessage) Duplicate() bool {
	return mqttmsg.msg.Duplicate()
}

func (mqttmsg *mqttMessage) Qos() byte {
	return mqttmsg.msg.Qos()
}

func (mqttmsg *mqttMessage) Retained() bool {
	return mqttmsg.msg.Retained()
}

func (mqttmsg *mqttMessage) Topic() string {
	return mqttmsg.msg.Topic()
}

func (mqttmsg *mqttMessage) MessageID() uint16 {
	return mqttmsg.msg.MessageID()
}

func (mqttmsg *mqttMessage) Payload() []byte {
	return mqttmsg.msg.Payload()
}

func (mqttmsg *mqttMessage) Connection() *MQTTConnection {
	return mqttmsg.mqttcn
}

func (mqttmsg *mqttMessage) Manager() *MQTTManager {
	return mqttmsg.mqttmmng
}

func (mqttmsg *mqttMessage) Ack() {
	mqttmsg.msg.Ack()
}

func (mqttcn *MQTTConnection) Fprint(w io.Writer) {
	if mqttcn != nil && w != nil {
		enc := json.NewEncoder(w)
		iorw.Fprint(w, "{")
		iorw.Fprint(w, "\"ClientID\":")
		enc.Encode(mqttcn.ClientId)
		iorw.Fprint(w, ",")

		iorw.Fprint(w, "\"broker\":")
		enc.Encode(mqttcn.broker)
		iorw.Fprint(w, ",")

		iorw.Fprint(w, "\"port\":")
		enc.Encode(mqttcn.port)
		iorw.Fprint(w, ",")

		iorw.Fprint(w, "\"user\":")
		enc.Encode(mqttcn.user)
		iorw.Fprint(w, ",")

		iorw.Fprint(w, "\"password\":")
		enc.Encode(mqttcn.password)
		iorw.Fprint(w, ",")

		iorw.Fprint(w, "\"autoack\":")
		enc.Encode(mqttcn.autoack)
		iorw.Fprint(w, ",")

		iorw.Fprint(w, "\"status\":")
		if mqttcn.IsConnected() {
			iorw.Fprint(w, "\"connected\"")
		} else {
			if mqttcn.IsConnected() {
				iorw.Fprint(w, "\"disconnected\"")
			}
		}
		iorw.Fprint(w, "}")
	}
}

func (mqttcn *MQTTConnection) String() (s string) {
	if mqttcn != nil {
		pr, pw := io.Pipe()
		ctx, ctxcancel := context.WithCancel(context.Background())
		go func() {
			defer pw.Close()
			ctxcancel()
			mqttcn.Fprint(pw)
		}()
		<-ctx.Done()
		s, _ = iorw.ReaderToString(pr)
	}
	return
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
				func() {
					var mqttmsg *mqttMessage = &mqttMessage{msg: msg, mqttcn: mqttcn, mqttmmng: mqttcn.mqttmngr}
					defer func() {
						mqttmsg.mqttcn = nil
						mqttmsg.msg = nil
						mqttmsg.tokenpath = ""
						mqttmsg.mqttmmng = nil
						mqttmsg = nil
					}()
					mqttcn.mqttmngr.messageReceived(mqttcn, mqttcn.ClientId, mqttmsg)
				}()
			}
		}
		tkn := mqttcn.pahomqtt.Subscribe(topic, qos, messagePubHandler)
		tkn.Wait()
		err = tkn.Error()
	}
	return err
}

func (mqttcn *MQTTConnection) Unsubscribe(topic ...string) (err error) {
	if mqttcn != nil && mqttcn.pahomqtt != nil {

		tkn := mqttcn.pahomqtt.Unsubscribe(topic...)
		tkn.Wait()
		err = tkn.Error()
	}
	return err
}

func init() {

}
