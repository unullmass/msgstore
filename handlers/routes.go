package handlers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/unullmass/msg-store/constants"
	"github.com/unullmass/msg-store/models"
)

var (
	fullSearchPath = fmt.Sprintf("%s/:%s/:%s/:%s", constants.RootPrefix, constants.TsPath, constants.KeyPath, constants.ValuePath)
	// searchByTimestamp  = fmt.Sprintf("%s/%s/:%s", constants.TsPath)
	// searchByKey        = fmt.Sprintf("%s/%s/:%s", constants.TsPath)
	// searchByValue      = fmt.Sprintf("%s/%s/:%s", constants.TsPath)
	// searchByKeyValue   = fmt.Sprintf("%s/%s/:%s", constants.TsPath)
	documentPath = fmt.Sprintf("%s/%s/:%s", constants.RootPrefix, constants.DocPath, constants.DocIdPath)

	dc *DocumentController
)

func setDocRoutes(r *gin.Engine, db *gorm.DB, wrc *chan *models.Document) {
	dc = &DocumentController{
		Db:           db,
		WriteRowChan: wrc,
	}
	r.PUT(documentPath, dc.NewDocumentHandler)
	r.GET(fullSearchPath, dc.SearchDocumentHandler)
	r.GET(documentPath, dc.RetrieveDocumentHandler)
}

// SetRoutes sets routes for all backend APIs
func SetRoutes(r *gin.Engine, db *gorm.DB, wrc *chan *models.Document) {
	setDocRoutes(r, db, wrc)
}
