package books

import (
	"fmt"
	"github.com/go-chi/chi"
	"html/template"
	"net/http"
	"studentbookef/config"
	"studentbookef/controller/misc"
	"studentbookef/domain"
	"studentbookef/io/book_io"
	location2 "studentbookef/io/location"
	user2 "studentbookef/io/user"
)

func BookAdmin(app *config.Env) http.Handler {
	r := chi.NewRouter()
	r.Get("/", Home(app))
	//r.Get("/", Home(app))
	return r
}

func Home(app *config.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userEmail := app.Session.GetString(r.Context(), "userEmail")
		result := misc.CheckAdmin(userEmail)
		if result == false {
			http.Redirect(w, r, "/user/login", 301)
			return
		}
		books, err := book_io.ReadBooks()
		if err != nil {
			fmt.Println(err, " error reading books")
		}

		type PageData struct {
			Books     []domain.Book
			UserBooks []UserBook
		}
		data := PageData{books, GetBookAndUsers()}
		files := []string{
			app.Path + "admin/admin_book/books.html",
			app.Path + "admin/template/sidebar.html",
			app.Path + "admin/template/topbar.html",
		}
		ts, err := template.ParseFiles(files...)
		if err != nil {
			app.ErrorLog.Println(err.Error())
			return
		}
		err = ts.Execute(w, data)
		if err != nil {
			app.ErrorLog.Println(err.Error())
		}
	}
}

type UserBook struct {
	User     domain.User
	Books    domain.Book
	Post     domain.BookPost
	Location domain.Location
}

func GetBookAndUsers() []UserBook {
	var userBooks []UserBook
	//Getting users with their books
	bookPosts, err := book_io.ReadBookPosts()
	if err != nil {
		fmt.Println(err, "error reading books")
		return userBooks
	}
	for _, bookpost := range bookPosts {
		book, err := book_io.ReadBook(bookpost.BookId)
		if err != nil {
			fmt.Println(err, "error reading book")
		}
		user, err := user2.ReadUser(bookpost.Email)
		if err != nil {
			fmt.Println(err, "error reading user")
		}
		location, err := location2.ReadLocation(bookpost.LocationId)
		if err != nil {
			fmt.Println(err, "error reading location")
		}
		userBook := UserBook{user, book, bookpost, location}
		userBooks = append(userBooks, userBook)
	}
	return userBooks
}
