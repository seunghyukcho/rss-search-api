package rsschannel_test

import (
	"github.com/DATA-DOG/go-sqlmock"
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
