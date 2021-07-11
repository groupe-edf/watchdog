package middleware

import (
	"bufio"
	"compress/gzip"
	"errors"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
)

// Compression middleware
type Compression struct {
}

func (middleware *Compression) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	L:
		for _, encoding := range strings.Split(r.Header.Get("Accept-Encoding"), ",") {
			switch strings.TrimSpace(encoding) {
			case "gzip":
				w.Header().Set("Content-Encoding", "gzip")
				w.Header().Add("Vary", "Accept-Encoding")
				gzipWriter := gzip.NewWriter(w)
				defer gzipWriter.Flush()
				defer gzipWriter.Close()
				w = &WrappedReponseWriter{
					originalWriter:    w,
					compressionWriter: gzipWriter,
				}
				break L
			}
		}
		next.ServeHTTP(w, r)
	})
}

type CompressionWriter interface {
	io.Writer
	Reset(io.Writer)
}

type WrappedReponseWriter struct {
	originalWriter    http.ResponseWriter
	compressionWriter CompressionWriter
}

func (writer *WrappedReponseWriter) Header() http.Header {
	return writer.originalWriter.Header()
}

func (writer *WrappedReponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := writer.originalWriter.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("compress response does not implement http.Hijacker")
	}
	writer.compressionWriter.Reset(ioutil.Discard)
	return h.Hijack()
}

func (writer *WrappedReponseWriter) Write(d []byte) (int, error) {
	headerMap := writer.originalWriter.Header()
	if headerMap.Get("Content-Type") == "" {
		headerMap.Set("Content-Type", http.DetectContentType(d))
	}
	headerMap.Del("Content-Length")
	return writer.originalWriter.Write(d)
}

func (writer *WrappedReponseWriter) WriteHeader(statusCode int) {
	writer.originalWriter.Header().Del("Content-Length")
	writer.originalWriter.WriteHeader(statusCode)
}
