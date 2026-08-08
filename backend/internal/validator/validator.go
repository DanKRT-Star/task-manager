package validator

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

var Validate = validator.New()

// FormatValidationError chuyển lỗi validator sang message dễ đọc
func FormatValidationError(err error) string {
	var messages []string
	for _, e := range err.(validator.ValidationErrors) {
		messages = append(messages, fmt.Sprintf("%s failed on '%s' rule", e.Field(), e.Tag()))
	}
	return strings.Join(messages, "; ")
}