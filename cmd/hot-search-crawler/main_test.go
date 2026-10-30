package main

import (
	"testing"
	"time"
)

func TestDaemonBackoff(t *testing.T) {
	cases := map[int]time.Duration{
		1:  time.Minute,
		2:  2 * time.Minute,
		3:  4 * time.Minute,
		4:  8 * time.Minute,
		5:  15 * time.Minute,
		20: 15 * time.Minute,
	}
	for failures, want := range cases {
		if got := daemonBackoff(failures); got != want {
			t.Errorf("daemonBackoff(%d) = %v, want %v", failures, got, want)
		}
	}
}
