package ingest

import (
	"context"
	"time"
)

func Retry(ctx context.Context, attempts int, delay time.Duration, action func(context.Context) error) error {
	if attempts < 1 {
		attempts = 1
	}
	var err error
	for i := 0; i < attempts; i++ {
		if err = action(ctx); err == nil {
			return nil
		}
		if i+1 == attempts {
			break
		}
		timer := time.NewTimer(delay)
		select {
		case <-timer.C:
		case <-ctx.Done():
			if !timer.Stop() {
				<-timer.C
			}
			return ctx.Err()
		}
	}
	return err
}
