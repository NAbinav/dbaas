# Database as a Service (DBaaS)

A RESTful API service that provides simplified database operations for developers. This service allows you to perform CRUD operations on your database tables through intuitive HTTP endpoints with query parameter filtering.

## Features

- **Simple Database Operations**: Create, read, update, and delete operations through REST endpoints
- **Query Parameter Filtering**: Filter data using URL parameters with comparison operators
- **Google OAuth2 Authentication**: Secure authentication through Google accounts
- **API Key Management**: Generate API keys for service access
- **Table Management**: Create and manage database tables dynamically
- **Developer-Friendly**: Simplified interface for rapid application development

## Authentication

### API Key Generation
After authentication, generate your API key:
```
GET /newApiKey
```

### Using API Keys
Include your API key in the `X-API-KEY` header for all requests to access protected endpoints. All database operations require valid API key authentication.

```bash
curl -H "X-API-KEY: your_api_key" \
     https://dboss.brogramiz.info/users/name
```

## API Endpoints

### Authentication Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/login` | Initiate Google OAuth2 login |
| `GET` | `/callback` | OAuth2 callback handler |
| `GET` | `/newApiKey` | Generate new API key |

### Table Management

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/create/{table_name}` | Create new table |
| `DELETE` | `/delete/{table_name}` | Delete table |

### Data Operations

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/{table_name}/{column}` | Retrieve data with optional filtering |
| `POST` | `/{table_name}` | Insert new records |
| `PUT` | `/{table_name}` | Update records with filtering |
| `DELETE` | `/{table_name}` | Delete records with filtering |

## Query Parameters

The service supports advanced filtering through URL query parameters using comparison operators:

### Supported Operators

| Operator | Description | Example |
|----------|-------------|---------|
| `_gt` | Greater than | `age_gt=25` |
| `_gte` | Greater than or equal to | `price_gte=100` |
| `_lt` | Less than | `date_lt=2024-01-01` |
| `_lte` | Less than or equal to | `weight_lte=75` |
| `_eq` | Equal to (default if no operator specified) | `status_eq=active` |
| `_ne` | Not equal to | `category_ne=deprecated` |

### Examples

#### Basic Filtering
```http
GET /users/username?age_gt=25&weight_lte=75
```
Returns usernames from the users table where age > 25 and weight ≤ 75.

#### Multiple Conditions
```http
GET /products/name?price_gte=100&price_lte=500&category_eq=electronics
```
Returns product names where price is between 100-500 and category is electronics.

#### Single Condition
```http
GET /orders/order_id?status_eq=completed
```
Returns order IDs where status equals "completed".

#### Update with Filtering
```http
PUT /users?age_gt=18&status_eq=active
```
Updates records in the users table where age > 18 and status equals "active".

#### Delete with Filtering
```http
DELETE /orders?status_eq=cancelled&date_lt=2024-01-01
```
Deletes orders where status equals "cancelled" and date is before 2024-01-01.

## Table Schema Definition

When creating a table, you must provide the table schema in the request body as JSON. The schema defines the column names and their data types.

### Request Body Format
```json
{
  "id": "int",
  "name": "string",
  "email": "string",
  "age": "int",
  "created_at": "datetime",
  "is_active": "boolean"
}
```

### Supported Data Types

| Type | Description | SQL Mapping |
|------|-------------|-------------|
| `string` | Variable-length text | `VARCHAR(255)` |
| `text` | Long text content | `TEXT` |
| `int` | Integer numbers | `INT` |
| `big int` | Large integer numbers | `BIGINT` |
| `decimal` | Decimal numbers | `DECIMAL(10,2)` |
| `float` | Floating point numbers | `FLOAT` |
| `boolean` | True/false values | `BOOLEAN` |
| `date` | Date values | `DATE` |
| `time` | Time values | `TIME` |
| `datetime` | Date and time | `DATETIME` |
| `timestamp` | Timestamp values | `TIMESTAMP` |
| `json` | JSON data | `JSON` |
| `uuid` | UUID values | `CHAR(36)` |
| `auto-increment int` | Auto-incrementing ID | `SERIAL` |

### Example: Creating a Users Table
```bash
curl -X POST "https://dboss.brogramiz.info/create/users" \
     -H "X-API-KEY: your_api_key" \
     -H "Content-Type: application/json" \
     -d '{
       "id": "auto-increment int",
       "name": "string",
       "email": "string",
       "age": "int",
       "created_at": "datetime",
       "is_active": "boolean"
     }'
```

## Data Insertion

When inserting data into a table, provide the record data as JSON in the request body.

### Request Body Format
```json
{
  "name": "John Doe",
  "email": "john@example.com",
  "age": 30,
  "created_at": "2024-01-15 10:30:00",
  "is_active": true
}
```

### Example: Inserting a User
```bash
curl -X POST "https://dboss.brogramiz.info/users" \
     -H "X-API-KEY: your_api_key" \
     -H "Content-Type: application/json" \
     -d '{
       "name": "John Doe",
       "email": "john@example.com",
       "age": 30,
       "created_at": "2024-01-15 10:30:00",
       "is_active": true
     }'
```

### Inserting Multiple Records
```bash
curl -X POST "https://dboss.brogramiz.info/users" \
     -H "X-API-KEY: your_api_key" \
     -H "Content-Type: application/json" \
     -d '[
       {
         "name": "John Doe",
         "email": "john@example.com",
         "age": 30,
         "is_active": true
       },
       {
         "name": "Jane Smith",
         "email": "jane@example.com",
         "age": 25,
         "is_active": true
       }
     ]'
```

## Getting Started

The service is live and ready to use at: **https://dboss.brogramiz.info**

No installation required - simply start using the API endpoints directly.

## Usage Flow

1. **API Key**: Generate your API key at `https://dboss.brogramiz.info/newApiKey`
2. **Table Creation**: Create tables using `POST https://dboss.brogramiz.info/create/{table_name}`
3. **Data Operations**: Perform CRUD operations using your API key
4. **Query Data**: Use query parameters for filtered data retrieval

## Security

- All database operations require valid API key authentication
- Google OAuth2 provides secure user authentication
- Middleware enforces authentication on protected endpoints

## Support

For issues, questions, or feature requests, please refer to the project documentation or contact the development team.

## License

This project is licensed under the MIT License. See the LICENSE file for details.
