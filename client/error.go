// Package error provides the interface for Warranted specific errors.
package client

import (
	"fmt"
)

// WarrantedError provides information about an unsuccessful request.
type WarrantedError struct {
	ErrorCode string `json:"errorCode"`
	Message   string `json:"message"`
}

func (err *WarrantedError) Error() string {
	return fmt.Sprintf("Status: Error - ApiError %s: %s ",
		err.ErrorCode, err.Message)
}
