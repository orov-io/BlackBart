package server

import (
	"context"
	"database/sql"
	"net/http"
	"sync"

	"cloud.google.com/go/profiler"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	badger "github.com/dgraph-io/badger/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"google.golang.org/appengine"
)

// expose relevant release modes
const (
	ReleaseMode = gin.ReleaseMode
	DebugMode   = gin.DebugMode
	Production  = "PROD"
)

var onceService sync.Once
var onceDB sync.Once
var onceInternalDB sync.Once
var onceRedis sync.Once
var instance *Service

// Service models the service and upstream needed capabilities
type Service struct {
	options   *Options
	db        *sql.DB
	dbx       *sqlx.DB
	redisPool *redis.Pool
	service   *gin.Engine
	log       *logrus.Logger
	firebase  *firebase.App
	badger    *badger.DB
}

// Init initializes a service if there are no other service initialized
func Init(options *Options) error {
	if instance != nil {
		return ServiceAlreadyInitializeError()
	}

	var err error

	onceService.Do(func() {

		instance, err = newService(options)

	})

	return err
}

func newService(options *Options) (*Service, error) {
	service := &Service{
		options: options,
	}

	service.init()

	return service, nil
}

func (s *Service) init() error {
	if s.options == nil {
		GetLogger().Warn("Warning: Starting new service with no plugins")
		return nil
	}

	var err error

	s.initLogger()
	err = s.initDB()
	if err != nil && !IsNoDatabaseOptionsError(err) {
		GetLogger().WithError(err).Fatal("Can't connect to provided database")
	} else if IsNoDatabaseOptionsError(err) {
		GetLogger().Debug("Database config not provided. Skipping db initialization")
	}

	err = s.initAuth()
	if err != nil && !IsNoFirebaseOptionsError(err) {
		GetLogger().WithError(err).Fatal("Can't connect to firebase auth system")
	} else if IsNoFirebaseOptionsError(err) {
		GetLogger().Debug("Firebase config not provided. Skipping auth initialization")
	}

	err = s.initRouter()
	if err != nil && !IsNoGinOptionsError(err) {
		GetLogger().WithError(err).Fatal("Can't initialize GIN router")
	} else if IsNoGinOptionsError(err) {
		GetLogger().Debug("Gin config not provided. Default router initialized")
	}

	err = s.initRedisPool()
	if err != nil && !IsNoRedisOptionsError(err) {
		GetLogger().WithError(err).Fatal("Can't connect to provided redis server")
	} else if IsNoRedisOptionsError(err) {
		GetLogger().Debug("Redis config not provided. Skipping redis initialization")
	}

	s.initProfiler()

	err = s.initInternalDB()
	GetLogger().Debug("Initializing badger")
	if err != nil && !IsNoInternalDatabaseOptionsError(err) {
		GetLogger().WithError(err).Fatal("Can't initialize the internal DB")
	} else if IsNoInternalDatabaseOptionsError(err) {
		GetLogger().Debug("InternalDB not required. Skipping badger initialization")
	}

	return nil
}

func (s *Service) initRouter() error {
	s.service = gin.New()
	if s.options.gin == nil {
		return NoGinOptionsError()
	}
	s.service.Use(s.options.gin.Middleware...)
	s.service.Use(cors.New(s.options.gin.Cors))
	return nil
}

func (s *Service) initProfiler() {
	if s.options.service.Profiler {
		profiler.Start(profiler.Config{})
	}
}

func (s *Service) initLogger() {
	setLogger(s.options.logger)
	s.log = GetLogger()
}

// GetService returns the service if initialized.
func GetService() (*Service, error) {
	if instance == nil {
		return nil, ServiceNotYetInitializeError()
	}

	return instance, nil
}

// GetDB returns the main sql database connection.
func (s *Service) GetDB() (*sql.DB, error) {
	if s.db == nil {
		return nil, DatabaseNotYetInitializeError()
	}

	return s.db, nil
}

// GetRedisPool returns the main sql database connection.
func (s *Service) GetRedisPool() (*redis.Pool, error) {
	if s.redisPool == nil {
		return nil, NewRedisNotYetInitializedError()
	}

	return s.redisPool, nil
}

// GetDBx returns the sqlx database connection wrapper.
func (s *Service) GetDBx() (*sqlx.DB, error) {
	if s.dbx == nil {
		return nil, DatabaseNotYetInitializeError()
	}

	return s.dbx, nil
}

// IsUsingDB returns if the service has a functional database.
func (s *Service) IsUsingDB() bool {
	return s.db != nil
}

// GetInternalDB returns the badger database
func (s *Service) GetInternalDB() (*badger.DB, error) {
	if s.badger == nil {
		return nil, NewInternalDBNotYetInitializeError()
	}

	return s.badger, nil
}

// GetAuthClient returns a instance of the attached auth client
func (s *Service) GetAuthClient() (*auth.Client, error) {
	if s.firebase == nil {
		return nil, FirebaseNotAlreadyInitializedError()
	}
	return s.firebase.Auth(context.Background())
}

// GetFirebaseApp returns a instance of the attached firebase app
func (s *Service) GetFirebaseApp() (*firebase.App, error) {
	if s.firebase == nil {
		return nil, FirebaseNotAlreadyInitializedError()
	}
	return s.firebase, nil
}

// GetLogger returns the actual logger
func (s *Service) GetLogger() *logrus.Logger {
	return log
}

// StartDefaultService returns an initialized service attached to the Ping handler.
func StartDefaultService() (*Service, error) {
	GetLogger().Infof("Starting service with default handlers/middleware")

	options := NewOptions().WithDefaultOptions()

	return initService(options)
}

// StartService returns an initialized service with that is not attached to
// any Handler.
func StartService(options *Options) (*Service, error) {
	GetLogger().Info("Starting service with no handlers/middleware attached")
	return initService(options)
}

func initService(options *Options) (*Service, error) {
	err := Init(options)
	if err != nil {
		return nil, err
	}
	return GetService()
}

// Run attaches the router to a http.Server and starts listening and serving HTTP requests.
// It is a shortcut for http.ListenAndServe(addr, router)
// Note: this method will block the calling goroutine indefinitely unless an error happens.
func (s *Service) Run(addr ...string) error {
	GetLogger().Infof("Running with gin router")
	return s.service.Run(addr...)
}

// RunAppEngine initializes the main appengine routine and starts listening and serving HTTP requests.
// It performs this opperation by attaching the router to a http.handler
// Note: this method will block the calling goroutine indefinitely unless an error happens.
func (s *Service) RunAppEngine() error {
	GetLogger().Infof("Running on appengine env")
	http.Handle("/", s.service)
	appengine.Main()
	return nil
}

// SetMode sets the router mode: gin.ReleaseMode or gin.DebugMode
func (s *Service) SetMode(mode string) {
	gin.SetMode(mode)
}

// Group declare a new group on the service router
func (s *Service) Group(relativePath string, handlers ...gin.HandlerFunc) *gin.RouterGroup {
	return s.service.Group(relativePath, handlers...)
}

// CloseAll closes all opened database connections
func (s *Service) CloseAll() error {
	if s.dbx != nil {
		if err := s.dbx.Close(); err != nil {
			return err
		}
	}

	if s.redisPool != nil {
		if err := s.redisPool.Close(); err != nil {
			return err
		}
	}

	if s.badger != nil {
		if err := s.db.Close(); err != nil {
			return err
		}
	}

	return nil
}
