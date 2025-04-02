# TalUp

## Installation Guide

### Prerequisites

- Go (latest version recommended)
- [PostgreSQL](https://www.postgresql.org/download/) (if using a database)

### Clone the Repository

```sh
git clone https://github.com/heyoSSam/TalUp-backend
cd your-repo
```

### Install Dependencies

```sh
go mod tidy
```

### Setup Environment Variables

Create a `.env` file in the project root and configure it as follows:

```ini
DATABASE_URL=""
ROBERTA_URL=""
PORT=""
JWT_SECRET=""
```

### Run the Project

```sh
go run main.go
```

### API Documentation

You can generate Swagger documentation by running:

```sh
swag init
```

Then, access the documentation at:

```
http://localhost:8080/swagger/index.html
```

