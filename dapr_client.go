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
	daprPortEnvVarName = "DAPR_GRPC_PORT"
	daprHostEnvVarName = "DAPR_GRPC_HOST"
	retryInterval      = 5 * time.Second // 定义重试间隔
)

var (
	daprPort = "50001"
	daprHost = "127.0.0.1"
)

func init() {
	if p := os.Getenv(daprPortEnvVarName); p != "" {
		daprPort = p
	}
	if h := os.Getenv(daprHostEnvVarName); h != "" {
		daprHost = h
	}
}

var once = sync.Once{}

func GetDaprClient() client.Client {
	once.Do(func() {

		for {
			maxRequestBodySize := 16384
			var opts []grpc.CallOption

			headerBuffer := 1
			opts = append(opts, grpc.MaxCallRecvMsgSize((maxRequestBodySize+headerBuffer)*1024*1024))
			conn, err := grpc.NewClient(net.JoinHostPort(daprHost,
				daprPort),
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
