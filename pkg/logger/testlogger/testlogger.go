package testlogger

import (
	"context"
	"testing"

	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"

	"github.com/freiheit-com/kuberpult/pkg/logger"
)

func Wrap(ctx context.Context, inner func(ctx context.Context)) *observer.ObservedLogs {
	config, obs := observer.New(zap.DebugLevel)
	log := zap.New(config)
	defer func() {
		if err := log.Sync(); err != nil {
			panic(err)
		}
	}()
	inner(logger.WithLogger(ctx, log))
	return obs
}

func AssertEmpty(t *testing.T, logs *observer.ObservedLogs) {
	l := logs.All()
	if len(l) != 0 {
		t.Errorf("expected no logs, but got %d\nexample: \t%#v",len(l), l[0])
	}
}

func AssertLogs(t *testing.T, logs *observer.ObservedLogs, tests ...func(observer.LoggedEntry)) {
	l := logs.All()
	if len(l) > len(tests) {
		t.Errorf("expected %d logs, but got %d\nexample: \t%#v",len(tests),len(l), l[len(tests)])
	} else if len(l) < len(tests) {
		t.Errorf("expected %d logs, but got %d",len(tests),len(l))
	} else {
		for i, j := range l {
			tests[i](j)
		}
	}
}
