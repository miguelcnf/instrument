package instrument

import (
	"fmt"
	"time"
)

const flushPeriod = 5 * time.Second

type Sinker interface {
	Timer(timer Timer)
	Counter(counter Counter)
	Gauge(gauge Gauge)
	Flush()
	Shutdown()
}

type Timer struct {
	name  string
	value time.Duration
}

type Counter struct {
	name  string
	value int
}

type Gauge struct {
	name  string
	value float64
}

type Sink struct {
	flush    chan bool
	input    chan interface{}
	quit     chan bool
	timers   []Timer
	counters []Counter
	gauges   []Gauge
}

func (s *Sink) Timer(timer Timer) {
	s.input <- timer
}

func (s *Sink) Counter(counter Counter) {
	s.input <- counter
}

func (s *Sink) Gauge(gauge Gauge) {
	s.input <- gauge
}

func (s *Sink) Flush() {
	s.flush <- true
}

func (s *Sink) print() {
	fmt.Println("buffer: flushing measurements")
	s.printTimers()
	s.printCounters()
	s.printGauges()
}

func (s *Sink) Shutdown() {
	s.quit <- true
}

func NewStdoutSink() *Sink {
	sink := &Sink{}
	sink.quit, sink.input, sink.flush = sink.process()
	return sink
}

func (s *Sink) printTimers() {
	for _, timer := range s.timers {
		val := float64(timer.value) / float64(time.Millisecond)
		fmt.Printf("measurement: %v; type: timer; value: %vms\n", timer.name, val)
	}
	s.timers = []Timer{}
}

func (s *Sink) printCounters() {
	for _, counter := range s.counters {
		fmt.Printf("measurement: %v; type: counter; value: %v\n", counter.name, counter.value)
	}
	s.counters = []Counter{}
}

func (s *Sink) printGauges() {
	for _, gauge := range s.gauges {
		fmt.Printf("measurement: %v; type: gauge; value: %v\n", gauge.name, gauge.value)
	}
	s.gauges = []Gauge{}
}

func (s *Sink) process() (quit chan bool, input chan interface{}, flush chan bool) {
	ticker := time.NewTicker(flushPeriod)

	flush = make(chan bool)
	// Buffered input channel
	input = make(chan interface{}, 10)
	quit = make(chan bool)

	go func() {
		for {
			select {
			case rcv := <-input:
				s.receive(rcv)
			case <-flush:
				s.print()
			case <-ticker.C:
				s.print()
			case <-quit:
				ticker.Stop()
				// Not closing the input channel to avoid panic from producer go routines
				close(quit)
				s.print()
				return
			}
		}
	}()

	return
}

func (s *Sink) receive(measurement interface{}) {
	switch measurement.(type) {
	case Timer:
		s.timers = append(s.timers, measurement.(Timer))
	case Counter:
		s.counters = append(s.counters, measurement.(Counter))
	case Gauge:
		s.gauges = append(s.gauges, measurement.(Gauge))
	}
}
