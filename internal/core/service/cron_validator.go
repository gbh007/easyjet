package service

import (
	"fmt"

	"github.com/robfig/cron/v3"
)

// validateCronExpression validates a cron expression using robfig's parser.
// Returns nil if the expression is valid or empty, otherwise returns an error.
// FIXME(ai-shit): перейти на playground validator.
func validateCronExpression(expr string) error {
	// Empty string is valid (means no schedule)
	if expr == "" {
		return nil
	}

	// Use standard cron parser to validate
	_, err := cron.ParseStandard(expr)
	if err != nil {
		return fmt.Errorf("invalid cron expression: %w", err)
	}

	return nil
}
