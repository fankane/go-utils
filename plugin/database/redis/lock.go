package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/fankane/go-utils/str"
	rds "github.com/redis/go-redis/v9"
)

type DistributeLock interface {
	Lock(key string) (bool, error)
	Release() (bool, error)
}

type rdsLock struct {
	cli    *rds.Client
	key    string
	uuid   string
	expire time.Duration
}

type Params struct {
	expire time.Duration
}

type Option func(params *Params)

var (
	lockScript = rds.NewScript(`
if redis.call("GET", KEYS[1]) == ARGV[1] then
    redis.call("SET", KEYS[1], ARGV[1], "PX", ARGV[2])
    return "OK"
else
    return redis.call("SET", KEYS[1], ARGV[1], "NX", "PX", ARGV[2])
end`)

	delScript = rds.NewScript(`
if redis.call("GET", KEYS[1]) == ARGV[1] then
    return redis.call("DEL", KEYS[1])
else
    return 0
end`)

	defaultParam = &Params{
		expire: time.Second, //默认过期时间1秒
	}
)

func NewRdsLock(cli *rds.Client, opts ...Option) DistributeLock {
	if cli == nil {
		return nil
	}
	p := defaultParam
	for _, opt := range opts {
		opt(p)
	}
	return &rdsLock{cli: cli, expire: p.expire}
}

func WithExpire(expire time.Duration) Option {
	return func(params *Params) {
		params.expire = expire
	}
}

func (l *rdsLock) Lock(key string) (bool, error) {
	if key == "" {
		return false, fmt.Errorf("key can not be empty")
	}
	uuid := str.UUID()
	resp, err := lockScript.Eval(context.Background(), l.cli, []string{key}, []string{uuid,
		fmt.Sprintf("%d", l.expire.Milliseconds())}).Result()
	if err == rds.Nil {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("error on acquiring lock for %s, %s", key, err)
	} else if resp == nil {
		return false, nil
	}
	reply, ok := resp.(string)
	if ok && reply == "OK" {
		l.uuid = uuid
		l.key = key
		return true, nil
	}
	return false, nil
}

func (l *rdsLock) Release() (bool, error) {
	if l.key == "" || l.uuid == "" {
		return false, nil
	}
	resp, err := delScript.Eval(context.Background(), l.cli, []string{l.key}, []string{l.uuid}).Result()
	if err != nil {
		return false, err
	}
	reply, ok := resp.(int64)
	if !ok {
		return false, nil
	}
	return reply == 1, nil
}
