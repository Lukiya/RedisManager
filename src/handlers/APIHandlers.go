package handlers

import (
	"encoding/json"
	"io/ioutil"
	"strconv"
	"sync"

	"github.com/syncfuture/go/sredis"

	"github.com/go-redis/redis/v7"

	"github.com/Lukiya/redismanager/src/go/core"
	"github.com/kataras/iris/v12"
	u "github.com/syncfuture/go/util"
)

const (
	_defaultMatch = "*"
)

func getClient(ctx iris.Context) (r redis.Cmdable) {
	if core.ClusterClient == nil {
		dbStr := ctx.FormValueDefault("db", "0")
		db, err := strconv.Atoi(dbStr)
		if u.LogError(err) {
			return nil
		}
		return core.DBs[db]
	} else {
		return core.ClusterClient
	}
}

// GetKeys GET /api/keys
func GetKeys(ctx iris.Context) {
	match := ctx.FormValueDefault("match", _defaultMatch)

	entries := make([]*core.RedisEntry, 0)
	if core.ClusterClient != nil {
		mtx := new(sync.Mutex)
		core.ClusterClient.ForEachMaster(func(client *redis.Client) error {
			mtx.Lock() // ForEachMaster is running concurrently, has to lock
			defer mtx.Unlock()

			nodeKeys := sredis.GetAllKeys(client, match, 100)
			for _, key := range nodeKeys {
				keyEntry := core.NewRedisEntry(client, key)
				entries = append(entries, keyEntry)
			}

			return nil
		})
	} else {
		client := getClient(ctx)
		if client == nil {
			return
		}
		nodeKeys := sredis.GetAllKeys(client, match, 100)
		for _, key := range nodeKeys {
			keyEntry := core.NewRedisEntry(client, key)
			entries = append(entries, keyEntry)
		}
	}

	bytes, err := json.Marshal(entries)
	if u.LogError(err) {
		return
	}

	ctx.ContentType(core.ContentTypeJson)
	ctx.Write(bytes)
}

// GetDBs GET /api/dbs
func GetDBs(ctx iris.Context) {
	dbCount := len(core.DBs)
	dbs := make([]int, dbCount)
	for i := 0; i < dbCount; i++ {
		dbs[i] = i
	}

	bytes, err := json.Marshal(dbs)
	if u.LogError(err) {
		return
	}

	ctx.ContentType(core.ContentTypeJson)
	ctx.Write(bytes)
}

// // GetDBCount GET /api/db/count
// func GetDBCount(ctx iris.Context) {
// 	dbcount := strconv.Itoa(len(core.DBs))
// 	ctx.WriteString(dbcount)
// }

// GetConfigs Get /api/configs
func GetConfigs(ctx iris.Context) {
	ctx.ContentType(core.ContentTypeJson)
	bytes, err := ioutil.ReadFile("configs.json")
	if err != nil {
		ctx.WriteString(err.Error())
	} else {
		ctx.Write(bytes)
	}
}

// GetEntry Get /api/entry?key={0}&field={1}
func GetEntry(ctx iris.Context) {
	key := ctx.FormValue("key")
	if key == "" {
		ctx.WriteString("key is missing in query")
		return
	}
	field := ctx.FormValue("field")

	client := getClient(ctx)
	if client == nil {
		return
	}

	entry := core.NewRedisEntry(client, key)
	entry.GetValue(field)

	bytes, err := json.Marshal(entry)
	if u.LogError(err) {
		return
	}
	ctx.ContentType(core.ContentTypeJson)
	ctx.Write(bytes)
}

// GetHashElements Get /api/hash?key={0}
func GetHashElements(ctx iris.Context) {
	key := ctx.FormValue("key")
	if key == "" {
		ctx.WriteString("key is missing in query")
		return
	}

	client := getClient(ctx)
	if client == nil {
		return
	}

	v, err := client.HGetAll(key).Result()
	if u.LogError(err) {
		return
	}

	ctx.ContentType(core.ContentTypeJson)
	if len(v) > 0 {
		bytes, err := json.Marshal(v)
		if u.LogError(err) {
			return
		}

		ctx.Write(bytes)
	} else {
		ctx.WriteString("[]")
	}
}

// GetListElements Get /api/list?key={0}
func GetListElements(ctx iris.Context) {
	key := ctx.FormValue("key")
	if key == "" {
		ctx.WriteString("key is missing in query")
		return
	}

	client := getClient(ctx)
	if client == nil {
		return
	}

	v, err := client.LRange(key, 0, -1).Result()
	if u.LogError(err) {
		return
	}

	ctx.ContentType(core.ContentTypeJson)
	if len(v) > 0 {
		bytes, err := json.Marshal(v)
		if u.LogError(err) {
			return
		}

		ctx.Write(bytes)
	} else {
		ctx.WriteString("[]")
	}
}

// GetSetElements Get /api/set?key={0}
func GetSetElements(ctx iris.Context) {
	key := ctx.FormValue("key")
	if key == "" {
		ctx.WriteString("key is missing in query")
		return
	}

	client := getClient(ctx)
	if client == nil {
		return
	}

	v, err := client.SMembers(key).Result()
	if u.LogError(err) {
		return
	}

	ctx.ContentType(core.ContentTypeJson)
	if len(v) > 0 {
		bytes, err := json.Marshal(v)
		if u.LogError(err) {
			return
		}

		ctx.Write(bytes)
	} else {
		ctx.WriteString("[]")
	}
}

// GetZSetElements Get /api/zset?key={0}
func GetZSetElements(ctx iris.Context) {
	key := ctx.FormValue("key")
	if key == "" {
		ctx.WriteString("key is missing in query")
		return
	}

	client := getClient(ctx)
	if client == nil {
		return
	}

	v, err := client.ZRangeByScoreWithScores(key, &redis.ZRangeBy{
		Min: "-inf",
		Max: "+inf",
		// Offset: 0,
		// Count:  2,
	}).Result()
	if u.LogError(err) {
		return
	}

	ctx.ContentType(core.ContentTypeJson)
	if len(v) > 0 {
		bytes, err := json.Marshal(v)
		if u.LogError(err) {
			return
		}

		ctx.Write(bytes)
	} else {
		ctx.WriteString("[]")
	}
}
