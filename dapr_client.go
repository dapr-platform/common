package common

import (
	"fmt"
	"os"
	"time"

	"github.com/dapr/go-sdk/client"
)

var daprClient client.Client

const (
	daprPortDefault    = "50001"
	daprPortEnvVarName = "DAPR_GRPC_PORT"
	retryInterval      = 5 * time.Second // 定义重试间隔
)

func getDaprPort() string {
	port := os.Getenv(daprPortEnvVarName)
	if port == "" {
		port = daprPortDefault
	}
	return port
}

func createDaprClient(port string) (client.Client, error) {
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
	return c, err
}

func GetDaprClient() (client.Client, error) {
	if daprClient != nil {
		return daprClient, nil
	}
	port := getDaprPort()
	var err error
	daprClient, err = createDaprClient(port)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Dapr client after retries: %w", err)
	}
	return daprClient, nil
}
