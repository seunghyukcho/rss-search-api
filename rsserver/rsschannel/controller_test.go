package rsschannel_test

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mmcdole/gofeed"
	"github.com/shhj1998/rss-search-api/rsserver/rsschannel"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTable_Create1(t *testing.T) {
	conn, mock, _ := sqlmock.New()
	defer conn.Close()

	mock.ExpectBegin()
	mock.ExpectExec(`INSERT INTO Channel`).
		WithArgs("https://media.daum.net/syndication/society.rss").
		WillReturnResult(sqlmock.NewResult(1, 1))
	mock.ExpectCommit()

	controller := rsschannel.Table{Connection: conn}
	err := controller.Create("https://media.daum.net/syndication/society.rss")

	assert.Nil(t, err)
}

func TestTable_Create2(t *testing.T) {
	conn, mock, _ := sqlmock.New()
	defer conn.Close()

	mock.ExpectBegin()
	mock.ExpectRollback()

	controller := rsschannel.Table{Connection: conn}
	err := controller.Create("http://media.daum.net/syndication/society")

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "http error: 400 Bad Request")
}

func TestTable_Create3(t *testing.T) {
	conn, mock, _ := sqlmock.New()
	defer conn.Close()

	mock.ExpectBegin()
	mock.ExpectRollback()

	controller := rsschannel.Table{Connection: conn}
	err := controller.Create("https://www.naver.com")

	assert.NotNil(t, err)
	assert.Equal(t, err.Error(), "Failed to detect feed type")
}

func TestTable_Select(t *testing.T) {
	var actualChannels, expectedChannels []*rsschannel.Schema
	expectedChannels = append(expectedChannels, &rsschannel.Schema{ID: 1, RSSLink: "rsslink01", Feed: gofeed.Feed{Title: "title01", Description: "description01", Link: "link01"}})
	expectedChannels = append(expectedChannels, &rsschannel.Schema{ID: 2, RSSLink: "rsslink02", Feed: gofeed.Feed{Title: "title02", Description: "description02", Link: "link02"}})

	db, mock, _ := sqlmock.New()
	defer db.Close()

	mockRows := sqlmock.NewRows([]string{"channel_id", "title", "description", "site_link", "rss_link"}).
		AddRow("1", "title01", "description01", "link01", "rsslink01").
		AddRow("2", "title02", "description02", "link02", "rsslink02")

	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT channel_id, title, description, site_link, rss_link FROM Channel`).
		WillReturnRows(mockRows)
	mock.ExpectCommit()

	controller := rsschannel.Table{Connection: db}
	err := controller.Select(&actualChannels)

	assert.Nil(t, err)
	assert.Equal(t, len(actualChannels), 2)

	for idx := range actualChannels {
		assert.Equal(t, expectedChannels[idx], actualChannels[idx])
	}
}

func TestTable_SelectWithItems(t *testing.T) {
	var actualChannels, expectedChannels []*rsschannel.Schema
	expectedChannels = append(expectedChannels, &rsschannel.Schema{ID: 1, RSSLink: "rsslink01", Feed: gofeed.Feed{Title: "title01", Description: "description01", Link: "link01"}})
	expectedChannels = append(expectedChannels, &rsschannel.Schema{ID: 2, RSSLink: "rsslink02", Feed: gofeed.Feed{Title: "title02", Description: "description02", Link: "link02"}})

	expectedChannels[0].Items = append(expectedChannels[0].Items, &gofeed.Item{GUID: "media:news:20190708154212517",
		Title: "I think it is a good day", Link: "link01", Description: "so do I!", Published: "2019-07-08 06:42:12"})
	expectedChannels[1].Items = append(expectedChannels[1].Items, &gofeed.Item{GUID: "media:news:20190708154523636",
		Title: "It's raining outside...", Link: "link02", Description: "I'm getting tired...", Published: "2019-07-09 07:42:12"})

	db, mock, _ := sqlmock.New()
	defer db.Close()

	channelRows := sqlmock.NewRows([]string{"channel_id", "title", "description", "site_link", "rss_link"}).
		AddRow("1", "title01", "description01", "link01", "rsslink01").
		AddRow("2", "title02", "description02", "link02", "rsslink02")

	itemRows1 := sqlmock.NewRows([]string{"guid", "title", "link", "description", "published", "creator", "url", "length", "type"}).
		AddRow("media:news:20190708154212517", "I think it is a good day", "link01", "so do I!", "2019-07-08 06:42:12", nil, nil, nil, nil)
	itemRows2 := sqlmock.NewRows([]string{"guid", "title", "link", "description", "published", "creator", "url", "length", "type"}).
		AddRow("media:news:20190708154523636", "It's raining outside...", "link02", "I'm getting tired...", "2019-07-09 07:42:12", nil, nil, nil, nil)

	mock.ExpectBegin()
	mock.ExpectQuery(`SELECT channel_id, title, description, site_link, rss_link FROM Channel`).
		WillReturnRows(channelRows)
	mock.ExpectQuery("SELECT").WithArgs(1, 1).
		WillReturnRows(itemRows1)
	mock.ExpectQuery("SELECT").WithArgs(2, 1).
		WillReturnRows(itemRows2)
	mock.ExpectCommit()

	controller := rsschannel.Table{Connection: db}
	err := controller.SelectWithItems(&actualChannels, nil, 1)

	assert.Nil(t, err)
	assert.Equal(t, len(actualChannels), 2)

	for idx := range actualChannels {
		assert.Equal(t, expectedChannels[idx], actualChannels[idx])
	}
}
