package service

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
	"share.ac.cn/common"
)

var clientManager = NewWebsocketManager() //管理者

func StartWebSocket() {
	go clientManager.start()
}

// WebSocketFunc gin 处理 websocket handler
func WebSocketFunc(ctx *gin.Context) {
	upGrader := websocket.Upgrader{
		// cross origin domain
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
		// 处理 Sec-WebSocket-Protocol Header
		Subprotocols: []string{ctx.GetHeader("Sec-WebSocket-Protocol")},
	}

	conn, err := upGrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		common.Log.Errorf("websocket connect error: %s", ctx.Param("channel"))
		http.NotFound(ctx.Writer, ctx.Request)
		return
	}
	common.Log.Info("升级协议", "ua:", ctx.Request.Header["User-Agent"], "referer:", ctx.Request.Header["Referer"])

	client := NewClient(ctx.DefaultQuery("group", "local"), conn.RemoteAddr().String(), conn)

	common.Log.Infof("websocket 建立连接: %s", conn.RemoteAddr().String())

	clientManager.RegisterClient(client)

	go client.Read()
	go client.Write()
}
