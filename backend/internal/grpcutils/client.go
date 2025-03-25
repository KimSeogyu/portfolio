package grpcutils

import (
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

// NewGrpcClient creates a new gRPC client.
//
// Parameters:
//   - target: the target to connect to
//   - opts: the options to pass to the client
//
// Returns:
//   - a new ClientConn instance
//   - an error if the client fails to initialize
//
// Details:
//   - The client will use the insecure credentials
//   - The client will use the keepalive parameters
//   - The client will use the options passed in the opts parameter
func NewGrpcClient(target string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	var dialOptions []grpc.DialOption

	dialOptions = append(dialOptions,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:    10 * time.Second, // send pings every 10 seconds if there is no activity
			Timeout: 5 * time.Second,  // wait 5 second for ping ack before considering the connection dead
		}),
	)

	dialOptions = append(dialOptions, opts...)

	conn, err := grpc.NewClient(target, dialOptions...)
	if err != nil {
		return nil, err
	}

	return conn, nil
}
