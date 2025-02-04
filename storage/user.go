package storage

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type (
	UserData struct {
		Id        int    `db:"id"`
		UserName  string `db:"user_name"`
		Email     string `db:"email"`
		Phone     string `db:"phone"`
		Gender    string `db:"gender"`
		CreatedAt int64  `db:"created_at"`
		Password  string `db:"password"`
	}
	UserDataPending struct {
		UserName string `db:"user_name"`
		Email    string `db:"email"`
		Password string `db:"password"`
	}
	UserStorage struct {
		mySql       *sqlx.DB
		redisClient *redis.Client
	}
)

func NewUserStorage(mySql *sqlx.DB, redisClient *redis.Client) *UserStorage {
	return &UserStorage{
		mySql:       mySql,
		redisClient: redisClient,
	}
}

func (u *UserStorage) CreateUser(userDataPending UserDataPending) error {
	ctx := context.Background()
	u.mySql.MustExec("DELETE FROM t_user_unverified WHERE email = ?", userDataPending.Email)
	u.redisClient.Del(ctx, "ver-token:"+userDataPending.Email)
	_, err := u.mySql.NamedExec("INSERT INTO t_user (user_name, email, password) values (:user_name, :email, :password)", userDataPending)
	return err
}

func (u *UserStorage) FindUserByEmail(email string) *UserData {
	user := &UserData{}
	err := u.mySql.Get(user, "SELECT * FROM t_user WHERE email = ?", email)
	if err != nil {
		user = nil
	}
	return user
}

func (u *UserStorage) FindUnverifiedUserByEmail(email string) *UserDataPending {
	user := &UserDataPending{}
	err := u.mySql.Get(user, "SELECT user_name, email, password FROM t_user_unverified WHERE email = ? and NOW() <= created_at + INTERVAL 5 MINUTE", email)
	if err != nil {
		user = nil
	}
	return user
}

func (u *UserStorage) CreateUnverifiedUser(userDataPending UserDataPending) error {
	_, err := u.mySql.NamedExec(`
		INSERT INTO t_user_unverified (user_name, email, password)
		VALUES (:user_name, :email, :password)
		ON DUPLICATE KEY UPDATE user_name = :user_name, password = :password, created_at = CURRENT_TIME
	`, userDataPending)
	return err
}

func (u *UserStorage) CreateVerificationToken(email string, token string) error {
	ctx := context.Background()
	return u.redisClient.Set(ctx, "ver-token:"+token, email, time.Minute*5).Err()
}

func (u *UserStorage) FindUnverifiedEmailByToken(token string) (string, error) {
	ctx := context.Background()
	return u.redisClient.Get(ctx, "ver-token:"+token).Result()
}
