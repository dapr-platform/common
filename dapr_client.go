package common

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/dapr/go-sdk/client"
)

var daprClient client.Client

const (
	daprPortDefault    = "50001"
	daprPortEnvVarName = "DAPR_GRPC_PORT"
	retryInterval      = 5 * time.Second // 定义重试间隔
)

var once = sync.Once{}

func getDaprPort() string {
	port := os.Getenv(daprPortEnvVarName)
	if port == "" {
		port = daprPortDefault
	}
	return port
}

func createDaprClient(port string) client.Client {
	var c client.Client
	var err error
	for {
		c, err = client.NewClientWithPort(port)
		if err == nil {
			break
		}
		fmt.Printf("Failed to create Dapr client: %v, retrying in %v...\n", err, retryInterval)
		time.Sleep(retryInterval)
	}
	return c
}

func GetDaprClient() client.Client {
	once.Do(func() {
		port := getDaprPort()
		var err error
		for {
			daprClient, err = client.NewClientWithPort(port)
			if err == nil {
				break
			}
			fmt.Printf("Failed to create Dapr client: %v, retrying in %v...\n", err, retryInterval)
			time.Sleep(retryInterval)
		}

	})
	return daprClient

}
