package common

import (
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"net"
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

func GetDaprClient() client.Client {
	once.Do(func() {
		port := getDaprPort()
		var err error
		for {
			maxRequestBodySize := 16384
			var opts []grpc.CallOption

			headerBuffer := 1
			opts = append(opts, grpc.MaxCallRecvMsgSize((maxRequestBodySize+headerBuffer)*1024*1024))
			conn, err := grpc.Dial(net.JoinHostPort("127.0.0.1",
				port),
				grpc.WithDefaultCallOptions(opts...), grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				fmt.Printf("Failed to create Dapr client: %v, retrying in %v...\n", err, retryInterval)
				time.Sleep(retryInterval)
				continue
			} else {
				daprClient = client.NewClientWithConnection(conn)
				break
			}

		}

	})
	return daprClient

}
