package instrument

import (
	"errors"
	"time"
)

type Instrument struct {
	sink Sinker
}

func (i *Instrument) shutdown() {
	i.sink.Shutdown()
}

// Counter records a counter measurement with a string name and int value.
func (i *Instrument) Counter(measurement string, value int) {
	i.sink.Counter(Counter{measurement, value})
}

// Gauge records a gauge measurement with a string name and a float64 value.
func (i *Instrument) Gauge(measurement string, value float64) {
	i.sink.Gauge(Gauge{measurement, value})
}

// Timer records a timer measurement with a string name and the time elapsed since the received init Time and the
// current time.
func (i *Instrument) Timer(measurement string, init time.Time) {
	elapsed := time.Since(init)

	i.sink.Timer(Timer{measurement, elapsed})
}

// NewInstrument initialises a new Instrument with a default Stdout Sink.
//
// Recorded measurements are printed to stdout every 5 seconds.
func NewInstrument() *Instrument {
	return &Instrument{sink: NewStdoutSink()}
}

// NewInstrumentWithSinker initialises a new Instrument with a custom provided Sinker.
func NewInstrumentWithSinker(sink Sinker) (*Instrument, error) {
	if sink == nil {
		return nil, errors.New("sinker cannot be nil")
	}

	return &Instrument{sink: sink}, nil
}

// Now is a helper function wrapper of the current time.
func Now() time.Time {
	return time.Now()
}
