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
	quit     chan struct{}
	timers   []Timer
	counters []Counter
	gauges   []Gauge
}

func (s *Sink) Timer(timer Timer) {
	s.timers = append(s.timers, timer)
}

func (s *Sink) Counter(counter Counter) {
	s.counters = append(s.counters, counter)
}

func (s *Sink) Gauge(gauge Gauge) {
	s.gauges = append(s.gauges, gauge)
}

func (s *Sink) flush() {
	s.flushTimers()
	s.flushCounters()
	s.flushGauges()
}

func (s *Sink) flushTimers() {
	for _, timer := range s.timers {
		val := float64(timer.value) / float64(time.Millisecond)
		fmt.Printf("measurement: %v; type: timer; value: %vms\n", timer.name, val)
	}
	s.timers = []Timer{}
}

func (s *Sink) flushCounters() {
	for _, counter := range s.counters {
		fmt.Printf("measurement: %v; type: counter; value: %v\n", counter.name, counter.value)
	}
	s.counters = []Counter{}
}

func (s *Sink) flushGauges() {
	for _, gauge := range s.gauges{
		fmt.Printf("measurement: %v; type: gauge; value: %v\n", gauge.name, gauge.value)
	}
	s.gauges = []Gauge{}
}

func NewStdoutSink() *Sink {
	sink := &Sink{}
	sink.quit = sink.flushCycle()
	return sink
}

func (s *Sink) Shutdown() {
	close(s.quit)
	s.flush()
}

func (s *Sink) flushCycle() chan struct{} {
	ticker := time.NewTicker(flushPeriod)

	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				s.flush()
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()

	return quit
}
