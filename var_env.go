package common

import (
	"os"
)

var DBNAME = "thingsdb"
var DB_SERVICE_NAME = "db-service"
var CENTER_DB_SERVICE_NAME = "center-db-service"
var DB_SCHEMA = "public"
var PUBSUB_NAME = "pubsub"
var CENTER_PUBSUB_NAME = "center_pubsub"
var DAPR_STATESTORE_NAME = "statestore"
var GLOBAL_STATESTOR_NAME = "global-redis"
var CENTER_DAPR_STATESTORE_NAME = "center-statestore"
var CENTER_GLOBAL_STATESTOR_NAME = "center-global-redis"
var BASE_CONTEXT = ""
var USER_STATESTORE_KEY_PREFIX = "user#"
var FILE_SERVICE_NAME = "file-service"
var USER_EXPIRED_SECONDS = 7200
var RUNNING_MODE = "center_edge"
var RUNNING_MODE_EDGE = "edge"
var RUNNING_MODE_CENTER = "center"
var RUNNING_MODE_CENTER_EDGE = "center_edge"
var EDGE_ID = ""
var DEVICE_DATA_TOPIC = "deviceData"
var RESOURCE_CHANGE_TOPIC = "resourceChange"
var METHOD_INVOKE_TOPIC = "methodInvoke"
var CENTER_METHOD_INVOKE_TOPIC = "centerMethodInvoke"
var DB_UPSERT_TOPIC = "db_upsert_event"
var PROPERTY_SET_TOPIC = "property_set_event"
var CENTER_DB_UPSERT_TOPIC = "center_db_upsert_event"

func init() {
	if val := os.Getenv("DBNAME"); val != "" {
		DBNAME = val
		Logger.Info("DBNAME:", DBNAME)
	}
	if val := os.Getenv("BASE_CONTEXT"); val != "" {
		BASE_CONTEXT = val
		Logger.Info("BASE_CONTEXT:", BASE_CONTEXT)
	}
	if val := os.Getenv("PUBSUB_NAME"); val != "" {
		PUBSUB_NAME = val
		Logger.Info("PUBSUB_NAME:", PUBSUB_NAME)
	}
	if val := os.Getenv("DB_SERVICE_NAME"); val != "" {
		DB_SERVICE_NAME = val
		Logger.Info("DB_SERVICE_NAME:", DB_SERVICE_NAME)
	}
	if val := os.Getenv("DB_SCHEMA"); val != "" {
		DB_SCHEMA = val
		Logger.Info("DB_SCHEMA:", DB_SCHEMA)
	}
	if val := os.Getenv("DAPR_STATESTORE_NAME"); val != "" {
		DAPR_STATESTORE_NAME = val
		Logger.Info("DAPR_STATESTORE_NAME:", DAPR_STATESTORE_NAME)
	}
	if val := os.Getenv("GLOBAL_STATESTOR_NAME"); val != "" {
		GLOBAL_STATESTOR_NAME = val
		Logger.Info("GLOBAL_STATESTOR_NAME:", GLOBAL_STATESTOR_NAME)
	}
	if val := os.Getenv("RUNNING_MODE"); val != "" {
		RUNNING_MODE = val
		Logger.Info("RUNNING_MODE:", RUNNING_MODE)
	}
	if val := os.Getenv("EDGE_ID"); val != "" {
		EDGE_ID = val
		Logger.Info("EDGE_ID:", EDGE_ID)
	}

}
