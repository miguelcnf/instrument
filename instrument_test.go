package instrument

import (
	"testing"
	"time"
)

const timerSleep = "timer.sleep.function"
const timerSleepAsync = "timer.sleep.function.async"

var instrument *Instrument

type mockSink struct {
	receivedTimer Timer
	receivedCounter Counter
	receivedGauge Gauge
	shutdownIsCalled bool

}

func (m *mockSink) Timer(timer Timer) {
	m.receivedTimer = timer
}

func (m *mockSink) Counter(counter Counter) {
	m.receivedCounter = counter
}

func (m *mockSink) Gauge(gauge Gauge) {
	m.receivedGauge = gauge
}

func (m *mockSink) Shutdown() {
	m.shutdownIsCalled = true
}

func sleepFunction() {
	defer instrument.Timer(timerSleep, Now())

	time.Sleep(1*time.Second)
}

func sleepFunctionAsync(done chan bool) {
	defer instrument.Timer(timerSleepAsync, Now())

	time.Sleep(1*time.Second)

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

	if sink.receivedTimer.name != timerSleep {
		t.Errorf("expected received timer name to be: %v but was: %v", timerSleep, sink.receivedTimer.name)
	}

	if sink.receivedTimer.value.Seconds() < 1*time.Second.Seconds() {
		t.Errorf("expected receiver timer value to be at least 1 second but was: %vs", sink.receivedTimer.value)
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

	_ = <- done
	close(done)

	if sink.receivedTimer.name != timerSleepAsync {
		t.Errorf("expected received timer name to be: %v but was: %v", timerSleepAsync, sink.receivedTimer.name)
	}

	if sink.receivedTimer.value.Seconds() < 1*time.Second.Seconds() {
		t.Errorf("expected receiver timer value to be at least 1 second but was: %vs", sink.receivedTimer.value)
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
