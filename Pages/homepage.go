package Pages

import (
	"fmt"
	"net/http"
)

func Homepage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	html := `
	<!DOCTYPE html>
	<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<link href="https://stackpath.bootstrapcdn.com/bootstrap/5.1.3/css/bootstrap.min.css" rel="stylesheet">
		<link href="https://cdnjs.cloudflare.com/ajax/libs/mdb-ui-kit/6.4.0/mdb.min.css" rel="stylesheet">
		<title>ElecShare</title>
		<style>
			.container2 {
				max-width: 600px;
				margin: auto;
			}
			.form-group {
				margin-bottom: 15px;
			}
		</style>
	</head>
	<body>
	  ` + Navbar() + ` 
		<div class="container2">
		</div>
		<script src="https://code.jquery.com/jquery-3.6.0.min.js"></script>
		<script src="https://cdn.jsdelivr.net/npm/@popperjs/core@2.9.2/dist/umd/popper.min.js"></script>
		<script src="https://stackpath.bootstrapcdn.com/bootstrap/5.1.3/js/bootstrap.min.js"></script>	
	</body>
	</html>
	`
	fmt.Fprint(w, html)
}
