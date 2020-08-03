package admin

import (
	"fmt"
	"github.com/go-chi/chi"
	"html/template"
	"net/http"
	"studentbookef/config"
	"studentbookef/controller/admin/books"
	"studentbookef/controller/admin/department"
	"studentbookef/controller/admin/language"
	"studentbookef/controller/admin/users"
	"studentbookef/controller/misc"
)

func AdminController(env *config.Env) http.Handler {
	mux := chi.NewMux()
	mux.Handle("/", Home(env))
	mux.Mount("/user", users.UserAdmin(env))
	mux.Mount("/book", books.BookAdmin(env))
	mux.Mount("/department", department.DepartmentAdmin(env))
	mux.Mount("/language", language.LanguageAdmin(env))
	return mux
}

func Home(env *config.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		userEmail := env.Session.GetString(r.Context(), "userEmail")
		fmt.Println(userEmail, "<<< userEmail")
		result := misc.CheckAdmin(userEmail)
		if result == false {
			http.Redirect(w, r, "/user/login", 301)
			return
		}
		fmt.Println("we are in admin home page")
		files := []string{
			env.Path + "admin/index.html",
			env.Path + "admin/template/sidebar.html",
			env.Path + "admin/template/topbar.html",
			env.Path + "template/footer.html",
		}
		ts, err := template.ParseFiles(files...)
		if err != nil {
			env.ErrorLog.Println(err.Error())
			return
		}
		err = ts.Execute(w, nil)
		if err != nil {
			env.ErrorLog.Println(err.Error())
		}
	}

}
