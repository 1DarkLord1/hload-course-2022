package main

import (
	"net/http"
	"database/sql"

	"github.com/gin-gonic/gin"
)

type PutCreateRequest struct {
	Longurl string `json:"longurl"`
}

func createPutHandler(conn *sql.DB, c *gin.Context) {
	var body PutCreateRequest

	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		return
	}
	
	tinyurl, err := getTinyurl(conn, body.Longurl)

	if err == sql.ErrNoRows {
		err = insertLongurl(conn, body.Longurl)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}

		tinyurl, err = getTinyurl(conn, body.Longurl)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
			return
		}
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	tinyurlStr, err := uintToString(tinyurl)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"longurl": body.Longurl, "tinyurl": tinyurlStr})
}

func urlGetHandler(conn *sql.DB, c *gin.Context) {
	tinyurlStr := c.Params.ByName("url")

	tinyurl, err := stringToUint(tinyurlStr)

	if err != nil {
		c.JSON(404, gin.H{"message": err.Error()})
		return
	}

	longurl, err := getLongurl(conn, tinyurl)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": err.Error()})
		return
	}

	c.Redirect(302, longurl)
}
