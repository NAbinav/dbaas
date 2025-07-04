# Database as a Service (DBaaS)

A RESTful API service that provides simplified database operations for developers. This service allows you to perform CRUD operations on your database tables through intuitive HTTP endpoints with query parameter filtering.

## Features

- **Simple Database Operations**: Create, read, update, and delete operations through REST endpoints
- **Query Parameter Filtering**: Filter data using URL parameters with comparison operators
- **Google OAuth2 Authentication**: Secure authentication through Google accounts
- **API Key Management**: Generate and manage API keys for service access
- **Table Management**: Create and manage database tables dynamically
- **Developer-Friendly**: Simplified interface for rapid application development

## Authentication

### Google OAuth2 Login
Access the login endpoint to authenticate with your Google account:
```
GET /login
```

### API Key Generation
After authentication, generate your API key:
```
GET /newApiKey
```

### Using API Keys
Include your API key in requests to access protected endpoints. All database operations require valid API key authentication.

## API Endpoints

### Table Management

#### Create Table
```http
POST /create/{table_name}
```
Creates a new database table with the specified name.

#### Delete Table
```http
DELETE /delete/{table_name}
```
Removes the specified table from the database.

### Data Operations

#### Insert Data
```http
POST /{table_name}
```
Adds new records to the specified table.

#### Retrieve Data
```http
GET /{table_name}/{column}
```
Retrieves data from the specified table and column with optional filtering.

#### Update Data
```http
PUT /{table_name}
```
Updates existing records in the specified table. Supports filtering to update specific records.

#### Delete Records
```http
DELETE /{table_name}
```
Removes records from the specified table. Supports filtering to delete specific records.

## Query Parameters

The service supports advanced filtering through URL query parameters using comparison operators:

### Supported Operators

- `_gt`: Greater than
- `_gte`: Greater than or equal to
- `_lt`: Less than
- `_lte`: Less than or equal to
- `_eq`: Equal to (default if no operator specified)
- `_ne`: Not equal to

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

## Getting Started

### Prerequisites
- Go 1.18 or higher
- Database connection (configuration required)
- Google OAuth2 credentials

### Installation

1. Clone the repository
2. Install dependencies:
   ```bash
   go mod tidy
   ```

3. Configure your database connection in the `db` package
4. Set up Google OAuth2 credentials in the `auth` package
5. Run the service:
   ```bash
   go run main.go
   ```

The service will start on port 8081.

### Configuration

Before running the service, ensure you have:

1. **Database Configuration**: Configure your database connection parameters
2. **Google OAuth2 Setup**: 
   - Create a Google Cloud Console project
   - Enable Google+ API
   - Create OAuth2 credentials
   - Configure authorized redirect URIs to include your callback URL

## Usage Flow

1. **Authentication**: Visit `/login` to authenticate with Google
2. **API Key**: Generate your API key at `/newApiKey`
3. **Table Creation**: Create tables using `POST /create/{table_name}`
4. **Data Operations**: Perform CRUD operations using your API key
5. **Query Data**: Use query parameters for filtered data retrieval

## Security

- All database operations require valid API key authentication
- Google OAuth2 provides secure user authentication
- Middleware enforces authentication on protected endpoints

## Error Handling

The service returns appropriate HTTP status codes:
- `200`: Success
- `400`: Bad Request (invalid parameters)
- `401`: Unauthorized (invalid or missing API key)
- `404`: Not Found (table or resource doesn't exist)
- `500`: Internal Server Error

## Support

For issues, questions, or feature requests, please refer to the project documentation or contact the development team.

## License

This project is licensed under the MIT License. See the LICENSE file for details.
