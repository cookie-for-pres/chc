package chc

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
)

type Response struct {
	StatusCode int
	Body       string
	Headers    map[string]string
	Cookies    map[string]string
	Conn       net.Conn
}

var statusCodes = map[int]string{
	100: "Continue",
	101: "Switching Protocols",
	102: "Processing",
	103: "Early Hints",
	200: "OK",
	201: "Created",
	202: "Accepted",
	203: "Non-Authoritative Information",
	204: "No Content",
	205: "Reset Content",
	206: "Partial Content",
	207: "Multi-Status",
	208: "Already Reported",
	226: "IM Used",
	300: "Multiple Choices",
	301: "Moved Permanently",
	302: "Found",
	303: "See Other",
	304: "Not Modified",
	305: "Use Proxy",
	307: "Temporary Redirect",
	308: "Permanent Redirect",
	400: "Bad Request",
	401: "Unauthorized",
	402: "Payment Required",
	403: "Forbidden",
	404: "Not Found",
	405: "Method Not Allowed",
	406: "Not Acceptable",
	407: "Proxy Authentication Required",
	408: "Request Timeout",
	409: "Conflict",
	410: "Gone",
	411: "Length Required",
	412: "Precondition Failed",
	413: "Payload Too Large",
	414: "URI Too Long",
	415: "Unsupported Media Type",
	416: "Range Not Satisfiable",
	417: "Expectation Failed",
	418: "I'm a teapot",
	421: "Misdirected Request",
	422: "Unprocessable Entity",
	423: "Locked",
	424: "Failed Dependency",
	426: "Upgrade Required",
	428: "Precondition Required",
	429: "Too Many Requests",
	431: "Request Header Fields Too Large",
	451: "Unavailable For Legal Reasons",
	500: "Internal Server Error",
	501: "Not Implemented",
	502: "Bad Gateway",
	503: "Service Unavailable",
	504: "Gateway Timeout",
	505: "HTTP Version Not Supported",
	506: "Variant Also Negotiates",
	507: "Insufficient Storage",
	508: "Loop Detected",
	510: "Not Extended",
	511: "Network Authentication Required",
}

func (path *Path) HandleResponse(res *Response, req *Request) string {
	var response string
	response += fmt.Sprintf("%s %d %s\r\n", req.Protocol, res.StatusCode, statusCodes[res.StatusCode])
	for k, v := range res.Headers {
		response += fmt.Sprintf("%s: %s\r\n", k, v)
	}

	for k, v := range res.Cookies {
		response += fmt.Sprintf("Set-Cookie: %s=%s\r\n", k, v)
	}

	response += "\r\n"
	response += res.Body

	return response
}

func (response *Response) SetStatusCode(statusCode int) {
	response.StatusCode = statusCode
}

func (response *Response) SetJsonBody(body string) {
	response.Body = body

	if response.Headers == nil {
		response.Headers = make(map[string]string)
	}

	response.Headers["Content-Type"] = "application/json"
}

func (response *Response) SetBody(body string) {
	response.Body = body
}

func (response *Response) SetHeader(key string, value string) {
	if response.Headers == nil {
		response.Headers = make(map[string]string)
	}

	response.Headers[key] = value
}

func (response *Response) SetCookie(key string, value string) {
	if response.Cookies == nil {
		response.Cookies = make(map[string]string)
	}

	response.Cookies[key] = value
}

func (response *Response) SetRedirect(url string) {
	response.SetHeader("Location", url)
}

func (response *Response) LoadHtmlFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	response.Body = string(data)

	if response.Headers == nil {
		response.Headers = make(map[string]string)
	}

	response.Headers["Content-Type"] = "text/html"

	return nil
}

func (response *Response) LoadImageFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}

	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	response.Body = string(data)

	if response.Headers == nil {
		response.Headers = make(map[string]string)
	}

	response.Headers["Content-Type"] = "image/png"

	return nil
}

func (response *Response) GetImageBytes(filepath string) ([]byte, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	return data, nil
}
