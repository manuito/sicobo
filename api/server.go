package api

/*
 * REST services for library managment
 * A swagger-ui is included => http://localhost:8080/apidocs/?url=http://localhost:8080/apidocs.json
 */

import (
	"log"
	"net/http"
	"sicobo/application"
	"sicobo/store"

	"github.com/emicklei/go-restful"
	restfulspec "github.com/emicklei/go-restful-openapi"
	"github.com/go-openapi/spec"
)

func selectDb(request *restful.Request) (db string) {
	db = request.PathParameter("database")
	store.SelectActiveDatabase(db)
	return
}

// GET http://1.2.3.4:8080/api/v1/library/{database}
func getLibrarySpec(request *restful.Request, response *restful.Response) {
	selectDb(request)
	application.Debug("Get library spec")
	response.WriteEntity(store.ActiveSpec())
}

// GET http://1.2.3.4:8080/api/v1/library/{database}/books
func getLibraryBooks(request *restful.Request, response *restful.Response) {
	selectDb(request)
	application.Debug("Get library books")
	response.WriteEntity(store.ListExistingBooks())
}

// GET http://1.2.3.4:8080/api/v1/library/{database}/books/{isbn}
func getLibraryBook(request *restful.Request, response *restful.Response) {
	selectDb(request)
	isbn := request.PathParameter("isbn")
	application.Debug("Get library book with isbn", isbn)
	book, err := store.GetBook(isbn)
	if err != nil {
		response.WriteError(404, err)
	} else {
		response.WriteEntity(book)
	}
}

// POST http://1.2.3.4:8080/api/v1/library/{database}/books/{isbn}
func addBookToLibrary(request *restful.Request, response *restful.Response) {
	db := selectDb(request)
	isbn := request.PathParameter("isbn")
	application.Debug("Add book with isbn", isbn, " to library books", db)
	book, _ := store.AddNewBook(isbn)
	response.WriteEntity(book)
}

// Routes for top level services
func rootWebService() *restful.WebService {
	ws := new(restful.WebService)
	ws.
		Path("/api/v1/library").
		Consumes(restful.MIME_JSON).
		Produces(restful.MIME_JSON)

	tags := []string{"library"}

	ws.Route(ws.GET("/{database}").To(getLibrarySpec).
		Doc("Get library spec").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.PathParameter("database", "database name").DataType("string")).
		Writes(store.BookDatabaseSpec{}).
		Returns(200, "OK", store.BookDatabaseSpec{}))

	ws.Route(ws.GET("/{database}/books").To(getLibraryBooks).
		Doc("get library books").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.PathParameter("database", "database name").DataType("string")).
		Writes([]store.Book{}))

	ws.Route(ws.GET("/{database}/books/{isbn}").To(getLibraryBook).
		Doc("get library book for specified isbn").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.PathParameter("database", "database name").DataType("string")).
		Param(ws.PathParameter("isbn", "isbn code").DataType("string")).
		Writes(store.Book{}))

	ws.Route(ws.POST("/{database}/books/{isbn}").To(addBookToLibrary).
		Doc("get library books").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Param(ws.PathParameter("database", "database name").DataType("string")).
		Param(ws.PathParameter("isbn", "isbn code").DataType("string")).
		Writes(store.Book{}))

	return ws
}

// StartServer : Start the REST service for library managment
func StartServer() {

	restful.DefaultContainer.Add(rootWebService())

	config := restfulspec.Config{
		WebServices:                   restful.RegisteredWebServices(),
		APIPath:                       "/apidocs.json",
		PostBuildSwaggerObjectHandler: enrichSwaggerObject}
	restful.DefaultContainer.Add(restfulspec.NewOpenAPIService(config))

	// Swagger-ui
	http.Handle("/apidocs/", http.StripPrefix("/apidocs/", http.FileServer(http.Dir("./api/swagger-ui/dist-v2"))))

	cors := restful.CrossOriginResourceSharing{
		AllowedHeaders: []string{"Content-Type", "Accept"},
		AllowedMethods: []string{"GET", "PUT"},
		CookiesAllowed: false,
		Container:      restful.DefaultContainer}
	restful.DefaultContainer.Filter(cors.Filter)

	ip := application.State.OutBountIP.String()

	application.Info("Get the API using http://" + ip + ":8080/apidocs.json")
	application.Info("Open Swagger UI using http://" + ip + ":8080/apidocs/?url=http://" + ip + ":8080/apidocs.json")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

// enrichSwaggerObject : Add swagger def for service
func enrichSwaggerObject(swo *spec.Swagger) {
	swo.Info = &spec.Info{
		InfoProps: spec.InfoProps{
			Title:       "Library manager",
			Description: "Resources for managing a library",
			Contact: &spec.ContactInfo{
				Name: "elecomte",
				URL:  "https://www.elecomte.com",
			},
			License: &spec.License{
				Name: "WTFPL",
				URL:  "http://www.wtfpl.net/about/",
			},
			Version: "0.0.1",
		},
	}
	swo.Tags = []spec.Tag{
		spec.Tag{TagProps: spec.TagProps{
			Name:        "library",
			Description: "Processing library content"}}}
}
