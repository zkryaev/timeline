package middleware

import "net/http"

type respWriterCustom struct {
	http.ResponseWriter
	statusCode int
	header     http.Header
}

func (rw *respWriterCustom) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *respWriterCustom) Write(b []byte) (int, error) {
	if rw.statusCode == 0 { // Если WriteHeader не был вызван
		rw.WriteHeader(http.StatusOK)
	}
	return rw.ResponseWriter.Write(b)
}

func (rw *respWriterCustom) Header() http.Header {
	return rw.ResponseWriter.Header()
}
