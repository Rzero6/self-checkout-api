# 📦 API Documentation

Base URL: /api

## 🔐 Login

Endpoint:
POST /api/login
User Authentication and to get JWT Token.

Request Body:

```json
{
  "username":"string",
  "password":"string"
}
```

Response:

```json
{
  "success" :true,
  "message" :"Login berhasil",
  "token"   :"jwt_token_string"
}
```

Authorization:Bearer [token]

## 🛍️ Product API

### Get All Products

Endpoint: GET /api/products?page={page}&limit={limit}

Response:

```json
{
  "success":true,
  "message":"Get data Products berhasil",
  "data":[{
    "id":1,
    "name":"string",
    "price":0,
    "barcode":"string",
    }],
  "meta":{
    "page":1,
    "per_page":10,
    "total_items":100,
    "total_pages":10
    },
  "links":{
    "next":"/api/products?page=2&limit=10",
    "prev":null
    }
}
```

### Search Product by Barcode

Endpoint: GET /api/products/search

Response:

```json
{
  "success":true,
  "message":"Get data Product berhasil",
  "data":{
    "id":1,
    "name":"string",
    "price":0,
    "barcode":"string",
    }
}
```

### Get Random Product

Endpoint: GET /api/products/random

Response:

```json
{
  "success":true,
  "message":"Get data Product berhasil",
  "data":{
    "id":1,
    "name":"string",
    "price":0,
    "barcode":"string",
    }
}
```

### Create Product

Endpoint: POST /api/products

Headers:

```json
Authorization: "Bearer [token]"
```

Request Body:

```json
{
  "name"  :"string",
  "price" :0
}
```

Response:

```json
{
  "success":true,
  "message":"Create data Product berhasil",
  "data":{
    "id"      :1,
    "name"    :"string",
    "price"   :0,
    "barcode" :"string",
    }
}
```

## 🛒 Cart API

### Start Cart

Endpoint: POST /api/cart/

Response:

```json
{
  "success"     :true,
  "message"     :"Create Cart berhasil",
  "data"        :{
    "id"          :1,
    "session_id"  :"string",
    "status"      :"string",
    }
}
```

### Get Current Cart

Endpoint: GET /api/cart/

Headers:

```json
X-Session-ID: "uuid"
```

Response:

```json
{
  "success":true,
  "message":"Get Cart berhasil",
  "data":{
    "id"          :1,
    "session_id"  :"string",
    "status"      :"string",
    }
}
```

### Delete Cart

Endpoint: DELETE /api/cart/{cart_id}

Headers:

```json
X-Session-ID: "uuid"
```

Response:

```json
{
  "success":true,
  "message":"Delete Cart berhasil",
  "data":{
    "id"          :1,
    "session_id"  :"string",
    "status"      :"string",
    }
}
```

### Add Product to Cart

Endpoint: POST /api/cart/detail

Headers:

```json
X-Session-ID: "uuid"
```

Request Body:

```json
{
  "barcode"  :"string",
  "quantity" :0
}
```

Response:

```json
{
  "success":true,
  "message":"Create Detail berhasil",
  "data":{
    "id"          :1,
    "cart_id"     :1,
    "product_id"  :1,
    "Product_name":"string",
    "price"       :0,
    "quantity"    :0,
    "subtotal"    :0,
    }
}
```

### Update Cart Detail

Endpoint: PATCH /api/cart/detail

Headers:

```json
X-Session-ID: "uuid"
```

Request Body:

```json
{
  "id"        :1,
  "quantity"  :0
}
```

Response:

```json
{
  "success":true,
  "message":"Update Detail berhasil",
  "data":{
    "id"          :1,
    "cart_id"     :1,
    "product_id"  :1,
    "Product_name":"string",
    "price"       :0,
    "quantity"    :0,
    "subtotal"    :0,
    }
}
```

### Delete Cart Detail

Endpoint: DELETE /api/cart/detail/{id}

Headers:

```json
X-Session-ID: "uuid"
```

Response:

```json
{
  "success":true,
  "message":"Delete Detail berhasil",
  "data":{
    "id"          :1,
    "cart_id"     :1,
    "product_id"  :1,
    "Product_name":"string",
    "price"       :0,
    "quantity"    :0,
    "subtotal"    :0,
    }
}
```

### Get All Products Detail from Cart

Endpoint: DELETE /api/cart/details

Headers:

```json
X-Session-ID: "uuid"
```

Response:

```json
{
  "success":true,
  "message":"Update Detail berhasil",
  "data":{
            "id"          :1,
            "cart_id"     :1,
            "product_id"  :1,
            "Product_name":"string",
            "price"       :0,
            "quantity"    :0,
            "subtotal"    :0,
          },
          {
            "id"          :1,
            "cart_id"     :1,
            "product_id"  :1,
            "Product_name":"string",
            "price"       :0,
            "quantity"    :0,
            "subtotal"    :0,
          },
}
```

## 💳 Payment API

### Create Transaction

Endpoint: POST /api/transaction

Headers:

```json
X-Session-ID: "uuid"
```

Response:

```json
{
  "success":true,
  "message":"Create Transaction berhasil",
  "data":{
    "id"            :1,
    "cart_id"       :1,
    "order_id"      :"string",
    "status"        :"string",
    "amount"        :0,
    "payment_type"  :"string",
    "qris_link"     :"string",
    "expire_time"   :"string",
    }
}
```

### Get Transaction Status

Endpoint: GET /api/transaction/{order_id}

Headers:

```json
X-Session-ID: "uuid"
```

Response:

```json
{
  "success":true,
  "message":"Get Transaction berhasil",
  "data":{
    "id"            :1,
    "cart_id"       :1,
    "order_id"      :"string",
    "status"        :"string",
    "amount"        :0,
    "payment_type"  :"string",
    "qris_link"     :"string",
    "expire_time"   :"string",
    }
}
```

### Get Transaction Details

Endpoint: GET /api/transaction/{order_id}/details

Headers:

```json
X-Session-ID: "uuid"
```

Response:

```json
{
  "success":true,
  "message":"Get Transaction Details berhasil",
  "data":{
    "id"              :1,
    "transaction_id"  :1,
    "product_id"      :1,
    "product_name"    :"string",
    "price"           :0,
    "quantity"        :0,
    "subtotal"        :0,
    }
}
```

### Cancel Transaction

Endpoint: POST /api/payment/{order_id}/cancel

Headers:

```json
X-Session-ID: "uuid"
```

Response:

```json
{
  "success":true,
  "message":"Transaction cancelled",
  "data":{
    "id"            :1,
    "cart_id"       :1,
    "order_id"      :"string",
    "status"        :"string",
    "amount"        :0,
    "payment_type"  :"string",
    "qris_link"     :"string",
    "expire_time"   :"string",
    }
}
```

### Midtrans Notification (Webhook)

Endpoint: POST /api/transaction/notification
