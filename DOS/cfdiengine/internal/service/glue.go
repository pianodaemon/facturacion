package service

import (
	"os"

	"github.com/gin-gonic/gin"

	co "blaucorp.com/fiscal-engine/internal/controllers"
	hups "blaucorp.com/fiscal-engine/pkg/hookups"
)

// Engages the RESTful API
func Engage() (merr error) {

	defer func() {

		if r := recover(); r != nil {
			merr = r.(error)
		}
	}()

	getEnv := func(key, fallback string) string {
		if value, ok := os.LookupEnv(key); ok {
			return value
		}
		return fallback
	}

	fiscalEngineImplt := hups.NewFiscalEngine(getEnv("MONGO_URI", "mongodb://127.0.0.1:27017"))

	r := gin.Default()
	setUpHandlers(r, fiscalEngineImplt)
	r.Run() // listen and serve on 0.0.0.0:8080

	return nil
}

func setUpHandlers(r *gin.Engine, pm co.FiscalEngineInterface) {

	r.GET("/ping", co.Health)
	r.POST("/receipts", co.CreateReceipt(pm))
	r.PUT("/receipts/:id", co.EditReceipt(pm))
}
