package goauthredis

import (
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"

	"github.com/go-redis/redis"
)

// RedisClient represents the Redis client
type RedisClient struct {
	client *redis.Client
}

// StartRedis start the connection to redis and return the client
func StartRedis(redisAddr string, pass string) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     redisAddr,
		Password: pass,
		DB:       0,
	})

	pingRes := client.Ping()
	if pingRes.Err() != nil {
		return nil, pingRes.Err()
	}
	log.Println("Connected to Redis server")

	return &RedisClient{client}, nil
}

// CloseConn close the connection to Redis
func CloseConn(rc *RedisClient) error {
	return rc.client.Close()
}

// SaveSession save current session in Redis
func (rc *RedisClient) SaveSession(username string) (string, error) {
	sessionToken := uuid.New().String()
	res := rc.client.SetNX(sessionToken, username, 120*time.Second)
	return sessionToken, res.Err()
}

// RetreiveSession retreive user from token
func (rc *RedisClient) RetreiveSession(token string) (string, error) {
	res := rc.client.Get(token)
	if res.Err() != nil {
		return "", res.Err()
	}
	return res.Val(), nil
}

// DeleteSession delete the session from Redis
func (rc *RedisClient) DeleteSession(token string) (string, error) {
	res := rc.client.Del(token)
	if res.Err() != nil {
		return "", res.Err()
	}
	return res.String(), nil
}

// SetCookie set cookie for client
func (rc *RedisClient) SetCookie(w http.ResponseWriter, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     "go_auth_id",
		Value:    token,
		Expires:  time.Now().Add(180 * time.Second),
		HttpOnly: true,
		SameSite: 1,
	})
}

const redisUserListKey = "go_auth_users"

// SaveUsers save users to Redis
func (rc *RedisClient) SaveUsers(usernameArr *[]string) error {
	delRes := rc.client.Del(redisUserListKey)
	if err := delRes.Err(); err != nil {
		return err
	}

	for i := 0; i < len(*usernameArr); i++ {
		res := rc.client.SAdd(redisUserListKey, (*usernameArr)[i])
		if err := res.Err(); err != nil {
			return err
		}
	}

	expireRes := rc.client.Expire(redisUserListKey, 120*time.Second)
	return expireRes.Err()
}

// InsertUser insert new username into users list
func (rc *RedisClient) InsertUser(username string) error {
	isMemberRes := rc.client.SIsMember(redisUserListKey, username)
	if isMemberRes.Val() == true {
		return errors.New("User is already a member")
	}
	log.Printf("Added %s to Redis user list", username)
	addRes := rc.client.SAdd(redisUserListKey, username)
	return addRes.Err()
}

// RetreiveAllUsers retreive all users from redis user list
func (rc *RedisClient) RetreiveAllUsers() ([]string, error) {
	allUsers := rc.client.SMembers(redisUserListKey)
	if err := allUsers.Err(); err != nil {
		return nil, err
	}

	return allUsers.Val(), nil
}
