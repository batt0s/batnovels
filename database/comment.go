package database

import "gorm.io/gorm"

type Comment struct{}

type CommentRepo interface{}

type SqlCommentRepo struct {
	db *gorm.DB
}

func NewSqlCommentRepo(db *gorm.DB) *SqlCommentRepo {
	return &SqlCommentRepo{
		db: db,
	}
}
