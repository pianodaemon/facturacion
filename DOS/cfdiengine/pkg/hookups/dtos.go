package hookups

type (
	ReceiptItemDTO struct {
		ProductID         string   `json:"product_id" binding:"required"`
		ProductDesc       string   `json:"product_desc" binding:"required"`
		ProductQuantity   float64  `json:"product_quantity" binding:"required"`
		ProductUnitPrice  float64  `json:"product_unit_price" binding:"required"`
		ProductTransfers  []TaxDTO `json:"product_transfers"`
		ProductDeductions []TaxDTO `json:"product_deductions"`
		FiscalProductID   string   `json:"fiscal_product_id" binding:"required"`
		FiscalProductUnit string   `json:"fiscal_product_unit" binding:"required"`
	}

	ReceiptDTO struct {
		Owner            string           `json:"owner" binding:"required"`
		ReceptorRFC      string           `json:"receptor_rfc" binding:"required"`
		ReceptorDataRef  string           `json:"receptor_data_ref" binding:"required"`
		DocumentCurrency string           `json:"document_currency" binding:"required"`
		BaseCurrency     string           `json:"base_currency" binding:"required"`
		ExchangeRate     float64          `json:"exchange_rate" binding:"required"`
		Items            []ReceiptItemDTO `json:"items" binding:"required"`
		Purpose          string           `json:"purpose" binding:"required"`
		PaymentWay       string           `json:"payment_way" binding:"required"`
		PaymentMethod    string           `json:"payment_method" binding:"required"`
	}

	TaxDTO struct {
		Rate         float64 `json:"rate" binding:"required"`
		FiscalFactor string  `json:"fiscal_factor" binding:"required"`
		FiscalType   string  `json:"fiscal_type" binding:"required"`
		Transfer     bool    `json:"transfer" binding:"required"`
	}
)
