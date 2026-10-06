package config

import (
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	"github.com/syrlramadhan/api-bendahara-inovdes/helper"
)

func ConnectToDatabase() (db *sql.DB, err error) {
	var w http.ResponseWriter
	err = godotenv.Load()
	if err != nil {
		helper.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return nil, err
	}

	dbName := os.Getenv("DB_NAME")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	mysql := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPass, dbHost, dbPort, dbName)
	db, err = sql.Open("mysql", mysql)
	if err != nil {
		helper.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return db, err
	}

	err = db.Ping()
	if err != nil {
		helper.WriteJSONError(w, http.StatusInternalServerError, err.Error())
		return db, err
	}

	db.SetConnMaxLifetime(time.Minute * 3)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	// Auto-migration for member table: jabatan and status
	migrateMemberTable(db)

	return db, nil
}

func migrateMemberTable(db *sql.DB) {
	var colCount int
	err := db.QueryRow("SELECT COUNT(*) FROM information_schema.COLUMNS WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = 'member' AND COLUMN_NAME = 'jabatan'").Scan(&colCount)
	if err == nil && colCount == 0 {
		_, _ = db.Exec("ALTER TABLE member ADD COLUMN jabatan VARCHAR(20) DEFAULT 'anggota' AFTER nama")
		_, _ = db.Exec("UPDATE member SET jabatan = CASE WHEN LOWER(status) = 'bph' THEN 'bph' ELSE 'anggota' END")
		_, _ = db.Exec("UPDATE member SET status = CASE WHEN LOWER(status) IN ('inactive', 'nonaktif') THEN 'nonaktif' ELSE 'aktif' END")
		_, _ = db.Exec("ALTER TABLE member MODIFY COLUMN status VARCHAR(20) DEFAULT 'aktif'")
	}
}
