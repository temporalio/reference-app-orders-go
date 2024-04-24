package temporalutil

import "context"

// WorkerInterruptFromContext returns a channel that will be closed when the context is done.
func WorkerInterruptFromContext(ctx context.Context) <-chan any {
	intCh := make(chan any)

	go func() {
		defer close(intCh)
		<-ctx.Done()
	}()

	return intCh
}
