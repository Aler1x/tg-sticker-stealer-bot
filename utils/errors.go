package utils

// BotError represents a custom error with i18n support
type BotError struct {
	Message   string
	I18nKey   string
	ErrorCode string
}

func (e *BotError) Error() string {
	return e.Message
}

// NewBotError creates a new BotError and logs it
func NewBotError(message, i18nKey, errorCode string) *BotError {
	Logger("error", message, map[string]any{"errorCode": errorCode})
	return &BotError{
		Message:   message,
		I18nKey:   i18nKey,
		ErrorCode: errorCode,
	}
}

// FailFast panics if error is not nil
func FailFast(err error) {
	if err != nil {
		Logger("error", "Error", map[string]any{"error": err})
		panic(err)
	}
}
