package query

import (
	"db-querier/storage"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

type TParam struct {
	Name string `json:"name" toml:"name"`
	Type string `json:"type" toml:"name"`
}
type TParams []TParam

func (ts TParams) NameTypePairs() map[string]string {
	pairs := make(map[string]string, len(ts))
	for _, param := range ts {
		pairs[param.Name] = param.Type
	}
	return pairs
}

type TQuery struct {
	Database string  `json:"database" toml:"database"`
	Name     string  `json:"name" toml:"name"`
	Params   TParams `json:"params" toml:"params"`
	Sql      string  `json:"sql" toml:"sql"`
}

type mysqlService struct {
	connMap map[string]*storage.DB
	qm      map[string]TQuery
	lock    sync.RWMutex
}

func (s mysqlService) get(name string) *storage.DB {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return s.connMap[name]
}

func (s mysqlService) List() []interface{} {
	return viper.Get("query").([]interface{})
}

func (s mysqlService) Query(name string, params ...interface{}) (*sqlx.Rows, error) {
	q, ok := s.qm[name]
	if !ok {
		return nil, errors.New(fmt.Sprintf("query `%s` is not configured", name))
	}
	if len(params) != len(q.Params) {
		return nil, errors.New(fmt.Sprintf("TParam length not match, expect %d but got %d", len(q.Params), len(params)))
	}
	db := s.get(q.Database)
	if db == nil {
		return nil, errors.New(fmt.Sprintf("database `%s` is not configured", q.Database))
	}

	rows, err := db.Query(s.formatSql(q.Sql), params...)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (s mysqlService) Initial() {}
func (s mysqlService) Close() {
	s.lock.Lock()
	defer s.lock.Unlock()
	for name := range s.connMap {
		storage.Close(name)
	}
}

func (s mysqlService) QueryParams(name string) TParams {
	if q, ok := s.qm[name]; ok {
		return q.Params
	}
	return nil
}

func (s mysqlService) formatSql(sql string) string {
	sql = strings.TrimSpace(sql)
	sql = strings.TrimRight(sql, ";")
	sql += " LIMIT 50;"
	return sql
}

func NewMysqlService() *mysqlService {
	var qs []TQuery
	_ = viper.UnmarshalKey("query", &qs)
	qm := make(map[string]TQuery, len(qs))
	connMap := make(map[string]*storage.DB)
	for _, q := range qs {
		qm[q.Name] = q
		connMap[q.Database] = storage.Get(q.Database)
	}
	return &mysqlService{
		connMap: connMap,
		qm:      qm,
		lock:    sync.RWMutex{},
	}
}
