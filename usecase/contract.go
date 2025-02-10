package usecase

import (
	"amartha-test/pkg/logruslogger"
	"database/sql"
	"encoding/json"
	"errors"
	"time"

	"github.com/go-redis/redis/v7"
)

// ContractUC ...
type ContractUC struct {
	ReqID     string
	DB        *sql.DB
	Tx        *sql.Tx
	Redis     *redis.Client
	EnvConfig map[string]string
}

// GetFromRedis get value from redis by key
func (uc ContractUC) GetFromRedis(key string, cb interface{}) error {
	ctx := "ContractUC.GetFromRedis"

	res, err := uc.Redis.Get(key).Result()
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "redis_get", uc.ReqID)
		return err
	}

	if res == "" {
		logruslogger.Log(logruslogger.WarnLevel, "", ctx, "redis_empty", uc.ReqID)
		return errors.New("[Redis] Value of " + key + " is empty.")
	}

	err = json.Unmarshal([]byte(res), &cb)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "json_unmarshal", uc.ReqID)
		return err
	}

	return err
}

// StoreToRedisExp save data to redis with key and exp time
func (uc ContractUC) StoreToRedisExp(key string, val interface{}, duration string) error {
	ctx := "ContractUC.StoreToRedisExp"

	dur, err := time.ParseDuration(duration)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "parse_duration", uc.ReqID)
		return err
	}

	b, err := json.Marshal(val)
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "json_marshal", uc.ReqID)
		return err
	}

	err = uc.Redis.Set(key, string(b), dur).Err()
	if err != nil {
		logruslogger.Log(logruslogger.WarnLevel, err.Error(), ctx, "redis_set", uc.ReqID)
		return err
	}

	return err
}
