package pagination

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncodeCursor(t *testing.T) {
	expected1 := "bmV4dDoxMjM0NTY3ODk="
	actual1 := EncodeCursor("123456789", "next")
	assert.Equal(t, expected1, actual1, "Encode Cursor1: invalid encode value")

	expected2 := "cHJldjoxMjM0NTY3ODk="
	actual2 := EncodeCursor("123456789", "prev")
	assert.Equal(t, expected2, actual2, "Encode Cursor2: invalid encode value")

	expected3 := ""
	actual3 := EncodeCursor("", "next")
	assert.Equal(t, expected3, actual3, "Encode Cursor3: invalid encode value")
}

func TestDecodeCursor(t *testing.T) {
	expectedDir1 := DirectionNext
	expectedCur1 := "123456789"

	actualDir1, actualCur1, actualErr1 := DecodeCursor("bmV4dDoxMjM0NTY3ODk=")

	assert.Equal(t, expectedDir1, actualDir1, "case 1: direction are not equal. it should be")
	assert.Equal(t, expectedCur1, actualCur1, "case 1: cursor are not equal. it should be")
	assert.NoError(t, actualErr1, "case 1: an error occured. there should be none")

	//
	expectedDir2 := Direction("")
	expectedCur2 := ""

	actualDir2, actualCur2, actualErr2 := DecodeCursor("YWZ0ZXI6cG9wb3BvcA==")

	assert.Equal(t, expectedDir2, actualDir2, "case 2: direction are not equal. it should be ")
	assert.Equal(t, expectedCur2, actualCur2, "case 2: cursor are not equal. it should be ")
	assert.Error(t, actualErr2, "case 2: an error occured. there should be none")

	_, _, err3 := DecodeCursor("YWZ0ZXI6cG9wb3BvcA==")
	assert.ErrorIs(t, err3, ErrDirectionInvalid, "case 3: an error did not occured. there should have been")

	_, _, err4 := DecodeCursor("YWZ0ZXIrK3BvcG9wb3A=")
	assert.ErrorIs(t, err4, ErrCursorInvalid, "case 4: an error did not occured. there should have been")
}

func Test_getStructFieldValue(t *testing.T) {
	example1 := struct {
		ID int
	}{1}

	example2 := struct {
		Title string
	}{"title of the page is otyoung"}

	expectedValue1 := 1
	expectedValue2 := "title of the page is otyoung"

	actual1 := getStructFieldValue(example1, "ID")
	actual2 := getStructFieldValue(example2, "Title")
	actual3 := func() {
		getStructFieldValue(example1, "FieldDoesNotExist")
	}

	assert.Equal(t, expectedValue1, actual1, "it should be same")
	assert.Equal(t, expectedValue2, actual2, "it should be same")
	assert.Panics(t, actual3, "it should panic")
}

func TestNewCursor(t *testing.T) {
	type Books struct {
		ID    string
		Title string
	}

	books := []Books{
		{ID: "ISBN1", Title: "48 Laws of power -1"},
		{ID: "ISBN2", Title: "48 Laws of power -2"},
		{ID: "ISBN3", Title: "48 Laws of power -3"},
		{ID: "ISBN4", Title: "48 Laws of power -4"},
	}

	expectedRecord1 := books
	expectedCursor1 := Cursor{
		Total:       len(books),
		HasPrevPage: false,
		HasNextPage: false,
		Start:       books[0].ID,
		End:         books[4-1].ID,
	}

	cursor1, records1 := NewCursor(books, true, 4, "ID")
	assert.Equal(t, expectedCursor1, cursor1, "cursor are not same. it should 1")
	assert.Equal(t, expectedRecord1, records1, "records are not same. it should 1")

	// case 2
	limit := 2
	expectedRecord2 := books[:limit-1]
	expectedCursor2 := Cursor{
		Total:       limit,
		HasPrevPage: true,
		HasNextPage: true,
		Start:       books[0].ID,
		End:         books[limit-1].ID,
	}

	cursor2, records2 := NewCursor(books, false, 2, "ID")

	assert.Equal(t, expectedCursor2, cursor2, "cursor are not same. it should 2")
	assert.Equal(t, expectedRecord2, records2, "records are not same. it should 2")

	assert.Panics(t,
		func() {
			NewCursor(books, true, 4, "FieldDoesNotExist")
		},
		"new cursor did not panic. it should")
}
