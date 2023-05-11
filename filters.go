package pagination

import (
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
)

var (
	SettingsDefaultLimit     = 200
	SettingsDefaultCursor    = "0"
	SettingsDefaultDirection = DirectionNext
	SettingsCursorSeperator  = ":"
	ErrCursorInvalid         = fmt.Errorf("invalid cursor: should be 'direction%sindex'", SettingsCursorSeperator)
)

type Direction string

const (
	DirectionNext Direction = "next"
	DirectionPrev Direction = "prev"
)

func (d Direction) String() string {
	return string(d)
}

func (d Direction) Validate() Direction {
	switch d {
	case DirectionPrev, DirectionNext:
		return d
	default:
		return SettingsDefaultDirection
	}
}

type Filters struct {
	// IfFirstPage helps determine if to show next or prev page
	IsFirstPage bool
	// Direction if next or prev. Used when making query if greater or less than
	Direction Direction
	// Cursor to begin read from.
	Cursor string
	// LimitPerPage as requested by user. and addition of 1 is made so we do all query at a go
	Limit int
}

// NewFilters decodes the Base64 enccode cursor and also
// generates filter params to use for the database query
func NewFilters(pageCursor string, userLimit int, dir Direction) Filters {
	userLimit = userLimit + 1
	if userLimit < 1 || userLimit > SettingsDefaultLimit {
		userLimit = SettingsDefaultLimit + 1
	}

	if pageCursor == "" {
		return Filters{
			IsFirstPage: true,
			Direction:   dir.Validate(),
			Cursor:      SettingsDefaultCursor,
			Limit:       userLimit,
		}
	}

	return Filters{
		IsFirstPage: false,
		Direction:   dir.Validate(),
		Cursor:      pageCursor,
		Limit:       userLimit,
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

	return Direction(split[0]).Validate(), split[1], nil
}
