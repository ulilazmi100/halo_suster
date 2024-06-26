# Halo Suster

## 🌄 Background

Halo Suster is a highly performant health record backend application designed to efficiently manage patient medical records. Developed as part of Project Sprint Batch 2 Week 3, this project showcases the ability to deliver high-quality, scalable applications quickly, with rigorous testing and load testing.

Utilizing Golang with labstack/echo framework, without using ORM just pure SQL, without external caching (this is part of the challange) with pgx/sqlx database library.

## Additional Information

### MAIN/PGX version
In the main folder it's using pgx library directly to communicate with PostgreSQL database without using sqlx.

### SQLX Version

Besides the main folder, there is also an sqlx version available in the `halo_suster_sqlx` folder, providing an alternative implementation using sqlx for database interactions.


---

## 🚀 Getting Started

### Prerequisites

- Go (1.22 or later)
- PostgreSQL
- [K6](https://k6.io/docs/get-started/installation/) for load testing
- WSL (Windows Subsystem for Linux) if on Windows

### Installation and Setup

1. **Clone the repository:**

    ```bash
    git clone https://github.com/ulilazmi100/halo_suster.git
    cd halo_suster
    ```

2. **Set up environment variables:**

    Create a `.env` file or export these variables in your shell:

    ```bash
    export DB_NAME=your_db_name
    export DB_PORT=5432
    export DB_HOST=localhost
    export DB_USERNAME=your_db_user
    export DB_PASSWORD=your_db_password
    export DB_PARAMS="sslmode=disable"  # or "sslrootcert=rds-ca-rsa2048-g1.pem&sslmode=verify-full" for production
    export JWT_SECRET=your_jwt_secret
    export BCRYPT_SALT=10  # or higher for production

    # S3 to upload, all uploaded files will be available for only a day
    export AWS_ACCESS_KEY_ID=
    export AWS_SECRET_ACCESS_KEY=
    export AWS_S3_BUCKET_NAME=projectsprint-bucket-public-read
    export AWS_REGION=ap-southeast-1
    ```

3. **Run database migrations:**

    ```bash
    migrate -database "postgres://$DB_USERNAME:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?$DB_PARAMS" -path db/migrations up
    ```

4. **Start the server:**

    ```bash
    go run main.go
    ```

5. **Optional: Rollback migrations if needed:**

    ```bash
    migrate -database "postgres://$DB_USERNAME:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?$DB_PARAMS" -path db/migrations down
    ```

### Running with Docker

1. **Build Docker image:**

    ```bash
    docker build -t halo_suster .
    ```

2. **Run Docker container:**

    ```bash
    docker run -p 8080:8080 --env-file .env halo_suster
    ```

---

## 🧪 Testing

### Prerequisites

- [K6](https://k6.io/docs/get-started/installation/)
- Linux environment (WSL/MacOS)

### Running Tests

1. **Set environment variable:**

    ```bash
    export BASE_URL=http://localhost:8080
    ```

2. **Run the server.**

3. **Navigate to the `tests` folder and run the tests:**

    #### Regular testing

    ```bash
    BASE_URL=http://localhost:8080 make run
    ```

    #### Load testing

    Prepare the NIP Generator:

    ```bash
    # this requires at minimum, go 1.22.3
    # run in separate terminal
    make run-go
    ```

    Then run:

    ```bash
    BASE_URL=http://localhost:8080 make run-load-test
    ```

    #### Regular testing with debug (output to `output.txt`)

    ```bash
    BASE_URL=http://localhost:8080 make run-debug
    ```

### Tests Result
- My personal test results can be seen in the `test_results` folder

---

## Database Migration

Use [golang-migrate](https://github.com/golang-migrate/migrate) for managing database migrations:

- Create migration:

    ```bash
    migrate create -ext sql -dir db/migrations add_user_table
    ```

- Execute migration:

    ```bash
    migrate -database "postgres://$DB_USERNAME:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?$DB_PARAMS" -path db/migrations up
    ```

- Rollback migration:

    ```bash
    migrate -database "postgres://$DB_USERNAME:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?$DB_PARAMS" -path db/migrations down
    ```

---

## 📝 Requirements

For detailed functional requirements, please refer to the [Notion page](https://openidea-projectsprint.notion.site/HaloSuster-be1866776fe84c2d8d9eac08ce09b7a5).

---

## 👥 Contributing

1. Fork the repository.
2. Create your feature branch (`git checkout -b feature/AmazingFeature`).
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`).
4. Push to the branch (`git push origin feature/AmazingFeature`).
5. Open a pull request.

---

## 📝 License

This project is licensed under the MIT License.

---

## 📚 Resources

- **Notion:** [Halo Suster's Requirements' Notion Page](https://openidea-projectsprint.notion.site/HaloSuster-be1866776fe84c2d8d9eac08ce09b7a5)
- **Tests:** [Project Sprint Batch 2 Week 3 Test Cases](https://github.com/nandanugg/EniQiloStoreTestCasesPSW2B2?tab=readme-ov-file#for-load-testing)
- **Migrations:** [Golang Migration](https://github.com/golang-migrate/migrate)

---

## 📞 Contact

[Muhammad Ulil 'Azmi](https://github.com/ulilazmi100) - [@M_Ulil_Azmi](https://twitter.com/M_Ulil_Azmi)

Project Link: [https://github.com/ulilazmi100/halo_suster](https://github.com/ulilazmi100/halo_suster)

---

## Note

Please be aware that the Amazon service subscription may have already ended, which could result in pipeline failures in GitHub.

---