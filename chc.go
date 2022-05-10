package chc

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type CHC struct {
	Routes  []*Route
	Logging bool
}

var colors = map[string]string{
	"reset":   "\x1b[0m",
	"bold":    "\x1b[1m",
	"under":   "\x1b[4m",
	"black":   "\x1b[30m",
	"red":     "\x1b[1;1m\x1b[31m",
	"green":   "\x1b[1;1m\x1b[32m",
	"yellow":  "\x1b[1;1m\x1b[33m",
	"blue":    "\x1b[1;1m\x1b[34m",
	"magenta": "\x1b[1;1m\x1b[35m",
	"cyan":    "\x1b[1;1m\x1b[36m",
	"white":   "\x1b[37m",
}

func logRequest(request *Request, statusCode int) {
	currentTime := time.Now().Format("2006-01-02 15:04:05")
	fmt.Printf("%srequest%s - %d %s %s %s %s\n", colors["magenta"], colors["reset"], statusCode, request.Method, request.URL, request.Protocol, currentTime)
}

// Create a new CHC object
func NewCHC() *CHC {
	return &CHC{
		Routes:  make([]*Route, 0),
		Logging: true,
	}
}

// Turn Request Logging on or off (if using NewCHC function default is on, otherwise default is off)
func (chc *CHC) RequestLogging(logging bool) {
	chc.Logging = logging
}

// Load Environment Variables from the given file
func (chc *CHC) LoadEnv(filePath string) {
	err := godotenv.Load(filePath)
	if err != nil {
		fmt.Printf("%serror%s - could not load environment variables from %s\n", colors["red"], colors["reset"], filePath)
		return
	}

	fmt.Printf("%sinfo%s - loaded environment variables from %s\n", colors["green"], colors["reset"], filePath)
}

// Get Environment Variable from the given key
func (chc *CHC) GetEnv(key string) string {
	return os.Getenv(key)
}

// Start the CHC server on the given address and port
func (chc *CHC) Listen(host string, port int) {
	ln, err := net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
	if err != nil {
		panic(err)
	}

	fmt.Printf("%sready%s - listening on %s:%d\n", colors["green"], colors["reset"], host, port)

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

			request := chc.parseRequest(string(buf[:n]), conn)
			for _, route := range chc.Routes {
				if route.Path == request.URL {
					for _, method := range route.Methods {
						if method == request.Method {
							response := chc.handleRoute(route, request)
							if chc.Logging {
								logRequest(request, response.StatusCode)
							}
							return
						}
					}

					response := "HTTP/1.1 405 Method Not Allowed\r\n"
					response += "Content-Type: text/plain\r\n\r\n"
					response += "Method Not Allowed"
					response += "\r\n"

					fmt.Fprint(conn, response)
					if chc.Logging {
						logRequest(request, 405)
					}
					return
				}
			}

			response := "HTTP/1.1 404 Not Found\r\n"
			response += "Content-Type: text/plain\r\n\r\n"
			response += "404 Not Found"
			response += "\r\n"

			fmt.Fprint(conn, response)

			if chc.Logging {
				logRequest(request, 404)
			}
		}(conn)
	}
}
