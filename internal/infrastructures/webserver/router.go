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
	route := r.server.app.Group("/")

	admin := route.Group("/admin")
	admin.Post("/questions", r.adminQuestionHandler.PostCreateQuestion)
	admin.Get("/questions", r.adminQuestionHandler.GetListQuestions)
	admin.Delete("/questions/:id", r.adminQuestionHandler.DeleteQuestion)
	admin.Get("/questions/:id", r.adminQuestionHandler.GetQuestionDetail)

	route.Get("/questions/random", r.userQuestionHandler.GetRandomQuestion)
	route.Post("/questions/:id/verify", r.userQuestionHandler.CheckAnswer)
}
