package commands

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pocketbase/dbx"
	"github.com/pocketbase/pocketbase/daos"
)

func OpenDao() (*daos.Dao, error) {
	db, err := sql.Open("sqlite3", FlagDBPath)
	if err != nil {
		return nil, err
	}
	open := dbx.NewFromDB(db, "sqlite3")
	if err != nil {
		return nil, err
	}

	return daos.New(open), nil
}
