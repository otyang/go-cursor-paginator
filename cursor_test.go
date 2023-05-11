package pagination

import (
	"reflect"
	"testing"
)

type Entry struct {
	ID   string
	Text string
}

var __records = []Entry{
	{ID: "1", Text: "Title 1"},
	{ID: "2", Text: "Title 2"},
	{ID: "3", Text: "Title 3"},
	{ID: "4", Text: "Title 4"},
	{ID: "5", Text: "Title 5"},
	{ID: "6", Text: "Title 6"},
	{ID: "7", Text: "Title 7"},
	{ID: "8", Text: "Title 8"},
	{ID: "9", Text: "Title 9"},
	{ID: "10", Text: "Title 10"},
}

var __records2 = []Entry{
	{ID: "1", Text: "Title 1"},
	{ID: "2", Text: "Title 2"},
	{ID: "3", Text: "Title 3"},
	{ID: "4", Text: "Title 4"},
}

var __dir = SettingsDefaultDirection

func TestNewCursor(t *testing.T) {
	__USERLIMIT := 10
	__case1_filter := NewFilters("", __USERLIMIT, __dir)
	__case1_nextPage_cursor := ""
	__case1_prevPage_cursor := ""
	_want_case1 := []Entry{
		{ID: "1", Text: "Title 1"},
		{ID: "2", Text: "Title 2"},
		{ID: "3", Text: "Title 3"},
		{ID: "4", Text: "Title 4"},
		{ID: "5", Text: "Title 5"},
		{ID: "6", Text: "Title 6"},
		{ID: "7", Text: "Title 7"},
		{ID: "8", Text: "Title 8"},
		{ID: "9", Text: "Title 9"},
		{ID: "10", Text: "Title 10"},
	}

	// ccase 2
	__USERLIMIT2 := 3
	__case2_filter := NewFilters("2", __USERLIMIT2, __dir)
	_want_case2 := []Entry{
		{ID: "1", Text: "Title 1"},
		{ID: "2", Text: "Title 2"},
		{ID: "3", Text: "Title 3"},
	}
	__case2_prevPage_cursor := "1"
	__case2_nextPage_cursor := "4"

	type args struct {
		entries           []Entry
		isFirstPage       bool
		userLimit         int
		cursorStructField string
	}

	tests := []struct {
		name  string
		args  args
		want  Cursor
		want1 []Entry
	}{
		{
			name: "case 1 length of records <= userlimit, on firstpage, so there should be no next & prev page",
			args: args{
				entries:           __records,
				isFirstPage:       __case1_filter.IsFirstPage,
				userLimit:         __USERLIMIT,
				cursorStructField: "ID",
			},
			want: Cursor{
				Total:       len(_want_case1),
				HasPrevPage: false,
				HasNextPage: false,
				Start:       __case1_prevPage_cursor,
				End:         __case1_nextPage_cursor,
			},
			want1: _want_case1,
		},

		{
			name: "case 2 length of records <= userlimit, not on firstpage, so there should be no next, only prev page",
			args: args{
				entries:           __records2,
				isFirstPage:       __case2_filter.IsFirstPage,
				userLimit:         __USERLIMIT2,
				cursorStructField: "ID",
			},
			want: Cursor{
				Total:       len(_want_case2),
				HasPrevPage: true,
				HasNextPage: true,
				Start:       __case2_prevPage_cursor,
				End:         __case2_nextPage_cursor,
			},
			want1: _want_case2,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := NewCursor(tt.args.entries, tt.args.isFirstPage, tt.args.userLimit, tt.args.cursorStructField)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCursor() got = %v, want %v", got, tt.want)
			}
			if !reflect.DeepEqual(got1, tt.want1) {
				t.Errorf("NewCursor() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
