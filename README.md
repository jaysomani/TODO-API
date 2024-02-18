# Todo API Project

This is a Todo API project developed using Golang and ScyllaDB, providing basic CRUD operations and pagination functionality for the list endpoint.

## Objective

The objective of this project was to create a robust Todo API that allows users to manage their tasks efficiently. The API supports creating, reading, updating, and deleting Todo items for individual users. Additionally, pagination functionality is implemented to retrieve Todo items in batches, with support for filtering based on status.

## Requirements

- Set up a Golang project and integrate ScyllaDB as the database for storing Todo items. Ensure items are stored user-wise.
- Implement endpoints for CRUD operations for Todo items, with properties including id, user_id, title, description, status, created, and updated.
- Implement a paginated list endpoint to retrieve Todo items.
- Support filtering based on Todo item status (e.g., pending, completed).

## Basic flow of the project
![basic flow drawio](https://github.com/jaysomani/TODO-API/assets/69755312/d2969231-4c1d-4b94-a074-7b89df3b64fe)


## Basic architecture diagram 
![ARCHI drawio](https://github.com/jaysomani/TODO-API/assets/69755312/7033f11f-3d9a-4c2a-a092-546e185d2e51)


## Getting Started

To clone and use this project on your local machine, follow these steps:

1. Clone the repository:

```bash
git clone https://github.com/jaysomani/todo-api.git
```
```bash
cd todo-api
```
```bash
go mod download
```
```bash
go run main.go
```

It should ben running on 
```bash
http://localhost:8080
```
