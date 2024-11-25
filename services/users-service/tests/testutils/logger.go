package testutils

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
)

func NewTestLogger(t zaptest.TestingT) *zap.Logger {
	return zaptest.NewLogger(t)
}
