package server

import (
	"context"
	"database/sql"
	"os"
	"strconv"

	badger "github.com/dgraph-io/badger/v2"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	"github.com/sirupsen/logrus"
)

const envKey = "ENV"
const badgerFlagKey = "ENABLE_BADGER"

// Options store all service configuration options
type Options struct {
	db         *DBOptions
	redis      *RedisOptions
	logger     *LoggerOptions
	firebase   *FirebaseOptions
	gin        *GinOptions
	service    *ServiceOptions
	internalDB *InternalDBOptions

	Context context.Context
}

// NewOptions returns an empty options object.
func NewOptions() *Options {
	return &Options{}
}

// DB sets the service database configuration
func (o *Options) DB(dbOptions *DBOptions) {
	GetLogger().Trace("DB options on DB() :")
	o.db = dbOptions
}

// Redis sets the service redis database configuration
func (o *Options) Redis(redisOptions *RedisOptions) {
	o.redis = redisOptions
}

// Logger sets the service logger configuration
func (o *Options) Logger(loggerOptions *LoggerOptions) {
	o.logger = loggerOptions
}

// Gin sets the service logger configuration
func (o *Options) Gin(ginOptions *GinOptions) {
	o.gin = ginOptions
}

// Firebase sets the service firebase app configuration
func (o *Options) Firebase(firebaseOptions *FirebaseOptions) {
	o.firebase = firebaseOptions
}

// Service sets the service options configuration
func (o *Options) Service(serviceOptions *ServiceOptions) {
	o.service = serviceOptions
}

// InternalDB sets the service options configuration
func (o *Options) InternalDB(internalDBOptions *InternalDBOptions) {
	o.internalDB = internalDBOptions
}

// WithDefaultOptions attach default configuration to the options struct and returns
// a pointer with the default configuration options.
func (o *Options) WithDefaultOptions() *Options {
	o.Logger(DefaultLoggerOptions())
	o.DB(DefaultDBOptions())
	o.Redis(DefaultRedisOptions())
	o.Context = context.Background()
	o.Gin(DefaultGinOptions())
	o.Firebase(DefaultFirebaseOptions())
	o.Service(DefaultServiceOptions())
	o.InternalDB(DefaultInternalDBOptions())

	return o
}

// DBOptions stores database configuration.
type DBOptions struct {
	MigrationDir string
	Host         string
	User         string
	SSLMode      string
	MainDatabase string
	Password     string

	db *sql.DB
}

const (
	migrationDirKey     = "DATABASE_MIGRATIONS_DIR" //carefull: Will be relative from calling file.
	databaseHostKey     = "DATABASE_HOST"
	databasePasswordKey = "DATABASE_PASSWORD"
	databaseUserKey     = "DATABASE_USER"
	databaseSSLModeKey  = "DATABASE_SSL_MODE"
	mainDatabaseKey     = "SERVICE_DATABASE_NAME"
)

// NewDBOptions returns a pointer to a new empty DBOptions struct
func NewDBOptions() *DBOptions {
	return &DBOptions{}
}

// DefaultDBOptions initializes and returns a DBOptions struct with default values.
// This values must be environment variables
func DefaultDBOptions() *DBOptions {

	if !databaseEnvIsSetting() {
		return nil
	}

	return &DBOptions{
		MigrationDir: os.Getenv(migrationDirKey),
		Host:         os.Getenv(databaseHostKey),
		User:         os.Getenv(databaseUserKey),
		SSLMode:      os.Getenv(databaseSSLModeKey),
		MainDatabase: os.Getenv(mainDatabaseKey),
		Password:     os.Getenv(databasePasswordKey),
	}
}

func databaseEnvIsSetting() bool {
	return envExist(migrationDirKey) &&
		envExist(databaseHostKey) &&
		envExist(databaseUserKey) &&
		envExist(databaseSSLModeKey) &&
		envExist(mainDatabaseKey) &&
		envExist(databasePasswordKey)
}

// WithInjectedDB can be used to database dependency injection
func (dbo *DBOptions) WithInjectedDB(db *sql.DB) *DBOptions {
	dbo.db = db
	return dbo
}

// GetInjectedDB returns the injected database
func (dbo *DBOptions) GetInjectedDB() *sql.DB {
	return dbo.db
}

// RedisOptions stores redis database configuration.
type RedisOptions struct {
	Address  string
	Password string

	injectedPool *redis.Pool
}

const (
	redisAddressKey  = "REDIS_ADDRESS"
	redisPasswordKey = "REDIS_PASSWORD"
)

// DefaultRedisOptions initializes and returns a DBOptions struct with default values.
// This values must be environment variables
func DefaultRedisOptions() *RedisOptions {
	if !redisEnvIsSetting() {
		return nil
	}

	return &RedisOptions{
		Address:  os.Getenv(redisAddressKey),
		Password: os.Getenv(redisPasswordKey),
	}
}

func redisEnvIsSetting() bool {
	return envExist(redisAddressKey) && envExist(redisPasswordKey)
}

// WithInjectedPool can be used to redis pool dependency injection
func (redisOptions *RedisOptions) WithInjectedPool(pool *redis.Pool) *RedisOptions {
	redisOptions.injectedPool = pool
	return redisOptions
}

// GetInjectedPool returns the injected redis pool
func (redisOptions *RedisOptions) GetInjectedPool() *redis.Pool {
	return redisOptions.injectedPool
}

// LoggerOptions stores logrus initialization options
type LoggerOptions struct {
	Env    string
	Level  logrus.Level
	Format logrus.Formatter
}

// DefaultLoggerOptions returns a LoggerOptions filled based on the environment.
func DefaultLoggerOptions() *LoggerOptions {
	env := os.Getenv(envKey)
	level, format := getDefaultLoggerConfByEnv(env)
	return &LoggerOptions{
		Env:    env,
		Level:  level,
		Format: format,
	}
}

const (
	firebaseBucketKey         = "FIREBASE_BUCKET"
	firebaseBucketFileNameKey = "FIREBASE_BUCKET_FILE_NAME"
	firebaseConfigPathKey     = "FIREBASE_CONFIG_PATH"
)

// FirebaseOptions stores firebase configurations to authorization/authentication
type FirebaseOptions struct {
	bucket     string
	name       string
	configPath string
}

// NewFirebaseOptions returns an empty FirebaseOptionsObject
func NewFirebaseOptions() *FirebaseOptions {
	return &FirebaseOptions{}
}

// FromFile fills the FirebaseObject with the path to the firebase.json config file.
// You can rename the file to other.json or whatever you want, only be sure that
// the path include the file name.
func (fbo *FirebaseOptions) FromFile(path string) *FirebaseOptions {
	fbo.configPath = path
	return fbo
}

// FromBucket fills the FirebaseObject with the gcloud storage bucket that stores
// the firebase.json. If you save the file with other name, you can provide it
// in the name parameter. In any other case, left it to blank string: ""
func (fbo *FirebaseOptions) FromBucket(bucket string, name string) *FirebaseOptions {
	fbo.bucket = bucket
	fbo.name = name
	return fbo
}

// DefaultFirebaseOptions returns a FirebaseOptions fills with the bucket/path
// found on environment variables.
func DefaultFirebaseOptions() *FirebaseOptions {
	switch true {
	case envExist(firebaseBucketKey):
		bucket := os.Getenv(firebaseBucketKey)
		name := os.Getenv(firebaseBucketFileNameKey)
		return NewFirebaseOptions().FromBucket(bucket, name)
	case envExist(firebaseConfigPathKey):
		path := os.Getenv(firebaseConfigPathKey)
		return NewFirebaseOptions().FromFile(path)

	default:
		return NewFirebaseOptions()
	}
}

// GinOptions stores gin router configuration
type GinOptions struct {
	Middleware []gin.HandlerFunc
	Cors       cors.Config
}

// NewGinOptions returns an empty GinOptions struct
func NewGinOptions() *GinOptions {
	return &GinOptions{}
}

// DefaultGinOptions returns a GinOptions fills with the default middleware and
// CORS
func DefaultGinOptions() *GinOptions {
	return &GinOptions{
		Middleware: getDefaultMiddleware(),
		Cors:       getDefaultCors(),
	}
}

func getDefaultMiddleware() []gin.HandlerFunc {
	return []gin.HandlerFunc{
		gin.Logger(),
		gin.Recovery(),
	}
}

func getDefaultCors() cors.Config {
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "PUT", "POST", "HEAD", "DELETE", "PATCH", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "authorization"}
	return config
}

const (
	serviceNameKey    = "SERVICE_NAME"
	serviceVersionKey = "SERVICE_VERSION"
	servicePathKey    = "SERVICE_BASE_PATH"
)

// ServiceOptions stores global service info params.
type ServiceOptions struct {
	Name     string
	Version  string
	Path     string
	Profiler bool
}

// NewServiceOptions returns an empty ServiceOptions struct
func NewServiceOptions() *ServiceOptions {
	return &ServiceOptions{}
}

// DefaultServiceOptions returns a ServiceOptions fills with Name and version
// founds on the  env variables.
func DefaultServiceOptions() *ServiceOptions {
	env := os.Getenv(envKey)
	profiler := false
	if env != "local" {
		profiler = true
	}
	return &ServiceOptions{
		Name:     GetEnvOrDefaultString(serviceNameKey),
		Version:  GetEnvOrDefaultString(serviceVersionKey),
		Path:     GetEnvOrDefaultString(servicePathKey),
		Profiler: profiler,
	}
}

// InternalDBOptions stores options to badger, the provided internal database.
type InternalDBOptions struct {
	BadgerOptions badger.Options
	enabled       bool
}

// DefaultInternalDBOptions returns a InternalDBOptions fills with Name and version
// founds on the  env variables.
func DefaultInternalDBOptions() *InternalDBOptions {
	// Don't care about error. If is not nil, we want to set enable to false.
	enable, _ := strconv.ParseBool(os.Getenv(badgerFlagKey))
	opt := badger.DefaultOptions("").WithInMemory(true)

	return &InternalDBOptions{
		BadgerOptions: opt,
		enabled:       enable,
	}
}
