package sdnotify

import (
	"context"
	"net"
	"sync"
	"testing"
	"time"

	"go.uber.org/zap/zaptest"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// scriptedHealthServer implements grpc_health_v1.HealthServer.Watch by
// pushing a scripted sequence of ServingStatus values, then holding the
// stream open until the client cancels.
type scriptedHealthServer struct {
	grpc_health_v1.UnimplementedHealthServer
	script []grpc_health_v1.HealthCheckResponse_ServingStatus
}

func (s *scriptedHealthServer) Check(_ context.Context, _ *grpc_health_v1.HealthCheckRequest) (*grpc_health_v1.HealthCheckResponse, error) {
	return &grpc_health_v1.HealthCheckResponse{Status: grpc_health_v1.HealthCheckResponse_SERVING}, nil
}

func (s *scriptedHealthServer) Watch(_ *grpc_health_v1.HealthCheckRequest, stream grpc_health_v1.Health_WatchServer) error {
	for _, st := range s.script {
		if err := stream.Send(&grpc_health_v1.HealthCheckResponse{Status: st}); err != nil {
			return err
		}
		// Small pause so the consumer can dedup correctly even if statuses repeat.
		time.Sleep(10 * time.Millisecond)
	}
	<-stream.Context().Done()
	return nil
}

// startScriptedHealthServer spins up an in-process gRPC server on a random
// loopback port and returns its endpoint. Cleaned up on test end.
func startScriptedHealthServer(t *testing.T, script []grpc_health_v1.HealthCheckResponse_ServingStatus) string {
	t.Helper()
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	srv := grpc.NewServer()
	grpc_health_v1.RegisterHealthServer(srv, &scriptedHealthServer{script: script})

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		_ = srv.Serve(lis)
	}()
	t.Cleanup(func() {
		srv.GracefulStop()
		wg.Wait()
	})
	return lis.Addr().String()
}

func TestDeepHealthcheck_StreamsStatusUpdates(t *testing.T) {
	rec := startNotifyRecorder(t)

	// Script: SERVING -> SERVING (dedup'd) -> NOT_SERVING -> SERVING.
	endpoint := startScriptedHealthServer(t, []grpc_health_v1.HealthCheckResponse_ServingStatus{
		grpc_health_v1.HealthCheckResponse_SERVING,
		grpc_health_v1.HealthCheckResponse_SERVING,
		grpc_health_v1.HealthCheckResponse_NOT_SERVING,
		grpc_health_v1.HealthCheckResponse_SERVING,
	})

	hc, err := dialHealthClient(endpoint, "", zaptest.NewLogger(t))
	if err != nil {
		t.Fatalf("dialHealthClient: %v", err)
	}
	t.Cleanup(func() { _ = hc.Close() })

	ext := newSDNotify(&Config{
		DeepHealthcheck:           true,
		HealthcheckV2GRPCEndpoint: endpoint, // bypass sibling lookup
	}, zaptest.NewLogger(t))
	ext.host = nil // not needed because we set HealthcheckV2GRPCEndpoint
	ext.hc = hc

	// Manually drive the watcher goroutine so we don't need to wire up a Host.
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	ext.watchCancel = cancel
	ext.watchDone = done
	updates := hc.Watch(ctx)
	go func() {
		defer close(done)
		for u := range updates {
			ext.sendStatusLine(u.Line)
		}
	}()

	// Expect three distinct STATUS= lines: SERVING, NOT_SERVING, SERVING.
	if !waitFor(t, 3*time.Second, func() bool { return rec.countPrefix("STATUS=") >= 3 }) {
		t.Fatalf("expected >=3 STATUS= lines, got %v", rec.all())
	}

	got := rec.all()
	var statuses []string
	for _, m := range got {
		if len(m) > len("STATUS=") && m[:len("STATUS=")] == "STATUS=" {
			statuses = append(statuses, m[len("STATUS="):])
		}
	}
	want := []string{"SERVING", "NOT_SERVING", "SERVING"}
	if len(statuses) != len(want) {
		t.Fatalf("want %v STATUS= lines, got %v", want, statuses)
	}
	for i, w := range want {
		if statuses[i] != w {
			t.Fatalf("status[%d]=%q, want %q (all=%v)", i, statuses[i], w, statuses)
		}
	}

	// Clean shutdown.
	cancel()
	<-done
}

func TestHealthClient_CheckOneShot(t *testing.T) {
	endpoint := startScriptedHealthServer(t, nil)
	hc, err := dialHealthClient(endpoint, "", zaptest.NewLogger(t))
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	t.Cleanup(func() { _ = hc.Close() })

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	sum, err := hc.Check(ctx)
	if err != nil {
		t.Fatalf("Check: %v", err)
	}
	if sum.Status != grpc_health_v1.HealthCheckResponse_SERVING {
		t.Fatalf("want SERVING, got %v", sum.Status)
	}
	if sum.Line != "SERVING" {
		t.Fatalf("want line %q, got %q", "SERVING", sum.Line)
	}
}
