package common

//message 指发送到UI的各类消息
var WEB_MESSAGE_TOPIC_NAME = "web"

var COMMON_MESSAGE_TYPE_PONG = "pong"

var COMMON_MESSAGE_KEY_TYPE = "type"
var COMMON_MESSAGE_KEY_MESSAGE = "message"
var COMMON_MESSAGE_KEY_TO_ID = "to_id"
var COMMON_MESSAGE_KEY_MARK = "mark"
var COMMON_MESSAGE_KEY_BUSINESS_TYPE = "business_type"
var COMMON_MESSAGE_KEY_CONNECT_ID = "connect_id"

type CommonMessage map[string]interface{}
