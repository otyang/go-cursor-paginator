package pagination

import (
	"reflect"
	"testing"
)

func TestEncodeCursor(t *testing.T) {
	type args struct {
		cursor    string
		direction Direction
	}

	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Encode Cursor - empty cursor next",
			args: args{
				cursor:    "",
				direction: DirectionNext,
			},
			want: "",
		},
		{
			name: "Encode Cursor - empty cursor prev",
			args: args{
				cursor:    "nexcto",
				direction: DirectionPrev,
			},
			want: "cHJldjpuZXhjdG8=",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := EncodeCursor(tt.args.cursor, tt.args.direction); got != tt.want {
				t.Errorf("EncodeCursor() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewFilters(t *testing.T) {
	type args struct {
		pageCursor string
		limit      int
	}
	tests := []struct {
		name    string
		args    args
		want    Filters
		wantErr bool
	}{
		{
			name: "Filters - Generate | limit less than 1 & no cursor",
			args: args{
				pageCursor: "",
				limit:      -3,
			},
			want: Filters{
				IsFirstPage: true,
				Direction:   DirectionNext,
				Cursor:      SettingsDefaultCursor,
				Limit:       SettingsDefaultLimit + 1,
			},
			wantErr: false,
		},

		{
			name: "Filters - Generate | limit less than 1 & base64 cursor",
			args: args{
				pageCursor: "cursor1234567890",
				limit:      -3,
			},
			want: Filters{
				IsFirstPage: false,
				Direction:   DirectionNext,
				Cursor:      "cursor1234567890",
				Limit:       SettingsDefaultLimit + 1,
			},
			wantErr: false,
		},

		{
			name: "Filters - Generate | wrong cursor seperator",
			args: args{
				pageCursor: "bmV4dF9fY3Vyc29yMTIzNDU2Nzg5MA==",
				limit:      100,
			},
			want: Filters{
				IsFirstPage: false,
				Cursor:      "bmV4dF9fY3Vyc29yMTIzNDU2Nzg5MA==",
				Direction:   SettingsDefaultDirection,
				Limit:       101,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewFilters(tt.args.pageCursor, tt.args.limit, SettingsDefaultDirection)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewFilters() = %v, want %v", got, tt.want)
			}
		})
	}
}
