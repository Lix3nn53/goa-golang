package error

import (
	"database/sql"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

// Generic Errors
var (
	// ErrNotFound error will be returned when not found a model
	ErrNotFound = errors.New("Not found")
	// InvalidPaymentMethod error will be returned when payment method is not available
	InvalidPaymentMethod = errors.New("Invalid Payment Method")
)

//ParseError Method to parse all the methods to specific HTTP Status
func ParseError(err error) int {
	switch {
	case errors.Is(err, sql.ErrNoRows):
		return http.StatusNotFound
	case strings.Contains(err.Error(), "Validate"):
		return http.StatusBadRequest
	case errors.Is(err, ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, InvalidPaymentMethod):
		return http.StatusServiceUnavailable
	}

	return http.StatusInternalServerError
}
