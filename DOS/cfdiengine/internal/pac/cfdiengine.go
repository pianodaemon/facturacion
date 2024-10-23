package pac

import (
	"fmt"
)

type CFDIEngine interface {
	DoFact(sourceID string) ([]byte, error)
}

type EngineCode uint8

type CFDIEngineError struct {
	Code    EngineCode
	Message string
}

const (
	SUCCESS EngineCode = iota

	DBIssue      = 250
	StorageIssue = 251
	ReqMalformed = 252
	PACConnIssue = 253
	UnknownIssue = 254
)

// Implement the error interface for CFDIError
func (e *CFDIEngineError) Error() string {
	return fmt.Sprintf("Error %d: %s", e.Code, e.Message)
}
