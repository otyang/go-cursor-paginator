package pagination

import (
	"encoding/base64"
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var (
	SettingsCursorSeperator = ":"
	ErrCursorInvalid        = fmt.Errorf("invalid cursor: should be 'direction%sindex'", SettingsCursorSeperator)
	ErrDirectionInvalid     = fmt.Errorf("invalid direction: should be 'next' or 'prev'")
)

type Direction string

const (
	DirectionNext Direction = "next"
	DirectionPrev Direction = "prev"
)

func (d Direction) String() string {
	return string(d)
}

func (d Direction) IsValid() bool {
	switch d {
	case DirectionPrev, DirectionNext:
		return true
	default:
		return false
	}
}

func EncodeCursor(cursor string, direction Direction) string {
	if len(cursor) == 0 {
		return ""
	}
	return base64.StdEncoding.EncodeToString([]byte((direction.String() + SettingsCursorSeperator + cursor)))
}

func DecodeCursor(pageCursor string) (dir Direction, cursor string, err error) {
	b, err := base64.StdEncoding.DecodeString(pageCursor)
	if err != nil {
		return "", "", errors.New("error base64 decoding: " + err.Error())
	}

	__cursor := string(b)
	split := strings.Split(__cursor, SettingsCursorSeperator)
	if len(split) != 2 {
		return "", "", ErrCursorInvalid
	}

	direction := Direction(split[0])
	if !direction.IsValid() {
		return "", "", ErrDirectionInvalid
	}

	return direction, split[1], nil
}

// Cursor holds pointers to the first and last items on a page.
type Cursor struct {
	Total       int
	HasPrevPage bool
	HasNextPage bool
	Start       any
	End         any
}

func NewCursor[Records any](entries []Records, isFirstPage bool, limit int, cursorStructField string) (Cursor, []Records) {
	results := entries
	total := len(entries)

	if len(entries) > limit {
		total = limit
		results = entries[0 : limit-1] // remember lenOfMapKeys has indexOf -1
	}

	var start, end any

	if len(entries) < 1 {
		start = ""
		end = ""
	}
	if len(entries) > 0 {
		start = getStructFieldValue(entries[0], cursorStructField)
		end = getStructFieldValue(entries[len(entries)-1], cursorStructField)
	}
	if len(entries) > limit {
		end = getStructFieldValue(entries[limit-1], cursorStructField)
	}

	return Cursor{
		Total:       total,
		HasPrevPage: !isFirstPage,
		HasNextPage: len(entries) > limit,
		Start:       start,
		End:         end,
	}, results
}

func getStructFieldValue(vStruct any, fieldName string) any {
	v := reflect.ValueOf(vStruct).FieldByName(fieldName)
	if !v.IsValid() {
		panic("invalid struct field. it doesn't exist, so cant be return for cursor")
	}
	return v.Interface()
}
