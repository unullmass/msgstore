package main

import (
	"fmt"
	"log"
	"time"

	"github.com/unullmass/msg-store/handlers"
	"github.com/unullmass/msg-store/models"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := gorm.Open("sqlite3", "msg.db")
	db.LogMode(true)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	// init models
	db.AutoMigrate(models.Document{}, models.Attribute{})

	r := gin.New()

	// setup logging middleware
	r.Use(gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {

		// set custom logging format
		return fmt.Sprintf("%s - [%s] \"%s %s %s %d %s \"%s\" %s\"\n",
			param.ClientIP,
			param.TimeStamp.Format(time.RFC1123),
			param.Method,
			param.Path,
			param.Request.Proto,
			param.StatusCode,
			param.Latency,
			param.Request.UserAgent(),
			param.ErrorMessage,
		)
	}))
	r.Use(gin.Recovery())
	handlers.SetRoutes(r, db)
	r.Run()
}
