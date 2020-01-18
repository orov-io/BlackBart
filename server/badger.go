package server

import (
	"fmt"

	badger "github.com/dgraph-io/badger/v2"
)

// GetInternalDB returns the internal badger DB
func GetInternalDB() (*badger.DB, error) {
	service, err := GetService()
	if err != nil {
		GetLogger().
			WithError(err).
			Warn("Can't retrieve internal database. Does you call server.Init()??")
		return nil, err
	}

	badger, err := service.GetInternalDB()
	if err != nil {
		GetLogger().
			WithError(err).
			Warn("Can't retrieve internal database. Does you add badger options??")
	}

	return badger, err
}

func (s *Service) initInternalDB() error {
	if !mustInitializeInternalDB(s.options) {
		return NoInternalDatabaseOptionsError()
	}

	if s.badger != nil {
		return InternalDatabaseAlreadyInitializeError()
	}

	var err error
	onceDB.Do(func() {
		s.badger, err = badger.Open(
			s.options.internalDB.BadgerOptions,
		)
	})

	if err == nil {
		GetLogger().Info("Badger DB configured")
	}

	return err
}

func mustInitializeInternalDB(options *Options) bool {
	return options.internalDB != nil && options.internalDB.enabled
}

// NoInternalDatabaseOptions is used when user try to get the internal
// database (badger) and it is not initialized
type NoInternalDatabaseOptions struct{}

func (e *NoInternalDatabaseOptions) Error() string {
	return "Can't initialize Database. Database config is not supplied"
}

// NoInternalDatabaseOptionsError returns a new NoInternalDatabaseOptionsError error
func NoInternalDatabaseOptionsError() error {
	return &NoInternalDatabaseOptions{}
}

// IsNoInternalDatabaseOptionsError checks if the error is a NoDatabaseOptions error
func IsNoInternalDatabaseOptionsError(err error) bool {
	_, ok := err.(*NoInternalDatabaseOptions)
	return ok
}

// InternalDatabaseAlreadyInitialize is used when user try to initialize the database and
//it is already initialized
type InternalDatabaseAlreadyInitialize struct{}

func (e *InternalDatabaseAlreadyInitialize) Error() string {
	return "Error initializing service. Database is already initialized"
}

// InternalDatabaseAlreadyInitializeError returns a new InternalDatabaseAlreadyInitialize error
func InternalDatabaseAlreadyInitializeError() error {
	return &InternalDatabaseAlreadyInitialize{}
}

// IsInternalDatabaseAlreadyInitializeError checks if the error is a InternalDatabaseAlreadyInitialize error
func IsInternalDatabaseAlreadyInitializeError(err error) bool {
	_, ok := err.(*InternalDatabaseAlreadyInitialize)
	return ok
}

// InternalDBNotYetInitializeError is used when badger database is not yet initialized.
type InternalDBNotYetInitializeError struct{}

func (e *InternalDBNotYetInitializeError) Error() string {
	return fmt.Sprintf("badger is already initialized")
}

// NewInternalDBNotYetInitializeError returns a new InternalDBNotYetInitializeErrorError error.
func NewInternalDBNotYetInitializeError() error {
	return &InternalDBNotYetInitializeError{}
}

// IsInternalDBNotYetInitializeError checks if the error is a InternalDBNotYetInitializeError error.
func IsInternalDBNotYetInitializeError(err error) bool {
	_, ok := err.(*InternalDBNotYetInitializeError)
	return ok
}
