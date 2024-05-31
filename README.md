# Price Tracker REST API

Price Tracker REST API is a Go-based web service that tracks product prices from ebuyer.com e-commerce website. It allows users to monitor price changes over time and fetches the latest prices daily using a scheduled task.

## Features

- Track product prices from specified URLs.
- Schedule daily price checks at a specific time.
- RESTful endpoints to manage tracked products and their price histories.
- MongoDB integration for storing product and price data.

## Prerequisites

- Go 1.16 or later
- MongoDB

## Installation

1. **Clone the repository:**

    ```bash
    git clone https://github.com/KemalBekir/price-tracker-rest-api-go.git
    cd price-tracker-rest-api-go
    ```

2. **Install Go dependencies:**

    ```bash
    go mod tidy
    ```

3. **Set up environment variables:**

    Create a `.env` file in the root directory and configure your MongoDB connection string:

    ```env
    MONGODB_URI=mongodb://localhost:27017
    DB_NAME=priceTracker
    ```

## Running the Application

### Backend

1. **Run the backend server:**

    ```bash
    go run main.go
    ```

    The server will start on `http://localhost:5000`.

### Frontend

If you have a frontend application (assuming it runs on `http://localhost:5173`), make sure it is running.

## API Endpoints

### Product Management

- **Get All Products**

    ```
    GET /
    ```

- **Get Product by ID**

    ```
    GET /:id
    ```

- **Add New Product**

    ```
    POST /scrape
    Body: {
      "url": "string"
    }
    ```
    
## Scheduled Tasks

The application uses a cron job to update the prices of all tracked products daily at 23:00.

The cron job is configured in the `main.go` file:

```go
c := cron.New()
_, err = c.AddFunc("0 23 * * *", func() {
    err := services.UpdatePrices(searchesCollection, pricesCollection)
    if err != nil {
        log.Printf("Error updating pricing: %v", err)
    } else {
        log.Println("Price update completed successfully.")
    }
})
c.Start()
```

## Code Structure
- main.go: Entry point of the application, initializes the server and the cron job.
- internal/db: Database connection and setup.
- internal/router: HTTP router and endpoints.
- internal/services: Business logic, including the price update logic.
- internal/handler: HTTP handlers for managing products and scraping.
- model: Data models for MongoDB collections.