package pagination

import (
	"reflect"
)

// Cursor holds pointers to the first and last items on a page.
type Cursor struct {
	Total       int
	HasPrevPage bool
	HasNextPage bool
	Start       any
	End         any
}

func NewCursor[Records any](entries []Records, isFirstPage bool, userLimit int, cursorStructField string) (Cursor, []Records) {
	// No result at all
	if len(entries) == 0 {
		return Cursor{
			Total:       0,
			HasPrevPage: false,
			HasNextPage: false,
			Start:       "",
			End:         "",
		}, nil
	}

	// Lets ensure the records if more, is not greater than an index
	if len(entries) > userLimit {
		entries = entries[:userLimit+1]
	}

	// start := // entries[0].ID
	start := getStructFieldByName(entries[0], cursorStructField)

	if len(entries) < userLimit+1 {
		if isFirstPage {
			return Cursor{
					Total:       len(entries),
					HasPrevPage: false,
					HasNextPage: false,
					Start:       "",
					End:         "",
				},
				entries
		}

		return Cursor{
				Total:       len(entries),
				HasPrevPage: true,
				HasNextPage: false,
				Start:       start,
				End:         "",
			},
			entries
	}

	// end := entries[userLimit].ID
	end := getStructFieldByName(entries[userLimit], cursorStructField)

	if isFirstPage {
		return Cursor{
				Total:       userLimit,
				HasPrevPage: false,
				HasNextPage: true,
				Start:       "",
				End:         end,
			},
			entries[:userLimit]
	}

	return Cursor{
			Total:       len(entries[:userLimit]),
			HasPrevPage: true,
			HasNextPage: true,
			Start:       start,
			End:         end,
		},
		entries[:userLimit]
}

func getStructFieldByName(vStruct any, fieldName string) any {
	v := reflect.ValueOf(vStruct).FieldByName(fieldName)
	if !v.IsValid() {
		panic("invalid struct field. it doesn't exist, so cant be return for cursor")
	}
	return v.Interface()
}
