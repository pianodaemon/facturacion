package service

import (
	"github.com/gorilla/mux"
	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"

	"immortalcrab.com/interview/internal/rsapi"
	dal "immortalcrab.com/interview/internal/storage"
)

var apiSettings rsapi.RestAPISettings

func init() {

	envconfig.Process("rsapi", &apiSettings)
}

// Engages the RESTful API
func Engage(logger *logrus.Logger) (merr error) {

	defer func() {

		if r := recover(); r != nil {
			merr = r.(error)
		}
	}()

	{
		/* Probes are configured here
		   with anonymous adaptive functions */
		apiSettings.Probes = append(apiSettings.Probes, func() error {

			var errDB error = dal.ProbeConn()

			if errDB != nil {
				return errDB
			}

			logger.Println("Database connection sucessfully verified!")

			return errDB
		})
	}

	var err error = nil

	{

		/* The connection of both components occurs through
		   the router glue and its adaptive functions */
		glue := func(api *rsapi.RestAPI) *mux.Router {

			router := mux.NewRouter()

			return router
		}

		api := rsapi.NewRestAPI(logger, &apiSettings, glue)

		api.PowerOn()
	}

	return err
}
