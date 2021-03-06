package api

import (
	"encoding/json"
	"net/url"

	"github.com/Lukiya/redismanager/src/go/core"
	"github.com/Lukiya/redismanager/src/go/rmr"
	"github.com/Lukiya/redismanager/src/go/shared"
	"github.com/syncfuture/go/serr"
	"github.com/syncfuture/go/u"
	"github.com/syncfuture/host"
)

func getDB(ctx host.IHttpContext) (rmr.IRedisDB, error) {
	serverID := ctx.GetParamString("serverID")
	// nodeID := ctx.GetParamString("nodeID")
	db := ctx.GetParamInt("db")

	var err error

	server, err := core.Manager.GetServer(serverID)
	if err != nil {
		if serr.Is(err, shared.ConnectServerFailedError) {
			ctx.WriteJsonBytes(u.StrToBytes(`{"err":"` + err.Error() + `"}`))
		}
		return nil, err
	} else if server == nil {
		return nil, serr.Errorf("cannot find Server '%s'", serverID)
	}

	// node := Server.GetNode(nodeID)
	// if node == nil {
	// 	return nil, serr.Errorf("cannot find node '%s/%s'", serverID, nodeID)
	// }

	dB, err := server.GetDB(db)
	if host.HandleErr(err, ctx) {
		return nil, err
	}

	return dB, nil
}

func getKey(ctx host.IHttpContext) (*rmr.RedisKey, error) {
	dB, err := getDB(ctx)
	if host.HandleErr(err, ctx) {
		return nil, err
	}

	key := ctx.GetParamString("key")
	key, err = url.PathUnescape(key)
	if host.HandleErr(err, ctx) {
		return nil, err
	}

	redisKey, err := dB.GetKey(key)
	if host.HandleErr(err, ctx) {
		return nil, err
	}
	return redisKey, nil
}

func writeMsgResultError(ctx host.IHttpContext, mr *rmr.MsgResult, err error) bool {
	if err != nil {
		// ctx.StatusCode(iris.StatusInternalServerError)
		mr.MsgCode = err.Error()
		ctx.WriteString(err.Error())
		return true
	}
	return false
}

func writeMsgResult(ctx host.IHttpContext, mr *rmr.MsgResult, msg string) {
	mr.MsgCode = msg
	jsonBytes, err := json.Marshal(mr)
	u.LogError(err)
	ctx.WriteJsonBytes(jsonBytes)
}
