package Pages

func NavbarMember(userName string) string {
	return `
  <head>
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
  <link href="https://fonts.googleapis.com/css2?family=Oswald:wght@200..700&display=swap" rel="stylesheet">
  <style>
    .navbar-brand {
    letter-spacing: 3px;
    color: #232b47;
    font-family: 'Oswald', sans-serif;
    font-size: 2rem;
  }

  .nav-link {
    letter-spacing: 3px;
    color: #543310;
    font-family: 'Oswald', sans-serif;
    font-size: 1rem;
    margin: 4px;
  }

  .navbar-brand:hover {
    color: #232b47;
  }

  .navbar-scroll .nav-link,
  .navbar-scroll .fa-bars {
    color: #543310;
    margin-top: 10px;
  }

  .navbar-scrolled .nav-link,
  .navbar-scrolled .fa-bars {
    color: #543310;
    margin-top: 10px;
  }

  .navbar-scrolled {
    background-color: #F8EDE3;
  }
  </style>
  </head>
  <nav class="navbar navbar-expand-lg fixed-top navbar-scroll shadow-0" style="background-color: #F8EDE3;">
    <div class="container">
      <a class="navbar-brand" href="#">ElecShare</a>
      <button class="navbar-toggler ps-0" type="button" data-mdb-toggle="collapse"
        aria-expanded="false" aria-label="Toggle navigation">
        <span class="d-flex justify-content-start align-items-center">
          <i class="fas fa-bars"></i>
        </span>
      </button>
      <div class="collapse navbar-collapse" id="navbarExample01">
        <ul class="navbar-nav me-auto mb-2 mb-lg-0">
          <li class="nav-item active">
            <a class="nav-link px-3" href="/homepage">Home</a>
          </li>
          <li class="nav-item">
            <a class="nav-link px-3" href="/vehiclepage">Vehicle Service</a>
          </li>
          <li class="nav-item">
            <a class="nav-link px-3" href="/billpage">Billing</a>
          </li>
        </ul>
        <ul class="navbar-nav ms-auto mb-2 mb-lg-0">
          <li class="nav-item">
            <a class="nav-link pe-3" href="/profile">
              <i class="fa-solid fa-user"></i> ` + userName + `
            </a>
          </li>
          <li class="nav-item">
            <a class="nav-link pe-3" href="/homepage">
              Logout
            </a>
          </li>
        </ul>
      </div>
    </div>
  </nav>
  <script src="https://kit.fontawesome.com/13df13ab87.js" crossorigin="anonymous"></script>
    `
}