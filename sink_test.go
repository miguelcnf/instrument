package instrument

import (
	"os"
	"testing"
	"time"
)

func TestSink_ThreadSafety(t *testing.T) {
	// Suppress stdout in this test
	os.Stdout, _ = os.Open(os.DevNull)
	sink := NewStdoutSink()

	done := make(chan bool)
	go func() {
		for i := 1; i <= 50; i++ {
			// Flush every 10 iterations concurrently with previous writes
			if i%10 == 0 {
				sink.Flush()
			}

			go sink.Timer(Timer{"test", time.Duration(i) * time.Millisecond})

			done <- true
		}
	}()

	for i := 1; i <= 50; i++ {
		<-done
	}
	close(done)
	sink.Shutdown()
}
