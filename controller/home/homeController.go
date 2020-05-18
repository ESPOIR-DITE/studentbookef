package home

import (
	"errors"
	"fmt"
	"github.com/go-chi/chi"
	"html/template"
	"net/http"
	"studentbookef/config"
	"studentbookef/controller/user"
	"studentbookef/domain"
	"studentbookef/io"
	"studentbookef/io/book_io"
	location2 "studentbookef/io/location"
	"studentbookef/io/picture_io"
	user2 "studentbookef/io/user"
)

func Home(app *config.Env) http.Handler {
	r := chi.NewRouter()
	r.Get("/", homeHandler(app))

	return r
}

func homeHandler(app *config.Env) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		myUser := domain.User{}
		cessionData := user.Message{}
		var confirmYourRedistration string
		//just_login_message:=user.Message{}
		//just_login:=app.Session.GetString(r.Context(),"just_login") 		// we checking if the user just login?
		//if just_login!=""{
		//	just_login_message=user.GetMessage(just_login)
		//}
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
			if myUser.PhoneNumber == "" || myUser.Surname == "" {
				confirmYourRedistration = "Please consider to confirm your registration from your email"
			}
		}
		// read all the posts deatils and images
		getBookDetails := getBookDetails()
		type PageData struct {
			PageMessage       user.Message
			User              domain.User
			BookPostDetails   []homePosts
			RegistrationNotic string
		}
		data := PageData{cessionData, myUser, getBookDetails, confirmYourRedistration}
		files := []string{
			app.Path + "index.html",
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

type homePosts struct {
	Book       domain.Book
	Picture    domain.Picture
	Post       domain.BookPost
	Location   domain.Location
	Department domain.Department
	User       domain.User
}

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
