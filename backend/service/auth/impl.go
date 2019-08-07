package auth

import (
	"db-querier/storage"
	"db-querier/utils/token"
	"errors"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"time"
)

var (
	ErrUnMatch = errors.New("email or password is not correct")
)

func NewMysqlService() *mysqlService {
	db := storage.Get("auth")
	if db == nil {
		log.Panic().Msg("auth database not configured")
	}
	return &mysqlService{
		conn: db,
	}
}

type mysqlService struct {
	conn *storage.DB
}

func (u mysqlService) Login(email, password string) (obj Account, t string, err error) {
	q := "SELECT * FROM auth WHERE email = ?"
	err = u.conn.Select(&obj, q, email)
	if err != nil {
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(obj.Password), []byte(password))
	if err != nil {
		err = ErrUnMatch
		return
	}
	t, err = token.Create(obj.Email)
	if err != nil {
		return
	}
	q = "UPDATE auth SET last_login_at = ? WHERE email = ?"
	err = u.conn.Exec(q, time.Now(), email)
	return
}

func (u mysqlService) Create(username, password, email string) (obj Account, err error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		return obj, err
	}
	obj = Account{
		Username: username,
		Password: string(hash),
		Email:    email,
	}
	q := "INSERT INTO auth (username, password, email) VALUES (?, ?, ?)"
	err = u.conn.Exec(q, username, password, email)
	return obj, err
}

func (u mysqlService) Close() {
	storage.Close("auth")
}

func (u mysqlService) Initial() {
	SQL :=
		`CREATE TABLE IF NOT EXISTS auth
	(
		id            int(10) unsigned NOT NULL AUTO_INCREMENT,
		username      varchar(255)     NOT NULL,
		password      varchar(255)     NOT NULL,
		email         varchar(50)      NOT NULL,
		created_at    timestamp        NULL DEFAULT now(),
		last_login_at timestamp        NULL DEFAULT NULL,
		PRIMARY KEY (id),
		UNIQUE KEY e (email) USING HASH
	) ENGINE = InnoDB
	  DEFAULT CHARSET = utf8mb4`
	_ = u.conn.Exec(SQL)
}
