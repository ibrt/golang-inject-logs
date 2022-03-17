package logz_test

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/ibrt/golang-errors/errorz"
	"github.com/ibrt/golang-fixtures/fixturez"
	"github.com/stretchr/testify/require"

	"golang-inject-logs/logz"
)

func TestModule_Debug(t *testing.T) {
	c := fixturez.CaptureOutput()
	defer c.Close()

	ctx := getDefaultCfg()
	injector, releaser := logz.Initializer(ctx)
	defer releaser()

	logs := logz.Get(injector(ctx))
	require.NotNil(t, logs)

	logs.Debug("message: %v", logz.A("value"), logz.M("k", "v"))
	v := make(map[string]interface{})
	fixturez.RequireNoError(t, json.Unmarshal(c.GetErr(), &v))
	require.Equal(t, "v", v["k"])
	require.Equal(t, "debug", v["level"])
	require.Equal(t, "message: value", v["msg"])
	require.NotNil(t, v["time"])
}

func TestModule_Info(t *testing.T) {
	c := fixturez.CaptureOutput()
	defer c.Close()

	ctx := getDefaultCfg()
	injector, releaser := logz.Initializer(ctx)
	defer releaser()

	logs := logz.Get(injector(ctx))
	require.NotNil(t, logs)

	logs.Debug("message: %v", logz.A("value"), logz.M("k", "v"))
	v := make(map[string]interface{})
	fixturez.RequireNoError(t, json.Unmarshal(c.GetErr(), &v))
	require.Equal(t, "v", v["k"])
	require.Equal(t, "debug", v["level"])
	require.Equal(t, "message: value", v["msg"])
	require.NotNil(t, v["time"])
}

func TestModule_Warning(t *testing.T) {
	c := fixturez.CaptureOutput()
	defer c.Close()

	ctx := getDefaultCfg()
	injector, releaser := logz.Initializer(ctx)
	defer releaser()

	logs := logz.Get(injector(ctx))
	require.NotNil(t, logs)

	logs.Warning(errorz.Errorf("message: %v", errorz.A("value"), errorz.M("k", "v")))
	v := make(map[string]interface{})
	fixturez.RequireNoError(t, json.Unmarshal(c.GetErr(), &v))
	require.Equal(t, "v", v["k"])
	require.Equal(t, "warning", v["level"])
	require.Equal(t, "message: value", v["msg"])
	require.NotNil(t, v["time"])
}

func TestModule_Error(t *testing.T) {
	c := fixturez.CaptureOutput()
	defer c.Close()

	ctx := getDefaultCfg()
	injector, releaser := logz.Initializer(ctx)
	defer releaser()

	logs := logz.Get(injector(ctx))
	require.NotNil(t, logs)

	logs.Warning(errorz.Errorf("message: %v", errorz.A("value"), errorz.M("k", "v")))
	v := make(map[string]interface{})
	fixturez.RequireNoError(t, json.Unmarshal(c.GetErr(), &v))
	require.Equal(t, "v", v["k"])
	require.Equal(t, "warning", v["level"])
	require.Equal(t, "message: value", v["msg"])
	require.NotNil(t, v["time"])
}

func getDefaultCfg() context.Context {
	cfg := &logz.LogsConfig{
		SentryLevel:      logz.Debug,
		OutputLevel:      logz.Debug,
		OutputFormat:     logz.JSON,
		SentryDSN:        "",
		SentrySampleRate: 1,
		ReleaseTimeout:   1,
		Environment:      "environment",
		Release:          "release",
		ServerName:       "serverName",
	}

	return logz.NewConfigSingletonInjector(cfg)(context.Background())
}
