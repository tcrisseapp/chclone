package grpcservice

import (
	"google.golang.org/grpc"
)

var (
	retryPolicy = `{
            "methodConfig": [{
                "waitForReady": true,
                "retryPolicy": {
                    "MaxAttempts": 4,
                    "InitialBackoff": ".01s",
                    "MaxBackoff": ".01s",
                    "BackoffMultiplier": 1.0,
                    "RetryableStatusCodes": [ "UNAVAILABLE" ]
                }
            }]
        }`
)

// Dial will dial a gRPC service with defautl dial options
func Dial(target string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(target, grpc.WithInsecure(), grpc.WithDefaultServiceConfig(retryPolicy))
	if err != nil {
		return nil, err
	}
	return conn, nil
}
