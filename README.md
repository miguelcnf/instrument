# instrument

[![GoDoc](https://godoc.org/github.com/miguelcnf/instrument?status.svg)](https://godoc.org/github.com/miguelcnf/instrument)

Instrument is an overly simplified, self-contained and idiomatic library (calling it a library is such an over-statement) to provide code level instrumentation for go programs.

# Concepts

* No 3rd-party libraries
* Use native defer concept to measure timers
* Pluggable custom sinks to handle (export/push/discard) measurements

# Installing

Install by running:

```shell
go get github.com/miguelcnf/instrument
```

# Example

## Default Stdout Sink

Testing it with the default stdout sink.

```go
package main

import (
	"github.com/miguelcnf/instrument"
	"math/rand"
	"time"
)

var inst *instrument.Instrument

func sleepFunction() {
	defer inst.Timer("timer.sleepFunction", instrument.Now())

	r := rand.Intn(500)
	time.Sleep(time.Duration(r)*time.Millisecond)
}

func main() {
	// Initialise instrument
   	inst = instrument.NewInstrument()
   
   	// Generate timers
	for i := 0; i<3; i++ {
        	go sleepFunction()
    	}
    
    	// Sleep the main thread to allow the stdout sink to print recorded measurements
    	time.Sleep(10*time.Second)
}
```

The previous code should output the following before exiting:

```
measurement: timer.sleepFunction; type: timer; value: 50.922904ms
measurement: timer.sleepFunction; type: timer; value: 85.846804ms
measurement: timer.sleepFunction; type: timer; value: 87.511807ms
```

## Custom Sink

Initialise it with a custom provided sink.
 
The sink must implement the `Sinker` interface.

```go
package main

import (
	"github.com/miguelcnf/instrument"
)

var inst *instrument.Instrument

func main() {
	// Initialise instrument with a provided custom sink
	var err error
    	inst, err = instrument.NewInstrumentWithSinker(&CustomSinkImplementation{})
    	if err != nil {
        	// handle error
    	}
}
```

# Notice

This is an experiment and is in no way ready to be used in production.

# License

See the [LICENSE](LICENSE.md) file for license rights and limitations (MIT).


