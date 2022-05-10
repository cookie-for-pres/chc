package main

import "github.com/cookie-for-pres/chc"

func LoginController(req *chc.Request) *chc.Response {
	res := req.NewResponse()
	if req.Method == "GET" {
		err := res.LoadHtmlFile("public/login.html")
		if err != nil {
			res.SetStatusCode(500)
			res.SetStringBody(err.Error())

			return res
		}

		return res
	} else {
		formData, err := req.FormDataBody()
		if err != nil {
			res.SetStatusCode(500)
			res.SetStringBody(err.Error())

			return res
		}

		if formData == nil {
			res.SetStatusCode(500)
			res.SetStringBody("No form data")

			return res
		}

		if formData["username"] != "cookie" || formData["password"] != "cookie" {
			res.SetStatusCode(500)
			res.SetStringBody("Invalid username or password")
		}

		res.SetStatusCode(301)
		res.SetCookie("loggedin", "true")
		res.SetRedirect("/dashboard")

		return res
	}
}

func DashboardController(req *chc.Request) *chc.Response {
	res := req.NewResponse()

	if req.Cookies["loggedin"] != "true" {
		res.SetStatusCode(500)
		res.SetStringBody("Not logged in")
	} else {
		res.SetStatusCode(200)
		res.LoadHtmlFile("public/dashboard.html")
	}

	return res
}

func DuckController(req *chc.Request) *chc.Response {
	res := req.NewResponse()
	if req.Method == "GET" {
		res.SetStatusCode(200)
		err := res.LoadImageFile("public/ducks.png")
		if err != nil {
			res.SetStatusCode(500)
			res.SetStringBody(err.Error())
		}

		return res

	} else {
		res.SetStatusCode(200)
		res.SetHeader("Content-Disposition", "attachment; filename=ducks.png")
		res.SetHeader("Content-Type", "image/png")
		imagesBytes, err := res.GetImageBytes("public/ducks.png")
		if err != nil {
			res.SetStatusCode(500)
			res.SetStringBody(err.Error())
		}

		res.SetStringBody(string(imagesBytes))

		return res
	}
}

func main() {
	CHC := chc.NewCHC()

	CHC.AddRoute(&chc.Route{
		Path:       "/login",
		Methods:    []string{"GET", "POST"},
		Controller: LoginController,
	})

	CHC.AddRoute(&chc.Route{
		Path:       "/dashboard",
		Methods:    []string{"GET"},
		Controller: DashboardController,
	})

	CHC.AddRoute(&chc.Route{
		Path:       "/ducks",
		Methods:    []string{"GET", "POST"},
		Controller: DuckController,
	})

	CHC.Listen("127.0.0.1", 8080)
}
