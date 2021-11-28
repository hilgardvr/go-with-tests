package main

import (
	"bytes"
	"reflect"
	"testing"
	"time"
)

const sleep = "sleep"
const write = "write"

type SpySleeper struct {
	Calls int
}

func (s *SpySleeper)Sleep() {
	s.Calls++
}

type SpySleeperOperations struct {
	Calls 		[]string
}

func (s *SpySleeperOperations)Sleep() {
	s.Calls = append(s.Calls, sleep)
}

func (s *SpySleeperOperations)Write(p []byte) (n int, err error) {
	s.Calls = append(s.Calls, write)
	return
}

type SpyTime struct {
	durationSlept time.Duration
}

func (s *SpyTime)Sleep(duration time.Duration) {
	s.durationSlept = duration
}


func TestCountDown(t *testing.T) {

	t.Run("count sleep call & output", func(t *testing.T) {
		buffer := &bytes.Buffer{}
		spySleeper := &SpySleeper{}

		CountDown(buffer, spySleeper)

		got := buffer.String()
		want := "3\n2\n1\nGo!"

		if got != want {
			t.Errorf("got %s want %s", got, want)
		}

		if spySleeper.Calls != 4 {
			t.Errorf("not enough calls to spysleeper, wanted 4 got %d", spySleeper.Calls)
		}
	})
	t.Run("check order of operations", func(t *testing.T) {
		spySleeperOpertions := &SpySleeperOperations{}

		CountDown(spySleeperOpertions, spySleeperOpertions)

		want := []string{
			sleep,
			write,
			sleep,
			write,
			sleep,
			write,
			sleep,
			write,
		}

		if (!reflect.DeepEqual(spySleeperOpertions.Calls, want)) {
			t.Errorf("got %s want %s", spySleeperOpertions.Calls, want)
		}
	})

	t.Run("check order of operations", func(t *testing.T) {
		sleepTime := 5 * time.Second

		spyTime := &SpyTime{}

		sleeper := ConfigurableSleeper{sleepTime, spyTime.Sleep}
		sleeper.Sleep()

		if (spyTime.durationSlept != sleepTime) {
			t.Errorf("should have slept for %v but slept for %v", sleepTime, spyTime.durationSlept)
		}
	})
}
