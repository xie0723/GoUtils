package GoUtils

import (
	"fmt"
	"github.com/lxzan/gws"
	"sync"
	"time"
)

var logger = InitLogger()

type WebGwsSocket struct {
	conn *gws.Conn
	Host string
}

func NewWebGwsSocket(host string, conn *gws.Conn) *WebGwsSocket {
	return &WebGwsSocket{
		conn: conn,
		Host: host,
	}
}

func (w *WebGwsSocket) OnClose(socket *gws.Conn, err error) {
	logger.Errorf("onerror: err=%s\n", err.Error())
}

func (w *WebGwsSocket) OnPong(socket *gws.Conn, payload []byte) {
}

func (w *WebGwsSocket) OnOpen(socket *gws.Conn) {
	_ = socket.WriteString("ping")
}

func (w *WebGwsSocket) OnPing(socket *gws.Conn, payload []byte) {
	_ = socket.WritePong(payload)
}

func (w *WebGwsSocket) OnMessage(socket *gws.Conn, msg *gws.Message) {
	defer msg.Close()
	s := msg.Data.String()
	if s != "pong" {
		logger.Infof("recv: %s\n", msg.Data.String())
	}
	message := msg.Bytes()
	host := socket.RemoteAddr().String()
	logger.Infof("Received: %s\n %s ", message, host)
}

func (w *WebGwsSocket) SendMsg(data string) {
	//for {
	logger.Infof("Send: %s\n", data)
	w.conn.WriteString(data)
	//}
}

func (w *WebGwsSocket) HeartBeat() {
	// 启动客户端心跳发送
	ticker := time.NewTicker(30)
	go func() {
		for range ticker.C {
			w.conn.WriteString("ping")
		}
	}()
}

func CreateGwsClients() []*WebGwsSocket {
	var wg sync.WaitGroup
	var connections []*WebGwsSocket
	prodWsHosts := []string{
		"xx.xx.com",
	}
	for _, h := range prodWsHosts {
		wg.Add(1)
		go func(host string) {
			defer wg.Done()
			addr := fmt.Sprintf("ws://%s/ws/", host)
			conn, _, err := gws.NewClient(new(WebGwsSocket), &gws.ClientOption{
				Addr: addr,
				PermessageDeflate: gws.PermessageDeflate{
					Enabled:               true,
					ServerContextTakeover: true,
					ClientContextTakeover: true,
				},
			})

			if err != nil {
				logger.Errorf("创建ws连接错误:%s,错误信息:%s", addr, err)
			} else {
				logger.Infof("创建ws连接成功:%s", addr)
			}
			ws := NewWebGwsSocket(host, conn)
			ws.HeartBeat()
			connections = append(connections, ws)
			go conn.ReadLoop()
		}(h)
	}
	wg.Wait()
	return connections
}

func CreateWsClient(addr string) *gws.Conn {
	conn, _, err := gws.NewClient(new(WebGwsSocket), &gws.ClientOption{
		Addr: addr,
		PermessageDeflate: gws.PermessageDeflate{
			Enabled:               true,
			ServerContextTakeover: true,
			ClientContextTakeover: true,
		},
	})

	if err != nil {
		logger.Errorf("创建ws连接错误:%s,错误信息:%s", addr, err)
	} else {
		logger.Infof("创建ws连接成功:%s", addr)
	}
	return conn
}