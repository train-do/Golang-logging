package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"github.com/train-do/Golang-Restfull-API/database"
	"github.com/train-do/Golang-Restfull-API/handler"
	"github.com/train-do/Golang-Restfull-API/library"
	mid "github.com/train-do/Golang-Restfull-API/middleware"
	"github.com/train-do/Golang-Restfull-API/repository"
	"github.com/train-do/Golang-Restfull-API/service"
	"go.uber.org/zap"
)

func NewRouter() *chi.Mux {
	logger := library.InitLog()
	db, err := database.InitDB()
	if err != nil {
		logger.Fatal("Database Connection Failed", zap.Error(err))
		// log.Fatal(err)
	}

	hUser := handler.NewUserHandler()
	rBook := repository.NewBookRepository(db, logger)
	sBook := service.NewBookService(rBook, logger)
	hBook := handler.NewBookHandler(sBook, logger)
	rOrder := repository.NewOrderRepository(db)
	sOrder := service.NewOrderService(rOrder)
	hOrder := handler.NewOrderHandler(sOrder)
	rReview := repository.NewReviewRepository(db)
	sReview := service.NewReviewService(rReview)
	hReview := handler.NewReviewHandler(sReview)
	// mWare := mid.NewMiddleware(logger)
	router := chi.NewRouter()

	// router.Use(middleware.Logger)
	fileServer := http.FileServer(http.Dir("./static"))
	router.Handle("/static/*", http.StripPrefix("/static", fileServer))
	router.Get("/login", hUser.Login)
	router.Post("/login", hUser.Login)
	// router.Group(func(r chi.Router) {
	// 	r.Use(mWare.Authentication)
	// 	r.Use(mWare.MinddlewareLogger)
	// 	r.Get("/dashboard", hBook.Dashboard)
	// 	r.Get("/addBook", hBook.CreateBook)
	// 	r.Post("/addBook", hBook.CreateBook)
	// 	r.Get("/books", hBook.GetAllBook)
	// 	r.Get("/orders", hOrder.GetAllOrder)
	// 	r.Put("/discount/{id}", hBook.UpdateBook)
	// 	r.Get("/logout", hUser.Logout)
	// 	r.Post("/logout", hUser.Logout)
	// })
	router.Group(func(r chi.Router) {
		r.Use(mid.Authentication)
		r.Use(mid.MinddlewareLogger(logger))
		r.Get("/dashboard", hBook.Dashboard)
		r.Post("/addBook", hBook.CreateBook)
		r.Get("/books", hBook.GetAllBook)
		r.Get("/orders", hOrder.GetAllOrder)
		r.Put("/discount", hBook.UpdateBook)
		r.Get("/logout", hUser.Logout)
		r.Post("/logout", hUser.Logout)
	})
	router.Group(func(r chi.Router) {
		r.Get("/customer/order", hOrder.CreateOrder)
		r.Post("/customer/review", hReview.CreateReview)
	})
	return router
}
