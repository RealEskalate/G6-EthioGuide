# EthioGuide - Backend API Service

<div align="center">
  <img src="https://ethio-guide.vercel.app/images/ethioguide-symbol.png" alt="EthioGuide Logo" width="200">
</div>

[![Build Status](https://github.com/Ethio-Guide/ethio-guide-backend/actions/workflows/go-ci.yml/badge.svg)](https://github.com/Ethio-Guide/ethio-guide-backend/actions)

---

### **Repository Context**

This repository contains the backend service for the EthioGuide project. It is a Go-based RESTful API responsible for all data management, business logic, and authentication for the platform.

**For the full platform source code (including frontend and mobile), please see the main project hub:** [**EthioGuide Platform Repository**](https://github.com/RealEskalate/G6-EthioGuide)

---

## 📚 Table of Contents

- [Live API Documentation](#-live-api-documentation)
- [Tech Stack](#-tech-stack)
- [Getting Started (Local Development)](#-getting-started-local-development)
  - [Option 1: Running with Docker (Recommended & Easiest)](#option-1-running-with-docker-recommended--easiest)
  - [Option 2: Running with Go Locally (For Active Development)](#option-2-running-with-go-locally-for-active-development)
- [Running Tests](#-running-tests)
- [Project Structure (Clean Architecture)](#-project-structure-clean-architecture)
- [CI/CD Pipeline](#-cicd-pipeline)
- [Implemented Functionalities](#-implemented-functionalities)
- [Contributing](#-contributing)

---

## 🚀 Live API Documentation

The complete API reference, including all endpoints, request/response models, and authorization details, is available via our live Swagger UI.

<div align="center">
  <strong><a href="https://ethio-guide-backend.onrender.com/swagger/index.html">View the Live API Documentation</a></strong>
</div>

## 🛠️ Tech Stack

This service is built with a focus on performance, concurrency, and maintainability.

| Category       | Technology                                                                     |
|----------------|--------------------------------------------------------------------------------|
| **Language**   | [**Go (Golang)**](https://go.dev/) (v1.24+)                                      |
| **Framework**  | [**Gin**](https://gin-gonic.com/) (High-performance HTTP web framework)          |
| **Database**   | [**MongoDB**](https://www.mongodb.com/) (via the official Go driver)             |
| **Caching**    | [**Redis**](https://redis.io/) (For session management and caching)              |
| **AI/ML**      | [**Google Gemini API**](https://ai.google.dev/) (For chat and translation)       |
| **Embeddings** | [**Cohere API**](https://cohere.com/) (For semantic search vectors)              |
| **Container**  | [**Docker & Docker Compose**](https://www.docker.com/) (For local development)   |
| **Deployment** | [**Render**](https://render.com/)                                              |

## ⚙️ Getting Started (Local Development)

You have two primary options for running the backend service locally. Docker is the easiest and recommended method.

---

### **Option 1: Running with Docker (Recommended & Easiest)**

This method builds a container with the Go environment and all dependencies. **You do not need to have Go installed on your machine for this method.**

#### Prerequisites

*   [Git](https://git-scm.com/)
*   [Docker](https://www.docker.com/products/docker-desktop/) and [Docker Compose](https://docs.docker.com/compose/install/)

#### Installation Steps

1.  **Clone the repository:**
    ```sh
    git clone https://github.com/Ethio-Guide/ethio-guide-backend
    cd backend-repo
    ```

2.  **Configure Environment Variables:**
    Create a `.env` file for your configuration.
    ```sh
    cp .env.example .env
    ```
    Now, open the `.env` file and fill in your secret keys, especially `JWT_SECRET`, `GEMINI_API_KEY`, and `HF_EMBEDDING_API_KEY`.

3.  **Configure for Local Swagger Documentation:**
    To ensure the API docs work locally, you must make one small change before building the container. Open the file `./delivery/main.go` and modify the `@host` annotation:
    ```go
    // FROM:
    // @host      ethio-guide-backend.onrender.com

    // TO:
    // @host      localhost:8080
    ```
    > **Note:** This change is for local development only. Do not commit this file. You can easily revert it later with `git checkout -- ./delivery/main.go`.

4.  **Build and Run with Docker Compose:**
    This single command will build the Go application image, start the MongoDB and Redis containers, and run the API.
    ```sh
    docker-compose up --build
    ```

5.  **Access the API:**
    *   **API available at:** `http://localhost:8080`
    *   **Local Swagger Docs:** `http://localhost:8080/swagger/index.html`

---

### **Option 2: Running with Go Locally (For Active Development)**

This method is ideal if you are actively developing the backend, as it provides faster feedback loops without rebuilding a Docker image for every change.

#### Prerequisites

*   [Git](https://git-scm.com/)
*   [Go (Golang)](https://go.dev/doc/install) version 1.24 or higher
*   [Swag CLI](https://github.com/swaggo/swag) (for regenerating API docs)
*   **A running instance of MongoDB and Redis.**

#### Installation Steps

1.  **Clone & Configure:**
    Follow steps 1 and 2 from the Docker method above (clone repo, set up `.env` file).

2.  **Ensure Dependencies are Running:**
    You must have MongoDB and Redis running and accessible from your local machine.

    *   **Option A: Use Docker for Dependencies (Recommended)**
        This is the simplest way to get consistent database versions. Run the database and cache in the background:
        ```sh
        docker-compose up -d mongo redis
        ```
        The default settings in the `.env` file will connect to these containers.

    *   **Option B: Use a Local Installation**
        If you have MongoDB and Redis installed locally (e.g., via Homebrew or as system services), ensure they are started. The default `.env` settings (`mongodb://localhost:27017` and `localhost:6379`) should work for standard local installations.

    *   **Option C: Use Cloud Services**
        You can use a cloud provider like MongoDB Atlas for your database and a managed Redis service (e.g., from Upstash or your cloud provider).
        **Important:** You must update your `.env` file with the connection strings provided by your cloud service.
        ```ini
        # Example for MongoDB Atlas
        MONGO_URI="mongodb+srv://<user>:<password>@cluster.mongodb.net/?retryWrites=true&w=majority"

        # Example for a cloud Redis provider
        REDIS_URI="rediss://<user>:<password>@host:port"
        # Or if separate details are provided:
        # REDIS_ADDR="host:port"
        # REDIS_PASSWORD="your-password"
        ```

3.  **Install Go Dependencies:**
    ```sh
    go mod tidy
    ```

4.  **Configure and Generate Local Swagger Docs:**
    *   **Install Swag CLI:**
        ```sh
        go install github.com/swaggo/swag/cmd/swag@latest
        ```
    *   **Modify `main.go`:** Open `./delivery/main.go` and change the `@host` to `localhost:8080` as described in the Docker instructions.
    *   **Regenerate Docs:** Run the following command from the root directory:
        ```sh
        swag init -g ./delivery/main.go
        ```

5.  **Run the Go Application:**
    ```sh
    go run ./delivery/main.go
    ```
    The server will start on port `8080`.

> **Warning:** Remember to revert the changes to `main.go` before committing your code: `git checkout -- ./delivery/main.go`.

## 🧪 Running Tests

This project has a comprehensive test suite. To run all unit and integration tests, use the following command:
```sh
go test -v ./...
```
*   `-v`: Enables verbose output to see which tests are running.

**NOTE:** The test specific env variables will be used for tests (`MONGO_URI_TEST`, `REDDIS_ADDR_TEST` and `REDDIS_DB_TEST`).

## 🏗️ Project Structure (Clean Architecture)

The codebase is organized following the principles of **Clean Architecture** to ensure separation of concerns, testability, and maintainability.

*   `/config`: Handles loading and parsing of environment variables.
*   `/delivery`: The presentation layer. Contains controllers for handling HTTP requests and the main router setup.
*   `/domain`: The core of the application. Contains business models, interfaces, and domain-specific logic. It has no external dependencies.
*   `/infrastructure`: Contains concrete implementations of external services like JWT, email, OAuth, and AI services.
*   `/repository`: The data access layer. Implements interfaces defined in `/domain` to interact with the MongoDB database.
*   `/usecase`: The application logic layer. Orchestrates the flow of data between the domain, repositories, and external services to execute business rules.

## 🔄 CI/CD Pipeline

This project utilizes a modern CI/CD pipeline to automate testing and deployment, ensuring code quality and rapid delivery. The process is split between GitHub Actions for integration and Render for deployment.

### **Continuous Integration (CI)**

*   **Platform:** GitHub Actions
*   **Workflow File:** `.github/workflows/Go CI.yml`

On every push or pull request to the `main` branch, the CI pipeline is automatically triggered to validate the code. The process includes:
1.  Setting up a clean Ubuntu environment.
2.  Starting MongoDB and Redis services within the CI runner for integration tests.
3.  Installing the specified version of Go.
4.  Downloading and caching Go module dependencies.
5.  Building the application to ensure there are no compilation errors.
6.  **Executing the full test suite** (`go test -v -race ./...`) against the live test database services.

This ensures that no broken code is merged into the main branch.

### **Continuous Deployment (CD)**

*   **Platform:** Render
*   **Trigger:** A successful push or merge to the `main` branch.

Continuous Deployment is managed directly by Render's auto-deploy feature, which is linked to this GitHub repository. The deployment flow is as follows:

1.  A developer pushes new code to the `main` branch.
2.  The CI pipeline on GitHub Actions runs first to verify the changes.
3.  Upon a successful push to `main`, Render detects the new commit.
4.  Render automatically triggers a new deployment. It pulls the latest code, builds a new application image using the `Dockerfile`, and deploys the new version, typically with zero downtime.

This creates a seamless, automated path from writing code to deploying it in the live production environment.

## ✨ Implemented Functionalities

This API provides a rich set of features, including:

#### User & Authentication
*   User registration and secure password login.
*   Social login via Google.
*   JWT-based authentication with access/refresh token rotation.
*   Email-based account verification and password reset flows.
*   Full CRUD for user profiles and preferences.

#### Organizations & Procedures
*   Admin-only creation of organization accounts.
*   Full CRUD operations for service procedures (by Admins or owning Orgs).
*   Advanced search and filtering for procedures (by name, cost, time, etc.).

#### Community & Feedback
*   Users can submit detailed feedback on any procedure.
*   Admins/Orgs can manage and respond to feedback.
*   A complete discussion forum (posts) with CRUD operations.

#### AI & Machine Learning
*   AI-powered chat guide using Google Gemini and RAG for accurate, context-aware answers about procedures.
*   On-the-fly JSON content translation to user's preferred language.
*   Semantic search capabilities powered by Cohere embeddings.

#### Checklists & Search
*   Users can create personal checklists to track their progress on procedures.
*   A powerful, application-wide search endpoint to find organizations and procedures.

## 🤝 Contributing

We welcome contributions! Please feel free to open an issue to report a bug or suggest a feature. If you would like to contribute code, please fork the repository and submit a pull request.

1.  Fork the Project
2.  Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3.  Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4.  Push to the Branch (`git push origin feature/AmazingFeature`)
5.  Open a Pull Request