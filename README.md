# Wallet

Wallet is a project for managing wallets and transaction history.

## Installation

Clone the repository:

```bash
git clone https://github.com/geejjoo/wallet.git
```

Navigate to the cloned directory:
```bash
cd wallet
```

Install dependencies:
```bash
go mod tidy
```

## Usage
### Running the Project
To run the project, execute the following command in your terminal:
```bash
docker compose up --build
```
By default, the project will be available at `http://localhost:8000`.


## Project Structure
- `cmd/`: Directory for main executable files.
- `pkg/`: Directory for packages used in the application.
  - `handler/`: Package containing HTTP request handlers.
  - `repository/`: Package containing the implementation of database operations.
  - `service/`: Package containing the business logic of the application.
- `configs/`: Directory for configuration files.
- `schema/`: Directory for database initialization.

