package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/pkg/errors"

	"github.com/unullmass/msg-store/cache"
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
	Id     *uuid.UUID `json:"id,omitempty"`
	Status string
}

// retrieveDocumentResponse is used to hold the retrieve document response payload
type retrieveDocumentResponse struct {
	Id        *uuid.UUID         `json:"id,omitempty"`
	Attrs     []models.Attribute `json:"attrs,omitempty"`
	Timestamp *int64             `json:"timestamp,omitempty"`
	Status    string
}

// documentSearchResponse is used to hold the search document response payload
type documentSearchResponse struct {
	Docs []uuid.UUID `json:"docs,omitempty"`
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
	StatusOk             = "ok"
)

const (
	decBase  = 10
	lenInt64 = 64
)

func convertTimestampStr(ts string) (*time.Time, error) {
	// convert timestamp to int64
	tint64, err := strconv.ParseInt(ts, decBase, lenInt64)
	if err != nil || tint64 < 0 {
		return nil, ErrInvalidTimestamp
	}

	t := time.Unix(tint64, 0)

	return &t, nil
}

func validateSearchParams(ts, key, value string) error {
	if _, err := convertTimestampStr(ts); err != nil {
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
	Db           *gorm.DB
	WriteRowChan chan *models.Document
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

	// update DocID on Attributes
	for i := range dcReq.Attrs {
		if dcReq.Attrs[i].ID == uuid.Nil {
			dcReq.Attrs[i].ID = uuid.New()
		}
		if dcReq.Attrs[i].DocumentID == uuid.Nil {
			dcReq.Attrs[i].DocumentID = dcReq.Id
		}
	}

	newDoc := models.Document{
		ID:         dcReq.Id,
		Attributes: dcReq.Attrs,
		Timestamp:  time.Unix(ts, 0),
	}

	// put to write chan
	dc.WriteRowChan <- &newDoc

	c.JSON(http.StatusCreated, documentCreateResponse{
		Id:     &newDoc.ID,
		Status: CreateSuccess,
	})

	// cache result
	_ = cache.ReadCache.Set(newDoc.ID.String(), newDoc, 1)
}

func RetrieveDocument(id uuid.UUID, db *gorm.DB) (*models.Document, error) {
	var doc models.Document

	// check if cached
	d, ok := cache.ReadCache.Get(id.String())
	if ok {
		doc = d.(models.Document)
	} else { // fetch from DB
		err := db.Preload("Attributes").First(&doc, "id = ?", id).Error
		if err != nil {
			return nil, err
		}
		// cache result
		_ = cache.ReadCache.Set(doc.ID.String(), doc, 1)
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
		c.JSON(http.StatusInternalServerError, retrieveDocumentResponse{
			Status: ErrDocNotFound.Error(),
		})
		return
	}

	ts := doc.Timestamp.Unix()

	c.JSON(http.StatusOK, retrieveDocumentResponse{
		Id:        &doc.ID,
		Timestamp: &ts,
		Attrs:     doc.Attributes,
		Status:    StatusOk,
	})
}

func (dc *DocumentController) SearchDocumentHandler(c *gin.Context) {

	// parse request params
	ts, _ := c.Params.Get(constants.TsPath)
	key, _ := c.Params.Get(constants.KeyPath)
	value, _ := c.Params.Get(constants.ValuePath)

	if err := validateSearchParams(ts, key, value); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
	}

	// convert timestamps
	ttime, _ := convertTimestampStr(ts)

	docids := []uuid.UUID{}

	rows, err := dc.Db.Model(&models.Document{}).Select("Documents.ID").
		Joins("JOIN Attributes on Documents.ID = Attributes.document_id").
		Where("Documents.timestamp = ?", *ttime).
		Where("Attributes.key = ? AND Attributes.value = ?", key, value).
		Limit(constants.DefaultMaxRecordsReturn).Rows()

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
	}

	defer rows.Close()

	for rows.Next() {
		var id uuid.UUID
		if err := rows.Scan(&id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
		}
		docids = append(docids, id)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, documentSearchResponse{})
	}

	c.JSON(http.StatusOK, documentSearchResponse{
		Docs: docids,
	})
}
