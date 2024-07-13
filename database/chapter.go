package database

import (
	"context"
	"time"

	"github.com/batt0s/batnovels/utils"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Chapter struct {
	ID        string         `gorm:"type:uuid;primary_key;" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Title     string         `gorm:"not null;size:128;" json:"title"`
	Content   string         `gorm:"type:text;" json:"content"`
	Slug      string         `gorm:"not null;unique;size:128;;" json:"slug"`
	ProjectID string         `json:"project_id"`
	Project   Project        `gorm:"foreignKey:ProjectID"`
}

type ChapterRepo interface {
	Find(ctx context.Context, id string) (Chapter, error)
	FindBySlug(ctx context.Context, slug string) (Chapter, error)
	Add(ctx context.Context, chapter Chapter) error
	Update(ctx context.Context, chapter Chapter) error
	Delete(ctx context.Context, chapter Chapter) error
	List(ctx context.Context, project_id string) ([]Chapter, error)
	ListBySlug(ctx context.Context, project_slug string) ([]Chapter, error)
}

type SqlChapterRepo struct {
	db *gorm.DB
}

func NewSqlChapterRepo(db *gorm.DB) *SqlChapterRepo {
	return &SqlChapterRepo{
		db: db,
	}
}

func (repo SqlChapterRepo) Find(ctx context.Context, id string) (Chapter, error) {
	select {
	case <-ctx.Done():
		return Chapter{}, ErrorOperationCanceled
	default:
		var chapter Chapter
		result := repo.db.First(&chapter, "id = ?", id)
		return chapter, result.Error
	}
}

func (repo SqlChapterRepo) FindBySlug(ctx context.Context, slug string) (Chapter, error) {
	select {
	case <-ctx.Done():
		return Chapter{}, ErrorOperationCanceled
	default:
		var chapter Chapter
		result := repo.db.Preload("Project").First(&chapter, "slug = ?", slug)
		return chapter, result.Error
	}
}

func (repo SqlChapterRepo) Add(ctx context.Context, chapter Chapter) error {
	select {
	case <-ctx.Done():
		return ErrorOperationCanceled
	default:
		if !chapter.IsValid() {
			return ErrorInvalidProject
		}
		chapter.ID = uuid.New().String()
		chapter.Slug = utils.Slugify(chapter.Title)
		result := repo.db.Create(&chapter)
		return result.Error
	}
}

func (repo SqlChapterRepo) Update(ctx context.Context, chapter Chapter) error {
	select {
	case <-ctx.Done():
		return ErrorOperationCanceled
	default:
		result := repo.db.Save(&chapter)
		return result.Error
	}
}

func (repo SqlChapterRepo) Delete(ctx context.Context, chapter Chapter) error {
	select {
	case <-ctx.Done():
		return ErrorOperationCanceled
	default:
		result := repo.db.Delete(&chapter)
		return result.Error
	}
}

func (repo SqlChapterRepo) List(ctx context.Context, project_id string) ([]Chapter, error) {
	select {
	case <-ctx.Done():
		return []Chapter{}, ErrorOperationCanceled
	default:
		var chapters []Chapter
		result := repo.db.Where("project_id = ?", project_id).Find(&chapters)
		return chapters, result.Error
	}
}

func (repo SqlChapterRepo) ListBySlug(ctx context.Context, project_slug string) ([]Chapter, error) {
	select {
	case <-ctx.Done():
		return []Chapter{}, ErrorOperationCanceled
	default:
		var chapters []Chapter
		result := repo.db.Joins("JOIN projects ON projects.id = chapters.project_id").Where("projects.slug = ?", project_slug).Find(&chapters)
		return chapters, result.Error
	}
}

func (c Chapter) IsValid() bool {
	if len(c.Title) > 128 || len(c.Title) < 3 {
		return false
	}
	if len(c.Content) < 64 {
		return false
	}
	return true
}
