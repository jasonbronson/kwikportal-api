package transport

import (
	"encoding/json"
	"net/http"
	"reflect"

	"log"

	"github.com/gin-gonic/gin"
)

func responseError(g *gin.Context, message error) {
	g.JSON(http.StatusInternalServerError, gin.H{
		"error": message,
	})
}
func responseSuccess(g *gin.Context, field, message string) {
	g.JSON(http.StatusCreated, gin.H{
		field: message,
	})
}
func responseData(g *gin.Context, data interface{}) {
	if data == nil {
		g.JSON(http.StatusOK, gin.H{})
		return
	}
	if isEmptyArray(data) {
		g.JSON(http.StatusOK, []interface{}{})
		return
	}
	d, err := json.Marshal(data)
	if err != nil {
		log.Println(err)
	}
	g.Data(http.StatusOK,
		"application/json",
		d,
	)
}

func isEmptyArray(data interface{}) bool {
	value := reflect.ValueOf(data)
	return value.Kind() == reflect.Slice && value.Len() == 0
}
