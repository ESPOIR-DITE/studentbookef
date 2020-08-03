package users

import (
	"fmt"
	"github.com/go-chi/chi"
	"github.com/go-chi/render"
	"html/template"
	"net/http"
	"studentbookef/config"
	"studentbookef/controller/misc"
	"studentbookef/domain"
	"studentbookef/io/user"
)

func UserAdmin(app *config.Env) http.Handler {
	r := chi.NewRouter()
	r.Get("/", HomeHandler(app))
	r.Get("/role", RoleHandler(app))
	r.Get("/gender", GenderHandler(app))
	r.Get("/getUser/{userId}", GetUserHandler(app))
	r.Post("/update_role", UpdateUserRoleHandler(app))
	r.Post("/create_role", CreateUserRoleHandler(app))
	r.Get("/delete_role/{roleId}", DeleteUserRoleHandler(app))
	return r
}

func DeleteUserRoleHandler(app *config.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userEmail := app.Session.GetString(r.Context(), "userEmail")
		result := misc.CheckAdmin(userEmail)
		if result == false {
			http.Redirect(w, r, "/user/login", 301)
			return
		}
		roleId := chi.URLParam(r, "roleId")
		userRole, err := user.ReadUserRole(roleId)
		if err != nil {
			fmt.Println(err, " error reading userRole")
		}

		if userRole.Id != "" {
			_, err := user.DeleteUserRole(userRole)
			if err != nil {
				fmt.Println(err, " error reading User Role")
			}
		}
		http.Redirect(w, r, "/director/user/role", 301)
		return
	}
}

func CreateUserRoleHandler(app *config.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userEmail := app.Session.GetString(r.Context(), "userEmail")
		result := misc.CheckAdmin(userEmail)
		if result == false {
			http.Redirect(w, r, "/user/login", 301)
			return
		}
		r.ParseForm()
		role := r.PostFormValue("role")
		description := r.PostFormValue("description")

		userRole := domain.UserRole{"", role, description}
		_, err := user.CreateUserRole(userRole)
		if err != nil {
			fmt.Println(err, " error creating user Role")
		}
		http.Redirect(w, r, "/director/user/role", 301)
	}
}

func UpdateUserRoleHandler(app *config.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		Email := app.Session.GetString(r.Context(), "userEmail")
		result := misc.CheckAdmin(Email)
		if result == false {
			http.Redirect(w, r, "/user/login", 301)
			return
		}
		r.ParseForm()
		roleId := r.PostFormValue("roleId")
		userEmail := r.PostFormValue("userEmail")
		fmt.Println(roleId, "<<<< RoleId")
		userAccount, err := user.ReadUserAccountWithEmail(userEmail)
		fmt.Println("userAccount: ", userAccount)
		if err != nil {
			fmt.Println(err, " error reading userAccount")
		}
		userAccountObject := domain.UserAccount{userAccount.Email, userAccount.Password, userAccount.AccountStatus, roleId, userAccount.Date}
		_, errx := user.UpdateUserAccount(userAccountObject)
		if errx != nil {
			fmt.Println(err, " error update userAccount")
		}
		http.Redirect(w, r, "/director/user", 301)
		return
	}
}

func GetUserHandler(app *config.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userEmail := app.Session.GetString(r.Context(), "userEmail")
		result := misc.CheckAdmin(userEmail)
		if result == false {
			http.Redirect(w, r, "/user/login", 301)
			return
		}
		var userObject domain.User
		userId := chi.URLParam(r, "userId")

		fmt.Println(userId, " <<<<UserId")
		if userId != "" {
			render.JSON(w, r, userObject)
			return
		}
		userObject, err := user.ReadUser(userId)
		if err != nil {
			http.Redirect(w, r, "/user/login", 301)
			return
		}

		render.JSON(w, r, userObject)
	}
}

func GenderHandler(app *config.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userEmail := app.Session.GetString(r.Context(), "userEmail")
		result := misc.CheckAdmin(userEmail)
		if result == false {
			http.Redirect(w, r, "/user/login", 301)
			return
		}
		files := []string{
			app.Path + "admin/admin_user/user_gender.html",
			app.Path + "admin/template/sidebar.html",
			app.Path + "admin/template/topbar.html",
		}
		ts, err := template.ParseFiles(files...)
		if err != nil {
			app.ErrorLog.Println(err.Error())
			return
		}
		err = ts.Execute(w, nil)
		if err != nil {
			app.ErrorLog.Println(err.Error())
		}
	}
}

func RoleHandler(app *config.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userEmail := app.Session.GetString(r.Context(), "userEmail")
		result := misc.CheckAdmin(userEmail)
		if result == false {
			http.Redirect(w, r, "/user/login", 301)
			return
		}

		userRole, err := user.ReadUserRoles()
		if err != nil {
			fmt.Println(err, " error reading users")
		}
		type PageData struct {
			Roles []domain.UserRole
		}
		data := PageData{userRole}
		files := []string{
			app.Path + "admin/admin_user/user_role.html",
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

func HomeHandler(app *config.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userEmail := app.Session.GetString(r.Context(), "userEmail")
		result := misc.CheckAdmin(userEmail)
		if result == false {
			http.Redirect(w, r, "/user/login", 301)
			return
		}
		Users, err := user.ReadUsers()
		if err != nil {
			fmt.Println(err, " error reading users")
		}
		userRoles, err := user.ReadUserRoles()
		if err != nil {
			fmt.Println(err, " error reading userRoles")
		}

		type PageData struct {
			Users        []domain.User
			UserRoleData []UserRoleData
			Roles        []domain.UserRole
		}
		data := PageData{Users, GetUserRole(), userRoles}
		files := []string{
			app.Path + "admin/admin_user/users.html",
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

//This type will help to aggregate necessary data for a user to be displayed in a table.
type UserRoleData struct {
	User        domain.User
	UserAccount domain.UserAccount
	Role        domain.UserRole
}

func GetUserRole() []UserRoleData {
	var userRoleDatas []UserRoleData
	var userRole domain.UserRole
	var userAccount domain.UserAccount
	users, err := user.ReadUsers()
	if err != nil {
		fmt.Println(err, " error reading Users")
		return userRoleDatas
	}
	for _, userdata := range users {
		if userdata.Email != "" {
			userAccount, err = user.ReadUserAccountWithEmail(userdata.Email)
			if err != nil {
				fmt.Println(err, " error reading UserAccount")
				//return userRoleData
			} else if userAccount.RoleId != "" {
				userRole, err = user.ReadUserRole(userAccount.RoleId)
				if err != nil {
					fmt.Println(err, " error reading UserRole")
					//return userRoleData
				}
			}

		}
		userRoleData := UserRoleData{userdata, userAccount, userRole}
		userRoleDatas = append(userRoleDatas, userRoleData)
		//cleaning the types
		userRole = domain.UserRole{}
		userAccount = domain.UserAccount{}
		userRoleData = UserRoleData{}
	}
	return userRoleDatas
}
