package item_test

import (
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/mmcdole/gofeed"
	"github.com/shhj1998/rss-search-api/rsserver"
	"github.com/shhj1998/rss-search-api/rsserver/item"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestController_GetItems(t *testing.T) {
	var answerItems []gofeed.Item
	answerItems = append(answerItems, gofeed.Item{GUID: "media:news:20190708154212517", Title: "I think it is a good day", Link: "link01", Description: "so do I!", Published: "2019-07-08 06:42:12"})
	answerItems = append(answerItems, gofeed.Item{GUID: "media:news:20190708154523636", Title: "It's raining outside...", Link: "link02", Description: "I'm getting tired...", Published: "2019-07-09 07:42:12"})

	conn, mock, _ := sqlmock.New()
	defer conn.Close()

	mockRows := sqlmock.NewRows([]string{"guid", "title", "link", "description", "published"}).
		AddRow("media:news:20190708154212517", "I think it is a good day", "link01", "so do I!", "2019-07-08 06:42:12").
		AddRow("media:news:20190708154523636", "It's raining outside...", "link02", "I'm getting tired...", "2019-07-09 07:42:12")

	mock.ExpectQuery("SELECT").WillReturnRows(mockRows)
	db := rsserver.DB{Connection: conn}
	controller := item.Controller{Table: &db}

	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)
	controller.GetItems(ctx)

	assert.Equal(t, w.Code, http.StatusOK)
	res, _ := ioutil.ReadAll(w.Body)
	ans, _ := json.Marshal(answerItems)

	assert.Equal(t, ans, res)
}
