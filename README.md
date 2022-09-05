# swagno-http
`net/http` & `gorilla/mux` middleware to serve Swagger Ui files and doc.json file which generated from [Swagno](https://github.com/go-swagno/swagno)

## Usage

1. Get [Swagno](https://github.com/go-swagno/swagno)
2. Create your endpoint array. See: https://github.com/go-swagno/swagno/blob/master/README.md#getting-started
3. Get swagno-http
```sh
go get github.com/go-swagno/swagno-http
```
4. Import swagno-http to your handler
```go
import "github.com/go-swagno/swagno-http/swagger"
```
5. **Be sure you created swagger instance and endpoints**
6. Create swagger handler
```go
// net-http
m.Handle("/swagger/", swagger.SwaggerHandler(sw.GenerateDocs()))

// gorilla/mux
r.PathPrefix("/swagger").Handler(swagger.SwaggerHandler(sw.GenerateDocs()))
```
7. Visit /swagger and /swagger/doc.json for confirmation
