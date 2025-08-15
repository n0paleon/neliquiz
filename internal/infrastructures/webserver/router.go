package webserver

import (
	categoryDelivery "NeliQuiz/internal/features/category/delivery"
	questionDelivery "NeliQuiz/internal/features/question/delivery"
)

type Router struct {
	categoryHandler *categoryDelivery.CategoryHandler
	questionHandler *questionDelivery.QuestionHandler
	server          *Server
}

func NewRouter(
	categoryHandler *categoryDelivery.CategoryHandler,
	questionHandler *questionDelivery.QuestionHandler,
	server *Server,
) *Router {
	return &Router{
		categoryHandler: categoryHandler,
		questionHandler: questionHandler,
		server:          server,
	}
}

func (r *Router) RegisterRoutes() {
	route := r.server.app.Group("/")

	admin := route.Group("/admin")
	// admin question routes
	admin.Post("/questions", r.questionHandler.CreateQuestion)
	admin.Get("/questions", r.questionHandler.GetListQuestion)
	admin.Delete("/questions/:id", r.questionHandler.DeleteQuestion)
	admin.Get("/questions/:id", r.questionHandler.GetQuestionDetail)
	admin.Put("/questions/:id", r.questionHandler.UpdateQuestionDetail)

	// category routes
	route.Get("/categories", r.categoryHandler.GetListCategories)

	route.Get("/questions/random", r.questionHandler.GetRandomQuestion)
	route.Post("/questions/:id/verify", r.questionHandler.PostVerifyAnswer)
}
