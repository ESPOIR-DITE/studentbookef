package controller

/**
* This is the main controller
* Every requests passe first here
* The purpose of this class is to direct the request(URL) coming from html to the respective controller classes
**/
import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"net/http"
	"studentbookef/config"
	"studentbookef/controller/admin"
	"studentbookef/controller/book"
	"studentbookef/controller/home"
	"studentbookef/controller/user"
)

func Controllers(env *config.Env) http.Handler {
	mux := chi.NewMux()
	mux.Use(middleware.RequestID)
	mux.Use(middleware.RealIP)
	mux.Use(middleware.Logger)
	mux.Use(env.Session.LoadAndSave)

	mux.Mount("/", home.Home(env))
	mux.Mount("/user", user.User(env))
	mux.Mount("/book", book.Book(env))
	mux.Mount("/director", admin.AdminController(env))

	fileServer := http.FileServer(http.Dir("./view/assets/"))
	// Use the mux.Handle() function to register the file server as the handler for
	// all URL paths that start with "/assets/". For matching paths, we strip the
	// "/static" prefix before the request reaches the file server.
	mux.Mount("/assets/", http.StripPrefix("/assets", fileServer))
	return mux
}
