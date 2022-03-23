package testlogz_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ibrt/golang-fixtures/fixturez"
	"github.com/stretchr/testify/require"

	"github.com/ibrt/golang-inject-logs/logz"
	"github.com/ibrt/golang-inject-logs/logz/testlogz"
)

func TestHelpers(t *testing.T) {
	fixturez.RunSuite(t, &Suite{})
	fixturez.RunSuite(t, &MockSuite{})
}

type Suite struct {
	*fixturez.DefaultConfigMixin
	Logs *testlogz.Helper
}

func (s *Suite) TestHelper(ctx context.Context, t *testing.T) {
	logs := logz.Get(ctx)
	require.NotNil(t, logs)
}

type MockSuite struct {
	*fixturez.DefaultConfigMixin
	Logs *testlogz.MockHelper
}

func (s *MockSuite) TestMockHelper(ctx context.Context, t *testing.T) {
	s.Logs.Mock.EXPECT().Debug(gomock.Any(), gomock.Eq(1), "message")
	logz.Get(ctx).Debug("message")
}
