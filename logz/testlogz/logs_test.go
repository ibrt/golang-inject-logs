package testlogz_test

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/ibrt/golang-fixtures/fixturez"
	"github.com/stretchr/testify/require"

	"golang-inject-logs/logz"
	"golang-inject-logs/logz/testlogz"
)

func TestHelpers(t *testing.T) {
	fixturez.RunSuite(t, &Suite{})
	fixturez.RunSuite(t, &MockSuite{})
}

type Suite struct {
	*fixturez.DefaultConfigMixin
	Logz *testlogz.Helper
}

func (s *Suite) TestHelper(ctx context.Context, t *testing.T) {
	logs := logz.Get(ctx)
	require.NotNil(t, logs)
}

type MockSuite struct {
	*fixturez.DefaultConfigMixin
	Logz *testlogz.MockHelper
}

func (s *MockSuite) TestMockHelper(ctx context.Context, t *testing.T) {
	s.Logz.Mock.EXPECT().Debug(gomock.Any(), gomock.Eq(1), "message")
	logz.Get(ctx).Debug("message")
}
