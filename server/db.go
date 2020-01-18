package server

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/eapache/go-resiliency/retrier"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // Here we initializes the database
	"github.com/pressly/goose"
)

const (
	pluginDatabaseDriver = "postgres"
	migrationDir         = "./migrations"
)

// GetDB returns the main sql database connection.
func GetDB() (*sql.DB, error) {
	service, err := GetService()
	if err != nil {
		GetLogger().
			WithError(err).
			Warn("Can't retrieve database. Does you call microserver.Init()??")
		return nil, err
	}

	db, err := service.GetDB()
	if err != nil {
		GetLogger().
			WithError(err).
			Warn("Can't retrieve database. Does you add a database options??")
	}

	return db, err
}

// GetDBx returns the sqlx database connection wrapper.
func GetDBx() (*sqlx.DB, error) {
	service, err := GetService()
	if err != nil {
		GetLogger().
			WithError(err).
			Warn("Can't retrieve database. Does you call microserver.Init()??")
		return nil, err
	}

	dbx, err := service.GetDBx()
	if err != nil {
		GetLogger().
			WithError(err).
			Warn("Can't retrieve database. Does you add a database options??")
	}

	return dbx, err
}

func (s *Service) initDB() error {
	if !mustInitializeDB(s.options) {
		return NoDatabaseOptionsError()
	}

	if s.db != nil {
		return DatabaseAlreadyInitializeError()
	}

	var err error
	onceDB.Do(func() {
		s.db, s.dbx, err = initializeDBFromOptions(s.options.db)
	})

	return err
}

func mustInitializeDB(options *Options) bool {
	return options.db != nil
}

func initializeDBFromOptions(options *DBOptions) (db *sql.DB, dbx *sqlx.DB, err error) {
	db, err = getDBFromOptions(options)
	if err != nil {
		return
	}
	dbx = sqlx.NewDb(db, pluginDatabaseDriver)
	assertDBExists(dbx, options)
	err = closeDBs(db, dbx)
	if err != nil {
		return
	}

	db, err = connectToDB(options)
	if err != nil {
		return
	}
	dbx = sqlx.NewDb(db, pluginDatabaseDriver)

	err = migrateDB(db)
	return
}

func getDBFromOptions(options *DBOptions) (db *sql.DB, err error) {
	if injectedDB := options.GetInjectedDB(); injectedDB != nil {
		db = injectedDB
		return
	}

	db, err = connectToDBServer(options)
	return
}

func assertDBExists(dbx *sqlx.DB, options *DBOptions) {
	if dbExists(dbx, options) {
		return
	}
	createDB(dbx, options)

	return
}

func dbExists(dbx *sqlx.DB, options *DBOptions) bool {
	query := "SELECT true from pg_database WHERE datname=$1"
	var exists bool
	err := dbx.Get(&exists, query, options.MainDatabase)
	return err == nil && exists
}

func createDB(dbx *sqlx.DB, options *DBOptions) {
	query := fmt.Sprintf("CREATE DATABASE %v;", options.MainDatabase)
	dbx.MustExec(query)
}

func connectToDBServer(options *DBOptions) (db *sql.DB, err error) {
	connectionParams := getServerConnectionString(options)
	ret := retrier.New(retrier.ExponentialBackoff(5, 1*time.Second), retrier.DefaultClassifier{})
	ret.Run(func() error {
		db, err = sql.Open(pluginDatabaseDriver, connectionParams)
		if err != nil {
			return err
		}
		err = db.Ping()
		return err
	})
	if err == nil {
		GetLogger().Info("database connection established")
	}
	return
}

func getServerConnectionString(options *DBOptions) string {
	return fmt.Sprintf(
		"user=%v password=%v sslmode=%v host=%v",
		options.User,
		options.Password,
		options.SSLMode,
		strings.Replace(options.Host, " ", "", -1),
	)
}

func connectToDB(options *DBOptions) (*sql.DB, error) {
	connectionParams := getDBConnectionString(options)
	db, err := sql.Open(pluginDatabaseDriver, connectionParams)
	err = db.Ping()
	if err != nil {
		GetLogger().WithError(err).Warn("Error Opening mainDB")
		return nil, err
	}
	return db, err
}

func getDBConnectionString(options *DBOptions) string {
	return fmt.Sprintf("dbname=%v %v", options.MainDatabase, getServerConnectionString(options))
}

func migrateDB(db *sql.DB) error {
	return goose.Up(db, getMigrationsPath())
}

func getMigrationsPath() string {
	return migrationDir
}

func closeDBs(db *sql.DB, dbx *sqlx.DB) error {
	err := db.Close()
	if err != nil {
		return err
	}

	err = dbx.Close()
	if err != nil {
		return err
	}

	return nil
}
