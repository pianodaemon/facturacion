package verification

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"blaucorp.com/fiscal-engine/internal/dal"
	"blaucorp.com/fiscal-engine/internal/dal/models"
)

func TestWithMongoDBContainer(t *testing.T) {

	ctx := context.Background()

	mongoC, err := setupMongoContainer(ctx)
	if err != nil {
		t.Fatalf("failed to start container: %s", err)
	}
	defer mongoC.Terminate(ctx)

	mongoURI, err := getMongoURI(ctx, mongoC)
	if err != nil {
		t.Fatalf("Failed to retrieve MongoDB URI: %s", err)
	}

	var client *mongo.Client
	err = dal.SetUpConnMongoDB(&client, mongoURI)
	if err != nil {
		t.Fatalf("failed to set up mongo client: %s", err)
	}
	defer client.Disconnect(ctx)

	db := client.Database("receipts_db")

	verifyReceiptStorage(t, db)
}

// Helper function to set up MongoDB container
func setupMongoContainer(ctx context.Context) (testcontainers.Container, error) {
	req := testcontainers.ContainerRequest{
		Image:        "mongo:6.0",
		ExposedPorts: []string{"27017/tcp"},
		WaitingFor:   wait.ForLog("Waiting for connections").WithStartupTimeout(30 * time.Second),
	}

	return testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
}

// Helper function to get MongoDB URI
func getMongoURI(ctx context.Context, mongoC testcontainers.Container) (string, error) {
	host, err := mongoC.Host(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get container host: %w", err)
	}

	port, err := mongoC.MappedPort(ctx, "27017")
	if err != nil {
		return "", fmt.Errorf("failed to get container port: %w", err)
	}

	return fmt.Sprintf("mongodb://%s:%s", host, port.Port()), nil
}

func verifyReceiptStorage(t *testing.T, db *mongo.Database) {

	receipt := models.Receipt{
		Owner:            "MX Seller SA de CV", // Example seller's name
		ReceptorRFC:      "XAXX010101000",      // Receptor RFC (generic for testing)
		ReceptorDataRef:  "Customer Ref 123",   // Some customer reference
		DocumentCurrency: "MXN",                // Mexican Peso (MXN)
		BaseCurrency:     "MXN",                // Base currency as MXN
		ExchangeRate:     1.0,                  // Exchange rate for same currency (MXN to MXN)
		SubtotalAmount:   1000.00,              // Total before taxes
		TotalTransfers:   160.00,               // Total IVA transfers
		TotalDeductions:  0.00,                 // No deductions
		TotalAmount:      1160.00,              // Total amount including IVA
		Purpose:          "Sale of goods",      // Purpose of the factura
		PaymentWay:       "PUE",                // Payment way: "PUE" (Pago en una sola exhibición - One-time payment)
		PaymentMethod:    "03",                 // Payment method: "03" (Bank transfer)
		Items: []models.ReceiptItem{
			{
				ProductID:        "PROD-001", // Internal product identifier
				ProductDesc:      "A limpia", // Interal product description
				ProductQuantity:  2,          // Two units sold
				ProductUnitPrice: 500,        // Price per unit
				ProductAmount:    1000,       // Total amount before taxes
				ProductTransfers: []models.Tax{
					{
						Base:         1000,   // Base amount for tax calculation
						Rate:         0.16,   // IVA rate 16% (AKA TasaOCuota sat)
						Amount:       160,    // Calculated tax amount (1000 * 0.16) (AKA Importe sat)
						FiscalFactor: "Tasa", // Factor (AKA TipoFactor sat)
						FiscalType:   "002",  // Tax type (AKA Impuesto sat)
						Transfer:     true,   // True if it's a transfer tax (AKA Traslado sat)
					},
				},
				ProductDeductions: nil,        // No deductions in this example
				FiscalProductID:   "81111600", // Mexican fiscal code for the product/service (AKA ClaveProdServ)
				FiscalProductUnit: "E48",      // SAT unit code (AKA ClaveUnidad)
			},
		},
	}

	// Call CreateReceipt function and check if it returns the generated ObjectID
	oid, err := dal.CreateReceipt(db, &receipt)
	assert.NoError(t, err)
	assert.False(t, oid.IsZero(), "Returned ObjectID should not be empty")

	// Get the test database and collection
	collection := db.Collection("receipts")

	// Verify the receipt was inserted into the collection
	var result models.Receipt
	err = collection.FindOne(context.TODO(), bson.M{"_id": oid}).Decode(&result)
	assert.NoError(t, err)

	// Assert that the inserted receipt matches the original one
	assert.Equal(t, receipt.Owner, result.Owner)
	assert.Equal(t, receipt.ReceptorRFC, result.ReceptorRFC)
	assert.Equal(t, receipt.TotalAmount, result.TotalAmount)
	assert.Equal(t, len(receipt.Items), len(result.Items))
}
