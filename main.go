package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/unullmass/msg-store/cache"
	"github.com/unullmass/msg-store/constants"
	"github.com/unullmass/msg-store/handlers"
	"github.com/unullmass/msg-store/models"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

const (
	helpText = `The following settings must be set in the environment before launch:
	DB_HOST
	DB_PORT
	DB_USERNAME
	DB_PASSWORD
	DB_SCHEMA
	`
)

//checkVarFound checks if a param is set in the environent
func checkVarFound(myvar *string, paramName string) bool {
	val, isFound := os.LookupEnv(paramName)
	if !isFound || strings.TrimSpace(val) == "" {
		fmt.Printf("%s is unset", paramName)
		return false
	}
	*myvar = val
	return true
}

func FlushAllInsertsToDB(db *gorm.DB, insertChan chan *models.Document, quit chan struct{}) {
	retryChan := make(chan *models.Document)

	for {
		select {
		case d := <-insertChan:
			if err := db.Create(d).Error; err != nil {
				// requeue for insertion
				retryChan <- d
			}
		case d := <-retryChan:
			if err := db.Create(d).Error; err != nil {
				log.Println("Error inserting Document ", d.ID)
			}
		case <-quit:
			// flush all data
			for d := range insertChan {
				if err := db.Create(d).Error; err != nil {
					log.Println("Error inserting Document ", d.ID)
				}
			}
			// close the insert channel
			close(insertChan)
			return
		}
	}
}

func main() {
	var (
		dbHost, dbUser, dbPass, dbSchema, dbPort string
	)
	defer cache.ReadCache.Clear()

	if !checkVarFound(&dbHost, constants.DbHostEnv) ||
		!checkVarFound(&dbUser, constants.DbUserEnv) ||
		!checkVarFound(&dbPass, constants.DbPassEnv) ||
		!checkVarFound(&dbSchema, constants.DbSchemaEnv) ||
		!checkVarFound(&dbPort, constants.DbPortEnv) {
		fmt.Println(helpText)
		os.Exit(1)
	}

	if _, err := strconv.Atoi(dbPort); err != nil {
		log.Fatal(errors.Wrap(err, "Invalid value for DB_PORT"))
	}

	dsn := fmt.Sprintf("%s%s:%s@%s:%s/%s?%s", constants.DbDSNBase,
		dbUser, dbPass, dbHost, dbPort, dbSchema, constants.DbInsecureConn)
	fmt.Println(dsn)
	db, err := gorm.Open(constants.DbType, dsn)
	if err != nil {
		log.Fatal(err)
	}
	//db.DB().SetMaxIdleConns(500)
	//db.DB().SetMaxOpenConns(1000)

	db.LogMode(true)

	defer db.Close()
	// init models
	db.AutoMigrate(models.Document{}, models.Attribute{})

	ic := make(chan *models.Document)
	quit := make(chan struct{})
	go FlushAllInsertsToDB(db, ic, quit)

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
	close(quit)
}
