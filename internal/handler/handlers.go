package handler

import (
	"fmt"
	"github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
	httpSwagger "github.com/swaggo/http-swagger"
	"go_kit_project/docs"
	_ "go_kit_project/docs"
	"go_kit_project/internal/db"
	"go_kit_project/internal/endpoints"
	"go_kit_project/internal/middleware"
	"go_kit_project/internal/repository"
	"go_kit_project/internal/service"
	"go_kit_project/internal/static"
	"gopkg.in/mgo.v2"
	"net/http"
	"time"
)

type App struct {
	Router           *mux.Router
	DB               *db.MongoConnection
	Logg             log.Logger
	PersonEndpoints  endpoints.PersonEndpoints
	PersonService    service.PersonService
	PersonRepository repository.PersonRepository
}

func (a *App) Run(addr string) error {
	err := http.ListenAndServe(addr, a.Router)
	return err
}

func (a *App) Initialize(user, password string) (err error) {
	fmt.Println(static.MsgResponseStartApplication)
	host := viper.GetString(static.MONGO_HOST) + ":" + viper.GetString(static.MONGO_PORT)
	dbs := viper.GetString(static.MONGO_DATABASE)
	info := &mgo.DialInfo{
		Addrs:    []string{host},
		Timeout:  60 * time.Hour,
		Database: dbs,
		Username: user,
		Password: password,
	}
	a.DB, _ = db.NewConnection(info)
	fmt.Println(static.MsgResponseConnectedMongoDB)
	muxObj := mux.NewRouter()
	muxObj.Use(middleware.CORS)
	a.Router = muxObj
	values := []interface{}{static.KeyType, static.SUCCESS, static.KeyURL, static.URLStartingNow, static.KeyMessage, static.MsgResponseStartingNow}
	middleware.LoggingOperation(a.Logg, values...)
	a.InitializeEndpoints()
	a.InitializeRoutes()
	a.InitializeSwagger()
	return err
}

// routing
func (a *App) InitializeRoutes() {
	var options []httptransport.ServerOption
	a.Router.PathPrefix(static.URLApi).Handler(httpSwagger.WrapHandler)
	a.Router.Methods(http.MethodPost).Path(static.URLCreatingOne).Handler(a.CreatePerson(options))
	a.Router.Methods(http.MethodPost).Path(static.URLUpdatingOne + "/{" + static.KeyId + "}").Handler(a.UpdatePerson(options))
	a.Router.Methods(http.MethodGet).Path(static.URLListingAll).Handler(a.ListPersons(options))
	a.Router.Methods(http.MethodGet).Path(static.URLGettingOne + "/{" + static.KeyId + "}").Handler(a.GetPerson(options))
	a.Router.Methods(http.MethodDelete).Path(static.URLDeletingOne + "/{" + static.KeyId + "}").Handler(a.DeletePerson(options))
}

// swagger
func (a *App) InitializeSwagger() {
	docs.SwaggerInfo.Title = static.MsgApiRestTitle
	docs.SwaggerInfo.Description = static.MsgApiRestDescription
	docs.SwaggerInfo.Version = static.MsgApiRestVersion1
	docs.SwaggerInfo.Host = viper.GetString(static.APP_HOST) + ":" + viper.GetString(static.APP_PORT)
	docs.SwaggerInfo.BasePath = static.URLStartingNow
	docs.SwaggerInfo.Schemes = []string{static.SchemaHttp}
}

// ENDPOINTS
func (a *App) InitializeEndpoints() {
	a.PersonRepository = repository.NewPersonRepository(a.DB, a.Logg)
	a.PersonService = service.NewPersonService(a.PersonRepository)
	a.PersonEndpoints = endpoints.MakePersonEndpoints(a.PersonService)
}
