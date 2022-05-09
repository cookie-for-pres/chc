package chc

import (
	"fmt"
	"net"
)

type CHC struct {
	Paths []*Path
}

func (chc *CHC) Listen(host string, port int) {
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		panic(err)
	}

	for {
		conn, err := ln.Accept()
		if err != nil {
			panic(err)
		}

		go func(conn net.Conn) {
			defer conn.Close()

			buf := make([]byte, 1024)
			n, err := conn.Read(buf)
			if err != nil {
				return
			}

			request := chc.ParseRequest(string(buf[:n]), conn)
			for _, path := range chc.Paths {
				if path.Path == request.URL {
					for _, method := range path.Methods {
						if method == request.Method {
							chc.HandlePath(path, request)
							return
						}
					}

					response := "HTTP/1.1 405 Method Not Allowed\r\n"
					response += "Content-Type: text/plain\r\n\r\n"
					response += "Method Not Allowed"
					response += "\r\n"

					fmt.Fprint(conn, response)
					return
				}
			}

			response := "HTTP/1.1 404 Not Found\r\n"
			response += "Content-Type: text/plain\r\n\r\n"
			response += "404 Not Found"
			response += "\r\n"

			fmt.Fprint(conn, response)
		}(conn)
	}
}
