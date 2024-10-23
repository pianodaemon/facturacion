package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	hups "blaucorp.com/fiscal-engine/pkg/hookups"
)

type (
	FiscalEngineInterface interface {
		DoCreateReceipt(dto *hups.ReceiptDTO) (string, error)
		DoEditReceipt(receiptID string, dto *hups.ReceiptDTO) error
	}
)

func CreateReceipt(fiscalEngineImplt FiscalEngineInterface) func(c *gin.Context) {

	return func(c *gin.Context) {

		dto := &hups.ReceiptDTO{}

		if errP := c.ShouldBind(dto); errP != nil {
			c.String(http.StatusBadRequest, "the body should be form of recipe type")
			return
		}

		name, err := fiscalEngineImplt.DoCreateReceipt(dto)
		if err != nil {
			panic(err.Error())
		}

		fmt.Println("-------> body: ", name)
		c.JSON(http.StatusOK, gin.H{
			"results": "ok",
		})
	}
}

func EditReceipt(fiscalEngineImplt FiscalEngineInterface) func(c *gin.Context) {

	return func(c *gin.Context) {

		dto := &hups.ReceiptDTO{}

		// Bind the JSON payload to the dto object
		if errP := c.ShouldBind(dto); errP != nil {
			c.String(http.StatusBadRequest, "the body should be in the form of receipt type")
			return
		}

		// Extract the receipt ID from the URL or request parameters
		receiptID := c.Param("id") // assumes the receipt ID is passed in the URL path as /receipts/:id
		if receiptID == "" {
			c.String(http.StatusBadRequest, "missing receipt ID")
			return
		}

		// Call the edit method in the fiscal engine
		err := fiscalEngineImplt.DoEditReceipt(receiptID, dto)
		if err != nil {
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		// Send a success response
		c.JSON(http.StatusOK, gin.H{
			"result": "receipt updated successfully",
		})
	}
}
