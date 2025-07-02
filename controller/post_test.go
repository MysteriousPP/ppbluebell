package controller

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestCreatePostHandler(t *testing.T) {
	r := gin.Default()
	url := "/api/v1/post"
	r.POST(url, CreatePostHandler)
	body := `{
		"title": "test",
		"content": "nihao",
		"community_id": 1
	}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, url, bytes.NewReader([]byte(body)))
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.Contains(t, w.Body.String(), "需要登录")
	// res := new(ResponseData)
	// if err := json.Unmarshal(w.Body.Bytes(), res); err != nil {
	// 	t.Fatalf("json.Unmarshal w.Body failed, err: %v\n", err)
	// }
	// assert.Equal(t, res.Code, CodeNeedLogin)
}
