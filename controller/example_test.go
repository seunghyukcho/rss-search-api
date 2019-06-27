package controller_test

import (
	"github.com/shhj1998/d2-fest-backend/controller"
	"testing"
)

func TestforTest(t *testing.T) {
	result := controller.Sample(1)

	if result != 1 {
		t.Error("Wrong Answer")
	}
}
