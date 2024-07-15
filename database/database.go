package database

import (
	"context"
	"database/sql"
	"log"
	"strings"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Database struct {
	DB *gorm.DB

	Users    UserRepo
	Projects ProjectRepo
	Chapters ChapterRepo
	Comments CommentRepo
}

func New(driver string, source string, config *gorm.Config) (*Database, error) {
	if strings.TrimSpace(source) == "" {
		return nil, ErrorDatabaseSourceInvalid
	}
	db := &Database{}
	if err := db.connect(driver, source, config); err != nil {
		log.Println("Failed to connect database.")
		return nil, err
	}
	if err := db.DB.AutoMigrate(&User{}, &Project{}, &Chapter{}); err != nil {
		log.Println("Failed to migrate database.")
		return nil, err
	}
	db.Users = NewSqlUserRepo(db.DB)
	db.Projects = NewSqlProjectRepo(db.DB)
	db.Chapters = NewSqlChapterRepo(db.DB)
	//db.Comments = NewSqlCommentRepo(db.db)
	return db, nil
}

func (db *Database) connect(driver string, source string, config *gorm.Config) error {
	if driver == "postgres" {
		return db.connect_postgres(source, config)
	} else if driver == "sqlite" {
		return db.connect_sqlite(source, config)
	}
	return ErrorDatabaseDriverInvalid
}

func (db *Database) connect_postgres(source string, config *gorm.Config) error {
	sqlDb, err := sql.Open("postgres", source)
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := sqlDb.PingContext(ctx); err != nil {
		return err
	}
	gdb, err := gorm.Open(postgres.New(
		postgres.Config{
			Conn: sqlDb,
		}), config)
	if err != nil {
		return err
	}
	db.DB = gdb
	return nil
}

func (db *Database) connect_sqlite(source string, config *gorm.Config) error {
	sqlDb := sqlite.Open(source)
	gdb, err := gorm.Open(sqlDb, config)
	if err != nil {
		return err
	}
	db.DB = gdb
	return nil
}
