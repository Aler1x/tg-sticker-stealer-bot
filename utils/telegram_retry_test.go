package utils

import (
	"fmt"
	"testing"
	"time"

	tg "gopkg.in/telebot.v4"
)

func TestFloodWaitFromString(t *testing.T) {
	wait, ok := FloodWait(fmt.Errorf("telegram: retry after 252 (429)"))
	if !ok {
		t.Fatal("expected flood wait")
	}
	if wait != 252*time.Second {
		t.Fatalf("got %v, want 252s", wait)
	}
}

func TestFloodWaitFromFloodError(t *testing.T) {
	err := tg.FloodError{RetryAfter: 8}
	wait, ok := FloodWait(err)
	if !ok {
		t.Fatal("expected flood wait")
	}
	if wait != 8*time.Second {
		t.Fatalf("got %v, want 8s", wait)
	}
}

func TestFloodWaitNonFlood(t *testing.T) {
	if _, ok := FloodWait(fmt.Errorf("telegram: Bad Request: STICKER_PNG_DIMENSIONS (400)")); ok {
		t.Fatal("did not expect flood wait")
	}
}

func TestCallWithFloodRetrySucceedsAfterWait(t *testing.T) {
	attempts := 0
	err := CallWithFloodRetry(func() error {
		attempts++
		if attempts == 1 {
			return fmt.Errorf("telegram: retry after 1 (429)")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if attempts != 2 {
		t.Fatalf("got %d attempts, want 2", attempts)
	}
}
