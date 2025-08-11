package util

import (
	"database/sql"
	"log"
)

func CommitOrRollBack(tx *sql.Tx) {
	err := recover()
	if err != nil {
		errRollBack := tx.Rollback()
		if errRollBack != nil {
			log.Fatal(errRollBack)
		}
		panic(err)
	} else {
		errCommit := tx.Commit()
		if errCommit != nil {
			log.Fatal(errCommit)
		}
	}
}