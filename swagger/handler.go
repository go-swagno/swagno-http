package swagger

import (
	"fmt"
	"net/http"

	swaggerUi "github.com/go-swagno/swagno-files"
	"golang.org/x/net/webdav"
)

type Config struct {
	Prefix string
}

var swaggerDoc string

var handler *webdav.Handler

var defaultConfig = Config{
	Prefix: "/swagger",
}

func SwaggerHandler(doc []byte, config ...Config) http.HandlerFunc {
	if len(config) != 0 {
		defaultConfig = config[0]
	}
	if swaggerDoc == "" {
		swaggerDoc = string(doc)
	}
	if handler == nil {
		handler = swaggerUi.Handler
	}

	return func(w http.ResponseWriter, r *http.Request) {
		prefix := defaultConfig.Prefix
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
