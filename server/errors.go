package server

import "fmt"

// ServiceNotYetInitialize is used when user try to get the service and it is
// not initialized
type ServiceNotYetInitialize struct{}

func (e *ServiceNotYetInitialize) Error() string {
	return "Error getting service. Service is not yet initialized"
}

// ServiceNotYetInitializeError returns a new ServiceNotYetInitialize error
func ServiceNotYetInitializeError() error {
	return &ServiceNotYetInitialize{}
}

// ServiceAlreadyInitialize is used when user try to initialize the service and it is
// already initialized
type ServiceAlreadyInitialize struct{}

func (e *ServiceAlreadyInitialize) Error() string {
	return "Error initializing service. Service is already initialized"
}

// ServiceAlreadyInitializeError returns a new ServiceAlreadyInitialize error
func ServiceAlreadyInitializeError() error {
	return &ServiceAlreadyInitialize{}
}

// DatabaseNotYetInitialize is used when user try to get the database and it is
// not initialized
type DatabaseNotYetInitialize struct{}

func (e *DatabaseNotYetInitialize) Error() string {
	return "Error getting database. Database is not yet initialized"
}

// DatabaseNotYetInitializeError returns a new DatabaseNotYetInitialize error
func DatabaseNotYetInitializeError() error {
	return &DatabaseNotYetInitialize{}
}

// DatabaseAlreadyInitialize is used when user try to initialize the database and
//it is already initialized
type DatabaseAlreadyInitialize struct{}

func (e *DatabaseAlreadyInitialize) Error() string {
	return "Error initializing service. Database is already initialized"
}

// DatabaseAlreadyInitializeError returns a new DatabaseAlreadyInitialize error
func DatabaseAlreadyInitializeError() error {
	return &DatabaseAlreadyInitialize{}
}

// IsDatabaseAlreadyInitializeError checks if the error is a DatabaseAlreadyInitialize error
func IsDatabaseAlreadyInitializeError(err error) bool {
	_, ok := err.(*DatabaseAlreadyInitialize)
	return ok
}

// NoDatabaseOptions is used when user try to get the database and it is
// not initialized
type NoDatabaseOptions struct{}

func (e *NoDatabaseOptions) Error() string {
	return "Can't initialize Database. Database config is not supplied"
}

// NoDatabaseOptionsError returns a new NoDatabaseOptionsError error
func NoDatabaseOptionsError() error {
	return &NoDatabaseOptions{}
}

// IsNoDatabaseOptionsError checks if the error is a NoDatabaseOptions error
func IsNoDatabaseOptionsError(err error) bool {
	_, ok := err.(*NoDatabaseOptions)
	return ok
}

// FirebaseNotAlreadyInitialized is used when user try to get the firebase app and it is
// not initialized
type FirebaseNotAlreadyInitialized struct{}

func (e *FirebaseNotAlreadyInitialized) Error() string {
	return "Can't initialize Firebase. Are firebase config supplied?"
}

// FirebaseNotAlreadyInitializedError returns a new FirebaseNotAlreadyInitializedError error
func FirebaseNotAlreadyInitializedError() error {
	return &FirebaseNotAlreadyInitialized{}
}

// IsFirebaseNotAlreadyInitializedError checks if the error is a FirebaseNotAlreadyInitializedError error
func IsFirebaseNotAlreadyInitializedError(err error) bool {
	_, ok := err.(*FirebaseNotAlreadyInitialized)
	return ok
}

// NoFirebaseOptions is used when user try to get the Firebase and it is
// not initialized
type NoFirebaseOptions struct{}

func (e *NoFirebaseOptions) Error() string {
	return "Can't initialize Firebase. Firebase config is not supplied"
}

// NoFirebaseOptionsError returns a new NoFirebaseOptionsError error
func NoFirebaseOptionsError() error {
	return &NoFirebaseOptions{}
}

// IsNoFirebaseOptionsError checks if the error is a NoFirebaseOptions error
func IsNoFirebaseOptionsError(err error) bool {
	_, ok := err.(*NoFirebaseOptions)
	return ok
}

// NoGinOptions is used when user try to get the Gin Options and there is no options.
type NoGinOptions struct{}

func (e *NoGinOptions) Error() string {
	return "Can't initialize Gin. Gin config is not supplied"
}

// NoGinOptionsError returns a new NoGinOptionsError error
func NoGinOptionsError() error {
	return &NoGinOptions{}
}

// IsNoGinOptionsError checks if the error is a NoGinOptions error
func IsNoGinOptionsError(err error) bool {
	_, ok := err.(*NoGinOptions)
	return ok
}

// RedisNotYetInitializedError is used when redis is not yet initialized
type RedisNotYetInitializedError struct{}

func (e *RedisNotYetInitializedError) Error() string {
	return fmt.Sprintf("Redis is not yet initialized")
}

// NewRedisNotYetInitializedError returns a new RedisNotYetInitializedErrorError error
func NewRedisNotYetInitializedError() error {
	return &RedisNotYetInitializedError{}
}

// IsRedisNotYetInitializedError checks if the error is a RedisNotYetInitializedError error
func IsRedisNotYetInitializedError(err error) bool {
	_, ok := err.(*RedisNotYetInitializedError)
	return ok
}

// NoRedisOptionsError is used when redis configuration is not supplied.
type NoRedisOptionsError struct{}

func (e *NoRedisOptionsError) Error() string {
	return fmt.Sprintf("No redis options provided")
}

// NewNoRedisOptionsError returns a new NoRedisOptionsErrorError error.
func NewNoRedisOptionsError() error {
	return &NoRedisOptionsError{}
}

// IsNoRedisOptionsError checks if the error is a NoRedisOptionsError error.
func IsNoRedisOptionsError(err error) bool {
	_, ok := err.(*NoRedisOptionsError)
	return ok
}

// RedisPoolAlreadyInitializedError is used when redis pool is already initialized.
type RedisPoolAlreadyInitializedError struct{}

func (e *RedisPoolAlreadyInitializedError) Error() string {
	return fmt.Sprintf("Redis is already initialized")
}

// NewRedisPoolAlreadyInitializedError returns a new RedisPoolAlreadyInitializedErrorError error.
func NewRedisPoolAlreadyInitializedError() error {
	return &RedisPoolAlreadyInitializedError{}
}

// IsRedisPoolAlreadyInitializedError checks if the error is a RedisPoolAlreadyInitializedError error.
func IsRedisPoolAlreadyInitializedError(err error) bool {
	_, ok := err.(*RedisPoolAlreadyInitializedError)
	return ok
}
