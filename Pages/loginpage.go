package Pages

import (
	"fmt"
	"net/http"
)

func Loginpage(w http.ResponseWriter, r *http.Request) {
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
		</style>
	</head>
	<body>
	  ` + Navbar() + ` 
	<section class="vh-100">
  	<div class="container-fluid">
    <div class="row">
      <div class="col-sm-6 text-black">

        <div class="d-flex align-items-center h-custom-2 px-5 ms-xl-4 mt-5 pt-5 pt-xl-0 mt-xl-n5">

          <form style="width: 23rem;">

            <h3 class="fw-normal mb-3 pb-3" style="letter-spacing: 1px; font-family: Oswald, sans serif;">Log in</h3>

            <div data-mdb-input-init class="form-outline mb-4">
              <input type="email" id="form2Example18" class="form-control form-control-lg" />
              <label class="form-label" for="form2Example18">Email address</label>
            </div>

            <div data-mdb-input-init class="form-outline mb-4">
              <input type="password" id="form2Example28" class="form-control form-control-lg" />
              <label class="form-label" for="form2Example28">Password</label>
            </div>

            <div class="pt-1 mb-4">
              <button data-mdb-button-init data-mdb-ripple-init class="btn btn-info btn-lg btn-block" type="button">Login</button>
            </div>

            <p class="small mb-5 pb-lg-2"><a class="text-muted" href="#!">Forgot password?</a></p>
            <p>Don't have an account? <a href="#!" class="link-info">Register here</a></p>

          </form>

        </div>

      </div>
      <div class="col-sm-6 px-0 d-none d-sm-block">
        <img src="https://a.storyblok.com/f/85281/1080x1440/2af3cc39d1/how-do-you-charge-an-electric-car__article_v1_header_vertical_3-4_mobile.png"
          alt="Login image" class="w-100 vh-100" style="object-fit: cover; object-position: left;">
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
