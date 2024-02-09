package revo

import (
	"testing"
	"time"
)

func TestComputeBackoff(t *testing.T) {
	backoffOne := computeBackoff(0, false)
	if backoffOne != 250*time.Millisecond {
		t.Errorf("Unexpected backoff time %d != %d", backoffOne.Milliseconds(), (250 * time.Millisecond).Milliseconds())
	}
	backoffTwo := computeBackoff(1, false)
	if backoffTwo != 500*time.Millisecond {
		t.Errorf("Unexpected backoff time %d != %d", backoffTwo.Milliseconds(), (500 * time.Millisecond).Milliseconds())
	}
	backoffThree := computeBackoff(2, false)
	if backoffThree != 1000*time.Millisecond {
		t.Errorf("Unexpected backoff time %d != %d", backoffThree.Milliseconds(), (1000 * time.Millisecond).Milliseconds())
	}
	maxBackoff := computeBackoff(10, false)
	if maxBackoff != 2000*time.Millisecond {
		t.Errorf("Unexpected backoff time %d != %d", maxBackoff.Milliseconds(), (2000 * time.Millisecond).Milliseconds())
	}
	overflow := computeBackoff(1000000, false)
	if overflow != 2000*time.Millisecond {
		t.Errorf("Unexpected backoff time %d != %d", overflow.Milliseconds(), (2000 * time.Millisecond).Milliseconds())
	}
}

func TestComputeBackoffWithRandom(t *testing.T) {
	randomRange := time.Duration(250)
	for i := 0; i < 10000; i++ {
		backoff := computeBackoff(0, true)
		min := (250 - randomRange) * time.Millisecond
		max := (250 + randomRange) * time.Millisecond
		if backoff < min || backoff > max {
			t.Fatalf("Unexpected backoff time %d <= (%d) <= %d", min, backoff.Milliseconds(), max)
		}
	}

	overflow := computeBackoff(1000000, true)
	if overflow != 2000*time.Millisecond {
		t.Fatalf("Unexpected backoff time %d != %d", overflow.Milliseconds(), (2000 * time.Millisecond).Milliseconds())
	}
}
