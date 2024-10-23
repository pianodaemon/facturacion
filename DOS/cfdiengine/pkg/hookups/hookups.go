package hookups

import (
	"errors"

	"blaucorp.com/fiscal-engine/internal/dal"
	"blaucorp.com/fiscal-engine/internal/dal/models"

	"go.mongodb.org/mongo-driver/mongo"
)

type (
	FiscalEngine struct {
		dbID string
		mcli *mongo.Client
	}
)

func NewFiscalEngine(mongoURI string) *FiscalEngine {

	pm := &FiscalEngine{}

	pm.dbID = "receipts_db"
	err := dal.SetUpConnMongoDB(&(pm.mcli), mongoURI)
	if err != nil {
		panic(err.Error())
	}

	return pm
}

func (self *FiscalEngine) DoCreateReceipt(dto *ReceiptDTO) (string, error) {

	if len(dto.Items) == 0 {
		return "", errors.New("receipt must have at least one item")
	}

	var subtotalAmount, totalTransfers, totalDeductions, totalAmount float64

	// Convert ReceiptItemDTO to ReceiptItem and calculate totals
	items := make([]models.ReceiptItem, len(dto.Items))
	for i, itemDTO := range dto.Items {
		productAmount := itemDTO.ProductQuantity * itemDTO.ProductUnitPrice
		subtotalAmount += productAmount

		// Convert product transfers and deductions, and calculate tax amounts
		transfers := convertTaxes(itemDTO.ProductTransfers, productAmount)
		deductions := convertTaxes(itemDTO.ProductDeductions, productAmount)

		// Calculate total transfers and deductions for the item
		for _, tax := range transfers {
			totalTransfers += tax.Amount
		}
		for _, tax := range deductions {
			totalDeductions += tax.Amount
		}

		items[i] = models.ReceiptItem{
			ProductID:         itemDTO.ProductID,
			ProductDesc:       itemDTO.ProductDesc,
			ProductQuantity:   itemDTO.ProductQuantity,
			ProductUnitPrice:  itemDTO.ProductUnitPrice,
			ProductAmount:     productAmount,
			ProductTransfers:  transfers,
			ProductDeductions: deductions,
			FiscalProductID:   itemDTO.FiscalProductID,
			FiscalProductUnit: itemDTO.FiscalProductUnit,
		}
	}

	// Calculate total amount
	totalAmount = subtotalAmount + totalTransfers - totalDeductions

	// Create the model
	receipt := &models.Receipt{
		Owner:            dto.Owner,
		ReceptorRFC:      dto.ReceptorRFC,
		ReceptorDataRef:  dto.ReceptorDataRef,
		DocumentCurrency: dto.DocumentCurrency,
		BaseCurrency:     dto.BaseCurrency,
		ExchangeRate:     dto.ExchangeRate,
		SubtotalAmount:   subtotalAmount,
		TotalTransfers:   totalTransfers,
		TotalDeductions:  totalDeductions,
		TotalAmount:      totalAmount,
		Items:            items,
		Purpose:          dto.Purpose,
		PaymentWay:       dto.PaymentWay,
		PaymentMethod:    dto.PaymentMethod,
	}

	db := self.mcli.Database(self.dbID)
	id, err := dal.CreateReceipt(db, receipt)
	if err != nil {
		return "", err
	}
	return id.Hex(), err
}

func (self *FiscalEngine) DoEditReceipt(receiptID string, dto *ReceiptDTO) error {

	if len(dto.Items) == 0 {
		return errors.New("receipt must have at least one item")
	}

	db := self.mcli.Database(self.dbID)
	existingReceipt, err := dal.GetReceiptByID(db, receiptID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return errors.New("receipt not found")
		}
		return err
	}

	// Update the fields of the existing receipt with the new dto values
	var subtotalAmount, totalTransfers, totalDeductions, totalAmount float64

	// Convert ReceiptItemDTO to ReceiptItem and calculate totals
	items := make([]models.ReceiptItem, len(dto.Items))
	for i, itemDTO := range dto.Items {
		productAmount := itemDTO.ProductQuantity * itemDTO.ProductUnitPrice
		subtotalAmount += productAmount

		// Convert product transfers and deductions, and calculate tax amounts
		transfers := convertTaxes(itemDTO.ProductTransfers, productAmount)
		deductions := convertTaxes(itemDTO.ProductDeductions, productAmount)

		// Calculate total transfers and deductions for the item
		for _, tax := range transfers {
			totalTransfers += tax.Amount
		}
		for _, tax := range deductions {
			totalDeductions += tax.Amount
		}

		items[i] = models.ReceiptItem{
			ProductID:         itemDTO.ProductID,
			ProductDesc:       itemDTO.ProductDesc,
			ProductQuantity:   itemDTO.ProductQuantity,
			ProductUnitPrice:  itemDTO.ProductUnitPrice,
			ProductAmount:     productAmount,
			ProductTransfers:  transfers,
			ProductDeductions: deductions,
			FiscalProductID:   itemDTO.FiscalProductID,
			FiscalProductUnit: itemDTO.FiscalProductUnit,
		}
	}

	// Calculate total amount
	totalAmount = subtotalAmount + totalTransfers - totalDeductions

	// Update the fields of the existing receipt
	existingReceipt.Owner = dto.Owner
	existingReceipt.ReceptorRFC = dto.ReceptorRFC
	existingReceipt.ReceptorDataRef = dto.ReceptorDataRef
	existingReceipt.DocumentCurrency = dto.DocumentCurrency
	existingReceipt.BaseCurrency = dto.BaseCurrency
	existingReceipt.ExchangeRate = dto.ExchangeRate
	existingReceipt.SubtotalAmount = subtotalAmount
	existingReceipt.TotalTransfers = totalTransfers
	existingReceipt.TotalDeductions = totalDeductions
	existingReceipt.TotalAmount = totalAmount
	existingReceipt.Items = items
	existingReceipt.Purpose = dto.Purpose
	existingReceipt.PaymentWay = dto.PaymentWay
	existingReceipt.PaymentMethod = dto.PaymentMethod

	// Save the updated receipt back to the database
	err = dal.EditReceipt(db, existingReceipt)
	if err != nil {
		return err
	}

	return nil
}

// convertTaxes converts a slice of TaxDTO to a slice of Tax, calculates amounts using Base and Rate
func convertTaxes(taxesDTO []TaxDTO, base float64) []models.Tax {
	taxes := make([]models.Tax, len(taxesDTO))
	for i, taxDTO := range taxesDTO {
		taxAmount := base * taxDTO.Rate // Calculate tax amount using Base * Rate
		taxes[i] = models.Tax{
			Base:         base, // Base comes from ProductAmount
			Rate:         taxDTO.Rate,
			FiscalFactor: taxDTO.FiscalFactor,
			FiscalType:   taxDTO.FiscalType,
			Transfer:     taxDTO.Transfer,
			Amount:       taxAmount, // Calculated tax amount
		}
	}
	return taxes
}
