// internal message 内部服务间各类消息
package common

import (
	"encoding/json"
)

var INTERNAL_MESSAGE_TOPIC_NAME = "internal_msg"

var INTERNAL_MESSAGE_TYPE_USER_LOGIN = "user_login"
var INTERNAL_MESSAGE_TYPE_USER_LOGOUT = "user_logout"
var INTERNAL_MESSAGE_TYPE_WEB_CONNECT = "web_connect"
var INTERNAL_MESSAGE_TYPE_WEB_DISCONNECT = "web_disconnect"
var INTERNAL_MESSAGE_TYPE_SYS_LOG = "sys_log"

type InternalMessage map[string]any

var INTERNAL_MESSAGE_KEY_TYPE = "type"
var INTERNAL_MESSAGE_KEY_USER_ID = "user_id"
var INTERNAL_MESSAGE_KEY_MARK = "mark"
var INTERNAL_MESSAGE_KEY_BUSINESS_TYPE = "business_type"
var INTERNAL_MESSAGE_KEY_CONNECT_ID = "connect_id"
var INTERNAL_MESSAGE_KEY_MESSAGE = "message"

func (m InternalMessage) GetType() string {
	val := m[INTERNAL_MESSAGE_KEY_TYPE]
	return val.(string)
}
func (m InternalMessage) SetType(t string) {
	m[INTERNAL_MESSAGE_KEY_TYPE] = t
}

func (m InternalMessage) FromStruct(s any) (err error) {
	buf, err := json.Marshal(s)
	if err != nil {
		return err
	}
	return json.Unmarshal(buf, &m)
}

type SyslogMessage struct {
	Type   string `json:"type"`
	UserId string `json:"user_id"`
	Ip     string `json:"ip"`
	Action int    `json:"action"`
	Info   string `json:"info"`
}

func (s *SyslogMessage) FromInternalMessage(m InternalMessage) {
	s.Type = m["type"].(string)
	s.UserId = m["user_id"].(string)
	s.Ip = m["ip"].(string)
	s.Action = int(m["action"].(float64))
	s.Info = m["info"].(string)
}
