package utils

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"time"

	tg "gopkg.in/telebot.v4"
)

const (
	maxFloodRetries   = 5
	defaultFloodWait  = 5 * time.Second
	maxFloodWait      = 5 * time.Minute
	stickerUploadPace = 400 * time.Millisecond
)

var retryAfterRe = regexp.MustCompile(`(?i)retry after (\d+)`)

// StickerUploadPace is the delay between successful addStickerToSet calls.
func StickerUploadPace() time.Duration {
	return stickerUploadPace
}

// CallWithFloodRetry runs fn and, on Telegram flood control (429), waits for
// retry_after then retries. Transient network errors use a short backoff.
func CallWithFloodRetry(fn func() error) error {
	var lastErr error

	for attempt := 1; attempt <= maxFloodRetries; attempt++ {
		err := fn()
		if err == nil {
			return nil
		}
		lastErr = err

		wait, ok := FloodWait(err)
		if !ok {
			return err
		}

		if wait > maxFloodWait {
			wait = maxFloodWait
		}
		if wait < time.Second {
			wait = defaultFloodWait
		}

		Logger("warn", "Telegram flood control, waiting before retry", map[string]any{
			"attempt":    attempt,
			"retryAfter": wait.String(),
			"error":      err.Error(),
		})
		time.Sleep(wait)
	}

	return fmt.Errorf("max flood retries exceeded: %w", lastErr)
}

// FloodWait reports whether err is a Telegram flood/429 and how long to wait.
func FloodWait(err error) (time.Duration, bool) {
	if err == nil {
		return 0, false
	}

	var floodErr tg.FloodError
	if errors.As(err, &floodErr) && floodErr.RetryAfter > 0 {
		return time.Duration(floodErr.RetryAfter) * time.Second, true
	}

	// FloodError is a value type and may not unwrap; also match string form.
	if floodErr, ok := err.(tg.FloodError); ok && floodErr.RetryAfter > 0 {
		return time.Duration(floodErr.RetryAfter) * time.Second, true
	}

	matches := retryAfterRe.FindStringSubmatch(err.Error())
	if len(matches) == 2 {
		seconds, parseErr := strconv.Atoi(matches[1])
		if parseErr == nil && seconds > 0 {
			return time.Duration(seconds) * time.Second, true
		}
	}

	return 0, false
}
