package chc

type Path struct {
	Path    string
	Methods []string
	Type    string
	Headers map[string]string
	File    string
	Handler func(request *Request) (response *Response)
}

func (chc *CHC) AddPath(path *Path) {
	chc.Paths = append(chc.Paths, path)
}

func (chc *CHC) AddPaths(paths []*Path) {
	for _, path := range paths {
		chc.AddPath(path)
	}
}

func (chc *CHC) HandlePath(path *Path, request *Request) {
	isMethod := false
	for _, method := range path.Methods {
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
		return
	}

	if path.Handler != nil {
		response := path.Handler(request)
		message := path.HandleResponse(response, request)
		request.Conn.Write([]byte(message))
	} else {
		response := "HTTP/1.1 404 Not Found\r\n"
		response += "Content-Type: text/plain\r\n\r\n"
		response += "Please add a handler for this path"
		response += "\r\n"
	}

}
