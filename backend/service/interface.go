package service

import (
	"db-querier/service/auth"
	"db-querier/service/query"

	"github.com/jmoiron/sqlx"
)

type BaseService interface {
	Initial()
	Close()
}

type QueryService interface {
	BaseService
	Query(name string, params ...interface{}) (*sqlx.Rows, error)
	QueryParams(name string) query.TParams
	List() []interface{}
}

type AuthService interface {
	BaseService
	Create(username, password, email string) (auth.Account, error)
	Login(email, password string) (user auth.Account, t string, err error)
}

var (
	Q QueryService
	// Auth  AuthService
)

func Initialize() {
	Q = query.NewMysqlService()
	Q.Initial()
	// Auth = auth.NewMysqlService()
	// Auth.Initial()
}

func Close() {
	Q.Close()
	// Auth.Close()
}
