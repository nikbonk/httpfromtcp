package response

import (
	"io"
	"strconv"

	"github.com/nikbonk/httpfromtcp/internal/headers"
)

type StatusCode int

const (
	StatusCodeSuccess             StatusCode = 200
	StatusCodeBadRequest          StatusCode = 400
	StatusCodeInternalServerError StatusCode = 500
)

func WriteStatusLine(w io.Writer, statusCode StatusCode) error {
	statusLine := "HTTP/1.1 "
	reasonPhrase := ""
	statusCodeString := ""

	statusCodeString = strconv.Itoa(int(statusCode))
	switch statusCode {
	case StatusCodeSuccess:
		reasonPhrase = "OK"
	case StatusCodeBadRequest:
		reasonPhrase = "Bad Request"
	case StatusCodeInternalServerError:
		reasonPhrase = "Internal Server Error"
	default:
		// Empty because we don't know the status code
	}
	_, err := io.WriteString(w, statusLine+statusCodeString+" "+reasonPhrase+"\r\n")
	return err
}

func GetDefaultHeaders(contentLen int) headers.Headers {
	return headers.Headers{
		"Content-Length": strconv.Itoa(contentLen),
		"Connection":     "close",
		"Content-Type":   "text/plain",
	}
}

func WriteHeaders(w io.Writer, headers headers.Headers) error {
	for key, value := range headers {
		_, err := io.WriteString(w, key+": "+value+"\r\n")
		if err != nil {
			return err
		}
	}
	_, err := io.WriteString(w, "\r\n")
	return err
}
