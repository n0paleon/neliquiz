package webserver

import "NeliQuiz/internal/delivery/http"

type Router struct {
	adminQuestionHandler *http.AdminQuestionHandler
	userQuestionHandler  *http.UserQuestionHandler
	server               *Server
}

func NewRouter(
	adminQuestionHandler *http.AdminQuestionHandler,
	userQuestionHandler *http.UserQuestionHandler,
	server *Server,
) *Router {
	return &Router{
		adminQuestionHandler: adminQuestionHandler,
		userQuestionHandler:  userQuestionHandler,
		server:               server,
	}
}

func (r *Router) RegisterRoutes() {
	route := r.server.app.Group("/api")

	admin := route.Group("/admin")
	admin.Post("/create-question", r.adminQuestionHandler.PostCreateQuestion)
	admin.Get("/get-questions", r.adminQuestionHandler.GetListQuestions)
	admin.Post("/delete-question", r.adminQuestionHandler.DeleteQuestion)

	route.Get("/get-random", r.userQuestionHandler.GetRandomQuestion)
	route.Post("/verify-answer", r.userQuestionHandler.CheckAnswer)
}
