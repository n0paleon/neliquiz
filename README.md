# NeliQuiz

NeliQuiz is a multiplayer quiz game built with Go, featuring a modular architecture, REST API, real-time game server, and PostgreSQL database. It supports solo and broadcast game modes, and is ready for deployment with Docker Compose.

## Features

- RESTful API for managing questions and categories
- Real-time multiplayer quiz game server (Pitaya-based)
- Modular, clean architecture (domain, usecase, delivery layers)
- PostgreSQL database with migration scripts
- Docker Compose for easy deployment
- Unit tests for core logic

## Getting Started

### Prerequisites

- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)
- (Optional) Go 1.24+ for local development
- [Golang Migrate](https://github.com/golang-migrate/migrate)
- [AWS Cognito](https://aws.amazon.com/cognito/) (Optional) for production usecase
- [AWS API Gateway](https://aws.amazon.com/api-gateway/) (Optional) for production usecase

### Running with Docker Compose

1. **Clone the repository:**
   ```sh
   git clone https://github.com/n0paleon/neliquiz.git
   cd neliquiz
   ```
2. **Migrate database schema:**
   ```sh
   make migrate-up
   ```
3. **Run service:**
    ```sh
   docker-compose up -d
   ```
4. **Access the Backend:**
    - API: http://localhost:3000/
    - Game server: ws://localhost:3250 
    - Postman collection: [Postman Collection](https://turu-developer.postman.co/workspace/Turu-Developer-Workspace~c575b45e-1a7e-4bd8-8026-b9c87dbf6eae/collection/32865624-58030ff2-6ba7-4dd0-b9e3-7cd0f6467538?action=share&source=copy-link&creator=32865624)
