#Controller
This the entry point of requests coming from the browser.  
The following are activities that classes of this package do:
* Call the web pages  
   
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
		
* Send the data to an HTML page    

    
        type PageData struct {  //creating an Object that you will send to the page
			User domain.User
		}
		data := PageData{user}  //instantiate that object
		 
		 err = ts.Execute(w, data) //send it
         		if err != nil {
         			app.ErrorLog.Println(err.Error())
         		}
         		
* Collect data from an HTML page from a form 
    
        r.ParseForm()           //Now we grabbing the contents of the form by call the name of the input(html)
    		name := r.PostFormValue("name") // "name" id the name of the field in the form
    		email := r.PostFormValue("email")
    		

* collect data from URL 
    
        r.Get("/userAccount/{code}", userAccountHandler(app))
        code := chi.URLParam(r, "code")