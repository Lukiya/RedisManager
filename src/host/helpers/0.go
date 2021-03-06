package helpers

import (
	"strings"
)

func IsJson(str string) bool {
	if str == "" {
		return false
	}

	str = strings.TrimSpace(str)

	return strings.HasPrefix(str, "{") && strings.HasSuffix(str, "}")
}

// func CheckServer(Server *rmr.RedisServer, ctx host.IHttpContext) bool {
// 	if Server == nil {
// 		host.HandleErr(serr.New("Server cannot be nil"), ctx)
// 		return true
// 	}
// 	return false
// }

// func GetClient(ctx host.IHttpContext) (r redis.UniversalClient) {
// 	dbStr := ctx.GetFormStringDefault("db", "0")
// 	db, err := strconv.Atoi(dbStr)
// 	u.LogError(err)

// 	// proxy := core.Manager.GetSelectedServer()
// 	return proxy.GetClient(db)
// }

// func handleError(ctx host.IHttpContext, err error) bool {
// 	if err != nil {
// 		// ctx.SetStatusCode(iris.StatusInternalServerError)
// 		ctx.WriteString(err.Error())
// 		return true
// 	}
// 	return false
// }

// func writeErrorString(ctx host.IHttpContext, errStr string) {
// 	// ctx.SetStatusCode(iris.StatusInternalServerError)
// 	ctx.WriteString(errStr)
// }

// func writeMsgResultError(ctx host.IHttpContext, mr *core.MsgResult, err error) bool {
// 	if err != nil {
// 		// ctx.SetStatusCode(iris.StatusInternalServerError)
// 		mr.MsgCode = err.Error()
// 		ctx.WriteString(err.Error())
// 		return true
// 	}
// 	return false
// }

// func writeMsgResult(ctx host.IHttpContext, mr *core.MsgResult, msg string) {
// 	mr.MsgCode = msg
// 	jsonBytes, err := json.Marshal(mr)
// 	u.LogError(err)
// 	ctx.Write(jsonBytes)
// }
