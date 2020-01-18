package server

import (
	"github.com/gomodule/redigo/redis"
)

const (
	hmset         = "HMSET"
	hGetAll       = "HGETALL"
	redisSet      = "SET"
	redisGet      = "GET"
	ping          = "PING"
	redisProtocol = "tcp"
	defaultPath   = "."
)

// GetRedisPool returns the main redis pool.
func GetRedisPool() (*redis.Pool, error) {
	service, err := GetService()
	if err != nil {
		GetLogger().
			WithError(err).
			Warn("Can't retrieve redisPool. Does you call microserver.Init()??")
		return nil, err
	}

	pool, err := service.GetRedisPool()
	if err != nil {
		GetLogger().
			WithError(err).
			Warn("Can't retrieve redisPool. Does you add a database options??")
	}

	return pool, err
}

// GetRedisConn returns a connection from the redis pool.
func GetRedisConn() (redis.Conn, error) {
	pool, err := GetRedisPool()

	return pool.Get(), err
}

func (s *Service) initRedisPool() error {
	if !mustInitializeRedis(s.options) {
		return NewNoRedisOptionsError()
	}

	if s.redisPool != nil {
		return NewRedisPoolAlreadyInitializedError()
	}

	var err error
	onceRedis.Do(func() {
		s.redisPool, err = initializeRedisPoolFromOptions(s.options.redis)
	})

	return err
}

func mustInitializeRedis(options *Options) bool {
	return options.redis != nil
}

func initializeRedisPoolFromOptions(options *RedisOptions) (pool *redis.Pool, err error) {
	pool, err = getRedisPoolFromOptions(options)
	if err != nil {
		return
	}

	return
}

func getRedisPoolFromOptions(options *RedisOptions) (pool *redis.Pool, err error) {
	if injectedPool := options.GetInjectedPool(); injectedPool != nil {
		pool = injectedPool
		return
	}

	pool, err = connectToRedisServer(options)
	return
}

func connectToRedisServer(options *RedisOptions) (pool *redis.Pool, err error) {
	pool = &redis.Pool{
		Dial: getDialToRedis(options),
	}

	pingConn, err := GetRedisConn()
	if err != nil {
		return nil, err
	}
	defer pingConn.Close()
	pingToRedis(pingConn)
	return
}

func getDialToRedis(options *RedisOptions) func() (redis.Conn, error) {
	return func() (redis.Conn, error) {
		conn, err := redis.Dial(redisProtocol, options.Address)
		if options.Password == "" {
			return conn, err
		}
		if err != nil {
			return nil, err
		}
		if _, err := conn.Do("AUTH", options.Password); err != nil {
			conn.Close()
			return nil, err
		}
		return conn, nil
	}
}

func pingToRedis(conn redis.Conn) {
	_, err := redis.String(conn.Do(ping))
	if err != nil {
		GetLogger().Fatalf("Can't stablish connection to redis server: %v", err)
	}
	GetLogger().Info("Redis connection ready")
}
