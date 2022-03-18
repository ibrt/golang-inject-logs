//go:generate go run github.com/golang/mock/mockgen@v1.6.0 -source ../logs.go -destination ./mocklogz/mocks.go -package mocklogz

package testlogz

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ibrt/golang-fixtures/fixturez"

	"github.com/ibrt/golang-inject-logs/logz"
	"github.com/ibrt/golang-inject-logs/logz/testlogz/mocklogz"
)

var (
	_ fixturez.BeforeSuite = &Helper{}
	_ fixturez.AfterSuite  = &Helper{}
	_ fixturez.BeforeTest  = &MockHelper{}
	_ fixturez.AfterTest   = &MockHelper{}
)

// Helper is a test helper for Logs.
type Helper struct {
	releaser func()
}

// BeforeSuite implements fixturez.BeforeSuite.
func (f *Helper) BeforeSuite(ctx context.Context, _ *testing.T) context.Context {
	cfg := &logz.Config{
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

// MockHelper is a test helper for Logs.
type MockHelper struct {
	Mock *mocklogz.MockLogs
	ctrl *gomock.Controller
}

// BeforeTest implements fixtures.BeforeTest.
func (f *MockHelper) BeforeTest(ctx context.Context, t *testing.T) context.Context {
	f.ctrl = gomock.NewController(t)
	f.Mock = mocklogz.NewMockLogs(f.ctrl)
	return logz.NewSingletonInjector(f.Mock)(ctx)
}

// AfterTest implements fixtures.AfterTest.
func (f *MockHelper) AfterTest(_ context.Context, _ *testing.T) {
	f.ctrl.Finish()
	f.ctrl = nil
	f.Mock = nil
}
