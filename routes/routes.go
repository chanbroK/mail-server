package routes

import (
	"encoding/json"
	"io"
	"net/http"
)

type Router struct {
	*http.ServeMux
}

func NewRouter() *Router {
	return &Router{
		ServeMux: http.NewServeMux(),
	}
}

func (r *Router) Post(pattern string, handler func(c *Context)) {
	r.ServeMux.HandleFunc(pattern, func(writer http.ResponseWriter, request *http.Request) {
		// post method 만 수행
		if request.Method == http.MethodPost {
			body, err := io.ReadAll(request.Body)
			if err != nil {
				panic(err)
			}

			handler(&Context{
				RawBody: body,
				Writer:  writer,
			})
		}
	})
}

type Context struct {
	RawBody []byte
	Writer  http.ResponseWriter
}

func (c *Context) JSON(v any) {
	c.Writer.WriteHeader(http.StatusOK)
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.Header().Set("Content-Encoding", "gzip")

	err := json.NewEncoder(c.Writer).Encode(v)
	if err != nil {
		panic(err)
	}
}

func (c *Context) BindBody(ref any) {
	err := json.Unmarshal(c.RawBody, ref)
	if err != nil {
		panic("invalid body")
	}
}
