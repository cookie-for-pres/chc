package chc

type Route struct {
	Path        string
	Methods     []string
	Headers     map[string]string
	File        string
	Controller  func(request *Request, response Response) (response_ *Response)
	Middlewares []func(response Response) (response_ *Response, next bool)
}

func (chc *CHC) handleRoute(route *Route, request *Request) *Response {
	isMethod := false
	for _, method := range route.Methods {
		if method == request.Method {
			isMethod = true
			break
		}
	}

	if !isMethod {
		response := "HTTP/1.1 405 Method Not Allowed\r\n"
		response += "Content-Type: text/plain\r\n\r\n"
		response += "Method Not Allowed"
		response += "\r\n"

		request.Conn.Write([]byte(response))
		return &Response{StatusCode: 405}
	}

	if route.Controller != nil {
		var newResponse Response
		for _, middleware := range route.Middlewares {
			response, isNext := middleware(newResponse)
			newResponse = *response
			if !isNext {
				break
			}
		}
		response := route.Controller(request, newResponse)
		// fmt.Println(newResponse)
		message := route.handleResponse(&newResponse, request)
		request.Conn.Write([]byte(message))
		return response
	} else {
		response := "HTTP/1.1 404 Not Found\r\n"
		response += "Content-Type: text/plain\r\n\r\n"
		response += "Please add a handler for this path"
		response += "\r\n"

		request.Conn.Write([]byte(response))
		return &Response{StatusCode: 404}
	}
}

// Create a new route object
func (chc *CHC) NewRoute(path string, methods []string, headers map[string]string, file string, controller func(request *Request, response Response) (response_ *Response)) *Route {
	return &Route{
		Path:       path,
		Methods:    methods,
		Headers:    headers,
		File:       file,
		Controller: controller,
	}
}

// Add a route to CHC
func (chc *CHC) AddRoute(route *Route) {
	chc.Routes = append(chc.Routes, route)
}

// Add multiple routes to CHC
func (chc *CHC) AddRoutes(routes []*Route) {
	for _, route := range routes {
		chc.AddRoute(route)
	}
}

// Add a middleware to a route
func (route *Route) AddMiddleware(middleware func(response Response) (response_ *Response, next bool)) {
	route.Middlewares = append(route.Middlewares, middleware)
}

// Add multiple middlewares to a route
func (route *Route) AddMiddlewares(middlewares []func(response Response) (response_ *Response, next bool)) {
	for _, middleware := range middlewares {
		route.AddMiddleware(middleware)
	}
}
