package toybox

import (
	"context"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"syscall"
	"testing"
	"time"

	"github.com/inoth/toybox/conf"
	"github.com/inoth/toybox/registry"
)

// --- Mock Transport ---

type mockTransport struct {
	name     string
	started  atomic.Bool
	stopped  atomic.Bool
	startCh  chan struct{} // closed when Start has been entered
	stopCh   chan struct{} // closed by Stop to unblock Start
	stopErr  error
	startErr error
}

func newMockTransport(name string) *mockTransport {
	return &mockTransport{
		name:    name,
		startCh: make(chan struct{}),
		stopCh:  make(chan struct{}),
	}
}

func (m *mockTransport) Start(ctx context.Context) error {
	m.started.Store(true)
	close(m.startCh)
	if m.startErr != nil {
		return m.startErr
	}
	// Block until Stop is called or context is cancelled.
	select {
	case <-ctx.Done():
		return nil
	case <-m.stopCh:
		return nil
	}
}

func (m *mockTransport) Stop(_ context.Context) error {
	m.stopped.Store(true)
	select {
	case <-m.stopCh:
	default:
		close(m.stopCh)
	}
	return m.stopErr
}

func (m *mockTransport) TransportName() string {
	return m.name
}

// --- Mock Transport with Endpoint ---

type mockEndpointTransport struct {
	mockTransport
	endpoint string
}

func newMockEndpointTransport(name, endpoint string) *mockEndpointTransport {
	return &mockEndpointTransport{
		mockTransport: mockTransport{
			name:    name,
			startCh: make(chan struct{}),
			stopCh:  make(chan struct{}),
		},
		endpoint: endpoint,
	}
}

func (m *mockEndpointTransport) Endpoint() (string, error) {
	return m.endpoint, nil
}

// --- Mock Registrar ---

type mockRegistrar struct {
	mu           sync.Mutex
	registered   []*registry.ServiceInstance
	deregistered []*registry.ServiceInstance
	registerErr  error
}

func (m *mockRegistrar) Register(_ context.Context, svc *registry.ServiceInstance) error {
	if m.registerErr != nil {
		return m.registerErr
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.registered = append(m.registered, svc)
	return nil
}

func (m *mockRegistrar) Deregister(_ context.Context, svc *registry.ServiceInstance) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.deregistered = append(m.deregistered, svc)
	return nil
}

// --- Mock ConfigMate ---

type mockConfigMate struct{}

func (m *mockConfigMate) PrimitiveDecode(_ ...conf.ConfigureMatcher) error {
	return nil
}

// --- Panic Transport ---

type panicTransport struct {
	startCh chan struct{}
}

func newPanicTransport() *panicTransport {
	return &panicTransport{startCh: make(chan struct{})}
}

func (p *panicTransport) Start(_ context.Context) error {
	close(p.startCh)
	panic("transport exploded")
}

func (p *panicTransport) Stop(_ context.Context) error {
	return nil
}

// --- Slow Stop Transport ---

type slowStopTransport struct {
	startCh chan struct{}
	delay   time.Duration
}

func newSlowStopTransport(delay time.Duration) *slowStopTransport {
	return &slowStopTransport{
		startCh: make(chan struct{}),
		delay:   delay,
	}
}

func (s *slowStopTransport) Start(ctx context.Context) error {
	close(s.startCh)
	<-ctx.Done()
	return nil
}

func (s *slowStopTransport) Stop(_ context.Context) error {
	time.Sleep(s.delay)
	return nil
}

// ============================================================
// Tests
// ============================================================

func TestNew_DefaultOptions(t *testing.T) {
	app := New()
	if app == nil {
		t.Fatal("New() returned nil")
	}
	if app.id == "" {
		t.Error("expected non-empty id")
	}
	if app.stopTimeout != defaultStopTimeout {
		t.Errorf("expected stopTimeout=%v, got %v", defaultStopTimeout, app.stopTimeout)
	}
	if len(app.sigs) == 0 {
		t.Error("expected default signals to be set")
	}
}

func TestNew_WithOptions(t *testing.T) {
	tr := newMockTransport("http")
	reg := &mockRegistrar{}
	cfg := &mockConfigMate{}

	app := New(
		WithServer(tr),
		WithRegistrar(reg),
		WithConfig(cfg),
		WithServiceInfo("test-svc", "v1.0"),
		WithMetadata(map[string]string{"env": "test"}),
		WithStopTimeout(3*time.Second),
	)

	if len(app.transports) != 1 {
		t.Fatalf("expected 1 transport, got %d", len(app.transports))
	}
	if app.serviceName != "test-svc" {
		t.Errorf("expected serviceName=test-svc, got %s", app.serviceName)
	}
	if app.serviceVersion != "v1.0" {
		t.Errorf("expected serviceVersion=v1.0, got %s", app.serviceVersion)
	}
	if app.stopTimeout != 3*time.Second {
		t.Errorf("expected stopTimeout=3s, got %v", app.stopTimeout)
	}
	if app.metadata["env"] != "test" {
		t.Errorf("expected metadata env=test")
	}
}

func TestRun_StartAndCancelStop(t *testing.T) {
	tr := newMockTransport("http")

	app := New(
		WithServer(tr),
		WithStopTimeout(2*time.Second),
	)

	errCh := make(chan error, 1)
	go func() {
		errCh <- app.Run()
	}()

	// Wait for transport to start.
	select {
	case <-tr.startCh:
	case <-time.After(3 * time.Second):
		t.Fatal("transport did not start in time")
	}

	if !tr.started.Load() {
		t.Error("expected transport to be started")
	}

	// Cancel to trigger graceful shutdown.
	app.cancel()

	select {
	case err := <-errCh:
		if err != nil {
			t.Fatalf("Run() returned error: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Run() did not return in time after cancel")
	}

	if !tr.stopped.Load() {
		t.Error("expected transport to be stopped")
	}
}

func TestRun_SignalStop(t *testing.T) {
	tr := newMockTransport("http")

	app := New(
		WithServer(tr),
		WithStopTimeout(2*time.Second),
	)

	errCh := make(chan error, 1)
	go func() {
		errCh <- app.Run()
	}()

	// Wait for transport to start.
	select {
	case <-tr.startCh:
	case <-time.After(3 * time.Second):
		t.Fatal("transport did not start in time")
	}

	// Brief sleep to let signal handler goroutine be scheduled.
	time.Sleep(50 * time.Millisecond)

	// Send signal to trigger shutdown.
	syscall.Kill(syscall.Getpid(), syscall.SIGINT)

	select {
	case err := <-errCh:
		if err != nil {
			t.Fatalf("Run() returned error: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Run() did not return in time after signal")
	}

	if !tr.stopped.Load() {
		t.Error("expected transport to be stopped after signal")
	}
}

func TestRun_MultipleTransports(t *testing.T) {
	tr1 := newMockTransport("http")
	tr2 := newMockTransport("grpc")

	app := New(
		WithServer(tr1),
		WithServer(tr2),
		WithStopTimeout(2*time.Second),
	)

	errCh := make(chan error, 1)
	go func() {
		errCh <- app.Run()
	}()

	// Wait for both transports to start.
	select {
	case <-tr1.startCh:
	case <-time.After(3 * time.Second):
		t.Fatal("transport 1 did not start in time")
	}
	select {
	case <-tr2.startCh:
	case <-time.After(3 * time.Second):
		t.Fatal("transport 2 did not start in time")
	}

	app.cancel()

	select {
	case err := <-errCh:
		if err != nil {
			t.Fatalf("Run() returned error: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Run() did not return in time")
	}

	if !tr1.stopped.Load() || !tr2.stopped.Load() {
		t.Error("expected both transports to be stopped")
	}
}

func TestRun_WithRegistrar(t *testing.T) {
	tr := newMockEndpointTransport("http", "127.0.0.1:8080")
	reg := &mockRegistrar{}

	app := New(
		WithServer(tr),
		WithRegistrar(reg),
		WithServiceInfo("test-svc", "v1"),
		WithStopTimeout(2*time.Second),
	)

	errCh := make(chan error, 1)
	go func() {
		errCh <- app.Run()
	}()

	<-tr.startCh

	// Give time for registration to complete (Register is called synchronously
	// in Run() after launching goroutines, so a small sleep suffices).
	time.Sleep(200 * time.Millisecond)

	reg.mu.Lock()
	regCount := len(reg.registered)
	reg.mu.Unlock()

	if regCount != 1 {
		t.Fatalf("expected 1 registered service, got %d", regCount)
	}

	reg.mu.Lock()
	svc := reg.registered[0]
	reg.mu.Unlock()
	if svc.Name != "test-svc" {
		t.Errorf("expected service name test-svc, got %s", svc.Name)
	}
	if len(svc.Endpoints) != 1 || svc.Endpoints[0] != "127.0.0.1:8080" {
		t.Errorf("unexpected endpoints: %v", svc.Endpoints)
	}

	app.cancel()

	select {
	case err := <-errCh:
		if err != nil {
			t.Fatalf("Run() returned error: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Run() did not return in time")
	}

	reg.mu.Lock()
	deregCount := len(reg.deregistered)
	reg.mu.Unlock()
	if deregCount != 1 {
		t.Fatalf("expected 1 deregistered service, got %d", deregCount)
	}
}

func TestRun_RegisterError(t *testing.T) {
	tr := newMockTransport("http")
	reg := &mockRegistrar{registerErr: fmt.Errorf("register failed")}

	app := New(
		WithServer(tr),
		WithRegistrar(reg),
		WithServiceInfo("test-svc", "v1"),
		WithStopTimeout(2*time.Second),
	)

	errCh := make(chan error, 1)
	go func() {
		errCh <- app.Run()
	}()

	select {
	case err := <-errCh:
		if err == nil {
			t.Fatal("expected error from registration failure")
		}
		if err.Error() != "service register: register failed" {
			t.Errorf("unexpected error: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Run() did not return in time after registration failure")
	}

	// Transport should NOT have been started since registration fails before cycle.
	if tr.started.Load() {
		t.Error("expected transport to not be started after registration failure")
	}
}

func TestRun_TransportPanicRecovery(t *testing.T) {
	ptr := newPanicTransport()

	app := New(
		WithServer(ptr),
		WithStopTimeout(2*time.Second),
	)

	errCh := make(chan error, 1)
	go func() {
		errCh <- app.Run()
	}()

	select {
	case err := <-errCh:
		if err == nil {
			return // Panic recovered, clean exit
		}
		// Panic converts to error containing "transport panic"
		t.Logf("Run() returned (expected after panic recovery): %v", err)
	case <-time.After(5 * time.Second):
		t.Fatal("Run() did not return after transport panic")
	}
}

func TestRun_StopTimeout(t *testing.T) {
	tr := newSlowStopTransport(5 * time.Second) // Stop takes 5s

	app := New(
		WithServer(tr),
		WithStopTimeout(500*time.Millisecond), // But timeout is 500ms
	)

	errCh := make(chan error, 1)
	go func() {
		errCh <- app.Run()
	}()

	<-tr.startCh
	app.cancel()

	select {
	case err := <-errCh:
		if err == nil {
			t.Fatal("expected timeout error from stopAll")
		}
		t.Logf("Got expected timeout error: %v", err)
	case <-time.After(10 * time.Second):
		t.Fatal("Run() did not return despite stop timeout")
	}
}

func TestRun_TransportStopError(t *testing.T) {
	tr := newMockTransport("http")
	tr.stopErr = fmt.Errorf("stop failed")

	app := New(
		WithServer(tr),
		WithStopTimeout(2*time.Second),
	)

	errCh := make(chan error, 1)
	go func() {
		errCh <- app.Run()
	}()

	<-tr.startCh
	app.cancel()

	select {
	case err := <-errCh:
		if err == nil {
			t.Fatal("expected error from transport stop failure")
		}
		t.Logf("Got expected stop error: %v", err)
	case <-time.After(5 * time.Second):
		t.Fatal("Run() did not return in time")
	}
}

func TestRun_NoTransports(t *testing.T) {
	app := New(
		WithStopTimeout(1 * time.Second),
	)

	errCh := make(chan error, 1)
	go func() {
		errCh <- app.Run()
	}()

	// With no transports, only signal handler + stopAll are running.
	time.Sleep(100 * time.Millisecond)
	app.cancel()

	select {
	case err := <-errCh:
		if err != nil {
			t.Fatalf("Run() returned error: %v", err)
		}
	case <-time.After(3 * time.Second):
		t.Fatal("Run() did not return in time with no transports")
	}
}

func TestRun_WithConfig(t *testing.T) {
	tr := newMockTransport("http")
	cfg := &mockConfigMate{}

	app := New(
		WithServer(tr),
		WithConfig(cfg),
		WithStopTimeout(2*time.Second),
	)

	errCh := make(chan error, 1)
	go func() {
		errCh <- app.Run()
	}()

	<-tr.startCh
	app.cancel()

	select {
	case err := <-errCh:
		if err != nil {
			t.Fatalf("Run() returned error: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Run() did not return in time")
	}
}

func TestBuildServiceInstance(t *testing.T) {
	tr := newMockEndpointTransport("http", "0.0.0.0:8080")

	app := New(
		WithServer(tr),
		WithServiceInfo("my-svc", "v2"),
		WithMetadata(map[string]string{"region": "cn"}),
	)

	svc := app.buildServiceInstance()
	if svc == nil {
		t.Fatal("expected non-nil service instance")
	}
	if svc.Name != "my-svc" {
		t.Errorf("expected name=my-svc, got %s", svc.Name)
	}
	if svc.Version != "v2" {
		t.Errorf("expected version=v2, got %s", svc.Version)
	}
	if svc.Metadata["region"] != "cn" {
		t.Error("expected metadata region=cn")
	}
	if len(svc.Endpoints) != 1 || svc.Endpoints[0] != "0.0.0.0:8080" {
		t.Errorf("unexpected endpoints: %v", svc.Endpoints)
	}
}

func TestBuildServiceInstance_NoName(t *testing.T) {
	app := New()
	svc := app.buildServiceInstance()
	if svc != nil {
		t.Error("expected nil service instance when no service name is set")
	}
}

func TestRun_TransportStartError(t *testing.T) {
	tr := newMockTransport("http")
	tr.startErr = fmt.Errorf("bind: address already in use")

	app := New(
		WithServer(tr),
		WithStopTimeout(2*time.Second),
	)

	errCh := make(chan error, 1)
	go func() {
		errCh <- app.Run()
	}()

	select {
	case err := <-errCh:
		if err == nil {
			t.Fatal("expected error from transport start failure")
		}
		t.Logf("Got expected start error: %v", err)
	case <-time.After(5 * time.Second):
		t.Fatal("Run() did not return in time after start error")
	}
}

// --- Reloadable Mock Transport (supports multiple start/stop cycles) ---

type reloadableMockTransport struct {
	name       string
	started    atomic.Int32
	stopped    atomic.Int32
	firstStart sync.Once
	startCh    chan struct{} // closed on FIRST start only
	mu         sync.Mutex
	stopCh     chan struct{} // recreated per cycle
}

func newReloadableMockTransport(name string) *reloadableMockTransport {
	return &reloadableMockTransport{
		name:    name,
		startCh: make(chan struct{}),
	}
}

func (m *reloadableMockTransport) Start(ctx context.Context) error {
	m.started.Add(1)
	m.firstStart.Do(func() {
		close(m.startCh)
	})
	ch := make(chan struct{})
	m.mu.Lock()
	m.stopCh = ch
	m.mu.Unlock()
	select {
	case <-ctx.Done():
		return nil
	case <-ch:
		return nil
	}
}

func (m *reloadableMockTransport) Stop(_ context.Context) error {
	m.stopped.Add(1)
	m.mu.Lock()
	ch := m.stopCh
	m.mu.Unlock()
	if ch != nil {
		select {
		case <-ch:
		default:
			close(ch)
		}
	}
	return nil
}

func (m *reloadableMockTransport) TransportName() string {
	return m.name
}

// --- Mock ConfigMate with OnChange and Reload support ---

type mockReloadableConfig struct {
	reloadCount int
	mu          sync.Mutex
	onChange    []func()
}

func (m *mockReloadableConfig) PrimitiveDecode(_ ...conf.ConfigureMatcher) error {
	return nil
}

func (m *mockReloadableConfig) OnChange(fn func()) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.onChange = append(m.onChange, fn)
}

func (m *mockReloadableConfig) Reload() error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.reloadCount++
	return nil
}

// triggerChange simulates a config change by calling all registered onChange callbacks.
func (m *mockReloadableConfig) triggerChange() {
	m.mu.Lock()
	fns := make([]func(), len(m.onChange))
	copy(fns, m.onChange)
	m.mu.Unlock()
	for _, fn := range fns {
		fn()
	}
}

func (m *mockReloadableConfig) getReloadCount() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.reloadCount
}

// ============================================================
// Reload Tests
// ============================================================

func TestRun_SIGHUPReload(t *testing.T) {
	tr := newReloadableMockTransport("http")

	app := New(
		WithServer(tr),
		WithStopTimeout(2*time.Second),
	)

	errCh := make(chan error, 1)
	go func() {
		errCh <- app.Run()
	}()

	// Wait for transport to start.
	select {
	case <-tr.startCh:
	case <-time.After(3 * time.Second):
		t.Fatal("transport did not start in time")
	}

	// Brief sleep to let signal handler goroutine be scheduled.
	time.Sleep(50 * time.Millisecond)

	// Send SIGHUP to trigger reload.
	syscall.Kill(syscall.Getpid(), syscall.SIGHUP)

	// Wait for reload to take effect.
	time.Sleep(500 * time.Millisecond)

	// The transport should have been stopped (old generation drained) and restarted.
	if tr.stopped.Load() < 1 {
		t.Error("expected transport to be stopped during reload")
	}
	if tr.started.Load() < 2 {
		t.Error("expected transport to be started at least twice (initial + reload)")
	}

	// Process should still be running.
	select {
	case err := <-errCh:
		t.Fatalf("Run() returned unexpectedly: %v", err)
	case <-time.After(200 * time.Millisecond):
		// Good — still running.
	}

	// Now shut down.
	app.cancel()

	select {
	case err := <-errCh:
		if err != nil {
			t.Fatalf("Run() returned error: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Run() did not return in time after cancel")
	}
}

func TestRun_ConfigChangeReload(t *testing.T) {
	tr := newReloadableMockTransport("http")
	cfg := &mockReloadableConfig{}

	app := New(
		WithServer(tr),
		WithConfig(cfg),
		WithStopTimeout(2*time.Second),
	)

	errCh := make(chan error, 1)
	go func() {
		errCh <- app.Run()
	}()

	// Wait for transport to start.
	select {
	case <-tr.startCh:
	case <-time.After(3 * time.Second):
		t.Fatal("transport did not start in time")
	}

	time.Sleep(50 * time.Millisecond)

	// Simulate config change.
	cfg.triggerChange()

	// Wait for reload to complete.
	time.Sleep(500 * time.Millisecond)

	// The transport should have been stopped for reload and restarted.
	if tr.stopped.Load() < 1 {
		t.Error("expected transport to be stopped during config change reload")
	}
	if tr.started.Load() < 2 {
		t.Error("expected transport to be started at least twice")
	}

	// Config Reload() should have been called.
	if cfg.getReloadCount() < 1 {
		t.Error("expected config Reload() to be called")
	}

	// Process should still be running.
	select {
	case err := <-errCh:
		t.Fatalf("Run() returned unexpectedly: %v", err)
	case <-time.After(200 * time.Millisecond):
		// Good.
	}

	app.cancel()

	select {
	case err := <-errCh:
		if err != nil {
			t.Fatalf("Run() returned error: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Run() did not return in time after cancel")
	}
}

func TestRun_SIGHUPWithReloadableConfig(t *testing.T) {
	tr := newReloadableMockTransport("http")
	cfg := &mockReloadableConfig{}

	app := New(
		WithServer(tr),
		WithConfig(cfg),
		WithStopTimeout(2*time.Second),
	)

	errCh := make(chan error, 1)
	go func() {
		errCh <- app.Run()
	}()

	<-tr.startCh
	time.Sleep(50 * time.Millisecond)

	// Send SIGHUP.
	syscall.Kill(syscall.Getpid(), syscall.SIGHUP)

	time.Sleep(500 * time.Millisecond)

	// Verify config was explicitly reloaded (for SIGHUP-triggered reload).
	if cfg.getReloadCount() < 1 {
		t.Error("expected config Reload() to be called on SIGHUP")
	}

	app.cancel()

	select {
	case err := <-errCh:
		if err != nil {
			t.Fatalf("Run() returned error: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Run() did not return in time")
	}
}

func TestPIDFile(t *testing.T) {
	pidPath := t.TempDir() + "/test.pid"

	tr := newMockTransport("http")

	app := New(
		WithServer(tr),
		WithPIDFile(pidPath),
		WithStopTimeout(2*time.Second),
	)

	errCh := make(chan error, 1)
	go func() {
		errCh <- app.Run()
	}()

	<-tr.startCh
	time.Sleep(50 * time.Millisecond)

	// Verify PID file exists and contains correct PID.
	data, err := os.ReadFile(pidPath)
	if err != nil {
		t.Fatalf("failed to read pid file: %v", err)
	}
	expected := fmt.Sprintf("%d\n", os.Getpid())
	if string(data) != expected {
		t.Errorf("pid file content = %q, want %q", string(data), expected)
	}

	app.cancel()

	select {
	case err := <-errCh:
		if err != nil {
			t.Fatalf("Run() returned error: %v", err)
		}
	case <-time.After(5 * time.Second):
		t.Fatal("Run() did not return in time")
	}

	// Verify PID file is removed after shutdown.
	if _, err := os.Stat(pidPath); !os.IsNotExist(err) {
		t.Error("expected pid file to be removed after shutdown")
	}
}
