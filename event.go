package common

import (
	"context"
	dapr "github.com/dapr/go-sdk/client"
	"github.com/pkg/errors"
)

// event指告警类事件
var (
	PLATFORM_ALARM_TOPIC_NAME = "platform-alarm"
	EVENT_TOPIC_NAME          = "eventTopic"
	EVENT_DATA_TOPIC_NAME     = "eventDataTopic"
	EVENT_POINT_RW_META_TOPIC = "pointRWMetaTopic"
	EventArchivedFlag         = int32(1)

	EventTypePlatform = int32(1)
	EventTypeDevice   = int32(2)
	EventTypeSecurity = int32(3)

	EventSubTypeService      = int32(101)
	EventSubTypeInterface    = int32(102)
	EventSubTypeComunication = int32(103)

	EventStatusActive  = int32(1)
	EventStatusClosed  = int32(0)
	EventLevelCritical = int32(1)
	EventLevelMajor    = int32(2)
	EventLevelMinor    = int32(3)
	EventLevelWarning  = int32(4)
)

type Event struct {
	Dn          string    `json:"dn"`
	Title       string    `json:"title"`
	Type        int32     `json:"type"`
	Description string    `json:"description"`
	Status      int32     `json:"status"`
	Level       int32     `json:"level"`
	EventTime   LocalTime `json:"event_time"`
	CreateAt    LocalTime `json:"create_at"`
	UpdatedAt   LocalTime `json:"updated_at"`
	ObjectID    string    `json:"object_id"`
	ObjectName  string    `json:"object_name"`
	Location    string    `json:"location"`
}

type MethodInvokeInfo struct {
	Service string `json:"service"`
	Method  string `json:"method"` //"GET POST PUT DELETE"
	Url     string `json:"url"`
	Data    any    `json:"data"`
}

func PublishMethodInvokeMessage(ctx context.Context, client dapr.Client, event MethodInvokeInfo) error {

	err := client.PublishEvent(ctx, PUBSUB_NAME, METHOD_INVOKE_TOPIC, event)
	if err != nil {
		err = errors.Wrap(err, "PublishDbUpsertMessage")
	}
	return err
}

type DbUpsertEvent struct {
	Db         string `json:"db"`
	Schema     string `json:"schema"`
	Table      string `json:"table"`
	Keys       string `json:"keys"`
	Batch      bool   `json:"batch"`
	Ignorekeys string `json:"ignorekeys"`
	Data       any    `json:"data"`
}

func PublishDbUpsertMessage(ctx context.Context, client dapr.Client, table, keys, ignorekeys string, batch bool, data any) error {
	event := DbUpsertEvent{
		Db:         DBNAME,
		Schema:     DB_SCHEMA,
		Table:      table,
		Keys:       keys,
		Batch:      batch,
		Ignorekeys: ignorekeys,
		Data:       data,
	}

	err := client.PublishEvent(ctx, PUBSUB_NAME, "db_upsert_event", event)
	if err != nil {
		err = errors.Wrap(err, "PublishDbUpsertMessage")
	}
	return err
}
