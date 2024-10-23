### Create receipt

This command creates a new receipt for specific receptor.

```sh

pianodaemon@LAPTOP-4RSVIK4C:~$ curl -X POST http://localhost:8080/receipts \
-H "Content-Type: application/json" \
-d '{
  "owner": "viajes Ponchito",
  "receptor_rfc": "ABC123456789",
  "receptor_data_ref": "CustomerRef123",
  "document_currency": "MXN",
  "base_currency": "MXN",
  "exchange_rate": 1.0,
  "items": [
    {
      "product_id": "Product001",
      "product_desc": "Product Description",
      "product_quantity": 5,
      "product_unit_price": 100.50,
      "product_transfers": [
        {
          "rate": 0.16,
          "fiscal_factor": "IVA",
          "fiscal_type": "transfer",
          "transfer": true
        }
      ],
      "product_deductions": []
    }
  ],
  "purpose": "sale",
  "payment_way": "credit_card",
  "payment_method": "PUE"
}'
```

**Explanation**:

-   This `POST` request creates a new receipt featuring serie and folio `A1`, owned by `viajes Ponchito` for `MXN` currency.

### Update receipt

This command updates an existing receipt for specific receptor.

```sh

pianodaemon@LAPTOP-4RSVIK4C:~$ curl -X PUT http://localhost:8080/receipts/60d5cabcabc1234567890def \
-H "Content-Type: application/json" \
-d '{
  "owner": "UpdatedCompanyName",
  "receptor_rfc": "ABC123456789",
  "receptor_data_ref": "UpdatedCustomerRef123",
  "document_currency": "USD",
  "base_currency": "USD",
  "exchange_rate": 20.0,
  "items": [
    {
      "product_id": "Product002",
      "product_desc": "Updated Product Description",
      "product_quantity": 10,
      "product_unit_price": 200.75,
      "product_transfers": [
        {
          "rate": 0.16,
          "fiscal_factor": "IVA",
          "fiscal_type": "transfer",
          "transfer": true
        }
      ],
      "product_deductions": []
    }
  ],
  "purpose": "purchase",
  "payment_way": "bank_transfer",
  "payment_method": "PUE"
}'
```

**Explanation**:

- This PUT request updates the receipt
