package core

import (
	"strconv"

	"github.com/go-redis/redis/v7"
	u "github.com/syncfuture/go/util"
)

type SaveRedisEntryCommand struct {
	Editing *RedisEntry `json:"editing"`
	Backup  *RedisEntry `json:"backup"`
}
type RedisEntry struct {
	Key    string
	Type   string
	Field  string
	Value  string
	TTL    int64
	Length uint64
	IsNew  bool `json:"isNew"`
	client redis.Cmdable
}

func (x *RedisEntry) GetType() {
	x.Type = x.client.Type(x.Key).Val()
}

func (x *RedisEntry) GetTTL() {
	ttl := x.client.TTL(x.Key).Val().Seconds()
	if ttl < 0 {
		x.TTL = -1
	} else {
		x.TTL = int64(ttl)
	}
}

func (x *RedisEntry) GetLength() {
	var err error
	switch x.Type {
	case RedisType_String:
		x.Length, err = x.client.StrLen(x.Key).Uint64()
		break
	case RedisType_Hash:
		x.Length, err = x.client.HLen(x.Key).Uint64()
		break
	case RedisType_List:
		x.Length, err = x.client.LLen(x.Key).Uint64()
		break
	case RedisType_Set:
		x.Length, err = x.client.SCard(x.Key).Uint64()
		break
	case RedisType_ZSet:
		x.Length, err = x.client.ZCard(x.Key).Uint64()
		break
	}
	u.LogError(err)
}

func (x *RedisEntry) GetValue(field string) {
	var err error
	switch x.Type {
	case RedisType_String:
		x.Value, err = x.client.Get(x.Key).Result()
		break
	case RedisType_Hash:
		if field != "" {
			x.Value, err = x.client.HGet(x.Key, field).Result()
			if err == nil {
				x.Field = field
			}
		}
		break
	case RedisType_List:
		if field != "" {
			var index int64
			index, err = strconv.ParseInt(field, 10, 64)
			if err == nil {
				x.Field = field
				x.Value, err = x.client.LIndex(x.Key, index).Result()
			}
		}
		break
	case RedisType_Set:
		if field != "" {
			x.Field = field
			x.Value = field
		}
		break
	case RedisType_ZSet:
		if field != "" {
			var score float64
			score, err = x.client.ZScore(x.Key, field).Result()
			if err == nil {
				x.Field = strconv.FormatFloat(score, 'f', -1, 64)
				x.Value = field
			}
		}
		break
	}

	u.LogError(err)
}

func NewRedisEntry(client redis.Cmdable, key string) (r *RedisEntry) {
	r = new(RedisEntry)
	r.client = client
	r.Key = key
	r.GetType()
	r.GetLength()
	r.GetTTL()

	return r
}