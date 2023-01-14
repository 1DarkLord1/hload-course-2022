package main

import (
	"net/http"
	"database/sql"
	"github.com/gin-gonic/gin"
	"main/utils"
)

type PutCreateRequest struct {
	Longurl string `json:"longurl"`
}

type Service struct {
	Storage *LongurlTinyurlStorage
}

func (s *Service) init() {
	s.Storage = &LongurlTinyurlStorage{}
	s.Storage.init()
}

func (s *Service) createPutHandler(c *gin.Context) {
	var body PutCreateRequest

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	
	tinyurl, err := s.Storage.getTinyurl(body.Longurl)

	if err == sql.ErrNoRows {
		tinyurl, err = s.Storage.insertLongurl(body.Longurl)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	tinyurlStr, err := utils.UintToString(tinyurl)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"longurl": body.Longurl, "tinyurl": tinyurlStr})
}

func (s *Service) urlGetHandler(c *gin.Context) {
	tinyurlStr := c.Params.ByName("url")

	tinyurl, err := utils.StringToUint(tinyurlStr)

	if err != nil {
		c.JSON(404, gin.H{"message": err.Error()})
		return
	}

	longurl, err := s.Storage.getLongurl(tinyurl)

	if err == sql.ErrNoRows {
		c.JSON(404, gin.H{"message": err.Error()})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.Redirect(302, longurl)
}
