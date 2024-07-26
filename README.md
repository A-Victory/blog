
# Blog API

Welcome to the Blog API project! This RESTful API application is designed to manage a blog with features including user authentication, post management, and comment handling. This README provides comprehensive instructions for setup, usage, and testing.

## Features

- **User Management**: Register, log in, and retrieve user profile information.
- **Post Management**: Create, read, update, and delete blog posts.
- **Comment Management**: Create, read, update, and delete comments on posts.
- **Pagination**: Implemented for retrieving lists of posts and comments.
- **Search Functionality**: Added to the comments retrieval endpoint for filtering by content.
- **Authentication**: JWT-based user authentication for secure access.
- **Authorization**: Ensures users can only modify their own posts and comments.
- **Deployment**: Dockerized application with deployment configurations for AWS EC2.

## Setup and Installation

### Prerequisites

- [Docker](https://docs.docker.com/get-docker/) - For containerization.
- [Docker Compose](https://docs.docker.com/compose/install/) - For managing multi-container Docker applications.
- [Go](https://go.dev/doc/install) - For local development and building the application.

### Database Setup

1. **Database Configuration**:
   - Create a `.env` file in the root directory of the project.
   - Add the following environment variables to configure your database:

     ```env
     DB_URI=mysql://username:password@tcp(localhost:3306)/
     DB_NAME=testdb
     ```

### Building and Running

1. **Build Docker Image**:

   ```bash
   docker build -t blog-api -f blog-api.dockerfile .
   ```

2. **Run with Docker Compose**:

   ```bash
   docker-compose up --build
   ```

   This command will start the application along with the database and other required services. Access the API at `http://localhost:80`.

### Configuration

- **Environment Variables**: Configure required environment variables in the `.env` file or through your cloud provider's configuration settings.

### Running Tests

1. **Unit and Integration Tests**:

   Ensure your database is set up and accessible as per the configuration in your test environment.

   ```bash
   go test ./...
   ```

   This command will run all tests, including those for CRUD operations and pagination functionality.

## API Endpoints

### User Endpoints

- **POST** `/api/users/register` - Register a new user.
  - **Request Body**:

    ```json
    {
      "username": "john_doe",
      "email": "john@example.com",
      "password": "securepassword"
    }
    ```

- **POST** `/api/users/login` - Authenticate user and receive a JWT token.
  - **Request Body**:

    ```json
    {
      "username": "john_doe",
      "password": "securepassword"
    }
    ```

  - **Response**:

    ```json
    {
        "status": 200,
        "message": "login successful",
        "data": {
            "authorization": "your generated token is eyJhbGciOiJIU......"
        }
    }
    ```

- **GET** `/api/users/profile` - Get user profile information (Authenticated).
  - **Headers**: `Authorization: token`

### Post Endpoints

- **GET** `/api/posts` - Retrieve all posts (Paginated).
  - **Example Request**: `GET /api/posts?page=1&limit=10`

  - **Response**:

    ```
    {
     "status": 200,
    "message": "success",
    "data": {
        "posts": [
            {
                "id": 1,
                "title": "My first post",
                "content": "Once upon a time...",
                "authorId": 1,
                "createdAt": "2024-07-26 07:56:12",
                "updatedAt": "2024-07-26 07:56:12"
            },
            ....
        ] 
    }
    ```

- **GET** `/api/posts/{id}` - Retrieve a single post by ID.

- **POST** `/api/posts` - Create a new post (Authenticated).
  - **Request Body**:

    ```json
    {
      "title": "My First Post",
      "content": "This is the content of the post."
    }
    ```

- **PUT** `/api/posts/{id}` - Update a post by ID (Authenticated & Author only).
  - **Request Body**:

    ```json
    {
      "title": "Updated Post Title",
      "content": "Updated content."
    }
    ```

- **DELETE** `/api/posts/{id}` - Delete a post by ID (Authenticated & Author only).

### Comment Endpoints

- **GET** `/api/posts/{postId}/comments` - Retrieve all comments for a post (Paginated).
  - **Example Request**: `GET /api/posts/1/comments?page=1&limit=10`

- **POST** `/api/posts/{postId}/comments` - Create a new comment on a post (Authenticated).
  - **Request Body**:

    ```json
    {
      "content": "Great post!"
    }
    ```

- **PUT** `/api/comments/{id}` - Update a comment by ID (Authenticated & Author only).
  - **Request Body**:

    ```json
    {
      "content": "Updated comment."
    }
    ```

- **DELETE** `/api/comments/{id}` - Delete a comment by ID (Authenticated & Author only).

## Pagination and Search

- **Pagination**: Implemented for both posts and comments. Use query parameters `page` and `limit` to control pagination.
- **Search**: For post, you can use the `search` query parameter to filter posts by content.

## Deployment

### Docker Configuration

- **Dockerfile**: `blog-api.dockerfile` - Defines the application image. Ensure it includes all dependencies and configurations.
- **Docker Compose**: `docker-compose.yml` - Manages the application and database services.

### AWS EC2 Deployment

The application is configured to run on an AWS EC2 instance. For the binary application to run best on AWS, some changes would have to be made to both the app and on the server. 

## Testing and Postman Collection

For a comprehensive set of tests and sample data, refer to the Postman collection linked below:

- **[Postman Collection](https://documenter.getpostman.com/view/21553602/2sA3kYjLKM)**

This collection includes example requests and responses for all API endpoints, which will help you verify and test the functionality of the API.

## Notes

- Update the `.env` file or environment settings with the correct database URL and other configurations for different environments (development, staging, production).
- For further clarity, please contact me at vuictory.agbabune@gmail.com.

---

This application was created in response to interview request from Intel Region.