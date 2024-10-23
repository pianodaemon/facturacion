package dal

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"blaucorp.com/fiscal-engine/internal/dal/models"
)

// CreateReceipt inserts a new Receipt into MongoDB.
func CreateReceipt(db *mongo.Database, receipt *models.Receipt) (primitive.ObjectID, error) {
	collection := db.Collection("receipts")

	// Ensure that the receipt doesn't already have an ID (new documents should not have one)
	if receipt.ID != primitive.NilObjectID {
		return primitive.NilObjectID, errors.New("receipt already has an ID; cannot create")
	}

	// Generate a new ObjectID for the new receipt
	receipt.ID = primitive.NewObjectID()

	// Insert the new receipt document into the collection
	result, err := collection.InsertOne(context.TODO(), receipt)
	if err != nil {
		return primitive.NilObjectID, err
	}

	// Return the generated ObjectID
	oid, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return primitive.NilObjectID, fmt.Errorf("failed to retrieve generated ObjectID")
	}
	return oid, nil
}

// EditReceipt updates an existing Receipt in MongoDB.
func EditReceipt(db *mongo.Database, receipt *models.Receipt) error {
	collection := db.Collection("receipts")

	if receipt.ID == primitive.NilObjectID {
		return errors.New("receipt ID is required for editing")
	}

	// Filter to find the receipt by ID
	filter := bson.M{"_id": receipt.ID}

	// Update the fields of the receipt document
	update := bson.M{
		"$set": bson.M{
			"owner":             receipt.Owner,
			"receptor_rfc":      receipt.ReceptorRFC,
			"receptor_data_ref": receipt.ReceptorDataRef,
			"document_currency": receipt.DocumentCurrency,
			"base_currency":     receipt.BaseCurrency,
			"exchange_rate":     receipt.ExchangeRate,
			"subtotal_amount":   receipt.SubtotalAmount,
			"total_transfers":   receipt.TotalTransfers,
			"total_deductions":  receipt.TotalDeductions,
			"total_amount":      receipt.TotalAmount,
			"items":             receipt.Items, // Updating the items with taxes
			"purpose":           receipt.Purpose,
			"payment_way":       receipt.PaymentWay,
			"payment_method":    receipt.PaymentMethod,
		},
	}

	// Perform the update without upsert
	result, err := collection.UpdateOne(context.TODO(), filter, update, options.Update().SetUpsert(false))
	if err != nil {
		return err
	}

	// If no documents were matched, return an error
	if result.MatchedCount == 0 {
		return errors.New("no receipt found with the given ID")
	}

	return nil
}

// GetReceiptByID retrieves a Receipt by its ID (as a string) from MongoDB.
func GetReceiptByID(db *mongo.Database, receiptID string) (*models.Receipt, error) {
	collection := db.Collection("receipts")

	// Convert the string ID to ObjectID
	objID, err := primitive.ObjectIDFromHex(receiptID)
	if err != nil {
		return nil, errors.New("invalid receipt ID format")
	}

	// Create a context with a timeout for the database operation
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Prepare the filter to find the receipt by ID
	filter := bson.M{"_id": objID}
	receipt := &models.Receipt{}

	// Find the receipt in the collection
	err = collection.FindOne(ctx, filter).Decode(receipt)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("receipt not found")
		}
		return nil, err
	}

	return receipt, nil
}
