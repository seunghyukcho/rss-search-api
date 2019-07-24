package rsschannel_test

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mmcdole/gofeed"
	"github.com/shhj1998/rss-search-api/rsserver/rsschannel"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFetch(t *testing.T) {
	var actualChannels, expectedChannels []*rsschannel.Schema
	expectedChannels = append(expectedChannels, &rsschannel.Schema{ID: 1, RSSLink: "rsslink01", Feed: gofeed.Feed{Title: "title01", Description: "description01", Link: "link01"}})
	expectedChannels = append(expectedChannels, &rsschannel.Schema{ID: 2, RSSLink: "rsslink02", Feed: gofeed.Feed{Title: "title02", Description: "description02", Link: "link02"}})

	db, mock, _ := sqlmock.New()
	defer db.Close()

	mockRows := sqlmock.NewRows([]string{"channel_id", "title", "description", "site_link", "rss_link"}).
		AddRow("1", "title01", "description01", "link01", "rsslink01").
		AddRow("2", "title02", "description02", "link02", "rsslink02")

	mock.ExpectQuery(`SELECT`).WillReturnRows(mockRows)
	rows, _ := db.Query(`SELECT`)

	err := rsschannel.Fetch(rows, &actualChannels)
	assert.Nil(t, err)
	assert.Equal(t, len(actualChannels), 2)

	for idx := range actualChannels {
		assert.Equal(t, expectedChannels[idx], actualChannels[idx])
	}
}
