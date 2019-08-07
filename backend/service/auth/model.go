package auth

import (
	"github.com/go-sql-driver/mysql"
	"time"
)

type Account struct {
	ID          int64          `json:"id"`
	Username    string         `json:"username"`
	Password    string         `json:"password"`
	Email       string         `json:"email"`
	CreatedAt   time.Time      `json:"created_at"`
	LastLoginAt mysql.NullTime `json:"last_login_at"`
}
