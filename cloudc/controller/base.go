package controller

import (
	"cloudc/model"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"log"
	"net/http"
	"strings"
)

// Server 服务模型
type Server struct {
	DB         *gorm.DB
	Router     *gin.Engine
	ws         *websocket.Conn
	deviceList []*model.Device
	pGrader    websocket.Upgrader
}

func (server *Server) init() {
	server.pGrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
}

func (server *Server) Ping(c *gin.Context) {
	var err error
	var nDevice *model.Device
	server.ws, err = server.pGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	defer server.ws.Close()
	for {
		//读取ws中的数据
		mt, message, err := server.ws.ReadMessage()
		if err != nil {
			break
		}
		//客户端主动链接服务器 服务器返回握手相应pong
		if string(message) == "ping" {
			nDevice = &model.Device{}
			nDevice.Connection = server.ws
			message = []byte("pong")
		}
		//客户端返回设备信息
		if strings.HasPrefix(string(message), "device:") {
			nDevice.DeviceType = strings.Split(string(message), ":")[1]
			fmt.Println(nDevice.DeviceType)
			server.deviceList = append(server.deviceList, nDevice)
			message = []byte("fine")
		}
		//客户端返回手机序列号 持久化设备
		if strings.HasPrefix(string(message), "serial:") {
			nDevice.DeviceSerial = strings.Split(string(message), ":")[1]
			fmt.Println(nDevice.DeviceSerial)
			nDevice.Save(server.DB, nDevice.DeviceSerial)
			message = []byte("ok")
		}
		//执行任务的时候客户端会告诉服务器端设备状态
		if strings.HasPrefix(string(message), "status:") {
			status := strings.Split(string(message), ":")[1]
			if nDevice.Status != status {
				nDevice.Status = status
				nDevice.UpdateStatus(server.DB, status)
			}
		}
		//写入ws数据
		err = server.ws.WriteMessage(mt, message)
		if err != nil {
			break
		}
	}
}

func (server *Server) Exec(c *gin.Context) {
	codeParam := c.Request.URL.Query().Get("code")
	nameParam := c.Request.URL.Query().Get("name")
	fmt.Println(len(server.deviceList))
	for _, item := range server.deviceList {
		item.Connection.WriteMessage(websocket.TextMessage, []byte("code~"+codeParam+"~"+nameParam))
	}
}

// Initialize 初始化数据库
func (server *Server) Initialize(Dbdriver, DbUser, DbPassword, DbPort, DbHost, DbName string) {
	var err error
	if Dbdriver == "mysql" {
		DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", DbUser, DbPassword, DbHost, DbPort, DbName)
		fmt.Println(DBURL)
		server.DB, err = gorm.Open(Dbdriver, DBURL)
		if err != nil {
			fmt.Printf("Cannot connect to %s database", Dbdriver)
			log.Fatal("This is the error:", err)
		} else {
			fmt.Printf("We are connected to the %s database", Dbdriver)
		}
	} else {
		fmt.Println("Unknown Driver")
	}
	//数据库初始化修改
	server.DB.Debug().AutoMigrate(
		&model.Device{},
	)
	server.Router = gin.Default()
	server.initializeRoutes()
}

// Run 系统运行入口方法
func (s *Server) Run(addr string) {
	fmt.Println("service is runing")
	log.Fatal(http.ListenAndServe(addr, s.Router))
}
