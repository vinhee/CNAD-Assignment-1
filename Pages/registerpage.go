package Pages

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/mail"
	"os"
	"regexp"

	_ "github.com/go-sql-driver/mysql"
)

type user struct {
	Id         int    `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	Password   string `json:"password"`
	MemberTier string `json:"memberTier"`
}

func Registerpage(w http.ResponseWriter, r *http.Request) {
	var errMsg string
	var successfulMsg string
	r.ParseForm()

	dbUser, dbPassword, dbHost, dbName := os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_HOST"), os.Getenv("DB_NAME")
	dsn := fmt.Sprintf("%s:%s@tcp(%s)/%s", dbUser, dbPassword, dbHost, dbName)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Println("Database connection error:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	query := "SELECT * FROM Users"
	results, err := db.Query(query)
	if err != nil {
		log.Println("Database query error:", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer results.Close()

	userList := []user{}
	for results.Next() {
		var user user
		if err := results.Scan(&user.Id, &user.Name, &user.Email, &user.Password, &user.MemberTier); err != nil {
			log.Println("Row scan error:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		userList = append(userList, user)
	}
	if r.Method == "POST" {
		var user user
		log.Println("Parsed form values:", r.Form)
		checkEmail := r.FormValue("userEmail")
		validEmail := isEmail(checkEmail)
		validNum := isPhoneNumber(checkEmail)

		if validEmail || validNum {
			for _, checkCourse := range userList {
				if checkCourse.Email == checkEmail {
					errMsg = "This email/phone number already has an account!"
				}
			}
		} else {
			errMsg = "This email/phone number is invalid! Try again"
		}
		user.Email = checkEmail
		user.Name = r.FormValue("userName")
		user.Password = r.FormValue("userPassword")

		insertQuery := "INSERT INTO Users (Name, Email, Password, MemberTier) VALUES (?, ?, ?, ?)"
		_, err = db.Exec(insertQuery, user.Name, user.Email, user.Password, "Basic")
		if err != nil {
			log.Println("Database insert error:", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		if errMsg == "" {
			successfulMsg = "Account Created Successfully!"
		}

	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	html := `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<link href="https://stackpath.bootstrapcdn.com/bootstrap/5.1.3/css/bootstrap.min.css" rel="stylesheet">
		<link href="https://cdnjs.cloudflare.com/ajax/libs/mdb-ui-kit/6.4.0/mdb.min.css" rel="stylesheet">
		<link rel="preconnect" href="https://fonts.googleapis.com">
		<link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
		<link href="https://fonts.googleapis.com/css2?family=Oswald:wght@200..700&display=swap" rel="stylesheet">
		<title>ElecShare</title>
		<style>
			.container2 {
				max-width: 600px;
				margin: auto;
			}
			.form-group {
				margin-bottom: 15px;
			}
			.form-outline {
				position: relative;
				border: 2px solid #232b47;
				border-radius: 0.5rem;
				padding: 0;
			}
			.form-control {
				padding: 20px 15px 5px;
				border: none;
				outline: none; 
				border-radius: 0.5rem; 
			}
			.label {
				color: #232b47;
				margin-bottom: 5px;
				margin-left: 5px;
			}
			.form-label {
				position: absolute;
				top: 10px; 
				left: 15px; 
				background-color: white;
				padding: 0 5px; 
				z-index: 1;
			}
			.bg-image-vertical {
				position: relative;
				overflow: hidden;
				background-repeat: no-repeat;
				background-position: right center;
				background-size: auto 100%;
			}

			@media (min-width: 1025px) {
				.h-custom-2 {
				height: 100%;
			}
			.loginLine {
				font-family: 'Oswald', sans-serif;
				color: #232b47;
			}
		</style>
	</head>
	<body>
	  ` + Navbar() + ` 
	<section class="vh-100">
  	<div class="container-fluid">
    <div class="row">
	<div class="col-sm-6 text-black">

        <div class="d-flex align-items-center h-custom-2 px-5 ms-xl-4 mt-5 pt-7 pt-xl-0 mt-xl-n5">

          <form style="width: 23rem;" method="POST" action="/register">

            <h3 class="fw-normal mb-3 pb-1" style="letter-spacing: 1px; font-family: Oswald, sans serif;">Register</h3>
		`
	if successfulMsg != "" {
		html += "<p class=\"label\" style=\"color:green\">" + successfulMsg + "</p>"
	}

	html += `
			<p class="label">Name</p>
            <div data-mdb-input-init class="form-outline mb-2">
              <input type="name" id="name" name="userName" class="form-control form-control-lg" />
            </div>
		
			<p class="label">Email Address/Phone Number</p>
		`
	if errMsg != "" {
		html += "<p class=\"label\" style=\"color:red\">" + errMsg + "</p>"
	}

	html += `
            <div data-mdb-input-init class="form-outline mb-2">
              <input type="emailNum" id="email" name="userEmail" class="form-control form-control-lg" />
            </div>

			<p class="label">Password</p>
			<div data-mdb-input-init class="form-outline mb-2">
              <input type="password" id="password" name="userPassword" class="form-control form-control-lg" />
            </div>

            <div class="pt-1 mb-4 loginLine">
              <button action="/registerpage" method="POST" data-mdb-button-init data-mdb-ripple-init class="btn btn-info btn-lg btn-block" style="background-color:#232b47" type="submit">Register</button>
            </div>

            <p class="loginLine">Have an account? <a href="/login" class="link-info">Login here</a></p>

          </form>

        </div>

      </div>
		<div class="col-sm-6 px-0 d-none d-sm-block">
				<img src="https://a.storyblok.com/f/85281/1080x1440/2af3cc39d1/how-do-you-charge-an-electric-car__article_v1_header_vertical_3-4_mobile.png"
				alt="Register image" class="w-100 vh-100" style="object-fit: cover; object-position: right;"></img>
				</div>
		</div>
		</div>
	</section>
		<script src="https://code.jquery.com/jquery-3.6.0.min.js"></script>
		<script src="https://cdn.jsdelivr.net/npm/@popperjs/core@2.9.2/dist/umd/popper.min.js"></script>
		<script src="https://stackpath.bootstrapcdn.com/bootstrap/5.1.3/js/bootstrap.min.js"></script>	
		</body>
	</html>
`
	fmt.Fprint(w, html)
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
