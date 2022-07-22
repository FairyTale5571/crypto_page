package database

import (
	"database/sql"
	"os"

	"github.com/fairytale5571/crypto_page/pkg/logger"
	_ "github.com/go-sql-driver/mysql"
)

type DB struct {
	db     *sql.DB
	logger *logger.LoggerWrapper
}

func New(uri string) (*DB, error) {

	var err error
	dbConnection := DB{
		logger: logger.New("database"),
	}
	db, err := sql.Open("mysql", uri)
	if err != nil {
		dbConnection.logger.Fatalf("error open database: %v", err)
		return nil, err
	}
	db.SetMaxOpenConns(10)
	dbConnection.db = db

	dbConnection.startMigrate()

	version, err := dbConnection.Version()
	if err != nil {
		dbConnection.logger.Errorf("error version database: %v", err)
		return nil, err
	}
	dbConnection.logger.Infof("database version: %s", version)

	return &dbConnection, nil
}

func (db *DB) isMigrated(filename string) bool {
	var version string
	err := db.QueryRow("SELECT version FROM migrations WHERE version = ?", filename).Scan(&version)
	if err != nil {
		return false
	}
	if version != "" {
		return true
	}
	return false
}

func (db *DB) startMigrate() {

	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS discord_migrations (
	   id int UNSIGNED NOT NULL AUTO_INCREMENT,
	   version varchar(255) NOT NULL,
	   time datetime NULL DEFAULT current_timestamp,
	   PRIMARY KEY (id)
	);`)
	if err != nil {
		db.logger.Fatalf("error create migrations table: %v", err)
		return
	}

	files, err := os.ReadDir("migrates")
	if err != nil {
		db.logger.Errorf("error read migrates files: %v", err)
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		name := file.Name()
		if !db.isMigrated(name) {

			read, err := os.ReadFile("migrates/" + name)
			if err != nil {
				db.logger.Errorf("error read migrates file: %v", err)
			}
			_, err = db.Exec("INSERT INTO migrations (version, time) VALUES (?, now())", name)
			if err != nil {
				db.logger.Errorf("error migrate: %v", err)
			}
			_, err = db.Exec(string(read))
			if err != nil {
				db.logger.Errorf("error migrate: %v", err)
			}
		}
	}

}

func (db *DB) Close() error {
	return db.db.Close()
}

func (db *DB) Exec(query string, args ...any) (sql.Result, error) {
	return db.db.Exec(query, args...)
}

func (db *DB) Query(query string, args ...any) (*sql.Rows, error) {
	return db.db.Query(query, args...)
}

func (db *DB) QueryRow(query string, args ...any) *sql.Row {
	return db.db.QueryRow(query, args...)
}

func (db *DB) Prepare(query string) (*sql.Stmt, error) {
	return db.db.Prepare(query)
}

func (db *DB) Version() (string, error) {
	var version string
	err := db.QueryRow("SELECT VERSION()").Scan(&version)
	if err != nil {
		return "", err
	}
	return version, nil
}
