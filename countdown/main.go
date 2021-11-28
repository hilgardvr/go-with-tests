package main

import (
	"fmt"
	"io"
	"os"
	"time"
)

const finalWord = "Go!"
const countDownStart = 3

type Sleep interface {
	Sleep()
}

type ConfigurableSleeper struct {
	duration	time.Duration
	sleep 		func(time.Duration)
}

func (c *ConfigurableSleeper)Sleep() {
	c.sleep(c.duration)
}

func CountDown(w io.Writer, s Sleep) {
	for i := countDownStart; i > 0; i-- {
		s.Sleep()
		fmt.Fprintf(w, "%d\n", i)
	}
	s.Sleep()
	fmt.Fprint(w, finalWord)
}

func main() {
	sleeper := &ConfigurableSleeper{1 * time.Second, time.Sleep}
	CountDown(os.Stdout, sleeper)
}
