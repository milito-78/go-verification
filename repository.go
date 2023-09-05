package go_verification

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

type CodeRepositoryInterface interface {
	SaveCode(username, code, scope string, expiresTime time.Duration) (*VerificationCode, error)
	GetCode(username, scope string) (*VerificationCode, error)
	DeleteCode(username, scope string) bool
	DeleteAllCodes(username string) bool
}

type RedisConfig struct {
	Password string
	Prefix   string
	Addr     string
	DB       int
}

type RedisCodeRepository struct {
	client *redis.Client
	prefix string
	ctx    context.Context
}

func NewRedisCodeRepository(ctx context.Context, options RedisConfig) *RedisCodeRepository {
	client := redis.NewClient(&redis.Options{
		Addr:     options.Addr,
		Password: options.Password,
		DB:       options.DB,
	})
	err := client.Ping(ctx).Err()
	if err != nil {
		log.Fatalf("Cannot ping redis %s", err)
	}

	return &RedisCodeRepository{client: client, prefix: options.Prefix, ctx: ctx}
}

func (r RedisCodeRepository) SaveCode(username, code, scope string, expiresTime time.Duration) (*VerificationCode, error) {
	verification := &VerificationCode{
		ExpiredAt:   time.Now().Add(expiresTime),
		ExpiredTime: Duration(expiresTime),
		ExpireAfter: int(expiresTime.Seconds()),
		Username:    username,
		Scope:       scope,
		Code:        code,
	}

	data, _ := json.Marshal(&verification)
	if res := r.client.Set(r.ctx, r.createKeyScope(username, scope), data, expiresTime); res.Err() != nil {
		return nil, res.Err()
	}
	return verification, nil
}

func (r RedisCodeRepository) GetCode(username, scope string) (*VerificationCode, error) {
	res, err := r.client.Get(r.ctx, r.createKeyScope(username, scope)).Result()
	if err == redis.Nil {
		//fmt.Println("key2 does not exist")
		return nil, errors.New("does not exist")
	} else if err != nil {
		return nil, err
	} else {
		var data VerificationCode
		err := json.Unmarshal([]byte(res), &data)
		if err != nil {
			return nil, err
		}
		data.ExpireAfter = int(data.ExpiredAt.Sub(time.Now()).Seconds())
		return &data, nil
	}
}

func (r RedisCodeRepository) DeleteCode(username, scope string) bool {
	if err := r.client.Del(r.ctx, r.createKeyScope(username, scope)).Err(); err != nil {
		//log error
		return false
	}
	return true
}

func (r RedisCodeRepository) DeleteAllCodes(username string) bool {
	var cursor uint64
	for {
		keys, nextCursor, err := r.client.Scan(r.ctx, cursor, r.createKey(username), 50).Result()
		if err != nil {
			log.Printf("Error during scan keys : %s", err)
			return false
		}

		// Delete keys
		if len(keys) > 0 {
			_, delErr := r.client.Del(r.ctx, keys...).Result()
			if delErr != nil {
				log.Printf("Error during scan keys : %s", delErr)
				return false
			}
		}

		// Update the cursor for the next iteration
		cursor = nextCursor

		// If cursor is 0, the iteration is complete
		if cursor == 0 {
			break
		}
	}
	return true
}

func (r RedisCodeRepository) createKeyScope(username string, scope string) string {
	return r.prefix + ":" + scope + ":" + username
}

func (r RedisCodeRepository) createKey(username string) string {
	return r.prefix + ":*:" + username
}
