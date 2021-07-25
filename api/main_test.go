package main

import (
	"goa-golang/internal/logger"

	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"os"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var router *gin.Engine
var testLogger logger.Logger

func TestMain(m *testing.M) {
	r, l := setupRouter()
	router = r
	testLogger = l

	testLogger.Info("Do stuff BEFORE the tests!")
	exitVal := m.Run()
	testLogger.Info("Do stuff AFTER the tests!")

	os.Exit(exitVal)
}

func TestPingRoute(t *testing.T) {
	testLogger.Info("TestPingRoute")

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/ping", nil)
	router.ServeHTTP(w, req)

	excepted, err := json.Marshal(gin.H{
		"message": "pong",
	})
	assert.NoError(t, err)

	assert.Equal(t, 200, w.Code)
	assert.Equal(t, excepted, w.Body.Bytes())
}
