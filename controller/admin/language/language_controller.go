package language

import (
	"fmt"
	"github.com/go-chi/chi"
	"html/template"
	"net/http"
	"studentbookef/config"
	"studentbookef/domain"
	"studentbookef/io/language"
)

func LanguageAdmin(app *config.Env) http.Handler {
	r := chi.NewRouter()
	r.Get("/", Home(app))
	r.Get("/delete/{languageId}", DeleteDepartment(app))
	r.Post("/create", CreateDepartmentHandler(app))
	return r
}

func DeleteDepartment(app *config.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		languageId := chi.URLParam(r, "languageId")
		langage, err := language.ReadLanguage(languageId)
		if err != nil {
			fmt.Println(err, " error reading language")
		}

		if langage.Id != "" {
			_, err := language.DeleteLanguage(langage)
			if err != nil {
				fmt.Println(err, " error reading language")
			}
		}
		http.Redirect(w, r, "/director/language", 301)
		return
	}
}

func CreateDepartmentHandler(app *config.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		langage := r.PostFormValue("language")

		if langage != "" {
			laguageObject := domain.Language{"", langage}
			_, err := language.CreateLanguage(laguageObject)
			if err != nil {
				fmt.Println(err, " error creating language")
			}
		}
		http.Redirect(w, r, "/director/language", 301)
		return
	}
}

func Home(env *config.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		languages, err := language.ReadLanguages()
		if err != nil {
			fmt.Println(err, " error reading language")
		}
		type PageData struct {
			Languages []domain.Language
		}
		data := PageData{languages}
		files := []string{
			env.Path + "admin/language/language.html",
			env.Path + "admin/template/sidebar.html",
			env.Path + "admin/template/topbar.html",
			env.Path + "template/footer.html",
		}
		ts, err := template.ParseFiles(files...)
		if err != nil {
			env.ErrorLog.Println(err.Error())
			return
		}
		err = ts.Execute(w, data)
		if err != nil {
			env.ErrorLog.Println(err.Error())
		}
	}
}
