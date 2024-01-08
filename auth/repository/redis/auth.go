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

// Login session dont have expired time
func (r *redisAuthRepository) CreateSession(ctx context.Context, user domain.ResponseLoginDTO, token string) (session domain.ResponseLoginDTO, err error) {
	// Use pipeline to execute multiple commands in a single round trip
	pipe := r.Conn.Pipeline()

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
		"Expired": session.ExpiredDatetime.Format(time.RFC3339), // Convert time to string
	}

	// Set session data and expire for email
	log.Printf("HMSET %s: %v", session.Email, sessionData)
	pipe.HMSet(ctx, session.Email, sessionData)
	pipe.Expire(ctx, session.Email, sessionExpire)

	// Set token and expire
	log.Printf("SET %s %s", token, session.Email)
	pipe.Set(ctx, token, session.Email, sessionExpire)
	pipe.Expire(ctx, token, sessionExpire)

	// Execute all commands in the pipeline
	log.Info("Executing Redis pipeline...")
	_, err = pipe.Exec(ctx)
	if err != nil {
		log.Errorf("Error executing Redis pipeline: %v", err)
		return domain.ResponseLoginDTO{}, err
	}

	return
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

	if sessionExpire.Before(time.Now()) {
		err = errors.New("Expired")
		return domain.ResponseLoginDTO{}, err
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
