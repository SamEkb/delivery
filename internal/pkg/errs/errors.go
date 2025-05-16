package errs

import (
	"errors"
	"fmt"
)

var (
	ErrValueIsRequired = errors.New("value is required")
	ErrNotFound        = errors.New("not found")
	ErrValidation      = errors.New("validation error")
	ErrDatabase        = errors.New("database error")
	ErrBusiness        = errors.New("business logic error")
	ErrConflict        = errors.New("conflict error")
)

type ValueIsRequiredError struct {
	ParamName string
	Cause     error
}

func NewValueIsRequiredError(paramName string) *ValueIsRequiredError {
	return &ValueIsRequiredError{
		ParamName: paramName,
	}
}

func NewValueIsRequiredErrorWithCause(paramName string, cause error) *ValueIsRequiredError {
	return &ValueIsRequiredError{
		ParamName: paramName,
		Cause:     cause,
	}
}

func (e *ValueIsRequiredError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (cause: %v)", ErrValueIsRequired, e.ParamName, e.Cause)
	}
	return fmt.Sprintf("%s: %s", ErrValueIsRequired, e.ParamName)
}

func (e *ValueIsRequiredError) Unwrap() error {
	return ErrValueIsRequired
}

type NotFoundError struct {
	Resource string
	ID       string
	Cause    error
}

func NewNotFoundError(resource, id string) *NotFoundError {
	return &NotFoundError{
		Resource: resource,
		ID:       id,
	}
}

func NewNotFoundErrorWithCause(resource, id string, cause error) *NotFoundError {
	return &NotFoundError{
		Resource: resource,
		ID:       id,
		Cause:    cause,
	}
}

func (e *NotFoundError) Error() string {
	if e.ID != "" {
		if e.Cause != nil {
			return fmt.Sprintf("%s: %s with id '%s' (cause: %v)", ErrNotFound, e.Resource, e.ID, e.Cause)
		}
		return fmt.Sprintf("%s: %s with id '%s'", ErrNotFound, e.Resource, e.ID)
	}
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (cause: %v)", ErrNotFound, e.Resource, e.Cause)
	}
	return fmt.Sprintf("%s: %s", ErrNotFound, e.Resource)
}

func (e *NotFoundError) Unwrap() error {
	return ErrNotFound
}

type ValidationError struct {
	Field   string
	Value   interface{}
	Message string
	Cause   error
}

func NewValidationError(field, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}

func NewValidationErrorWithValue(field string, value interface{}, message string) *ValidationError {
	return &ValidationError{
		Field:   field,
		Value:   value,
		Message: message,
	}
}

func NewValidationErrorWithCause(field, message string, cause error) *ValidationError {
	return &ValidationError{
		Field:   field,
		Message: message,
		Cause:   cause,
	}
}

func (e *ValidationError) Error() string {
	if e.Value != nil {
		if e.Cause != nil {
			return fmt.Sprintf("%s: field '%s' with value '%v': %s (cause: %v)", ErrValidation, e.Field, e.Value, e.Message, e.Cause)
		}
		return fmt.Sprintf("%s: field '%s' with value '%v': %s", ErrValidation, e.Field, e.Value, e.Message)
	}
	if e.Cause != nil {
		return fmt.Sprintf("%s: field '%s': %s (cause: %v)", ErrValidation, e.Field, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: field '%s': %s", ErrValidation, e.Field, e.Message)
}

func (e *ValidationError) Unwrap() error {
	return ErrValidation
}

type DatabaseError struct {
	Operation string
	Entity    string
	Cause     error
}

func NewDatabaseError(operation, entity string, cause error) *DatabaseError {
	return &DatabaseError{
		Operation: operation,
		Entity:    entity,
		Cause:     cause,
	}
}

func (e *DatabaseError) Error() string {
	if e.Entity != "" {
		return fmt.Sprintf("%s: failed to %s %s (cause: %v)", ErrDatabase, e.Operation, e.Entity, e.Cause)
	}
	return fmt.Sprintf("%s: failed to %s (cause: %v)", ErrDatabase, e.Operation, e.Cause)
}

func (e *DatabaseError) Unwrap() error {
	return ErrDatabase
}

type BusinessError struct {
	Operation string
	Reason    string
	Cause     error
}

func NewBusinessError(operation, reason string) *BusinessError {
	return &BusinessError{
		Operation: operation,
		Reason:    reason,
	}
}

func NewBusinessErrorWithCause(operation, reason string, cause error) *BusinessError {
	return &BusinessError{
		Operation: operation,
		Reason:    reason,
		Cause:     cause,
	}
}

func (e *BusinessError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s - %s (cause: %v)", ErrBusiness, e.Operation, e.Reason, e.Cause)
	}
	return fmt.Sprintf("%s: %s - %s", ErrBusiness, e.Operation, e.Reason)
}

func (e *BusinessError) Unwrap() error {
	return ErrBusiness
}

type ConflictError struct {
	Resource string
	ID       string
	Reason   string
	Cause    error
}

func NewConflictError(resource, id, reason string) *ConflictError {
	return &ConflictError{
		Resource: resource,
		ID:       id,
		Reason:   reason,
	}
}

func NewConflictErrorWithCause(resource, id, reason string, cause error) *ConflictError {
	return &ConflictError{
		Resource: resource,
		ID:       id,
		Reason:   reason,
		Cause:    cause,
	}
}

func (e *ConflictError) Error() string {
	if e.ID != "" {
		if e.Cause != nil {
			return fmt.Sprintf("%s: %s with id '%s' - %s (cause: %v)", ErrConflict, e.Resource, e.ID, e.Reason, e.Cause)
		}
		return fmt.Sprintf("%s: %s with id '%s' - %s", ErrConflict, e.Resource, e.ID, e.Reason)
	}
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s - %s (cause: %v)", ErrConflict, e.Resource, e.Reason, e.Cause)
	}
	return fmt.Sprintf("%s: %s - %s", ErrConflict, e.Resource, e.Reason)
}

func (e *ConflictError) Unwrap() error {
	return ErrConflict
}

func IsValueRequired(err error) bool {
	return errors.Is(err, ErrValueIsRequired)
}

func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

func IsValidation(err error) bool {
	return errors.Is(err, ErrValidation)
}

func IsDatabase(err error) bool {
	return errors.Is(err, ErrDatabase)
}

func IsBusiness(err error) bool {
	return errors.Is(err, ErrBusiness)
}

func IsConflict(err error) bool {
	return errors.Is(err, ErrConflict)
}
