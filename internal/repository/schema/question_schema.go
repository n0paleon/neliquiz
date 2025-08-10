package schema

import (
	"NeliQuiz/internal/domain/entities"
	"NeliQuiz/pkg/utils"
	"database/sql/driver"
	"fmt"
	"github.com/bytedance/sonic"
	"gorm.io/gorm"
	"time"
)

type Question struct {
	ID             string `gorm:"primaryKey;type:char(26)"`
	Content        string
	Hit            int
	Options        Options    `gorm:"type:jsonb;default:'[]'"`
	Categories     []Category `gorm:"many2many:question_categories"`
	ExplanationURL string
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

func (s *Question) BeforeCreate(tx *gorm.DB) (err error) {
	s.ID = utils.GenerateULID()

	for i, _ := range s.Options {
		s.Options[i].ID = utils.GenerateULID()
	}

	return
}

func (s *Question) BeforeUpdate(tx *gorm.DB) (err error) {
	if s.ID == "" {
		s.ID = utils.GenerateULID()
	}

	for i, opt := range s.Options {
		if opt.ID == "" {
			s.Options[i].ID = utils.GenerateULID()
		}
	}

	return nil
}

func (s *Question) TableName() string {
	return "questions"
}

func (s *Question) ToEntity() *entities.Question {
	options := make([]entities.Option, len(s.Options))
	for i, option := range s.Options {
		options[i] = entities.Option{
			ID:        option.ID,
			Content:   option.Content,
			IsCorrect: option.IsCorrect,
		}
	}

	categories := make([]entities.Category, len(s.Categories))
	for i, category := range s.Categories {
		categories[i] = *category.ToEntity()
	}

	return &entities.Question{
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

func ToQuestionSchema(q *entities.Question) *Question {
	options := make(Options, len(q.Options))
	for i, option := range q.Options {
		options[i] = entities.Option{
			ID:        option.ID,
			Content:   option.Content,
			IsCorrect: option.IsCorrect,
		}
	}

	categories := make([]Category, len(q.Categories))
	for i, category := range q.Categories {
		categories[i] = *ToCategorySchema(&category)
	}

	return &Question{
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

type Options []entities.Option

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
