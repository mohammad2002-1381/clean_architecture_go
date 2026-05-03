#!/bin/bash

BASE_URL="http://localhost:9000/api/v1"

echo "========================================="
echo "1. Creating a new product..."
echo "========================================="
CREATE_RESPONSE=$(curl -s -X POST $BASE_URL/products \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Mechanical Keyboard",
    "description": "RGB mechanical gaming keyboard",
    "price": 149.99,
    "items": [
      {
        "name": "Keycap Puller",
        "description": "Tool for removing keycaps"
      },
      {
        "name": "USB Cable",
        "description": "Braided USB-C cable"
      }
    ]
  }')

echo $CREATE_RESPONSE | jq '.'

# Extract product ID (assuming it's in the response)
PRODUCT_ID=$(echo $CREATE_RESPONSE | jq -r '.id // 1')

echo -e "\n========================================="
echo "2. Getting product by ID: $PRODUCT_ID"
echo "========================================="
curl -s -X GET $BASE_URL/products/$PRODUCT_ID | jq '.'

echo -e "\n========================================="
echo "3. Getting all products"
echo "========================================="
curl -s -X GET $BASE_URL/products | jq '.'

echo -e "\n========================================="
echo "4. Updating product ID: $PRODUCT_ID"
echo "========================================="
curl -s -X PUT $BASE_URL/products/$PRODUCT_ID \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Premium Mechanical Keyboard",
    "description": "RGB mechanical gaming keyboard with Cherry MX switches",
    "price": 199.99
  }' | jq '.'

echo -e "\n========================================="
echo "5. Verifying the update"
echo "========================================="
curl -s -X GET $BASE_URL/products/$PRODUCT_ID | jq '.'

echo -e "\n========================================="
echo "6. Testing error: Invalid product ID"
echo "========================================="
curl -s -X GET $BASE_URL/products/invalid-id | jq '.'

echo -e "\n========================================="
echo "7. Testing error: Create product with invalid data"
echo "========================================="
curl -s -X POST $BASE_URL/products \
  -H "Content-Type: application/json" \
  -d '{
    "price": "invalid_price"
  }' | jq '.'

echo -e "\n========================================="
echo "8. Testing error: Update non-existent product"
echo "========================================="
curl -s -X PUT $BASE_URL/products/99999 \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Non-existent Product",
    "price": 99.99
  }' | jq '.'