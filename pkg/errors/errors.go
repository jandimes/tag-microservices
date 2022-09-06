package errors

import (
	"errors"
	"strings"
)

// ErrorWithContext is the error with other data context
type ErrorWithContext struct {
	Msg string
	Err error
}

// Error return the error in string
func (e *ErrorWithContext) Error() string {
	return e.Msg
}

// Wrap return the wrapped error
func (e *ErrorWithContext) Wrap() error {
	if e.Err == nil {
		return errors.New(e.Msg)
	}
	return e.Err
}

// QueryError is the error related to querying the database
type QueryError struct {
	ErrorWithContext
	DBName string
}

// NewQueryError is the constructor for QueryError
func NewQueryError(dbName string, message string, err error) *QueryError {
	return &QueryError{
		DBName:           dbName,
		ErrorWithContext: ErrorWithContext{Msg: message, Err: err},
	}
}

// Error return the error in string
func (e *QueryError) Error() string {
	msgs := []string{e.DBName + ": An error occurred in the database. Please try again."}
	if strings.TrimSpace(e.Msg) != "" {
		msgs = append(msgs, e.Msg)
	}

	return strings.Join(msgs, " ")
}
