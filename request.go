package chc

import (
	"encoding/json"
	"net"
	"strings"
)

type Request struct {
	Method   string
	URL      string
	Protocol string
	Body     string
	Params   map[string]string
	Headers  map[string]string
	Cookies  map[string]string
	Conn     net.Conn
}

func (chc *CHC) ParseCookies(cookies string) map[string]string {
	cookieMap := make(map[string]string)
	lines := strings.Split(cookies, "\n")
	cookieLine := ""

	for _, line := range lines {
		if !strings.Contains(line, "Cookie:") {
			continue
		} else {
			cookieLine = line
			break
		}
	}

	cookieLine = strings.TrimPrefix(cookieLine, "Cookie:")
	cookieLine = strings.TrimSpace(cookieLine)
	cookieLines := strings.Split(cookieLine, ";")

	for _, cookie := range cookieLines {
		kv := strings.Split(cookie, "=")
		if len(kv) != 2 {
			continue
		}

		cookieMap[kv[0]] = kv[1]
	}

	return cookieMap
}

func (chc *CHC) ParseHeaders(headers string) map[string]string {
	headerMap := make(map[string]string)

	headers = strings.TrimPrefix(headers, "")
	headers = strings.TrimSpace(headers)

	for _, header := range strings.Split(headers, "\n") {
		kv := strings.Split(header, ":")
		if len(kv) != 2 {
			continue
		}

		if kv[0] == "Cookie" {
			continue
		}

		headerMap[kv[0]] = kv[1]
	}

	return headerMap
}

func (chc *CHC) ParseParams(url string) (map[string]string, string) {
	params := make(map[string]string)
	if !strings.Contains(url, "?") && !strings.Contains(url, "=") && !strings.Contains(url, "&") {
		return params, url
	}

	newUrl := strings.Split(url, "?")[0]
	url = strings.TrimSpace(url)
	url = strings.Replace(url, "?", "", -1)
	url = strings.Replace(url, "&", "&", -1)
	url = strings.Replace(url, "=", "=", -1)

	for _, param := range strings.Split(url, "&") {
		kv := strings.Split(param, "=")
		if len(kv) != 2 {
			continue
		}

		params[kv[0]] = kv[1]
	}

	return params, newUrl
}

func (chc *CHC) ParseRequest(requestString string, conn net.Conn) *Request {
	request := &Request{}
	request.Conn = conn
	request.Headers = make(map[string]string)
	request.Cookies = make(map[string]string)

	body := strings.Split(requestString, "\r\n\r\n")[1]

	lines := strings.Split(requestString, "\n")
	for _, line := range lines {
		if line == "" || line == "\r" {
			continue
		}

		if strings.HasPrefix(line, "GET") {
			request.Method = "GET"
			if strings.Contains(line, "HTTP/1.1") {
				request.Protocol = "HTTP/1.1"
			} else if strings.Contains(line, "HTTP/1.0") {
				request.Protocol = "HTTP/1.0"
			}

			request.URL = strings.Replace(line, request.Method, "", 1)
			request.URL = strings.Replace(request.URL, request.Protocol, "", 1)
			request.URL = strings.TrimSpace(request.URL)
		} else if strings.HasPrefix(line, "POST") {
			request.Method = "POST"
			if strings.Contains(line, "HTTP/1.1") {
				request.Protocol = "HTTP/1.1"
			} else if strings.Contains(line, "HTTP/1.0") {
				request.Protocol = "HTTP/1.0"
			}

			request.URL = strings.Replace(line, request.Method, "", 1)
			request.URL = strings.Replace(request.URL, request.Protocol, "", 1)
			request.URL = strings.TrimSpace(request.URL)
		} else if strings.HasPrefix(line, "PUT") {
			request.Method = "PUT"
			if strings.Contains(line, "HTTP/1.1") {
				request.Protocol = "HTTP/1.1"
			} else if strings.Contains(line, "HTTP/1.0") {
				request.Protocol = "HTTP/1.0"
			}

			request.URL = strings.Replace(line, request.Method, "", 1)
			request.URL = strings.Replace(request.URL, request.Protocol, "", 1)
			request.URL = strings.TrimSpace(request.URL)
		} else if strings.HasPrefix(line, "DELETE") {
			request.Method = "DELETE"
			if strings.Contains(line, "HTTP/1.1") {
				request.Protocol = "HTTP/1.1"
			} else if strings.Contains(line, "HTTP/1.0") {
				request.Protocol = "HTTP/1.0"
			}

			request.URL = strings.Replace(line, request.Method, "", 1)
			request.URL = strings.Replace(request.URL, request.Protocol, "", 1)
			request.URL = strings.TrimSpace(request.URL)
		} else if strings.HasPrefix(line, "PATCH") {
			request.Method = "PATCH"
			if strings.Contains(line, "HTTP/1.1") {
				request.Protocol = "HTTP/1.1"
			} else if strings.Contains(line, "HTTP/1.0") {
				request.Protocol = "HTTP/1.0"
			}

			request.URL = strings.Replace(line, request.Method, "", 1)
			request.URL = strings.Replace(request.URL, request.Protocol, "", 1)
			request.URL = strings.TrimSpace(request.URL)
		} else if strings.HasPrefix(line, "HEAD") {
			request.Method = "HEAD"
			if strings.Contains(line, "HTTP/1.1") {
				request.Protocol = "HTTP/1.1"
			} else if strings.Contains(line, "HTTP/1.0") {
				request.Protocol = "HTTP/1.0"
			}

			request.URL = strings.Replace(line, request.Method, "", 1)
			request.URL = strings.Replace(request.URL, request.Protocol, "", 1)
			request.URL = strings.TrimSpace(request.URL)
		} else if strings.HasPrefix(line, "OPTIONS") {
			request.Method = "OPTIONS"
			if strings.Contains(line, "HTTP/1.1") {
				request.Protocol = "HTTP/1.1"
			} else if strings.Contains(line, "HTTP/1.0") {
				request.Protocol = "HTTP/1.0"
			}

			request.URL = strings.Replace(line, request.Method, "", 1)
			request.URL = strings.Replace(request.URL, request.Protocol, "", 1)
			request.URL = strings.TrimSpace(request.URL)
		}
	}

	cookies := chc.ParseCookies(requestString)
	headers := chc.ParseHeaders(requestString)
	params, newUrl := chc.ParseParams(request.URL)
	request.Cookies = cookies
	request.Headers = headers
	request.Params = params
	request.URL = newUrl
	request.Body = body

	return request
}

func (request *Request) NewResponse() *Response {
	response := &Response{}
	response.Headers = make(map[string]string)
	response.Cookies = make(map[string]string)
	response.Body = ""
	response.StatusCode = 200

	return response
}

func (request *Request) GetParam(key string) string {
	return request.Params[key]
}

func (request *Request) Json() (map[string]interface{}, error) {
	var jsonData map[string]interface{}
	err := json.Unmarshal([]byte(request.Body), &jsonData)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

func (request *Request) JsonArray() ([]map[string]interface{}, error) {
	var jsonData []map[string]interface{}
	err := json.Unmarshal([]byte(request.Body), &jsonData)
	if err != nil {
		return nil, err
	}

	return jsonData, nil
}

func (request *Request) FormData() (map[string]string, error) {
	formData := make(map[string]string)
	for _, pair := range strings.Split(request.Body, "&") {
		// fmt.Println(pair)
		kv := strings.Split(pair, "=")
		if len(kv) != 2 {
			continue
		}

		formData[kv[0]] = kv[1]
	}

	return formData, nil
}
