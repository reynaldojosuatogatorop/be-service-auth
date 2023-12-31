package redis

import (
	"be-service-auth/domain"
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
)

type redisAuthRepository struct {
	Conn *redis.Client
}

// NewRedisAuthRepository is constructor of Redis repository
func NewRedisAuthRepository(Conn *redis.Client) domain.AuthRedisRepo {
	return &redisAuthRepository{Conn}
}

func (r *redisAuthRepository) getToken(ctx context.Context, token string) (session domain.ResponseLoginDTO, err error) {
	log.Debug("Get session token " + token)
	data := r.Conn.Get(ctx, token)
	res, err := data.Result()
	if err != nil {
		return
	}
	if len(res) == 0 {
		err = errors.New("not found")
		return
	}
	err = data.Scan(&session.Email)
	if err != nil {
		return
	}

	session.Email = res
	session.Token = token

	return
}

// GetSession is check if any user already have login session
func (r *redisAuthRepository) GetSession(ctx context.Context, username string) (session domain.ResponseLoginDTO, err error) {
	token, err := r.Conn.HGet(ctx, username, "Token").Result()
	if err != nil {
		return
	}

	session = domain.ResponseLoginDTO{
		Token: token,
	}

	return
}

// Delete data from redis
func (r *redisAuthRepository) DeleteSession(ctx context.Context, token string) (err error) {
	session, err := r.getToken(ctx, token)
	if err != nil {
		return
	}

	res, err := r.Conn.Del(ctx, session.Email).Result()
	if err != nil {
		return
	}
	if res == 0 {
		return errors.New("not found")
	}
	res, err = r.Conn.Del(ctx, token).Result()
	if err != nil {
		return
	}
	if res == 0 {
		return errors.New("not found")
	}

	return
}

// Login session dont have expired time
func (r *redisAuthRepository) CreateSession(ctx context.Context, user domain.ResponseLoginDTO, token string) (session domain.ResponseLoginDTO, err error) {
	// var params []interface{}
	now := time.Now()
	sessionExpire := time.Duration(viper.GetInt("server.session_expire")) * time.Second

	log.Print(viper.GetInt("server.session_expire"))
	session = domain.ResponseLoginDTO{
		Token:           token,
		ID:              user.ID,
		Email:           user.Email,
		Role:            user.Role,
		ExpiredDatetime: now.Add(sessionExpire),
	}

	sessionData := map[string]interface{}{
		"Token":   token,
		"ID":      user.ID,
		"Email":   user.Email,
		"Role":    user.Role,
		"Expired": session.ExpiredDatetime,
	}

	// store data using SET command
	pipe := r.Conn.Pipeline()

	// Simpan data sesi dalam Redis menggunakan HSet
	r.Conn.HMSet(ctx, session.Email, sessionData)
	r.Conn.Expire(ctx, session.Email, sessionExpire)

	// Simpan token ke kunci sesi
	r.Conn.Set(ctx, token, session.Email, sessionExpire)
	r.Conn.Expire(ctx, token, sessionExpire)
	_, err = pipe.Exec(ctx)

	return
}

func (r *redisAuthRepository) GetAuth(ctx context.Context, token string) (response domain.ResponseLoginDTO, err error) {
	session, err := r.getToken(ctx, token)
	if err != nil {
		return
	}

	res, err := r.Conn.HGetAll(ctx, session.Email).Result()

	if err != nil {
		return
	}

	sessionExpire, err := time.Parse(time.RFC3339, res["Expired"])
	if err != nil {
		log.Error("Error convert expired")
		return
	}
	ID, err := strconv.Atoi(res["ID"])
	if err != nil {
		return
	}

	response.ID = ID
	response.Email = res["Email"]
	response.Token = token
	response.Role = res["Role"]
	response.ExpiredDatetime = sessionExpire

	return
}
