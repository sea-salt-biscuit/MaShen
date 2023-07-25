package net

import (
	"awesomeProject/XueLang/MaShen/utils"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

// websocket服务
type wsServer struct {
	wsConn       *websocket.Conn
	router       *router
	outChan      chan *WsMsgRsp // Response通道 （写队列）
	Seq          int64
	property     map[string]interface{} // 属性 key-value
	propertyLock sync.RWMutex           // 读写锁
}

func NewWsServer(wsConn *websocket.Conn) *wsServer {
	return &wsServer{
		wsConn:   wsConn,
		outChan:  make(chan *WsMsgRsp, 1000),
		property: make(map[string]interface{}),
		Seq:      0,
	}
}

func (w *wsServer) Router(router *router) {
	w.router = router
}
func (w *wsServer) SetProperty(key string, value interface{}) {
	w.propertyLock.Lock()
	defer w.propertyLock.Unlock()
	w.property[key] = value
}
func (w *wsServer) GetProperty(key string) (interface{}, error) {
	w.propertyLock.RLock()
	defer w.propertyLock.RUnlock()
	return w.property[key], nil
}
func (w *wsServer) RemoveProperty(key string) {
	w.propertyLock.Lock()
	defer w.propertyLock.Unlock()
	delete(w.property, key)
}
func (w *wsServer) Addr() string {
	return w.wsConn.RemoteAddr().String()
}
func (w *wsServer) Push(name string, data interface{}) {
	rsp := &WsMsgRsp{Body: &RspBody{Name: name, Msg: data, Seq: 0}}
	w.outChan <- rsp
}

// 通道一旦建立，那么收发消息就得一直监听才行
func (w *wsServer) Start() {
	// 启动读写数据的处理逻辑
	go w.readMsgLoop()
	go w.writeMsgLoop()
}

// 写消息
func (w *wsServer) writeMsgLoop() {
	for {
		select {
		case msg := <-w.outChan:
			fmt.Println(msg)
		}
	}
}

// 读消息
func (w *wsServer) readMsgLoop() {

	// 先读到客户端 发送过来的数据 然后 进行处理 然后 再回消息
	// 经过路由 实际处理程序
	defer func() { // 出了问题 程序进行关闭
		if err := recover(); err != nil {
			log.Fatal(err)
			w.Close()
		}
	}()
	for {
		_, data, err := w.wsConn.ReadMessage()
		fmt.Println("start read Message")
		if err != nil {
			log.Println("接受消息出现错误:", err)
			break
		}
		// 收到消息 解析消息 前端发送过来的消息 就是Json格式
		// 1. data 解压  unzip

		data, err = utils.UnZip(data)
		if err != nil {
			fmt.Println("解压数据出错，非法格式：", err)
			continue
		}

		// 2. 前端的消息 加密消息 进行解密
		secretKey, err := w.GetProperty("secretKey")
		if err == nil {
			key := secretKey.(string)
			d, err := utils.AesCBCEncrypt(data, []byte(key), []byte(key), openssl.ZEROS.PADDING)
		}
		// 3. data 转为 body

		body := &ReqBody{}

		json.Unmarshal(data, body)

	}
	w.Close()
}

// 关闭通道
func (w *wsServer) Close() {
	w.wsConn.Close()
}
