# Go E-Commerce Backend

## Features

- User registration & login with JWT authentication
- Product CRUD operations (admin only)
- Cart functionality (add, remove, view items)
- Order management (create, view, cancel, update orders)
- Stripe payment integration (intent + webhook)
- Role-based authorization (admin/user)


## Environment Variables

Create a `.env` file with:

```env
DATABASE_URL=your_postgres_connection_string
STRIPE_SECRET_KEY=your_stripe_secret_key
JWT_SECRET=your_jwt_secret_key
```

## API Endpoints
#### Auth Routes

| Method | Endpoint   | Description        |
|--------|------------|--------------------|
| POST   | /register  | User registration  |
| POST   | /login     | User login         |
---
#### Product Routes
| Method | Endpoint                      | Description               |
|--------|-------------------------------|---------------------------|
| GET    | /api/products                 | List all products         |
| GET    | /api/products/{id}           | Get product by ID         |
| POST   | /api/admin/products          | Add new product (admin)   |
| PUT    | /api/admin/products/{id}     | Update product (admin)    |
| DELETE | /api/admin/products/{id}     | Delete product (admin)    |
---
#### Cart Routes

| Method | Endpoint                      | Description             |
|--------|-------------------------------|-------------------------|
| POST   | /api/cart                     | Add item to cart        |
| GET    | /api/cart                     | View cart               |
| DELETE | /api/cart/{product_id}        | Remove item from cart   |
---
#### Order Routes
| Method | Endpoint                              | Description             |
|--------|----------------------------------------|-------------------------|
| POST   | /api/order                             | Create new order        |
| GET    | /api/orders                            | List user's orders      |
| GET    | /api/orders/{id}                       | View order details      |
| DELETE | /api/orders/{id}/cancel                | Cancel order            |
| PUT    | /api/admin/orders/{id}/status          | Update order status (admin) |
---
#### Payment Routes
| Method | Endpoint                 | Description                    |
|--------|--------------------------|--------------------------------|
| POST   | /api/create-payment-intent | Create Stripe payment intent |
| POST   | /api/webhook             | Stripe webhook endpoint        |


## Test Flow:

    Register/Login
    Add product (if admin)
    Add to cart
    Create order
    Call /create-payment-intent
    Simulate webhook /api/webhook