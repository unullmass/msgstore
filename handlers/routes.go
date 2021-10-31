package handlers

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/unullmass/msg-store/constants"
)

var (
	fullSearchPath = fmt.Sprintf("%s/:%s/:%s/:%s", constants.RootPrefix, constants.TsPath, constants.KeyPath, constants.ValuePath)
	// searchByTimestamp  = fmt.Sprintf("%s/%s/:%s", constants.TsPath)
	// searchByKey        = fmt.Sprintf("%s/%s/:%s", constants.TsPath)
	// searchByValue      = fmt.Sprintf("%s/%s/:%s", constants.TsPath)
	// searchByKeyValue   = fmt.Sprintf("%s/%s/:%s", constants.TsPath)
	retrievePath       = fmt.Sprintf("%s/%s/:%s", constants.RootPrefix, constants.DocPath, constants.DocIdPath)
	createDocumentPath = constants.RootPrefix

	dc *DocumentController
)

func setDocRoutes(r *gin.Engine, db *gorm.DB) {
	dc = &DocumentController{
		Db: db,
	}
	r.POST(createDocumentPath, dc.NewDocumentHandler)
	r.GET(fullSearchPath, dc.SearchDocumentHandler)
	r.GET(retrievePath, dc.RetrieveDocumentHandler)
}

// SetRoutes sets routes for all backend APIs
func SetRoutes(r *gin.Engine, db *gorm.DB) {
	setDocRoutes(r, db)
}