package sdnotify

import (
	"context"
	"errors"
	"fmt"
	"io"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// healthSummary is what the watcher goroutine consumes for each status change.
type healthSummary struct {
	Status grpc_health_v1.HealthCheckResponse_ServingStatus
	// Line is a short, human-readable form suitable for `STATUS=<Line>`
	// systemd notifications. Includes the watched service name when non-empty.
	Line string
}

// healthClient wraps grpc_health_v1.HealthClient with a Watch helper that
// owns its own reconnect/backoff loop and delivers results on a channel.
type healthClient struct {
	conn    *grpc.ClientConn
	client  grpc_health_v1.HealthClient
	service string // "" for overall, or pipeline name
	logger  *zap.Logger
}

// dialHealthClient opens an insecure gRPC connection to endpoint. The
// healthcheckv2 gRPC service is loopback by default; TLS is out of scope
// for this iteration.
func dialHealthClient(endpoint, service string, logger *zap.Logger) (*healthClient, error) {
	if endpoint == "" {
		return nil, errors.New("healthcheckv2 gRPC endpoint is empty")
	}
	conn, err := grpc.NewClient(endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("dial healthcheckv2 gRPC %q: %w", endpoint, err)
	}
	return &healthClient{
		conn:    conn,
		client:  grpc_health_v1.NewHealthClient(conn),
		service: service,
		logger:  logger,
	}, nil
}

// Close releases the underlying gRPC connection. Safe to call multiple times.
func (c *healthClient) Close() error {
	if c == nil || c.conn == nil {
		return nil
	}
	err := c.conn.Close()
	c.conn = nil
	return err
}

// Check is a one-shot unary health query.
func (c *healthClient) Check(ctx context.Context) (*healthSummary, error) {
	resp, err := c.client.Check(ctx, &grpc_health_v1.HealthCheckRequest{Service: c.service})
	if err != nil {
		return nil, err
	}
	return &healthSummary{Status: resp.GetStatus(), Line: formatLine(resp.GetStatus(), c.service)}, nil
}

// Watch opens a server-stream and emits a healthSummary on every status
// change until ctx is cancelled. Transient stream errors trigger a
// bounded exponential backoff retry inside the goroutine -- callers only
// see clean summaries on the channel, never errors.
//
// The returned channel is closed when ctx is done.
func (c *healthClient) Watch(ctx context.Context) <-chan healthSummary {
	out := make(chan healthSummary, 1)
	go func() {
		defer close(out)

		const (
			minBackoff = 200 * time.Millisecond
			maxBackoff = 10 * time.Second
		)
		backoff := minBackoff

		for {
			if ctx.Err() != nil {
				return
			}

			stream, err := c.client.Watch(ctx, &grpc_health_v1.HealthCheckRequest{Service: c.service})
			if err != nil {
				if ctx.Err() != nil {
					return
				}
				c.logger.Debug("healthcheckv2 Watch open failed; will retry",
					zap.Duration("backoff", backoff), zap.Error(err))
				if !sleep(ctx, backoff) {
					return
				}
				backoff = nextBackoff(backoff, maxBackoff)
				continue
			}

			// Reset backoff once we have a working stream.
			backoff = minBackoff

			for {
				resp, err := stream.Recv()
				if err != nil {
					if ctx.Err() != nil {
						return
					}
					if !errors.Is(err, io.EOF) {
						c.logger.Debug("healthcheckv2 Watch recv failed; reconnecting",
							zap.Error(err))
					}
					break // break inner loop, retry stream
				}
				select {
				case out <- healthSummary{Status: resp.GetStatus(), Line: formatLine(resp.GetStatus(), c.service)}:
				case <-ctx.Done():
					return
				}
			}
		}
	}()
	return out
}

func formatLine(s grpc_health_v1.HealthCheckResponse_ServingStatus, service string) string {
	if service == "" {
		return s.String()
	}
	return fmt.Sprintf("%s (%s)", s.String(), service)
}

func nextBackoff(cur, max time.Duration) time.Duration {
	d := cur * 2
	if d > max {
		return max
	}
	return d
}

// sleep is a context-aware time.Sleep; returns false if ctx fired first.
func sleep(ctx context.Context, d time.Duration) bool {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-t.C:
		return true
	case <-ctx.Done():
		return false
	}
}
