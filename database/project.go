package database

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Project struct {
	ID        string         `gorm:"type:uuid;primary_key;" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Title     string         `gorm:"not null;size:256;" json:"title"`
	Synopsis  string         `gorm:"not null;size:1024;" json:"synopsis"`
	Author    string         `gorm:"not null;size:128;" json:"author"`
	Status    string         `gorm:"not null;size:64;" json:"status"`
	Tags      string         `gorm:"not null;size:256;" json:"tags"` // , ile ayırarak
	Views     int32          `json:"views"`
	Image     string         `json:"image"`
	Slug      string         `gorm:"not null;unique;size:128;;" json:"slug"`
}

type ProjectRepo interface {
	Find(ctx context.Context, id string) (Project, error)
	FindBySlug(ctx context.Context, slug string) (Project, error)
	Add(ctx context.Context, project Project) (Project, error)
	Update(ctx context.Context, project Project) (Project, error)
	Delete(ctx context.Context, project Project) error
	List(ctx context.Context, limit int, offset int, orderby string) ([]Project, error)
}

type SqlProjectRepo struct {
	db *gorm.DB
}

func NewSqlProjectRepo(db *gorm.DB) *SqlProjectRepo {
	return &SqlProjectRepo{
		db: db,
	}
}

func (repo SqlProjectRepo) Find(ctx context.Context, id string) (Project, error) {
	select {
	case <-ctx.Done():
		return Project{}, ErrorOperationCanceled
	default:
		var project Project
		result := repo.db.First(&project, "id = ?", id)
		return project, result.Error
	}
}

func (repo SqlProjectRepo) FindBySlug(ctx context.Context, slug string) (Project, error) {
	select {
	case <-ctx.Done():
		return Project{}, ErrorOperationCanceled
	default:
		var project Project
		result := repo.db.First(&project, "slug = ?", slug)
		return project, result.Error
	}
}

func (repo SqlProjectRepo) Add(ctx context.Context, project Project) (Project, error) {
	select {
	case <-ctx.Done():
		return project, ErrorOperationCanceled
	default:
		if !project.IsValid() {
			return project, ErrorInvalidProject
		}
		project.ID = uuid.New().String()
		project.Slug = Slugify(project.Title)
		result := repo.db.Create(&project)
		return project, result.Error
	}
}

func (repo SqlProjectRepo) Update(ctx context.Context, project Project) (Project, error) {
	select {
	case <-ctx.Done():
		return project, ErrorOperationCanceled
	default:
		result := repo.db.Save(&project)
		return project, result.Error
	}
}

func (repo SqlProjectRepo) Delete(ctx context.Context, project Project) error {
	select {
	case <-ctx.Done():
		return ErrorOperationCanceled
	default:
		result := repo.db.Delete(&project)
		return result.Error
	}
}

func (repo SqlProjectRepo) List(ctx context.Context, limit int, offset int, orderby string) ([]Project, error) {
	select {
	case <-ctx.Done():
		return []Project{}, ErrorOperationCanceled
	default:
		var projects []Project
		result := repo.db.Limit(limit).Offset(offset).Order(orderby).Find(&projects)
		return projects, result.Error
	}
}

func (p Project) IsValid() bool {
	if len(p.Title) > 256 || len(p.Title) < 3 {
		return false
	}
	if len(p.Synopsis) > 1024 || len(p.Synopsis) < 64 {
		return false
	}
	return true
}
