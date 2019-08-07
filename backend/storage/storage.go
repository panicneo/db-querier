package storage

import (
	"github.com/jmoiron/sqlx"
	"github.com/rs/zerolog"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

type Storage struct {
	db    map[string]*DB
	mutex sync.RWMutex
}

var once sync.Once
var s *Storage

func register(name string, dialect string, dsn string) {
	conn, err := sqlx.Connect(dialect, dsn)
	if err != nil {
		log.Panic().Str("db", name).Err(err).Msg("database connect failed")
	}
	log.Info().Str("name", name).Msg("database connected")
	conn.SetMaxOpenConns(5)
	conn.SetMaxIdleConns(2)
	conn.SetConnMaxLifetime(time.Hour)
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.db[name] = newDB(conn)
}

func Get(name string) *DB {
	once.Do(func() {
		s = &Storage{
			db:    make(map[string]*DB),
			mutex: sync.RWMutex{},
		}
		var dbs []struct {
			Name    string
			Dialect string
			Dsn     string
		}
		_ = viper.UnmarshalKey("databases", &dbs)
		for _, db := range dbs {
			register(db.Name, db.Dialect, db.Dsn)
		}
	})

	s.mutex.RLock()
	defer s.mutex.RUnlock()
	db, ok := s.db[name]
	if !ok {
		return nil
	} else {
		return db
	}
}

func Close(name string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	_ = s.db[name].conn.Close()
	log.Info().Str("name", name).Msg("database closed")

}

func newDB(conn *sqlx.DB) *DB {
	return &DB{
		logger: &log.Logger,
		conn:   conn,
	}
}

type DB struct {
	logger *zerolog.Logger
	conn   *sqlx.DB
}

func (db DB) Query(query string, params ...interface{}) (rows *sqlx.Rows, err error) {
	start := time.Now()
	rows, err = db.conn.Queryx(query, params...)
	logWithError(err).
		Str("sql", query).
		Interface("args", params).
		Str("latency", time.Since(start).String()).
		Msg("SQL")
	return
}

func (db DB) Exec(query string, params ...interface{}) error {
	start := time.Now()
	result, err := db.conn.Exec(query, params...)
	lastInsertedId, _ := result.LastInsertId()
	rowsAffected, _ := result.RowsAffected()
	defer logWithError(err).
		Str("sql", query).
		Interface("args", params).
		Str("latency", time.Since(start).String()).
		Int64("last_inserted_id", lastInsertedId).
		Int64("rows_affected", rowsAffected).
		Msg("SQL")
	return nil
}
func (db DB) Select(dest interface{}, query string, args ...interface{}) (err error) {
	start := time.Now()
	err = db.conn.Select(dest, query, args...)
	defer logWithError(err).
		Str("sql", query).
		Interface("args", args).
		Str("lagency", time.Since(start).String()).
		Msg("SQL")
	return
}

func logWithError(err error) *zerolog.Event {
	if err != nil {
		return log.Error().Err(err)
	}
	return log.Info()
}
