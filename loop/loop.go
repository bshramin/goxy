package loop

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
)

func InfiniteLoop(ctx context.Context, name string, f func(context.Context) error, reDoDuration, failDuration time.Duration) context.Context {
	go func() {
		DoTasKWithRetry(ctx, name, f, failDuration)
		select {
		case <-ctx.Done():
			logrus.Debug(fmt.Sprintf("exit infinite loop: %s", name))
			return
		case <-time.After(reDoDuration):
		}
	}()
	return ctx
}

func DoTasKWithRetry(c context.Context, name string, f func(context.Context) error, failDuration time.Duration) {
	ctx, cl := context.WithCancel(c)
	defer func() {
		cl()
		e := recover()
		if e != nil {
			logrus.Errorf(fmt.Sprintf("%s failed", name))
		}
	}()
	for {
		err := f(ctx)
		if err == nil {
			return
		}
		select {
		case <-ctx.Done():
			logrus.Debug(fmt.Sprintf("exit infinite loop: %s", name))
			return
		case <-time.After(failDuration):
		}
	}
}
