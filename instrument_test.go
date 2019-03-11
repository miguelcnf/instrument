package instrument

import (
	"sync"
	"testing"
	"time"
)

const timerSleep = "timer.sleep.function"
const timerSleepAsync = "timer.sleep.function.async"

var instrument *Instrument

type mockSink struct {
	// Mutex to not trigger data races false positives
	lock             sync.Mutex
	rcvTimer         Timer
	shutdownIsCalled bool
}

func (m *mockSink) Flush() {
	panic("implement me")
}

func (m *mockSink) Gauge(gauge Gauge) {
	panic("implement me")
}

func (m *mockSink) Counter(counter Counter) {
	panic("implement me")
}

func (m *mockSink) Timer(timer Timer) {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.rcvTimer = timer
}

func (m *mockSink) RcvTimer() Timer {
	m.lock.Lock()
	defer m.lock.Unlock()
	return m.rcvTimer
}

func (m *mockSink) Shutdown() {
	m.shutdownIsCalled = true
}

func sleepFunction() {
	defer instrument.Timer(timerSleep, Now())

	time.Sleep(1 * time.Second)
}

func sleepFunctionAsync(done chan bool) {
	defer instrument.Timer(timerSleepAsync, Now())

	time.Sleep(1 * time.Second)

	done <- true
}

func TestTimer(t *testing.T) {
	sink := &mockSink{}

	var err error
	instrument, err = NewInstrumentWithSinker(sink)
	if err != nil {
		t.Fatalf("error creating instrument with sinker")
	}

	sleepFunction()

	if sink.RcvTimer().name != timerSleep {
		t.Errorf("expected received timer name to be: %v but was: %v", timerSleep, sink.RcvTimer().name)
	}

	if sink.RcvTimer().value.Seconds() < 1*time.Second.Seconds() {
		t.Errorf("expected receiver timer value to be at least 1 second but was: %vs", sink.RcvTimer().value)
	}

	instrument.shutdown()

	if !sink.shutdownIsCalled {
		t.Errorf("expected shutdown to be called but wasn't")
	}
}

func TestTimerAsync(t *testing.T) {
	sink := &mockSink{}

	var err error
	instrument, err = NewInstrumentWithSinker(sink)
	if err != nil {
		t.Fatalf("error creating instrument with sinker")
	}

	done := make(chan bool)
	go sleepFunctionAsync(done)

	_ = <-done
	close(done)

	if sink.RcvTimer().name != timerSleepAsync {
		t.Errorf("expected received timer name to be: %v but was: %v", timerSleepAsync, sink.RcvTimer().name)
	}

	if sink.RcvTimer().value.Seconds() < 1*time.Second.Seconds() {
		t.Errorf("expected receiver timer value to be at least 1 second but was: %vs", sink.RcvTimer().value)
	}

	instrument.shutdown()

	if !sink.shutdownIsCalled {
		t.Errorf("expected shutdown to be called but wasn't")
	}
}

func TestNewInstrumentWithNilSinker(t *testing.T) {
	_, err := NewInstrumentWithSinker(nil)

	if err == nil {
		t.Errorf("expected NewInstrumentWithSinker to fail when receiving a nil sinker")
	}
}
