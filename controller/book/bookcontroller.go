package book

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"html/template"
	"io/ioutil"
	"net/http"
	"strconv"
	"studentbookef/config"
	"studentbookef/controller/misc"
	"studentbookef/controller/user"
	"studentbookef/domain"
	"studentbookef/io"
	"studentbookef/io/book_io"
	language2 "studentbookef/io/language"
	location2 "studentbookef/io/location"
	"studentbookef/io/picture_io"
	user2 "studentbookef/io/user"
	"time"
)

func Book(app *config.Env) http.Handler {
	r := chi.NewRouter()
	r.Get("/", NewPostHandler(app))
	r.Post("/post_Book", PostBookHandler(app)) //this method receives signUp form
	r.Post("/post_book_location", PostLocationHandler(app))
	r.Get("/book_Image", BookImageHandler(app))
	r.Post("/post_book_Image", PostBookImage(app))
	r.Get("/location", LocationHandler(app))
	r.Get("/details/{bookId}", DetailsHandler(app))
	r.Get("/category", CategoryHandler(app))
	r.Get("/get_category/{departmentId}", GetCategoryHandler(app))
	r.Get("/get_mypost", GetPostHandler(app)) //This Method is called when a user want to see all his/her posts.
	r.Get("/post_edit/{bookId}", EditPostHandler(app))
	r.Post("/image_update", updatePicture(app))

	//r.Post("/post_book_location",PostBookLocation(app))
	return r
}

func updatePicture(app *config.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		//var status bool
		//we need to check the user if has created an account first
		email := app.Session.GetString(r.Context(), "userEmail")
		if email == "" {
			app.Session.Put(r.Context(), "userMessage", "post_error_need_to_signup")
			http.Redirect(w, r, "/login", 301)
			return
		}
		fmt.Println("voila we are in picture update process")
		//picture := domain.Picture{}
		r.ParseForm()
		fmt.Println(" reading the file")

		file, _, err := r.FormFile("file")
		pictureId := r.PostFormValue("pictureId")
		picturedescription := r.PostFormValue("picturedescription")
		bookId := r.PostFormValue("bookId")
		fmt.Println(file)
		fmt.Println(pictureId)

		if err != nil {
			fmt.Println(err, "<<<<<<>>>>>>>")
			app.Session.Put(r.Context(), "userMessage", "post_image_error")
			http.Redirect(w, r, "/book/post_edit/"+bookId, 301)
			return
		}
		fmt.Println(" read successful")

		//fmt.Println(" converting to []byte and into slice array", handler)
		reader := bufio.NewReader(file)
		//Converting the files into byteArrays
		content, _ := ioutil.ReadAll(reader)

		newpicture := domain.Picture{pictureId, content, picturedescription}
		_, errx := picture_io.UpdatePicture(newpicture)
		if errx != nil {
			fmt.Println(err, "<<<<<<>>>>>>>")
			app.Session.Put(r.Context(), "userMessage", "error_update_image")
			http.Redirect(w, r, "/book/post_edit/"+bookId, 301)
			return
		}

		// if all go well
		fmt.Println("bookImage creation Successful: ")
		app.Session.Put(r.Context(), "userMessage", "update_successful")
		http.Redirect(w, r, "/book/post_edit/"+bookId, 301)
		return

	}
}

func EditPostHandler(app *config.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var picture1 string
		var picture2 string
		var picture1Id string
		var picture2Id string
		var department domain.Department
		var bookPostObject domain.BookPost
		myUser := domain.User{}

		bookId := chi.URLParam(r, "bookId")
		if bookId == "" {
			app.Session.Put(r.Context(), "userMessage", "error_reading_book_details")
			http.Redirect(w, r, "/", 301)
			return
		}
		bookDepartment, err := book_io.ReadBookDepartment(bookId)
		if err != nil {
			app.InfoLog.Println(err, " reading bookDepartment")
		}
		if bookDepartment.DepartmentId != "" {
			department, err = io.ReadDepartment(bookDepartment.DepartmentId)
			if err != nil {
				app.InfoLog.Println(err, " reading department")
			}
		}

		bookPost, err := book_io.ReadWithBookId(bookId)
		if err != nil {
			app.InfoLog.Println(err, " reading bookPost")
		} else {
			/***
			Here I HAVE TRIED TO REPLACE LOCATIONID FIELD WITH FORMATTED TIME IN STRING TO ALLOW HTML TO READ TIME PROPERLY
			*/
			bookPostObject = domain.BookPost{bookPost.Id, bookPost.Email, bookPost.BookId, bookPost.Date, misc.FormatDateTime(bookPost.Date), bookPost.Description, bookPost.Description}
			fmt.Println(bookPostObject.LocationId, "<<<<<<Time")
		}
		book, err := book_io.ReadBook(bookId)
		if err != nil {
			app.InfoLog.Println(err, " reading book")
		}

		bookImage, err := book_io.ReadAllOfBookImage(bookId)
		if err != nil {
			app.InfoLog.Println(err)
		}

		pictures, err := picture_io.ReadAllOf(getBookImageArray(bookImage))
		if err != nil {
			app.InfoLog.Println(err, "  reading images")
		}
		/******
		Be careful here. the picture object has been reversed since the backend in the following way:
		i have returned base64 string on the place of pictureId and the picture id in th place of description
		*/
		for index, valeu := range pictures {
			if index == 0 {
				picture1 = valeu.Id
				picture1Id = valeu.Description
			} else {
				picture2 = valeu.Id
				picture2Id = valeu.Description
			}
		}
		location, err := location2.ReadLocation(bookPost.LocationId)
		if err != nil {
			app.InfoLog.Println(err, " Location")
		}

		user, err := user2.ReadUser(bookPost.Email)
		if err != nil {
			app.InfoLog.Println(err, " user")
		}
		departments, err := io.ReadDepartments()
		if err != nil {
			app.InfoLog.Println(err, " reading departments")
		}
		languages, err := language2.ReadLanguages()
		if err != nil {
			app.InfoLog.Println(err, " reading languages")
		}
		language := getBookLanguage(bookId)

		type PageData struct {
			Book        domain.Book
			Department  domain.Department
			Post        domain.BookPost
			Picture1    string
			Picture2    string
			Picture1Id  string
			Picture2Id  string
			User        domain.User
			BookOwner   domain.User
			Location    domain.Location
			Departments []domain.Department
			Languages   []domain.Language
			Language    domain.Language
		}
		data := PageData{book, department, bookPostObject, picture1, picture2, picture1Id, picture2Id, myUser, user, location, departments, languages, language}

		//we need to check the user if has created an account first
		//email := app.Session.GetString(r.Context(), "userEmail")
		//if email == "" {
		//	app.Session.Put(r.Context(), "userMessage", "post_error_need_to_signup")
		//	http.Redirect(w, r, "/login", 301)
		//	return
		//}

		files := []string{
			app.Path + "user/edit_post.html",
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
func getBookLanguage(bookId string) domain.Language {
	entity := domain.Language{}
	bookLanguage, err := book_io.ReadBookLanguage(bookId)
	if err != nil { //if an error, we return an empty object
		fmt.Println(err, "  error reading booklanguage")
		return entity
	}
	language, err := language2.ReadLanguage(bookLanguage.LanguageId)
	if err != nil { //if an error, we return an empty object
		fmt.Println(err, "  error reading language")
		return entity
	}
	return language
}

func GetPostHandler(app *config.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := app.Session.GetString(r.Context(), "userEmail") //We are checking for the user has login? Here if the user is logged in, we should redirect him to the login page
		myUser := domain.User{}
		/****
		Here we are checking if the user has logged in and we verify if the user is authentic
		*/
		if email != "" {
			err := errors.New("")
			myUser, err = user2.ReadUser(email)
			if err != nil {
				app.InfoLog.Println(err)
				app.Session.Put(r.Context(), "userMessage", "login_error_missing")
				http.Redirect(w, r, "/user/login", 301)
				return
			}
		} else {
			app.Session.Put(r.Context(), "userMessage", "login_error_missing")
			//app.Session.Put(r.Context(), "userMessage","just_login")
			http.Redirect(w, r, "/user/login", 301)
			return
		}
		userPosts := getUserBookDetails(email)
		type PageData struct {
			User            domain.User
			BookPostDetails []homePosts
		}
		data := PageData{myUser, userPosts}

		files := []string{
			app.Path + "user/user_post.html",
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

func GetCategoryHandler(app *config.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		myUser := domain.User{}
		departmentId := chi.URLParam(r, "departmentId")
		email := app.Session.GetString(r.Context(), "userEmail") //We are checking of the user has login?
		if email != "" {
			err := errors.New("")
			myUser, err = user2.ReadUser(email)
			if err != nil {
				app.InfoLog.Println(err)
			}
		}
		userDertement, err := io.ReadDepartment(departmentId)
		if err != nil {
			app.InfoLog.Println(err, "  reading userDepartment")
		}
		departments, err := io.ReadDepartments()
		if err != nil {
			app.InfoLog.Println(err, "  reading department")
		}

		bookOfDepartment := getBookOfDepartment(departmentId)
		homePost := getBookDetails()
		type PageData struct {
			User             domain.User
			Departements     []domain.Department
			BookPostDetails  []homePosts
			Department       domain.Department
			BookOfDepartment []homePostCategories
		}
		data := PageData{myUser, departments, homePost, userDertement, bookOfDepartment}
		fmt.Println("voila we are in")
		files := []string{
			app.Path + "book/category_of_department.html",
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

func CategoryHandler(app *config.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		myUser := domain.User{}

		email := app.Session.GetString(r.Context(), "userEmail") //We are checking of the user has login?
		if email != "" {
			err := errors.New("")
			myUser, err = user2.ReadUser(email)
			if err != nil {
				app.InfoLog.Println(err)
			}
		}
		departments, err := io.ReadDepartments()
		if err != nil {
			app.InfoLog.Println(err, "  reading department")
		}

		homePost := getBookDetails()
		type PageData struct {
			User            domain.User
			Departements    []domain.Department
			BookPostDetails []homePosts
		}
		data := PageData{myUser, departments, homePost}
		fmt.Println("voila we are in")
		files := []string{
			app.Path + "category.html",
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

//this method is just returning a list of the Image Ids
func getBookImageArray(bookImages []domain.BookImage) []string {
	valeus := []string{}

	for _, valeu := range bookImages {
		valeus = append(valeus, valeu.ImageId)
	}
	return valeus
}

/***
this method should collect all the information of one book and it two pictures
*/
func DetailsHandler(app *config.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var picture1 string
		var picture2 string
		var department domain.Department
		myUser := domain.User{}

		bookId := chi.URLParam(r, "bookId")
		if bookId == "" {
			app.Session.Put(r.Context(), "userMessage", "error_reading_book_details")
			http.Redirect(w, r, "/", 301)
			return
		}
		bookDepartment, err := book_io.ReadBookDepartment(bookId)
		if err != nil {
			app.InfoLog.Println(err, " reading bookDepartment")
		}
		if bookDepartment.DepartmentId != "" {
			department, err = io.ReadDepartment(bookDepartment.DepartmentId)
			if err != nil {
				app.InfoLog.Println(err, " reading department")
			}
		}

		bookPost, err := book_io.ReadWithBookId(bookId)
		if err != nil {
			app.InfoLog.Println(err, " reading bookPost")
		}
		book, err := book_io.ReadBook(bookId)
		if err != nil {
			app.InfoLog.Println(err, " reading book")
		}

		bookImage, err := book_io.ReadAllOfBookImage(bookId)
		if err != nil {
			app.InfoLog.Println(err)
		}

		pictures, err := picture_io.ReadAllOf(getBookImageArray(bookImage))
		if err != nil {
			app.InfoLog.Println(err, "  reading images")
		}
		for index, valeu := range pictures {
			if index == 0 {
				picture1 = valeu.Id
			} else {
				picture2 = valeu.Id
			}
		}
		location, err := location2.ReadLocation(bookPost.LocationId)
		if err != nil {
			app.InfoLog.Println(err, " Location")
		}

		user, err := user2.ReadUser(bookPost.Email)
		if err != nil {
			app.InfoLog.Println(err, " user")
		}
		type PageData struct {
			Book       domain.Book
			Department domain.Department
			Post       domain.BookPost
			Picture1   string
			Picture2   string
			User       domain.User
			BookOwner  domain.User
			Location   domain.Location
		}

		data := PageData{book, department, bookPost, picture1, picture2, myUser, user, location}

		//we need to check the user if has created an account first
		//email := app.Session.GetString(r.Context(), "userEmail")
		//if email == "" {
		//	app.Session.Put(r.Context(), "userMessage", "post_error_need_to_signup")
		//	http.Redirect(w, r, "/login", 301)
		//	return
		//}

		files := []string{
			app.Path + "book/single-book.html",
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

func NewPostHandler(app *config.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		myUser := domain.User{}
		cessionData := user.Message{}
		sessionType := app.Session.GetString(r.Context(), "userMessage")
		if sessionType != "" {
			cessionData = user.GetMessage(sessionType)
		}
		email := app.Session.GetString(r.Context(), "userEmail") //We are checking of the user has login?
		if email != "" {
			err := errors.New("")
			myUser, err = user2.ReadUser(email)
			if err != nil {
				app.InfoLog.Println(err)
			}
		}
		department, err := io.ReadDepartments()
		if err != nil {
			app.InfoLog.Println(err, "  error reading department")
		}
		language, err := language2.ReadLanguages()
		if err != nil {
			app.InfoLog.Println(err, "  error reading Language")
		}
		type PageData struct {
			PageMessage user.Message
			User        domain.User
			Departement []domain.Department
			Language    []domain.Language
		}
		data := PageData{cessionData, myUser, department, language}
		fmt.Println("voila we are in")
		files := []string{
			app.Path + "book/book_page.html",
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

func BookImageHandler(app *config.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("voila we are in")
		myUser := domain.User{}
		cessionData := user.Message{}
		sessionType := app.Session.GetString(r.Context(), "userMessage")
		//bookId := app.Session.GetString(r.Context(), "bookId")
		if sessionType != "" {
			cessionData = user.GetMessage(sessionType)
		}
		email := app.Session.GetString(r.Context(), "userEmail") //We are checking of the user has login?
		if email != "" {
			err := errors.New("")
			myUser, err = user2.ReadUser(email)
			if err != nil {
				app.InfoLog.Println(err)
			}
		}
		department, err := io.ReadDepartments()
		if err != nil {
			app.InfoLog.Println(err, "  error reading department")
		}
		type PageData struct {
			PageMessage user.Message
			User        domain.User
			Department  []domain.Department
		}
		data := PageData{cessionData, myUser, department}
		files := []string{
			app.Path + "book/book_image.html",
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

func LocationHandler(app *config.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		fmt.Println("voila we are in")

		myUser := domain.User{}
		cessionData := user.Message{}
		sessionType := app.Session.GetString(r.Context(), "userMessage")
		//bookId := app.Session.GetString(r.Context(), "bookId")
		if sessionType != "" {
			cessionData = user.GetMessage(sessionType)
		}
		email := app.Session.GetString(r.Context(), "userEmail") //We are checking of the user has login?
		if email != "" {
			err := errors.New("")
			myUser, err = user2.ReadUser(email)
			if err != nil {
				app.InfoLog.Println(err)
			}
		}
		type PageData struct {
			PageMessage user.Message
			User        domain.User
		}
		data := PageData{cessionData, myUser}

		files := []string{
			app.Path + "book/book_location.html",
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

func PostLocationHandler(app *config.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		/***
		The following are the data that we are retrieving from the cession
		-we have to put the correct key so that we can get what was saved in that key jsut like in queue data structure
		*/
		email := app.Session.GetString(r.Context(), "userEmail")
		bookId := app.Session.GetString(r.Context(), "bookId")
		//app.Session.Remove(r.Context(),"bookId")
		/***
		Here we are checking if the cession variables are empty, we redirect the user to /book/post link which is going to call the same HTML page.
		if both variables are not empty the program continues.
		*/
		if email == "" && bookId == "" {
			fmt.Println("missing either email or bookId", email, "<<<<<email||bookId>>>>", bookId)
			app.Session.Put(r.Context(), "userMessage", "post_error_need_to_signup")
			http.Redirect(w, r, "/book/post", 301)
			return
		}
		//Grabbing data from the HTML form
		r.ParseForm()
		description := r.PostFormValue("description")
		place := r.PostFormValue("place")
		placeDescription := r.PostFormValue("place_description")

		//This is the final, if all went well. we now need to create a post
		//todo we will need to implement map that can give us longitude an latitude
		if place != "" && placeDescription != "" {
			location := domain.Location{"", placeDescription, "", "", place}
			newLocation, err := location2.CreateLocation(location)
			if err != nil { // if an error occur
				app.InfoLog.Println(err, "there is an error when creating a post location")
				app.Session.Put(r.Context(), "userMessage", "post_error")
				http.Redirect(w, r, "/book/location", 301)
				return
			}
			post := domain.BookPost{"", email, bookId, time.Now(), newLocation.Id, "new", description}
			newPost, error := book_io.CreatBookPost(post)
			if error != nil { // if an error occur we delete what was created(rollback)
				_, err := location2.DeleteLocation(newLocation)
				if err != nil {
					app.InfoLog.Println(error, "there is an error when rolling back on location after failing to create location")
				}
				app.InfoLog.Println(error, "there is an error when creating a post location")
				app.Session.Put(r.Context(), "userMessage", "post_error")
				http.Redirect(w, r, "/book/location", 301)
				return
			} else { // This section executes when all the conditions in the ifs are OK.
				//todo create UserPost.
				myUserPost := domain.UserPost{newPost.Id, email}
				_, err := user2.CreateUserPost(myUserPost)
				if err != nil { //when an error when creating UserPost occurs
					_, err := location2.DeleteLocation(newLocation)
					if err != nil {
						app.InfoLog.Println(error, "there is an error when rolling back on location after failing to create location")
					}
					_, errr := book_io.DeleteBookPost(newPost)
					if errr != nil {
						app.InfoLog.Println(error, "there is an error when rolling back on bookPost after failing to create location")
					}
					app.InfoLog.Println(err, "  error creating UserPost")
					app.Session.Put(r.Context(), "userMessage", "post_error")
					http.Redirect(w, r, "/book/location", 301)
					return

				}
				app.Session.Put(r.Context(), "userMessage", "post_successful")
				http.Redirect(w, r, "/", 301)
				return
			}

		} else {
			fmt.Println("missing either place or placeDescription", place, "<<<<<place||placeDescription>>>>", placeDescription)
			//todo need to implement this cession error
			app.Session.Put(r.Context(), "userMessage", "post_error")
			http.Redirect(w, r, "/book/location", 301)
			return
		}

	}
}

func PostBookHandler(app *config.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		//we need to check the user if has created an account first
		email := app.Session.GetString(r.Context(), "userEmail")
		if email == "" {
			app.Session.Put(r.Context(), "userMessage", "post_error_need_to_signup")
			http.Redirect(w, r, "/login", 301)
			return
		}

		title := r.PostFormValue("title")
		language := r.PostFormValue("language")
		author := r.PostFormValue("author")
		edition := r.PostFormValue("edition")
		price, _ := strconv.ParseFloat(r.PostFormValue("price"), 64) // converting string into double

		fmt.Println(language, "<<< language")

		//Creating a book first
		if title != "" && language != "" && author != "" && edition != "" && price != 0.0 {
			book := domain.Book{"", title, language, edition, price, author}
			newBook, err := book_io.CreateBook(book)
			if err != nil { //if an error occurs we interrupt everything here and return an error message to the user
				app.Session.Put(r.Context(), "userMessage", "post_error")
				http.Redirect(w, r, "/book/", 301)
				return
			} else { // if all Good. we now create BookLanguage
				bookLanguage := domain.BookLanguage{book.Id, language}
				_, err := book_io.CreateBookLanguage(bookLanguage)
				if err != nil {
					app.ErrorLog.Println(err.Error(), " error creating book language")
					app.Session.Put(r.Context(), "userMessage", "post_error")
					http.Redirect(w, r, "/book/", 301)
					return
				}
				// we put the book id on the cession. final good condition
				app.Session.Put(r.Context(), "bookId", newBook.Id)
				http.Redirect(w, r, "/book/book_Image", 301)
				return
			}
		} else { //if one of the field is empty this should happen
			app.Session.Put(r.Context(), "userMessage", "post_empty_error")
			http.Redirect(w, r, "/book/", 301)
			return
		}

	}
}
func PostBookImage(app *config.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		//var status bool
		//we need to check the user if has created an account first
		email := app.Session.GetString(r.Context(), "userEmail")
		bookId := app.Session.GetString(r.Context(), "bookId")
		if email == "" {
			app.Session.Put(r.Context(), "userMessage", "post_error_need_to_signup")
			http.Redirect(w, r, "/login", 301)
			return
		}
		fmt.Println("voila we are in book PostBookHandler")
		picture := domain.Picture{}
		r.ParseForm()
		fmt.Println(" reading the file")

		file, _, err := r.FormFile("file")
		file1, _, err := r.FormFile("file1")
		department := r.PostFormValue("departmentId")
		imagedescription := r.PostFormValue("imagedescription")
		fmt.Println("Department: ", department, " and imagedescription: ", imagedescription)
		fmt.Println(" read successful")
		if err != nil {
			fmt.Println(err, "<<<<<<>>>>>>>")
			app.Session.Put(r.Context(), "userMessage", "post_image_error")
			http.Redirect(w, r, "/book/book_Image", 301)
			return
		}
		//fmt.Println(" converting to []byte and into slice array", handler)
		reader := bufio.NewReader(file)
		reader1 := bufio.NewReader(file1)
		content, _ := ioutil.ReadAll(reader)
		content1, _ := ioutil.ReadAll(reader1)

		//converting the file into an slice of byte
		sliceOfImage := [][]byte{content, content1}

		for index, pic := range sliceOfImage {
			fmt.Println("index: ", index)
			if index == 0 {
				picture, err = sendingPicture(pic, "fist picture")
			} else {
				picture, err = sendingPicture(pic, "second picture")
			}
			if err != nil { //if an error occurs when creating a picture here we return an error message and we delete the book and it picture
				fmt.Println(err, "<<<<<<picture>>>>>>>")
				app.Session.Put(r.Context(), "userMessage", "post_image_error")
				http.Redirect(w, r, "/book/book_Image", 301)
				return
			} else { //in case when there is no error when creating a picture, now we will need to create BookImageHandler
				if bookId != "" {
					bookImage := domain.BookImage{picture.Id, bookId, imagedescription}
					fmt.Println("bookImage: ", bookImage)
					newBookImage, err := book_io.CreatBookImage(bookImage)
					if err != nil { //if could not create a bookImage we rollback
						app.InfoLog.Println(err, "newBookImage fail")
						_, errr := picture_io.DeletePicture(picture)
						if errr != nil {
							app.InfoLog.Println(err, " when trying to delete  Picture: ", picture)
						}
						app.Session.Put(r.Context(), "userMessage", "post_error")
						http.Redirect(w, r, "/book/book_Image", 301)
						return
					} else {
						bookDepartment := domain.BookDepartment{bookId, department, imagedescription}
						_, err := book_io.CreateBookdepartment(bookDepartment)
						if err != nil {
							app.InfoLog.Println(err, "bookDepartment fail")
							_, errr := picture_io.DeletePicture(picture)
							if errr != nil {
								app.InfoLog.Println(err, " when trying to delete  Picture: ", picture)
							}
							_, err := book_io.DeleteBookImage(newBookImage)
							if err != nil {
								app.InfoLog.Println(err, " when trying to delete  bookPicture: ", newBookImage)
							}
							app.Session.Put(r.Context(), "userMessage", "post_error")
							http.Redirect(w, r, "/book/book_Image", 301)
							return
						}
					}
				} else {
					fmt.Println("Could not get BookId in the cession")
					app.Session.Put(r.Context(), "userMessage", "post_error")
					http.Redirect(w, r, "/book/book_Image", 301)
					return
				}
			}
		}
		// if all go well
		fmt.Println("bookImage creation Successful: ")
		app.Session.Put(r.Context(), "userMessage", "post_successful")
		http.Redirect(w, r, "/book/location", 301)
		return

	}
}

func sendingPicture(picture []byte, desc string) (domain.Picture, error) {
	newPicture := domain.Picture{"", picture, desc}
	result, err := picture_io.CreatePicture(newPicture)
	return result, err
}

type homePostCategories struct {
	Book       domain.Book
	Picture    domain.Picture
	Post       domain.BookPost
	Department domain.Department
}

func getBookOfDepartment(departmentId string) []homePostCategories {
	entity := []homePostCategories{}
	image := domain.Picture{}
	department := domain.Department{}

	bookDepartments, err := book_io.ReadAllOfBookDepartment(departmentId)
	if err != nil {
		fmt.Println(err, " there is an error when reading all the bookPosts")
		return entity
	}

	for _, bookDepartment := range bookDepartments {
		book, err := book_io.ReadBook(bookDepartment.BookId)
		if err != nil {
			fmt.Println(err, " there is an error when reading book")
		}

		bookImage, err := book_io.ReadBookImageWithBookId(book.Id)
		if err != nil {
			fmt.Println(err, " there is an error when reading all the bookImage")
		} else {
			image, err = picture_io.ReadFirstPicture(bookImage.ImageId)
			//fmt.Println(image.Id)
			if err != nil {
				fmt.Println(bookImage, "<<<<<<<bookImage")
				fmt.Println(image, "<<<<<<Image")
				fmt.Println(err, " there is an error when reading all the Image")
			}
		}
		post, err := book_io.ReadWithBookId(bookDepartment.BookId)
		if err != nil {
			fmt.Println(err, " there is an error when reading all the bookPosts")
			return entity
		}

		//fmt.Println(post.LocationId)
		newEntity := homePostCategories{book, image, post, department}
		entity = append(entity, newEntity)
		image = domain.Picture{}
		department = domain.Department{}
		newEntity = homePostCategories{}
	}
	return entity
}

//Thi is the type that we are going to return the home category
type homePosts struct {
	Book       domain.Book
	Picture    domain.Picture
	Post       domain.BookPost
	Location   domain.Location
	Department domain.Department
	User       domain.User
}

//This method returns all the posts so that we can populate the home category page
func getBookDetails() []homePosts {
	entity := []homePosts{}
	image := domain.Picture{}
	department := domain.Department{}

	// we first check all the posts
	posts, err := book_io.ReadBookPosts()
	if err != nil {
		fmt.Println(err, " there is an error when reading all the bookPosts")
		return entity
	}
	for _, post := range posts {
		book, err := book_io.ReadBook(post.BookId)
		if err != nil {
			fmt.Println(err, " there is an error when reading book")
		}

		location, err := location2.ReadLocation(post.LocationId)
		if err != nil {
			fmt.Println(err, " there is an error when reading all the location")
		}

		bookImage, err := book_io.ReadBookImageWithBookId(book.Id)
		if err != nil {
			fmt.Println(err, " there is an error when reading all the bookImage")
		} else {
			image, err = picture_io.ReadFirstPicture(bookImage.ImageId)
			//fmt.Println(image.Id)
			if err != nil {
				fmt.Println(bookImage, "<<<<<<<bookImage")
				fmt.Println(image, "<<<<<<Image")
				fmt.Println(err, " there is an error when reading all the Image")
			}
		}
		user, err := user2.ReadUser(post.Email)
		if err != nil {
			fmt.Println(err, " there is an error when reading all the user")
		}
		bookdepartment, err := book_io.ReadBookDepartment(book.Id)
		if err != nil {
			fmt.Println(err, " there is an error when reading all the bookdepartment")
		} else {
			department, err = io.ReadDepartment(bookdepartment.DepartmentId)
			if err != nil {
				fmt.Println(err, " there is an error when reading all the bookdepartment")
			}
		}
		//fmt.Println(post.LocationId)
		newEntity := homePosts{book, image, post, location, department, user}
		entity = append(entity, newEntity)
		image = domain.Picture{}
		department = domain.Department{}
		newEntity = homePosts{}
	}
	return entity
}

//This method should return all the post of a user
func getUserBookDetails(id string) []homePosts {
	entity := []homePosts{}
	image := domain.Picture{}
	department := domain.Department{}
	posts := []domain.BookPost{}
	/***
	  We first read all of the post of one user
	*/
	userPosts, err := user2.ReadAllOfUserPost(id)
	if err != nil {
		fmt.Println(err, " there is an error when reading all the userPosts")
		return entity
	}
	// we secondly check all the posts of that user if he/she has one.
	for _, post := range userPosts {
		myposts, err := book_io.ReadBookPost(post.PostId)
		if err != nil {
			fmt.Println(err, " there is an error when reading all the bookPosts")
			return entity
		}
		posts = append(posts, myposts)
	}

	//Now we loop through posts slice to get all the details of
	for _, post := range posts {
		book, err := book_io.ReadBook(post.BookId)
		if err != nil {
			fmt.Println(err, " there is an error when reading book")
		}

		location, err := location2.ReadLocation(post.LocationId)
		if err != nil {
			fmt.Println(err, " there is an error when reading all the location")
		}

		bookImage, err := book_io.ReadBookImageWithBookId(book.Id)
		if err != nil {
			fmt.Println(err, " there is an error when reading all the bookImage")
		} else {
			image, err = picture_io.ReadFirstPicture(bookImage.ImageId)
			//fmt.Println(image.Id)
			if err != nil {
				fmt.Println(bookImage, "<<<<<<<bookImage")
				fmt.Println(image, "<<<<<<Image")
				fmt.Println(err, " there is an error when reading all the Image")
			}
		}
		user, err := user2.ReadUser(post.Email)
		if err != nil {
			fmt.Println(err, " there is an error when reading all the user")
		}
		bookdepartment, err := book_io.ReadBookDepartment(book.Id)
		if err != nil {
			fmt.Println(err, " there is an error when reading all the bookdepartment")
		} else {
			department, err = io.ReadDepartment(bookdepartment.DepartmentId)
			if err != nil {
				fmt.Println(err, " there is an error when reading all the bookdepartment")
			}
		}
		//fmt.Println(post.LocationId)
		newEntity := homePosts{book, image, post, location, department, user}
		entity = append(entity, newEntity)
		image = domain.Picture{}
		department = domain.Department{}
		newEntity = homePosts{}
	}
	return entity
}
