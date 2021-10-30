package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/unullmass/msg-store/constants"
	"github.com/unullmass/msg-store/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
)

// documentCreateRequest is used to hold the unmarshalled DocumentCreate request payload
// The fields extracted from this are passed on to the Document model
type documentCreateRequest struct {
	Id        uuid.UUID
	Attrs     []models.Attribute
	Timestamp string
}

// documentCreateResponse is used to hold the response payload for Create Document flow
type documentCreateResponse struct {
	Id     *uuid.UUID
	Status string
}

// retrieveDocumentResponse is used to hold the retrieve document response payload
type retrieveDocumentResponse struct {
	Id        *uuid.UUID
	Attrs     []models.Attribute
	Timestamp *int64
	Status    string
}

// documentSearchResponse is used to hold the search document response payload
type documentSearchResponse struct {
	Docs []uuid.UUID
}

var (
	ErrEmptyAttrsList    = errors.New("document attributes list empty")
	ErrAttrsNotFound     = errors.New("attributes not found")
	Errbackend           = errors.New("db error")
	ErrInvalidId         = errors.New("invalid id")
	ErrInvalidTimestamp  = errors.New("invalid timestamp")
	ErrInvalidParams     = errors.New("invalid params")
	ErrDocNotFound       = errors.New("document not found")
	ErrDocCreateFailed   = errors.New("document create failed")
	ErrDocRetrieveFailed = errors.New("document retrieve failed")
	ErrDocSearchFailed   = errors.New("document search failed")
	ErrDcrParseFailed    = errors.New("document create request parse failed")
	CreateSuccess        = "created"
)

func convertTimestampStr(ts string) (*time.Time, error) {
	/* 	// convert string to big int
	   	b, ok := big.NewInt(0).SetString(ts, 10)
	   	if !ok {
	   		return ErrInvalidTimestamp
	   	}
	   	// perform range check < 0 or > maxint64
	   	z, _ := b.SetString(ts, 10)
	   	if z.Cmp(big.NewInt(0)) < 0 || z.Cmp(big.NewInt(math.MaxInt64)) > 0 {
	   		return ErrInvalidTimestamp
	   	} */
	// convert timestamp to int64
	tint64, err := strconv.ParseInt(ts, 10, 64)
	if err != nil {
		return nil, ErrInvalidTimestamp
	}
	if tint64 < 0 {
		return nil, ErrInvalidTimestamp
	}
	t := time.Unix(tint64, 0)

	return &t, nil
}

func validateSearchParams(startTs, endTs, key, value string) error {
	if _, err := convertTimestampStr(startTs); err != nil {
		return err
	}
	if _, err := convertTimestampStr(endTs); err != nil {
		return err
	}
	if err := models.ValidateAttribute(key, value); err != nil {
		return err
	}
	return nil
}

func (dc documentCreateRequest) IsValid() error {
	if len(dc.Attrs) == 0 {
		return ErrEmptyAttrsList
	}
	if dc.Id == uuid.Nil {
		return ErrInvalidId
	}
	if _, err := convertTimestampStr(dc.Timestamp); err != nil {
		return ErrInvalidTimestamp
	}

	return nil
}

type DocumentController struct {
	Db *gorm.DB
}

func (dc *DocumentController) NewDocumentHandler(c *gin.Context) {
	var dcReq documentCreateRequest
	if err := c.BindJSON(&dcReq); err != nil {
		c.JSON(http.StatusBadRequest, documentCreateResponse{
			Status: ErrDcrParseFailed.Error(),
		})
		return
	}

	// validate req
	if err := dcReq.IsValid(); err != nil {
		c.JSON(http.StatusBadRequest, documentCreateResponse{
			Status: err.Error(),
		})
		return
	}

	_, err := RetrieveDocument(dcReq.Id, dc.Db)
	// we need to change the UUID since to prevent overwrite
	if err == nil {
		dcReq.Id = uuid.New()
	}

	// convert timestamp to int64
	ts, _ := strconv.ParseInt(dcReq.Timestamp, 10, 64)

	newDoc := models.Document{
		ID:         dcReq.Id,
		Attributes: dcReq.Attrs,
		Timestamp:  time.Unix(ts, 0),
	}

	// persist to store
	if err := dc.Db.Create(&newDoc).Error; err != nil {
		c.JSON(http.StatusInternalServerError, documentCreateResponse{
			Status: errors.Wrap(err, ErrDocCreateFailed.Error()).Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, documentCreateResponse{
		Id:     &newDoc.ID,
		Status: CreateSuccess,
	})
}

func RetrieveDocument(id uuid.UUID, db *gorm.DB) (*models.Document, error) {
	doc := models.Document{}

	err := db.First(&doc, "id = ?", id).Error

	if err != nil {
		return nil, err
	}

	return &doc, nil
}

func (dc *DocumentController) RetrieveDocumentHandler(c *gin.Context) {
	id, err := uuid.Parse(c.Param(constants.DocIdPath))

	if err != nil || id == uuid.Nil {
		c.JSON(http.StatusBadRequest, retrieveDocumentResponse{
			Status: ErrInvalidId.Error(),
		})
		return
	}

	doc, err := RetrieveDocument(id, dc.Db)
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			c.JSON(http.StatusNotFound, retrieveDocumentResponse{
				Status: ErrDocNotFound.Error(),
			})
		} else {
			c.JSON(http.StatusInternalServerError, retrieveDocumentResponse{
				Status: ErrDocNotFound.Error(),
			})
		}
		return
	}

	ts := doc.Timestamp.Unix()

	c.JSON(http.StatusOK, retrieveDocumentResponse{
		Id:        &doc.ID,
		Timestamp: &ts,
	})
}

func (dc *DocumentController) SearchDocumentHandler(c *gin.Context) {
	docids := []uuid.UUID{}

	// parse request params
	startTs, _ := c.Params.Get(constants.StartTsPath)
	endTs, _ := c.Params.Get(constants.EndTsPath)
	key, _ := c.Params.Get(constants.KeyPath)
	value, _ := c.Params.Get(constants.ValuePath)

	if err := validateSearchParams(startTs, endTs, key, value); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	// convert timestamps
	stime, _ := convertTimestampStr(startTs)
	etime, _ := convertTimestampStr(endTs)

	if err := dc.Db.("Attributes").Find(&models.Document{}).
		Where("timestamp >= ? AND timestamp <= ?", *stime, *etime).
		Where("Attributes.key = ? AND Attributes.value = ?", key, value).
		Limit(constants.MaxRecordsReturn).Find(&docids).Error; err != nil {
		if gorm.IsRecordNotFoundError(err) {
			c.JSON(http.StatusNotFound, documentSearchResponse{})
		} else {
			c.JSON(http.StatusInternalServerError, documentSearchResponse{})
		}
	}

	c.JSON(http.StatusOK, documentSearchResponse{
		Docs: docids,
	})
}
