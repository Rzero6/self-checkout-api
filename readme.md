# 📦 API Documentation

Base URL: /api

## 🔐 Login

Endpoint: POST /api/login
Headers:

- Content-Type: application/json
  Body:
  {"username": "string","password": "string"}
  Response:
  {"status": true,"message": "string","token": "string"}

## 🛍️ Product API

### Get All Products

Endpoint: GET /api/products?page={page}&limit={limit}
Query Params:

- page: int
- limit: int
  Response:
  {"success": true,"message": "string","data": [{"id": int, "barcode": "string","name": "string","price": 0}]}

### Search Product by Barcode

Endpoint: GET /api/products/search
Query Params:

- barcode: string
  Response:
  {"status": true,"message": "string","data": {"id": int, "barcode": "string","name": "string","price": 0}}

### Get Random Product

Endpoint: GET /api/products/random
Response:
{"status": true,"message": "string","data": {"id": int, "barcode": "string","name": "string","price": 0}}

### Create Product

Endpoint: POST /api/products
Headers:

- Authorization: Bearer <token>
- Content-Type: application/json
  Body:
  {"name": "string","price": 0}
  Response:
  {"status": true,"message": "string","data": {"id": int, "barcode": "string","name": "string","price": 0}}

## 🛒 Cart API

### Start Cart

Endpoint: POST /api/cart/
Headers:
Response:
{
"success": true, "message": string,
"session_id": string
}

### Get Current Cart

Endpoint: GET /api/cart/
Headers:
- X-Session-ID: string
  Response:
  {
  "success": true,
  "message": string,
  "data": {
  "cart": {"id": int,"session_id": string, "status": string},
  "details": [{"id": int, "product_name": string, "price": int, "quantity":int, "subtotal": int}],
  },
  }

### Add Product to Cart

Endpoint: POST /api/cart/add
Headers:

- X-Session-ID: string
- Content-Type: application/json
  Body:
  {
  "cart_id": "<fill_here>",
  "product_id": "<fill_here>",
  "quantity": "<fill_here>"
  }
  Response:
  {
  "message": "<fill_here>"
  }

### Update Cart Detail

Endpoint: PATCH /api/cart/update
Headers:

- X-Session-ID: string
- Content-Type: application/json
  Body:
  {
  "detail_id": "<fill_here>",
  "quantity": "<fill_here>"
  }
  Response:
  {
  "message": "<fill_here>"
  }

### Delete Cart Detail

Endpoint: DELETE /api/cart/detail/:id
Headers:

- X-Session-ID: string
  Params:
- id: <fill_here>
  Response:
  {
  "message": "<fill_here>"
  }

### Delete All Products from Cart

Endpoint: DELETE /api/cart/:id
Headers:

- X-Session-ID: string
  Params:
- id: <fill_here>
  Response:
  {
  "message": "<fill_here>"
  }

## 💳 Payment API

### Create Transaction

Endpoint: POST /api/payment
Headers:

- Authorization: <fill_here>
- Content-Type: application/json
  Body:
  {
  "cart_id": "<fill_here>",
  "payment_method": "<fill_here>"
  }
  Response:
  {
  "order_id": "<fill_here>",
  "payment_url": "<fill_here>"
  }

### Get Transaction Status

Endpoint: GET /api/payment/:order_id
Headers:

- Authorization: <fill_here>
  Params:
- order_id: <fill_here>
  Response:
  {
  "status": "<fill_here>"
  }

### Cancel Transaction

Endpoint: PATCH /api/payment/:order_id/cancel
Headers:

- Authorization: <fill_here>
  Params:
- order_id: <fill_here>
  Response:
  {
  "message": "<fill_here>"
  }

### Midtrans Notification (Webhook)

Endpoint: POST /api/payment/notification
Headers:

- Content-Type: application/json
- X-Signature-Key: <fill_here>
  Body:
  {
  "order_id": "<fill_here>",
  "transaction_status": "<fill_here>",
  "fraud_status": "<fill_here>",
  "signature_key": "<fill_here>"
  }
  Response:
  {
  "message": "<fill_here>"
  }
