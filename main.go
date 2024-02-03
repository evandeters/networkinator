package main

import (
	"networkinator/models"
	"log"
    "os"

	"gorm.io/gorm"
	"github.com/gin-gonic/gin"
    "github.com/gorilla/websocket"
)

var (
    HostCount int
    agentClients = make(map[*websocket.Conn]bool)
    webClients = make(map[*websocket.Conn]bool)
    statusChan = make(chan string)
    agentChan = make(chan string)
    db = &gorm.DB{}

    tomlConf = &models.Config{}
    configPath = "config.conf"
)

func main() {
    models.ReadConfig(tomlConf, configPath)
    db = ConnectToDB()

    f, err := os.OpenFile("server.txt", os.O_RDWR | os.O_CREATE | os.O_APPEND, 0666)
    if err != nil {
        log.Fatalf("error opening file: %v", err)
    }
    defer f.Close()

    log.SetOutput(f)

 // gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	router.LoadHTMLGlob("templates/**/*")
	router.MaxMultipartMemory = 8 << 20 // 8 MiB
	router.Static("/assets", "./assets/")

	public := router.Group("/")
	addPublicRoutes(public)

    err = db.AutoMigrate(&models.Agent{})
	if err != nil {
		log.Fatalln(err)
	}

    go handleMsg()

    log.Fatalln(router.Run(":80"))
}
