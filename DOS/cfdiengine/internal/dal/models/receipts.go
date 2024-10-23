package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	Tax struct {
		Base         float64 `json:"base" bson:"base,omitempty"`
		Rate         float64 `json:"rate" bson:"rate,omitempty"`
		Amount       float64 `json:"amount" bson:"amount,omitempty"`
		FiscalFactor string  `json:"fiscal_factor" bson:"fiscal_factor,omitempty"`
		FiscalType   string  `json:"fiscal_type" bson:"fiscal_type,omitempty"`
		Transfer     bool    `json:"transfer" bson:"transfer,omitempty"`
	}

	Receipt struct {
		ID               primitive.ObjectID `json:"_id" bson:"_id,omitempty"`
		Owner            string             `json:"owner" bson:"owner,omitempty"`
		ReceptorRFC      string             `json:"receptor_rfc" bson:"receptor_rfc,omitempty"`
		ReceptorDataRef  string             `json:"receptor_data_ref" bson:"receptor_data_ref,omitempty"`
		DocumentCurrency string             `json:"document_currency" bson:"document_currency,omitempty"`
		BaseCurrency     string             `json:"base_currency" bson:"base_currency,omitempty"`
		ExchangeRate     float64            `json:"exchange_rate" bson:"exchange_rate,omitempty"`
		SubtotalAmount   float64            `json:"subtotal_amount" bson:"subtotal_amount,omitempty"`
		TotalTransfers   float64            `json:"total_transfers" bson:"total_transfers,omitempty"`
		TotalDeductions  float64            `json:"total_deductions" bson:"total_deductions,omitempty"`
		TotalAmount      float64            `json:"total_amount" bson:"total_amount,omitempty"`
		Items            []ReceiptItem      `json:"items" bson:"items,omitempty"`
		Purpose          string             `json:"purpose" bson:"purpose,omitempty"`
		PaymentWay       string             `json:"payment_way" bson:"payment_way,omitempty"`
		PaymentMethod    string             `json:"payment_method" bson:"payment_method,omitempty"`
	}

	ReceiptItem struct {
		ProductID         string  `json:"product_id" bson:"product_id,omitempty"`
		ProductDesc       string  `json:"product_desc" bson:"product_desc,omitempty"`
		ProductQuantity   float64 `json:"product_quantity" bson:"product_quantity,omitempty"`
		ProductUnitPrice  float64 `json:"product_unit_price" bson:"product_unit_price,omitempty"`
		ProductAmount     float64 `json:"product_amount" bson:"product_amount,omitempty"`
		ProductTransfers  []Tax   `json:"product_transfers" bson:"product_transfers,omitempty"`
		ProductDeductions []Tax   `json:"product_deductions" bson:"product_deductions,omitempty"`
		FiscalProductID   string  `json:"fiscal_product_id" bson:"fiscal_product_id,omitempty"`
		FiscalProductUnit string  `json:"fiscal_product_unit" bson:"fiscal_product_unit,omitempty"`
	}
)
