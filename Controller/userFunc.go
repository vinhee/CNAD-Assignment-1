package Controller

import (
	"crypto/rand"
	"encoding/base64"
	"html/template"
	"log"
	"net/http"
	"net/mail"
	"regexp"
	"strconv"
	"time"

	database "CNAD-Assignment-1/Database"

	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

func generateSecretKey() string {
	key := make([]byte, 32) // generates a random 32 byte
	if _, err := rand.Read(key); err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(key)
}

var store = sessions.NewCookieStore([]byte(generateSecretKey()))

func Logout(w http.ResponseWriter, r *http.Request) {
	nameSess, _ := store.Get(r, "cookieName")
	nameSess.Values["cookieName"] = nil
	emailSess, _ := store.Get(r, "cookieEmail")
	emailSess.Values["cookieEmail"] = nil
	idSess, _ := store.Get(r, "cookieID")
	idSess.Values["cookieID"] = nil

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func Loginpage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("Pages/UserManage/loginpage.html", "Pages/navbar.html")
		if err != nil {
			log.Println("Error executing template:", err)
			http.Error(w, "Error loading template", http.StatusInternalServerError)
			return
		}
		err = tmpl.ExecuteTemplate(w, "LoginPage", nil)
		if err != nil {
			log.Println("Error executing template:", err)
			http.Error(w, "Error executing template", http.StatusInternalServerError)
			return
		}
		return
	}
	if r.Method == http.MethodPost {
		var errMsg string
		r.ParseForm()
		checkEmail := r.FormValue("userEmail")
		checkPassword := r.FormValue("userPassword")
		log.Println("Email:", checkEmail)
		log.Println("Password:", checkPassword)

		userList, err := database.GetLoginUser()
		if err != nil {
			log.Println("Error retrieving users:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		userFound := false
		if checkEmail != "" && checkPassword != "" {
			for _, checkUser := range userList {
				if checkUser.Email == checkEmail {
					err := bcrypt.CompareHashAndPassword([]byte(checkUser.Password), []byte(checkPassword))
					log.Printf("%x", []byte(checkUser.Password))
					log.Printf("%x", []byte(checkPassword))
					if err == nil {
						userFound = true
						sessionName, _ := store.Get(r, "cookieName")
						sessionName.Values["userName"] = checkUser.Name
						sessionEmail, _ := store.Get(r, "cookieEmail")
						sessionEmail.Values["userEmail"] = checkEmail
						sessionID, _ := store.Get(r, "cookieID")
						sessionID.Values["userID"] = checkUser.Id
						sessionName.Save(r, w)
						sessionEmail.Save(r, w)
						sessionID.Save(r, w)
						http.Redirect(w, r, "/homemember", http.StatusSeeOther)
						return
					}
				}
			}
			if !userFound {
				errMsg = "Incorrect email/phone number or password, try again!"
				LoginError(w, errMsg)
				return
			}
		}
	}
}

func LoginError(w http.ResponseWriter, errMsg string) {
	tmpl, err := template.ParseFiles("Pages/UserManage/loginpage.html", "Pages/navbar.html")
	if err != nil {
		log.Println("Error parsing template:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	data := struct {
		ErrMsg string
		Navbar interface{}
	}{
		ErrMsg: errMsg,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	log.Println("Error message:", errMsg)
	if err := tmpl.ExecuteTemplate(w, "LoginPage", data); err != nil {
		log.Println("Error executing template:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func HomeMember(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "cookieName")
	userName, ok := session.Values["userName"].(string)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	}

	tmpl, err := template.ParseFiles("Pages/UserManage/homemember.html", "Pages/UserManage/navbarmember.html")
	if err != nil {
		log.Println("Error parsing template:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	data := struct {
		UserName  string
		NavMember interface{} // needed as navbar is accessing the username from homemember.html
	}{
		UserName: userName,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, data); err != nil {
		log.Println("Error executing template:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func Registerpage(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("Pages/UserManage/registerpage.html", "Pages/navbar.html")
		if err != nil {
			log.Println("Error executing template:", err)
			http.Error(w, "Error loading template", http.StatusInternalServerError)
			return
		}
		err = tmpl.ExecuteTemplate(w, "RegPage", nil)
		if err != nil {
			log.Println("Error executing template:", err)
			http.Error(w, "Error executing template", http.StatusInternalServerError)
			return
		}
		return
	}
	if r.Method == http.MethodPost {
		log.Println("Parsed form values:", r.Form)
		checkEmail := r.FormValue("userEmail")
		validEmail := isEmail(checkEmail)
		validNum := isPhoneNumber(checkEmail)

		userList, err := database.GetLoginUser()
		if err != nil {
			log.Println("Error retrieving users:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		var errMsg string
		var successMsg string
		var user database.User
		if r.Method == "POST" {
			if validEmail || validNum {
				for _, checkUser := range userList {
					if checkUser.Email == checkEmail {
						errMsg = "This email/phone number already has an account!"
						RegError(w, errMsg, "")
						return
					}
				}
			} else {
				errMsg = "This email/phone number is invalid! Try again"
				RegError(w, errMsg, "")
				return
			}
			user.Email = checkEmail
			user.Name = r.FormValue("userName")
			user.Password = r.FormValue("userPassword")

			hash, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
			if err != nil {
				log.Println("Hash error:", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			user.Password = string(hash)

			if err := database.InsertNewUser(user); err != nil {
				log.Println("Error inserting new user:", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
				return
			}
			successMsg = "Account created successfully!"
			RegSuccess(w, successMsg, "")
		}
	}
}

func isEmail(input string) bool {
	_, err := mail.ParseAddress(input)
	return err == nil
}

func isPhoneNumber(input string) bool {
	phoneRegex := `^(9|8|6)\d{7}$`
	re := regexp.MustCompile(phoneRegex)
	return re.MatchString(input)
}

func RegError(w http.ResponseWriter, errMsg string, successMsg string) {
	tmpl, err := template.ParseFiles("Pages/UserManage/registerpage.html", "Pages/navbar.html")
	if err != nil {
		log.Println("Error parsing template:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	data := struct {
		ErrMsg     string
		Navbar     interface{}
		SuccessMsg string
	}{
		ErrMsg:     errMsg,
		SuccessMsg: successMsg,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	log.Println("Error message:", errMsg)
	if err := tmpl.ExecuteTemplate(w, "RegPage", data); err != nil {
		log.Println("Error executing template:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func RegSuccess(w http.ResponseWriter, successMsg string, errMsg string) {
	tmpl, err := template.ParseFiles("Pages/UserManage/registerpage.html", "Pages/navbar.html")
	if err != nil {
		log.Println("Error parsing template:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	data := struct {
		SuccessMsg string
		Navbar     interface{}
		ErrMsg     string
	}{
		SuccessMsg: successMsg,
		ErrMsg:     errMsg,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	log.Println("Successful message:", successMsg)
	if err := tmpl.ExecuteTemplate(w, "RegPage", data); err != nil {
		log.Println("Error executing template:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func ProfilePage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("Pages/UserManage/profilepage.html", "Pages/UserManage/navbarmember.html")
	if err != nil {
		log.Println("Error parsing template:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	session, _ := store.Get(r, "cookieEmail")
	userEmail, _ := session.Values["userEmail"].(string)
	idSess, _ := store.Get(r, "cookieID")
	userID, _ := idSess.Values["userID"].(int)
	userList, err := database.GetUserDetail(userEmail)
	var userName string
	var userTier string
	var userBooking int
	for _, checkUser := range userList {
		if checkUser.Email == userEmail {
			userName = checkUser.Name
			userTier = checkUser.MemberTier
			userBooking = checkUser.Bookings
		}
	}

	bookList, err := database.GetUserBook(userID)
	if err != nil {
		log.Println("Error getting user booking:", err)
	}

	todayDate := time.Now()

	data := struct {
		UserName    string
		UserEmail   string
		UserTier    string
		UserBooking int
		BookList    []database.CarsBooking
		TodayDate   time.Time
		UserID      int
	}{
		UserName:    userName,
		UserEmail:   userEmail,
		UserTier:    userTier,
		UserBooking: userBooking,
		BookList:    bookList,
		TodayDate:   todayDate,
		UserID:      userID,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.ExecuteTemplate(w, "ProfPage", data); err != nil {
		log.Println("Error executing template:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}

}

func EditProfile(w http.ResponseWriter, r *http.Request) {
	var successMsg string
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("Pages/UserManage/editprofile.html", "Pages/UserManage/navbarmember.html")
		if err != nil {
			log.Println("Error parsing template:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		emailSess, _ := store.Get(r, "cookieEmail")
		userEmail, _ := emailSess.Values["userEmail"].(string)
		nameSess, _ := store.Get(r, "cookieName")
		userName, _ := nameSess.Values["userName"].(string)

		userList, err := database.GetUserDetail(userEmail)
		var userTier string
		var userPassword string
		for _, checkUser := range userList {
			if checkUser.Email == userEmail {
				userTier = checkUser.MemberTier
				userPassword = checkUser.Password
			}
		}

		data := struct {
			UserName     string
			UserEmail    string
			UserTier     string
			UserPassword string
			SuccessMsg   string
		}{
			UserName:     userName,
			UserEmail:    userEmail,
			UserPassword: userPassword,
			UserTier:     userTier,
			SuccessMsg:   "",
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		log.Println("User profile name", data.UserName)
		if err := tmpl.ExecuteTemplate(w, "EditProfile", data); err != nil {
			log.Println("Error executing template:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
	if r.Method == http.MethodPost {
		r.ParseForm()
		userEmail := r.FormValue("userEmail")
		userPassword := r.FormValue("userPassword")
		userName := r.FormValue("userName")
		userTier := r.FormValue("userTier")
		hash, err := bcrypt.GenerateFromPassword([]byte(userPassword), bcrypt.DefaultCost)
		if err != nil {
			log.Println("Hash error:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		userPassword = string(hash)
		if err := database.UpdateUser(userName, userEmail, userPassword, userTier); err != nil {
			log.Println("Error inserting new user:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		successMsg = "Account updated successfully!"

		tmpl, err := template.ParseFiles("Pages/UserManage/editprofile.html", "Pages/UserManage/navbarmember.html")
		if err != nil {
			log.Println("Error parsing template:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		data := struct {
			UserName     string
			UserEmail    string
			UserTier     string
			UserPassword string
			SuccessMsg   string
		}{
			UserName:     userName,
			UserEmail:    userEmail,
			UserPassword: userPassword,
			UserTier:     userTier,
			SuccessMsg:   successMsg,
		}

		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		log.Println("User  profile name", data.UserName)
		if err := tmpl.ExecuteTemplate(w, "EditProfile", data); err != nil {
			log.Println("Error executing template:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}

func CancelBooking(w http.ResponseWriter, r *http.Request) {
	bookingIDstr := r.FormValue("bookingID")
	bookingID, _ := strconv.Atoi(bookingIDstr)
	database.UpdateCancelled(bookingID)
	database.IncreaseBook(bookingID)
	booking, _ := database.GetBookingByID(bookingID)
	log.Print("UserID is: ", booking.UserID)
	database.IncreaseBook(booking.UserID)
	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}
