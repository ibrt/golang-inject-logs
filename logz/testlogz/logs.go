//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -source ../logs.go -destination ./mocks.go -package testlogz

package testlogz

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ibrt/golang-fixtures/fixturez"

	"golang-inject-logs/logz"
)

var (
	_ fixturez.BeforeSuite = &Helper{}
	_ fixturez.AfterSuite  = &Helper{}
	_ fixturez.BeforeTest  = &MockHelper{}
	_ fixturez.AfterTest   = &MockHelper{}
)

// Helper provides a test helper for logz using a real logger.
type Helper struct {
	releaser func()
}

// BeforeSuite implements fixturez.BeforeSuite.
func (f *Helper) BeforeSuite(ctx context.Context, _ *testing.T) context.Context {
	cfg := &logz.LogsConfig{
		SentryLevel:      logz.Debug,
		OutputLevel:      logz.Debug,
		OutputFormat:     logz.Text,
		SentryDSN:        "",
		SentrySampleRate: 1,
		ReleaseTimeout:   1,
		Environment:      "test",
		Release:          "test",
		ServerName:       "test",
	}

	ctx = logz.NewConfigSingletonInjector(cfg)(ctx)
	injector, releaser := logz.Initializer(ctx)
	f.releaser = releaser
	return injector(ctx)
}

// AfterSuite implements fixturez.AfterSuite.
func (f *Helper) AfterSuite(_ context.Context, _ *testing.T) {
	f.releaser()
	f.releaser = nil
}

// MockHelper provides a test helper for logz using a mock logger.
type MockHelper struct {
	Mock *MockLogs
	ctrl *gomock.Controller
}

// BeforeTest implements fixtures.BeforeTest.
func (f *MockHelper) BeforeTest(ctx context.Context, t *testing.T) context.Context {
	f.ctrl = gomock.NewController(t)
	f.Mock = NewMockLogs(f.ctrl)
	return logz.NewSingletonInjector(f.Mock)(ctx)
}

// AfterTest implements fixtures.AfterTest.
func (f *MockHelper) AfterTest(_ context.Context, _ *testing.T) {
	f.ctrl.Finish()
	f.ctrl = nil
	f.Mock = nil
}
