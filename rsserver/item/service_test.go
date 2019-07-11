package item_test

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mmcdole/gofeed"
	"github.com/shhj1998/rss-search-api/rsserver/item"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestItemController_Fetch(t *testing.T) {
	var itemList []*gofeed.Item
	var answerItems []*gofeed.Item
	answerItems = append(answerItems, &gofeed.Item{GUID: "media:news:20190708154212517", Title: "I think it is a good day", Link: "link01", Description: "so do I!", Published: "2019-07-08 06:42:12"})
	answerItems = append(answerItems, &gofeed.Item{GUID: "media:news:20190708154523636", Title: "It's raining outside...", Link: "link02", Description: "I'm getting tired...", Published: "2019-07-09 07:42:12"})

	db, mock, _ := sqlmock.New()
	defer db.Close()

	mockRows := sqlmock.NewRows([]string{"guid", "title", "link", "description", "published", "creator", "url", "length", "type"}).
		AddRow("media:news:20190708154212517", "I think it is a good day", "link01", "so do I!", "2019-07-08 06:42:12", nil, nil, nil, nil).
		AddRow("media:news:20190708154523636", "It's raining outside...", "link02", "I'm getting tired...", "2019-07-09 07:42:12", nil, nil, nil, nil)

	mock.ExpectQuery("select").WillReturnRows(mockRows)
	rows, _ := db.Query(`select`)

	err := item.Fetch(rows, &itemList)
	assert.Nil(t, err)
	assert.Equal(t, len(itemList), 2)

	for idx := range itemList {
		assert.Equal(t, itemList[idx], answerItems[idx])
	}
}
