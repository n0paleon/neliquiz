package solo

import (
	"NeliQuiz/internal/features/game/delivery/modes/solo/dto"
	questionDomain "NeliQuiz/internal/features/question/domain"
	"NeliQuiz/internal/infrastructures/gameserver"
	"context"
	pitaya "github.com/topfreegames/pitaya/v3/pkg"
	"github.com/topfreegames/pitaya/v3/pkg/component"
	"github.com/topfreegames/pitaya/v3/pkg/timer"
)

type Room struct {
	component.Base
	timer      *timer.Timer
	app        pitaya.Pitaya
	questionUC questionDomain.QuestionUseCase
}

func (r *Room) GetRandomQuestion(_ context.Context, req *dto.GetRandomQuestionRequest) (*dto.GetRandomQuestionResponse, error) {
	q, err := r.questionUC.GetRandomQuestion(req.Categories...)
	if err != nil {
		return nil, err
	}

	return dto.ToGetRandomQuestionResponse(q), nil
}

func (r *Room) VerifyAnswer(_ context.Context, req *dto.VerifyAnswerRequest) (*dto.VerifyAnswerResponse, error) {
	isCorrect, correctOption, explanationUrl, err := r.questionUC.CheckAnswer(req.QuestionID, req.SelectedOptionID)
	if err != nil {
		return nil, pitaya.Error(err, "U-000", map[string]string{
			"message": err.Error(),
		})
	}

	response := &dto.VerifyAnswerResponse{
		IsCorrect:      isCorrect,
		ExplanationURL: explanationUrl,
		CorrectOption: dto.VerifyAnswer_Option{
			OptionID: correctOption.ID,
			Content:  correctOption.Content,
		},
	}

	return response, nil
}

// Factory
func NewSoloRoom(questionUC questionDomain.QuestionUseCase, server *gameserver.Server) *Room {
	r := &Room{
		questionUC: questionUC,
		app:        server.App,
	}

	return r
}
