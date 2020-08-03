package user

import (
	"fmt"
	"github.com/go-chi/chi"
	"html/template"
	"net/http"
	"studentbookef/config"
	"studentbookef/controller/misc"
	"studentbookef/domain"
	"studentbookef/io/user"
	"time"
)

func User(app *config.Env) http.Handler {
	r := chi.NewRouter()
	r.Get("/", homeHandler(app))
	r.Get("/login", logInHandler(app))
	r.Post("/loginpost", LoginPostHandler(app))
	r.Get("/signup", SignUpHandler(app))
	r.Post("/register", RegisterHandler(app)) //this method receives signUp form
	r.Get("/userAccount/{code}", userAccountHandler(app))
	r.Post("/useraccount/register", RegisterUserAccount(app))
	r.Get("/logout", LogOutHandler(app))
	r.Get("/profile/{email}", UserProfileHandler(app))
	r.Post("/update_profile", UserUpdateProfileHandler(app))
	return r
}

func UserUpdateProfileHandler(app *config.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//myuser := domain.User{} //creating an empty object
		//clear the cession

		r.ParseForm() //Now we grabbing the contents of the form by call the name of the input(html)
		name := r.PostFormValue("name")
		email := r.PostFormValue("email")
		surname := r.PostFormValue("surname")
		cellphone := r.PostFormValue("cellphone")

		if name != "" && email != "" {
			userToUpdate := domain.User{email, name, surname, cellphone}
			_, err := user.UpdateUser(userToUpdate)
			if err != nil {
				fmt.Println(err, " error reading user")
				//app.Session.Put(r.Context(), "userEmail", email)
				app.Session.Put(r.Context(), "userMessage", "unknown_error")
				http.Redirect(w, r, "/", 301)
			} else {
				app.Session.Put(r.Context(), "userMessage", "Profile_successfully_updated")
				http.Redirect(w, r, "/", 301)
			}
		}
	}
}

func UserProfileHandler(app *config.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := chi.URLParam(r, "email")
		user, err := user.ReadUser(email)
		if err != nil {
			fmt.Println(err, " error reading user")
			//app.Session.Put(r.Context(), "userEmail", email)
			app.Session.Put(r.Context(), "userMessage", "unknown_error")
			http.Redirect(w, r, "/", 301)
		}
		type PageData struct {
			User domain.User
		}
		data := PageData{user}
		files := []string{
			app.Path + "user/profile.html",
			app.Path + "template/navigator.html",
			app.Path + "template/footer.html",
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

func LogOutHandler(app *config.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//clear the cession
		app.Session.Destroy(r.Context())
		http.Redirect(w, r, "/", 301)
		return
	}
}

func RegisterUserAccount(app *config.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		myuser := domain.User{} //creating an empty object
		//clear the cession
		app.Session.Destroy(r.Context())

		r.ParseForm() //Now we grabbing the contents of the form by call the name of the input(html)
		name := r.PostFormValue("name")
		email := r.PostFormValue("email")
		surname := r.PostFormValue("surname")
		cellphone := r.PostFormValue("cellphone")
		password1 := r.PostFormValue("password1")
		if email != "" {
			myuser = domain.User{email, name, surname, cellphone}
			userresult, err := user.UpdateUser(myuser)
			if err != nil { //when an error occurs when signing up
				fmt.Println(err, "errror in userUpdate")
				app.Session.Put(r.Context(), "userMessage", "account-confirmed_error")
				http.Redirect(w, r, "/user/signup", 301)
				return
			} else {
				userAccountObject := domain.UserAccount{userresult.Email, password1, "confirmed", "", time.Now()}
				_, err := user.UpdateUserAccount(userAccountObject)
				if err != nil { //when an error occurs when signing up
					fmt.Println(err, "errror in userUpdate")
					app.Session.Put(r.Context(), "userMessage", "account-confirmed_error")
					http.Redirect(w, r, "/user/signup", 301)
					return
				} else {
					app.Session.Remove(r.Context(), "userEmail")
					app.Session.Put(r.Context(), "userEmail", userresult.Email)
					app.Session.Put(r.Context(), "userMessage", "account-confirmed-successfully")
					http.Redirect(w, r, "/", 301)
					return
				}
				//Creating User Role

			}

		}
		app.Session.Put(r.Context(), "userMessage", "account-confirmed_error")
		http.Redirect(w, r, "/user/signup", 301)
		return
	}
}

type Message struct {
	Message string
	Class   string
}

func GetMessage(Type string) Message {
	switch Type {
	case "sign_up_error":
		text := "An error has occurred during sign up. You may have already signed up"
		return Message{text, "warning"}
	case "sign_up_success":
		text := "You have successfully sign up, please verify your email we have sent your an email with your temporary password"
		return Message{text, "info"}
	case "just_login":
		text := "Welcome back"
		return Message{text, "info"}
	case "login_error":
		text := "An error has occurred, Please try again with correct details"
		return Message{text, "warning"}
	case "post_error":
		text := "An error has occurred, Please try again later"
		return Message{text, "warning"}
	case "post_error_need_to_signup":
		text := "An error has occurred, Please try to sign in first and again later"
		return Message{text, "warning"}
	case "post_empty_error":
		text := "Please make sure that you have filled in all the required fields"
		return Message{text, "warning"}
	case "post_image_error":
		text := "Please try again uploading light pictures "
		return Message{text, "warning"}

	case "userAccount_error": // done by Taylor
		text := "An error has occurred. please check your input and try again"
		return Message{text, "warning"}

	case "userAccount_successful_added":
		text := "Thanks for your time, your account was successfully created"
		return Message{text, "info"}
	case "error_reading_book_details":
		text := "Sorry an error has occurred, please try again"
		return Message{text, "info"}
	case "login_error_missing":
		text := "Please login before checking your posts"
		return Message{text, "info"}
	case "error_update_image": //this error should be reported on edit post page
		text := "We are sorry an error has occurred while updating the picture. It may be due to the size of your picture"
		return Message{text, "warning"}
	case "update_successful": //this error should be reported on edit post page
		text := "Successful update of your book picture"
		return Message{text, "warning"}
	case "delete_successful": //this error should be reported on user_post page
		text := "Successful delete of your book post picture"
		return Message{text, "info"}
	case "register_error": //this error should be reported on user_post page
		text := "Unknown error when registering, please try again"
		return Message{text, "info"}
	case "account-confirmed-successfully": //this error should be reported on home page it means that the user has confirmed his email successfully
		text := "Unknown error when registering, please try again"
		return Message{text, "info"}
	case "account-confirmed_error": //this error should be reported on sign up or confirm registration page it means that the user failed to  confirmed his email
		text := "Unknown error when registering, please try again"
		return Message{text, "info"}
	case "unknown_error": //this error should be reported on home page it means that the something wait wrong when trying to get profile page
		text := "Unknown error, please try again"
		return Message{text, "info"}
	case "Profile_successfully_updated": //this message should be reported on home page it means that the user has update his profile successfully
		text := "Unknown error, please try again"
		return Message{text, "info"}
	}
	return Message{}
}

/****
When the user press submit button on sign up form this method will execute.
we will collect all the data in the form with r.ParseForm() method now we getting each input by passing the input name(html name).
we then create a user with only email and name other attributes will remain empty until when the user update his profile.
if an error occurs we will redirect the url address to /user/signup.
this Url will return a sign up page on user's interface with a proper error Message
But if there no errors, we will direct the user on home page with a notification Message for him/her to check the email to confirm registration.
*/
func LoginPostHandler(app *config.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		myuser := domain.UserAccount{}
		r.ParseForm()
		password := r.PostFormValue("password")
		email := r.PostFormValue("email")
		if password != "" || email != "" {
			myuser = domain.UserAccount{email, password, "", "", time.Now()}
			fmt.Println(myuser)
			resultUser, err := user.UserLog(myuser)
			if err != nil { //when an error occurs when signing up
				fmt.Println(err, " error reading resultUser")
				app.Session.Put(r.Context(), "userMessage", "sign_up_error")
				http.Redirect(w, r, "/user/login", 301)
				return
			}

			//Checking if the user is an admin
			isAdmin := misc.CheckAdmin(email)
			if isAdmin == true {
				app.Session.Put(r.Context(), "userEmail", email)
				app.Session.Put(r.Context(), "userMessage", "just_login")
				http.Redirect(w, r, "/director", 301)
				return
			}

			fmt.Println(resultUser, " result")
			if resultUser.Email != "" && resultUser.RoleId == "" {
				// If there is no error we save the login details in the cession so that we can authenticate the user during her/his cession period
				//And if the user is not a manager
				app.Session.Put(r.Context(), "userEmail", email)
				app.Session.Put(r.Context(), "userMessage", "just_login")
				http.Redirect(w, r, "/", 301)
				return
			} else {
				app.InfoLog.Println(err)
				app.Session.Put(r.Context(), "userMessage", "login_error")
				//app.Session.Put(r.Context(), "userMessage","just_login")
				http.Redirect(w, r, "/user/login", 301)
				return
			}

		} else {
			app.InfoLog.Println("error with password or username")
			app.Session.Put(r.Context(), "userMessage", "login_error")
			//app.Session.Put(r.Context(), "userMessage","just_login")
			http.Redirect(w, r, "/user/login", 301)
			return
		}

	}
}

/****
the first time when the user register on the system.
*/

func RegisterHandler(app *config.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		myuser := domain.User{} //creating an empty object
		r.ParseForm()           //Now we grabbing the contents of the form by call the name of the input(html)
		name := r.PostFormValue("name")
		email := r.PostFormValue("email")
		if email != "" {
			myuser = domain.User{email, name, "", ""}
			_, err := user.CreateUser(myuser)
			if err != nil { //when an error occurs when signing up
				app.Session.Put(r.Context(), "userMessage", "sign_up_error")
				http.Redirect(w, r, "/user/signup", 301)
				return
			}
			//creating user role
			app.Session.Put(r.Context(), "userEmail", email)
			app.Session.Put(r.Context(), "userMessage", "sign_up_success")
			http.Redirect(w, r, "/", 301)
			return

		}
	}
}

func SignUpHandler(app *config.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		Type := Message{}
		sessionType := app.Session.GetString(r.Context(), "userMessage") // we are checking what could be in the cession.
		app.Session.Remove(r.Context(), "userMessage")
		if sessionType != "" { //if there is something in the cession we need to check what it is
			Type = GetMessage(sessionType)
		}
		files := []string{
			app.Path + "user/sign_up.html",
			app.Path + "template/navigator.html",
			app.Path + "template/footer.html",
		}
		ts, err := template.ParseFiles(files...)
		if err != nil {
			app.ErrorLog.Println(err.Error())
			return
		}
		err = ts.Execute(w, Type)
		if err != nil {
			app.ErrorLog.Println(err.Error())
		}
	}
}

func logInHandler(app *config.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("voila we are in")
		files := []string{
			app.Path + "user/loginpage.html",
			app.Path + "template/navigator.html",
			app.Path + "template/footer.html",
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

func homeHandler(app *config.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		files := []string{
			app.Path + "user/loginpage.html",
			app.Path + "template/navigator.html",
			app.Path + "template/footer.html",
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

func userAccountHandler(app *config.Env) http.HandlerFunc { //done by Taylor
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("voila we are in")
		code := chi.URLParam(r, "code")
		userAccount, err := user.ReadWithpassword(code)
		if err != nil { //when an error occurs when
			fmt.Println(err, "error reading user with code")
			app.Session.Put(r.Context(), "userMessage", "register_error")
			http.Redirect(w, r, "/user/signup", 301)
			return
		} else if userAccount.Email == "" {
			app.Session.Put(r.Context(), "userMessage", "register_error")
			http.Redirect(w, r, "/user/signup", 301)
			return
		}
		user, err := user.ReadUser(userAccount.Email)
		if err != nil { //when an error occurs
			fmt.Println(err, "error reading user with email")
			app.Session.Put(r.Context(), "userMessage", "register_error")
			http.Redirect(w, r, "/user/signup", 301)
			return
		}
		if user.PhoneNumber != "" {
			app.Session.Put(r.Context(), "userEmail", user.Email)
			//app.Session.Put(r.Context(), "userMessage", "account-confirmed-successfully")
			http.Redirect(w, r, "/", 301)
			return
		}
		type PageData struct {
			User domain.User
		}

		data := PageData{user}
		files := []string{
			app.Path + "user/register.html",
			app.Path + "template/navigator.html",
			app.Path + "template/footer.html",
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
