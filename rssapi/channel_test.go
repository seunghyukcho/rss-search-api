package rssapi_test

import (
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/shhj1998/rss-search-api/rssapi"
	"github.com/shhj1998/rss-search-api/rsserver"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
	"testing"
)

func serverSetting() rssapi.Server {
	var ret rssapi.Server

	if err := godotenv.Load(); err != nil {
		panic(err)
	}

	name, address, id, password := os.Getenv("DB_NAME"), os.Getenv("DB_HOST"), os.Getenv("DB_ID"), os.Getenv("DB_PW")
	rssDB := rsserver.DB{}

	if err := rssDB.Open(name, address, id, password); err != nil {
		panic(err)
	}

	ret = rssapi.Server{DB: &rssDB}
	return ret
}

func BenchmarkServer_SelectChannels(b *testing.B) {
	server := serverSetting()
	defer server.DB.Close()

	gin.SetMode(gin.TestMode)
	c, _ := gin.CreateTestContext(httptest.NewRecorder())

	for i := 0; i < b.N; i++ {
		server.SelectChannels(c)
	}
}

func BenchmarkServer_SelectChannelsWithItems(b *testing.B) {
	server := serverSetting()
	defer server.DB.Close()

	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{URL: &url.URL{}}

	for i := 0; i < b.N; i++ {
		if server.SelectChannelsWithItems(c); w.Code != http.StatusOK {
			b.Fail()
		}
	}
}

func BenchmarkServer_SelectChannelsWithItems2(b *testing.B) {
	server := serverSetting()
	defer server.DB.Close()

	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	sampleURL, _ := url.Parse("https://sample/api/v1/channel/items?id=3")
	c, _ := gin.CreateTestContext(w)
	c.Request = &http.Request{URL: sampleURL}

	for i := 0; i < b.N; i++ {
		if server.SelectChannelsWithItems(c); w.Code != http.StatusOK {
			b.Fail()
		}
	}
}

func BenchmarkServer_SelectChannelsWithItems3(b *testing.B) {
	server := serverSetting()
	defer server.DB.Close()

	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	for cnt := 1; cnt <= 10; cnt++ {
		sampleURL, _ := url.Parse("https://sample/api/v1/channel/items/count/" + strconv.Itoa(cnt))
		c.Request = &http.Request{URL: sampleURL}
		b.Run("count"+strconv.Itoa(cnt), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				if server.SelectChannelsWithItems(c); w.Code != http.StatusOK {
					b.Fail()
				}
			}
		})
	}
}
