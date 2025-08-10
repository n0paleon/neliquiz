package webserver

import "NeliQuiz/internal/delivery/http"

type Router struct {
	adminQuestionHandler *http.AdminQuestionHandler
	adminCategoryHandler *http.AdminCategoryHandler
	userQuestionHandler  *http.UserQuestionHandler
	server               *Server
}

func NewRouter(
	adminQuestionHandler *http.AdminQuestionHandler,
	adminCategoryHandler *http.AdminCategoryHandler,
	userQuestionHandler *http.UserQuestionHandler,
	server *Server,
) *Router {
	return &Router{
		adminQuestionHandler: adminQuestionHandler,
		adminCategoryHandler: adminCategoryHandler,
		userQuestionHandler:  userQuestionHandler,
		server:               server,
	}
}

func (r *Router) RegisterRoutes() {
	route := r.server.app.Group("/")

	admin := route.Group("/admin")
	// question routes
	admin.Post("/questions", r.adminQuestionHandler.PostCreateQuestion)
	admin.Get("/questions", r.adminQuestionHandler.GetListQuestions)
	admin.Delete("/questions/:id", r.adminQuestionHandler.DeleteQuestion)
	admin.Get("/questions/:id", r.adminQuestionHandler.GetQuestionDetail)
	admin.Put("/questions/:id", r.adminQuestionHandler.PutQuestion)
	// category routes
	admin.Get("/categories", r.adminCategoryHandler.GetListCategories)

	route.Get("/questions/random", r.userQuestionHandler.GetRandomQuestion)
	route.Post("/questions/:id/verify", r.userQuestionHandler.CheckAnswer)
}
