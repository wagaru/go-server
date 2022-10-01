package http

import (
	"errors"
	"regexp"
)

var (
	ERR_NOT_SUPPORTED_REQUEST_HTTP_METHOD  = errors.New("Request: unsupported request method")
	ERR_NOT_SUPPORTED_REQUEST_HTTP_VERSION = errors.New("Request: unsupported request HTTP version")

	SUPPORTED_HTTP_VERSION = []string{"1.0", "1.1"}
)

type Request struct {
	Method  RequestMethod
	URI     string
	Version string
	Headers []*RequestHeader
	Body    string
}

type RequestHeader struct {
	Name  string
	Value string
}

type RequestMethod int

const (
	REQUEST_METHOD_GET RequestMethod = iota
	REQUEST_METHOD_POST
	REQUEST_METHOD_DELETE
	REQUEST_METHOD_PUT
	REQUEST_METHOD_PATCH
)

// GET /xsdfasfasdfasdf HTTP/1.1
func ParseReceivedData(str string, req *Request) error {
	// start line
	reg := regexp.MustCompile(`([a-zA-Z]+)\s(\/.*)\sHTTP\/(\d+\.\d+)`)
	matches := reg.FindStringSubmatch(str)
	if matches != nil {
		method, err := parseRequestMethod(matches[1])
		if err != nil {
			return err
		}
		req.Method = method
		req.URI = matches[2]
		version, err := parseRequestVersion(matches[3])
		if err != nil {
			return err
		}
		req.Version = version
		return nil
	}

	// header
	reg = regexp.MustCompile(`(.*):\s?(.*)`)
	matches = reg.FindStringSubmatch(str)
	if matches != nil {
		req.Headers = append(req.Headers, &RequestHeader{
			Name:  matches[1],
			Value: matches[2],
		})
		return nil
	}

	// body
	req.Body += str
	return nil
}

func parseRequestMethod(str string) (RequestMethod, error) {
	switch str {
	case "GET":
		return REQUEST_METHOD_GET, nil
	case "POST":
		return REQUEST_METHOD_POST, nil
	case "DELETE":
		return REQUEST_METHOD_DELETE, nil
	case "PUT":
		return REQUEST_METHOD_PUT, nil
	case "PATCH":
		return REQUEST_METHOD_PATCH, nil
	default:
		return RequestMethod(0), ERR_NOT_SUPPORTED_REQUEST_HTTP_METHOD
	}
}

func parseRequestVersion(str string) (string, error) {
	for _, v := range SUPPORTED_HTTP_VERSION {
		if v == str {
			return v, nil
		}
	}
	return "", ERR_NOT_SUPPORTED_REQUEST_HTTP_VERSION
}
