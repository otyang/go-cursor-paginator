package main

import (
	"fmt"

	paginator "github.com/otyang/go-cursor-paginator"
)

func main() {
	limitFromURL := 10
	cursorFromURL := ""

	dir, cursor, err := paginator.DecodeCursor("bmV4dDoxMjM0NTY3ODk=")
	if err != nil {
		panic("error decoding cursor: " + err.Error())
	}

	// perform db query with this
	results := ListResults(dir, cursor)

	// returning query
	// results

	cursor1, records1 := paginator.NewCursor(results, cursorFromURL == "", limitFromURL, "ID")

	fmt.Println("Total:", cursor1.Total)
	fmt.Println("	===========		")
	fmt.Println("HasPrevPage:", cursor1.HasPrevPage)
	fmt.Println("HasNextPage:", cursor1.HasNextPage)
	fmt.Println("	===========		")
	fmt.Println("Start:", cursor1.Start)
	fmt.Println("End:", cursor1.End)
	fmt.Println("	===========		")
	fmt.Println("Records:", records1)
}

type Books struct {
	ID    string
	Title string
}

func ListResults(dir paginator.Direction, cursor string) []Books {
	return []Books{
		{ID: "ISBN1", Title: "48 Laws of power -1"},
		{ID: "ISBN2", Title: "48 Laws of power -2"},
		{ID: "ISBN3", Title: "48 Laws of power -3"},
		{ID: "ISBN4", Title: "48 Laws of power -4"},
	}
}
