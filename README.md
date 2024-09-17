# Go REST API with MySQL

This is a simple REST API built using Go and MySQL. It includes CRUD operations (Create, Read, Update, Delete) for managing products in a MySQL database.

## Prerequisites

- Go (>= 1.16)
- MySQL
- Go modules for dependencies:
    - [github.com/gorilla/mux](https://pkg.go.dev/github.com/gorilla/mux) for routing.
    - [github.com/go-sql-driver/mysql](https://pkg.go.dev/github.com/go-sql-driver/mysql) for MySQL connection.

## Setup MySQL Database

1. Create a MySQL database and table using the following SQL commands:

   ```sql
   CREATE DATABASE productdb;

   USE productdb;

   CREATE TABLE products (
       id INT AUTO_INCREMENT PRIMARY KEY,
       name VARCHAR(255) NOT NULL,
       price DECIMAL(10, 2) NOT NULL,
       description TEXT
   );
   
2. Update the MySQL connection details in the main.go file:

go

```dsn := "root:password@tcp(127.0.0.1:3306)/productdb ```

Replace root:password with your actual MySQL username and password.


3. How to Run
   Clone the repository:


```
git clone https://github.com/yourusername/go-mysql-api.git
cd go-mysql-api 
```
4. Install dependencies:


```
go mod tidy
```
5. Run the application:



```
go run main.go
```

The server will run at http://localhost:8000.


6. API Endpoints

GET /products: Retrieve all products.
GET /products/{id}: Retrieve a single product by its ID.
POST /products: Create a new product.

Example request body:
json

```
{
"name": "Laptop",
"price": 1000.00,
"description": "A new laptop"
}
```

PUT : /products/{id}: 
Update a product by ID.


Example request body:
json

```
{
"name": "Updated Laptop",
"price": 1200.00,
"description": "An updated laptop"
}
```
DELETE /products/{id}: Delete a product by ID.
Example Requests Using Curl

Get all products:

```
curl http://localhost:8000/products
```
Get a product by ID:

```
curl http://localhost:8000/products/1
```

Create a new product:

```
curl -X POST -H "Content-Type: application/json" -d '{"name":"Laptop","price":1000,"description":"A new laptop"}' http://localhost:8000/products

```

Update a product:

```
curl -X PUT -H "Content-Type: application/json" -d '{"name":"Updated Laptop","price":1200,"description":"Updated"}' http://localhost:8000/products/1
```

Delete a product:
```
curl -X DELETE http://localhost:8000/products/1
```

test main_test.go
```
go test -v
```
