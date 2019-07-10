package rsserver

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/mmcdole/gofeed"
	"testing"
)

func TestItemController_GetItems(t *testing.T) {
	var itemList []*gofeed.Item

	db, mock, _ := sqlmock.New()
	defer db.Close()

	mockRows := sqlmock.NewRows([]string{"guid", "title", "link", "description", "published"}).
		AddRow("example1", "I think it is a good day", "link01", "so do I!", "2019-07-08 06:42:12").
		AddRow("example2", "It's raining outside...", "linke02", "I'm getting tired...", "2019-07-09 07:42:12")
	mock.ExpectQuery("select").WillReturnRows(mockRows)
	rows, _ := db.Query(`select`)

	err := FetchItems(rows, &itemList)
	if err != nil {
		t.Errorf("Error occured. It should be nil")
	}

	if len(itemList) != 2 {
		t.Error(len(itemList), 1)
	}
}
