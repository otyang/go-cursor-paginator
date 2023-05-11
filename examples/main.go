package main

import (
	"context"
	"database/sql"
	"fmt"
	"sort"

	"github.com/otyang/yasante/pkg/pagination"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/uptrace/bun/extra/bundebug"
)

func main() {
	ctx := context.Background()

	sqlite, err := sql.Open(sqliteshim.ShimName, "file::memory:?cache=shared")
	if err != nil {
		panic(err)
	}
	sqlite.SetMaxOpenConns(1)

	db := bun.NewDB(sqlite, sqlitedialect.New())
	db.AddQueryHook(bundebug.NewQueryHook(
		bundebug.WithVerbose(true),
		bundebug.FromEnv("BUNDEBUG"),
	))

	if err := resetDB(ctx, db); err != nil {
		panic(err)
	}

	page1, err := selectPrevPage(ctx, db, "100", 10)
	if err != nil {
		panic(err)
	}

	x, z := pagination.NewCursor(page1, true, 10, "ID")

	// fmt.Printf("%v, %v", x, z)
	fmt.Printf("%v", z)
	fmt.Println("total:", x.Total, ",end:", x.End, ",Hasnext:", x.HasNextPage, ",Hasprev:", x.HasPrevPage, ",start:", x.Start, z)

	// page2, cursor, err := selectNextPage(ctx, db, cursor.End)
	// if err != nil {
	// 	panic(err)
	// }

	// page3, cursor, err := selectNextPage(ctx, db, cursor.End)
	// if err != nil {
	// 	panic(err)
	// }

	// prevPage, _, err := selectPrevPage(ctx, db, cursor.Start)
	// if err != nil {
	// 	panic(err)
	// }

	// fmt.Println("page #1", page1)
	// fmt.Println("page #2", page2)
	// fmt.Println("page #3", page3)
	// fmt.Println("prev page", prevPage)
}

type Entry struct {
	ID   int64 `bun:",pk,autoincrement"`
	Text string
}

func (e Entry) String() string {
	return fmt.Sprint(e.ID)
}

func selectNextPage(ctx context.Context, db *bun.DB, cursor string, limit int) ([]Entry, error) {
	var entries []Entry
	if err := db.NewSelect().
		Model(&entries).
		Where("id > ?", cursor).
		OrderExpr("id ASC").
		Limit(limit).
		Scan(ctx); err != nil {
		return nil, err
	}
	fmt.Println("=======", entries)
	return entries, nil
}

func selectPrevPage(ctx context.Context, db *bun.DB, cursor string, limit int) ([]Entry, error) {
	var entries []Entry
	if err := db.NewSelect().
		Model(&entries).
		Where("id < ?", cursor).
		OrderExpr("id DESC").
		Limit(limit).
		Scan(ctx); err != nil {
		return nil, err
	}
	sort.Slice(entries, func(i, j int) bool {
		return entries[i].ID < entries[j].ID
	})
	return entries, nil
}

func resetDB(ctx context.Context, db *bun.DB) error {
	if err := db.ResetModel(ctx, (*Entry)(nil)); err != nil {
		return err
	}

	seed := make([]Entry, 100)

	for i := range seed {
		seed[i] = Entry{ID: int64(i + 1), Text: fmt.Sprintf("text %d", i)}
	}

	if _, err := db.NewInsert().Model(&seed).Exec(ctx); err != nil {
		return err
	}

	return nil
}
