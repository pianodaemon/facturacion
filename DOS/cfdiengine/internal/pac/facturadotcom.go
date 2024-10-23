package pac

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"blaucorp.com/fiscal-engine/internal/dal"
)

type (

	// Represents a factura.com PAC
	FacturaDotComEngine struct {
	}
)

const createEP = "%s/v4/cfdi40/create"

// Spawns an newer instance of the factura.com implementation
func NewFacturaDotComEngine() *FacturaDotComEngine {

	return &FacturaDotComEngine{}
}

func (self *FacturaDotComEngine) DoFact(sourceID string) ([]byte, error) {
	sandbox := "https://sandbox.factura.com/api"
	provider := dal.NewPgSQLPayloadProvider(sourceID)

	issueHeaders, err := provider.GetIssueHeaders()
	if err != nil {
		return nil, &CFDIEngineError{Code: DBIssue, Message: err.Error()}
	}

	buff, err := provider.DumpBuffer()
	if err != nil {
		return nil, &CFDIEngineError{Code: DBIssue, Message: err.Error()}
	}

	barr, err := hit(buff, sandbox, issueHeaders.FPlugin, issueHeaders.FApiKey, issueHeaders.FSecretKey)
	if err != nil {
		return nil, err
	}

	return barr, nil
}

func hit(body *bytes.Buffer,
	host string,
	headerFPlugin string,
	headerFApiKey string,
	headerFSecretKey string) ([]byte, error) {

	url := fmt.Sprintf(createEP, host)
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, &CFDIEngineError{Code: UnknownIssue, Message: err.Error()}
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("F-PLUGIN", headerFPlugin)
	req.Header.Set("F-Api-Key", headerFApiKey)
	req.Header.Set("F-Secret-Key", headerFSecretKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, &CFDIEngineError{Code: PACConnIssue, Message: err.Error()}
	}
	defer resp.Body.Close()

	barr, err := ioutil.ReadAll(resp.Body)

	var msg map[string]interface{}
	err = json.Unmarshal(barr, &msg)
	if err != nil {
		return nil, &CFDIEngineError{Code: UnknownIssue, Message: err.Error()}
	}

	if msg["response"] == "success" {
		return barr, nil
	}

	return nil, parseErrorMessage(barr)
}

// Parses the response containing an error message from PAC
// https://factura.com/apidocs/crear-cfdi-40.html#ejemplo-de-respuesta-de-error
func parseErrorMessage(responseBody []byte) error {

	var errMsgA struct {
		Response string `json:"response"`
		Message  string `json:"message"`
		XML      string `json:"xml"`
		UID      string `json:"uid"`
		UUID     string `json:"uuid"`
	}

	if err := json.Unmarshal(responseBody, &errMsgA); err == nil {
		return &CFDIEngineError{Code: ReqMalformed, Message: errMsgA.Message}
	}

	var errMsgB struct {
		Response string `json:"response"`
		Message  struct {
			Message       string  `json:"message"`
			MessageDetail string  `json:"messageDetail"`
			Data          *string `json:"data"` // Data can be null
			Status        string  `json:"status"`
		} `json:"message"`
		XMLError string `json:"xmlerror"`
	}

	if err := json.Unmarshal(responseBody, &errMsgB); err == nil {
		return &CFDIEngineError{Code: ReqMalformed, Message: errMsgB.Message.Message}
	}

	// If everyone has failed, then return an error indicating unknown structure
	return &CFDIEngineError{Code: ReqMalformed, Message: "Unsupported PAC's response"}
}
