package repository

import (
	categoryDomain "NeliQuiz/internal/features/category/domain"
	categoryRepo "NeliQuiz/internal/features/category/repository"
	questionDomain "NeliQuiz/internal/features/question/domain"
	"NeliQuiz/internal/shared/strutil"
	"database/sql/driver"
	"fmt"
	"github.com/bytedance/sonic"
	"gorm.io/gorm"
	"time"
)

type QuestionSchema struct {
	ID             string `gorm:"primaryKey;type:char(26)"`
	Content        string
	Hit            int
	Options        Options                       `gorm:"type:jsonb;default:'[]'"`
	Categories     []categoryRepo.CategorySchema `gorm:"many2many:question_categories;joinForeignKey:question_id;joinReferences:category_id;"`
	ExplanationURL string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (s *QuestionSchema) BeforeCreate(tx *gorm.DB) (err error) {
	s.ID = strutil.GenerateULID()

	for i, _ := range s.Options {
		s.Options[i].ID = strutil.GenerateULID()
	}

	return
}

func (s *QuestionSchema) BeforeUpdate(tx *gorm.DB) (err error) {
	if s.ID == "" {
		s.ID = strutil.GenerateULID()
	}

	for i, opt := range s.Options {
		if opt.ID == "" {
			s.Options[i].ID = strutil.GenerateULID()
		}
	}

	return nil
}

func (s *QuestionSchema) TableName() string {
	return "questions"
}

func (s *QuestionSchema) ToEntity() *questionDomain.Question {
	options := make([]questionDomain.Option, len(s.Options))
	for i, option := range s.Options {
		options[i] = questionDomain.Option{
			ID:        option.ID,
			Content:   option.Content,
			IsCorrect: option.IsCorrect,
		}
	}

	categories := make([]categoryDomain.Category, len(s.Categories))
	for i, c := range s.Categories {
		categories[i] = *c.ToEntity()
	}

	return &questionDomain.Question{
		ID:             s.ID,
		Content:        s.Content,
		Hit:            s.Hit,
		Options:        options,
		Categories:     categories,
		ExplanationURL: s.ExplanationURL,
		CreatedAt:      s.CreatedAt,
		UpdatedAt:      s.UpdatedAt,
	}
}

func ToQuestionSchema(q *questionDomain.Question) *QuestionSchema {
	options := make(Options, len(q.Options))
	for i, option := range q.Options {
		options[i] = questionDomain.Option{
			ID:        option.ID,
			Content:   option.Content,
			IsCorrect: option.IsCorrect,
		}
	}

	categories := make([]categoryRepo.CategorySchema, len(q.Categories))
	for i, c := range q.Categories {
		categories[i] = *categoryRepo.ToCategorySchema(&c)
	}

	return &QuestionSchema{
		ID:             q.ID,
		Content:        q.Content,
		Hit:            q.Hit,
		Options:        options,
		Categories:     categories,
		ExplanationURL: q.ExplanationURL,
		CreatedAt:      q.CreatedAt,
		UpdatedAt:      q.UpdatedAt,
	}
}

type Options []questionDomain.Option

func (o Options) Value() (driver.Value, error) {
	return sonic.Marshal(o)
}

func (o *Options) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB value: %v", value)
	}
	return sonic.Unmarshal(bytes, o)
}
