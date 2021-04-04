package rmr

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/syncfuture/go/sconv"
	"github.com/syncfuture/go/serr"
	"github.com/syncfuture/host"
)

type RedisNode struct {
	ID            string
	Addr          string
	password      string
	DBs           []*RedisDB
	clusterClient redis.UniversalClient
}

func NewStandaloneReidsNode(addr, pwd string) *RedisNode {
	r := &RedisNode{
		ID:       host.GenerateID(),
		Addr:     addr,
		password: pwd,
	}

	return r
}

func NewClusterRedisNode(addr string, clusterClient redis.UniversalClient) *RedisNode {
	r := &RedisNode{
		ID:            host.GenerateID(),
		Addr:          addr,
		clusterClient: clusterClient,
	}

	return r
}

func (x *RedisNode) LoadDBs() error {
	var err error
	x.DBs, err = x.GetDBs()
	return err
}

func (x *RedisNode) GetDBs() ([]*RedisDB, error) {
	if x.clusterClient != nil {
		db0 := NewRedisDB(0, x.clusterClient)
		return []*RedisDB{db0}, nil
	} else {
		db0Client := x.createStandaloneClient(0)
		db0 := NewRedisDB(0, db0Client)
		databases, err := db0Client.ConfigGet(context.Background(), "databases").Result()
		if err != nil {
			return nil, serr.WithStack(err)
		}
		dbcount := sconv.ToInt(databases[1])

		dbs := make([]*RedisDB, 0, dbcount)
		dbs = append(dbs, db0)

		for i := 1; i < dbcount; i++ { // skip db0 since it's already been added
			client := x.createStandaloneClient(i)
			dbs = append(dbs, NewRedisDB(i, client))
		}

		return dbs, nil
	}
}

func (x *RedisNode) createStandaloneClient(db int) redis.UniversalClient {
	return redis.NewClient(&redis.Options{
		Addr:     x.Addr,
		Password: x.password,
		DB:       db,
	})
}