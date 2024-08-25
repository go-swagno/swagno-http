package swagger

import (
	"fmt"
	"net/http"

	swaggerUi "github.com/go-swagno/swagno-files"
)

type Config struct {
	Prefix string
}

var defaultConfig = Config{
	Prefix: "/swagger",
}

type optFunc func(*Config)

func WithPrefix(prefix string) optFunc {
	return func(c *Config) {
		c.Prefix = prefix
	}
}

func SwaggerHandler(doc []byte, opts ...optFunc) http.HandlerFunc {
	conf := defaultConfig

	for _, opt := range opts {
		opt(&conf)
	}

	swaggerDoc := string(doc)

	handler := swaggerUi.Handler

	return func(w http.ResponseWriter, r *http.Request) {
		prefix := conf.Prefix
		handler.Prefix = prefix

		switch r.RequestURI {
		case prefix, prefix + "/":
			http.Redirect(w, r, prefix+"/index.html", http.StatusMovedPermanently)
		case prefix + "/doc.json":
			fmt.Fprint(w, swaggerDoc)
		default:
			handler.ServeHTTP(w, r)
		}
	}
}
