package highwinds

import (
	"fmt"
	"strings"
)

const (
	// ErrScopeIsNilStr is an error identifying that the config returned from the
	// remote host is empty or missing required portions
	ErrScopeIsNilStr = "Scope is nil on %s/%s/%d"
)

// ErrScopeIsNil is testing a method of having easier to manage dynamic errors
func ErrScopeIsNil(accountHash string, hostHash string, ScopeID int) error {
	return fmt.Errorf(ErrScopeIsNilStr, accountHash, hostHash, ScopeID)
}

// ErrSetState is an error response for failures during state setting
func ErrSetState(errs []error) error {
	if len(errs) > 0 {
		var stringedErrors []string
		for _, err := range errs {
			stringedErrors = append(stringedErrors, err.Error())
		}
		return fmt.Errorf("Error setting state: %v", strings.Join(stringedErrors, " / "))
	}
	return nil
}

// ErrSetModel is an error response for failures during model creation
func ErrSetModel(errs []error) error {
	if len(errs) > 0 {
		var stringedErrors []string
		for _, err := range errs {
			stringedErrors = append(stringedErrors, err.Error())
		}
		return fmt.Errorf("Error building models: %v", strings.Join(stringedErrors, " / "))
	}
	return nil
}
