package department

import (
	"fmt"
	"github.com/go-chi/chi"
	"html/template"
	"net/http"
	"studentbookef/config"
	"studentbookef/domain"
	"studentbookef/io"
)

func DepartmentAdmin(app *config.Env) http.Handler {
	r := chi.NewRouter()
	r.Get("/", Home(app))
	r.Get("/delete/{departmentId}", DeleteDepartment(app))
	r.Post("/create", CreateDepartmentHandler(app))
	return r
}

func DeleteDepartment(app *config.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		departmentId := chi.URLParam(r, "departmentId")
		department, err := io.ReadDepartment(departmentId)
		if err != nil {
			fmt.Println(err, " error reading deparment")
		}

		if departmentId != "" {
			_, err := io.DeleteDepartment(department)
			if err != nil {
				fmt.Println(err, " error reading deparment")
			}
		}
		http.Redirect(w, r, "/director/department", 301)
		return
	}
}

func CreateDepartmentHandler(app *config.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		department := r.PostFormValue("department")
		description := r.PostFormValue("description")

		if department != "" && description != "" {
			departmentObject := domain.Department{"", department, description}
			_, err := io.CreateDepartment(departmentObject)
			if err != nil {
				fmt.Println(err, " error creating department")
			}
		}
		http.Redirect(w, r, "/director/department", 301)
		return
	}
}

func Home(env *config.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		department, err := io.ReadDepartments()
		if err != nil {
			fmt.Println(err, " error reading department")
		}
		type PageData struct {
			Departments []domain.Department
		}
		data := PageData{Departments: department}
		files := []string{
			env.Path + "admin/department/department.html",
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
