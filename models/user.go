package models

import (
	"errors"
	"fmt"

	"github.com/go-redis/redis"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidLogin = errors.New("invalid login")
)

type User struct {
	key string
}

func NewUser(username string, hash []byte) (*User, error) {
	id, err := client.Incr("user:next-id").Result()
	if err != nil {
		return nil, err
	}
	key := fmt.Sprintf("user:%d", id)
	pipe := client.Pipeline() //set a bunch of commands to redis in one round-trip
	pipe.HSet(key, "id", id)
	pipe.HSet(key, "username", id)
	pipe.HSet(key, "hash", hash)
	pipe.HSet("user:by-username", username, id)
	_, err = pipe.Exec()
	if err != nil {
		return nil, err
	}
	return &User{key}, nil
}

func (user *User) Authenticate(password string) error {
	hash, err := user.GetHash()
	if err != nil {
		return err
	}
	err = bcrypt.CompareHashAndPassword(hash, []byte(password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return ErrInvalidLogin
	}
	return err
}
func AuthenticateUser(username, password string) error {
	user, err := GetUserByUserName(username)
	if err != nil {
		return err
	}
	return user.Authenticate(password)
}
func (user *User) GetUserName() (string, error) {
	return client.HGet(user.key, "username").Result()
}

func (user *User) GetHash() ([]byte, error) {
	return client.HGet(user.key, "hash").Bytes()
}
func GetUserByUserName(username string) (*User, error) {
	id, err := client.HGet("user:by-username", username).Int64()
	if err == redis.Nil {
		return nil, ErrUserNotFound
	} else if err != nil {
		return nil, err
	}
	key := fmt.Sprintf("user:%d", id)
	return &User{key}, nil
}

func RegisterUser(username, password string) error {
	cost := bcrypt.DefaultCost
	hash, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return err
	}
	_, err = NewUser(username, hash)
	return err
}
