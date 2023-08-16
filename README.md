# Chirpy Backend Go API

Chirpy is a social network similar to Twitter. This project hosts the backend API for Chirpy, allowing users to create chirps, automatically validate content, manage users, and handle authentication using JWT tokens.

## Features

- User Management: Create and authenticate users; update user details securely.
- Chirp Management: Create, retrieve, and delete chirps; associate with user IDs.
- Chirp Validation: Ensure chirp content meets criteria and handle profane words.
- Token-based Authentication: Generate JWT tokens, refresh tokens, and revoke them.
- Webhooks Integration: Process Polka webhooks for user events and upgrades.
- Error Handling: Provide clear error messages and appropriate status codes.

## Getting Started

Follow these steps to get your development environment up and running.

### Prerequisites

- [Go](https://golang.org/doc/install) installed on your machine.
- [Git](https://git-scm.com/downloads) installed.

### Setup

To set up the Chirpy Backend Go API in your local development environment, follow these steps:

1. Clone this repository:

   ```bash
   git clone https://github.com/abhishekghimire40/chirpy_go_server
   cd chirpy_go_server
   ```

2. Install dependencies:

   ```bash
    go mod download
   ```

3. Create a .env file:
   > store JWT_SECRET ,POLKA_API_KEY in .env
   > **NOTE**:You can create these keys yourself.
   ```env
   JWT_SECRET=<key>
   POLKA_API_KEY=<APIkey>
   ```
4. Build the project:

   ```bash
    go build
   ```

## API Endpoints

1. **Health Check** `GET /api/healthz`

   Returns a `200` status code and a JSON response:

   ```json
   {
     "status": "ok"
   }
   ```

2. **Create User** `POST /api/users`

   Example request body:

   ```json
   {
     "email": "johnDoe@gmail.com",
     "password": "password123"
   }
   ```

   Example response body:

   ```json
   {
     "id": 1,
     "email": "johnDoe@gmail.com",
     "is_chirpy_red": false
   }
   ```

3. **Login User** `POST /api/login`

   Example request body:

   ```json
   {
     "email": "johnDoe@gmail.com",
     "password": "password123"
   }
   ```

   Example response body:

   ```json
   {
     "id": 1,
     "email": "johnDoe@gmail.com",
     "is_chirpy_red": false,
     "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJjaGlycHktYWNjZXNzIiwic3ViIjoiMTAxIiwiZXhwIjoxNjkyMTcyMzkyLCJpYXQiOjE2OTIxNjkxOTJ9.1Wp4zSVxphWbA5KZwWJlgB6_samCxztx3-dnYRhdbnM",
     "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJjaGlycHktcmVmcmVzaCIsInN1YiI6IjEiLCJleHAiOjE2OTIyNTU1OTIsImlhdCI6MTY5MjE2OTE5Mn0._oQbOgG8Gzh2dkMfMlV9gzXPe0oV7QMK6XRUkz1Tyw0"
   }
   ```

4. **Update User** `PUT /api/users`
   Updates user email and password
   `Authentication Required`
   Request headers: `Authorization: Bearer <JWT_ACCESS_TOKEN>`
   Example request body:

   ```json
   {
     "email": "johnDoe@gmail.com",
     "password": "newpassword123" //new password to be updated
   }
   ```

   Example response body:
   `200 ok`

   ```json
   {
     "id": 1,
     "email": "johnDoe@gmail.com",
     "is_chirpy_red": false
   }
   ```

5. **Create Chirp** `POST /api/chirps`
   `Authentication Required`
   Request headers: `Authorization: Bearer <JWT_ACCESS_TOKEN>`
   Example request body:

   ```json
   {
     "body": "This is a example response chirp"
   }
   ```

   Example response body:
   `201 ok`

   ```json
   {
     "id": 1,
     "body": "This is a first chirp",
     "author_id": 1 //user_id
   }
   ```

6. **Get all Chirps of any user** `GET /api/chirps`
   Request Params: `author_id=<int> & sort="asc or desc"`

   Example response body:
   `200 ok`

   ```json
   [
     {
       "id": 1,
       "body": "This is a first chirp",
       "author_id": 1 //user_id
     }
   ]
   ```

7. **Get single chirp** `GET /api/chirps/{chirpID}`

   Example response body:
   `200 ok`

   ```json
   {
     "id": 1,
     "body": "This is a first chirp",
     "author_id": 1 //user_id
   }
   ```

8. **Delete chirp** `DELETE /api/chirps/{chirpID}`
   `Authentication Required`
   Request headers: `Authorization: Bearer <JWT_ACCESS_TOKEN>`

   > **NOTE**: user trying to delete should be author of that chirp

9. **Refresh Token** `POST /api/refresh`
   Generates new accesss token for user using previously provided refresh token
   `Authentication Required`
   Request headers:`Authorization: Bearer <JWT_REFRESH_TOKEN>`

   Example response body:
   `200 ok`

   ```json
   {
     "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJjaGlycHktYWNjZXNzIiwic3ViIjoiMSIsImV4cCI6MTY5MjE3NDE4NSwiaWF0IjoxNjkyMTcwNTg1fQ.ZFKTVYQO6EPb13UqmFi75yjZ7nfZ22wuWjHjhIDxaqo"
   }
   ```

10. **Revoke refresh Token** `POST /api/revoke`
    Revokes refresh token so that i cannot be used to generate access token again using same refresh token
    `Authentication Required`
    Request headers:`Authorization: Bearer <JWT_REFRESH_TOKEN>`

    Example Response:
    `200 ok `
    `response body: null`

11. **Upgrade User tier** `POST /polka/webhooks`
    Upgrades user to red(premium) tier. This request is made sample for intgration with stripe like api when payment is proccessed it sends a request with api key to upgrade user

    `Authentication Required`
    Request headers:`Authorization: ApiKey <APIKEY>`

    Example Reqeust body:

    ```json
    {
      "event": "user.upgraded",
      "data": {
        "user_id": 3
      }
    }
    ```

    Example Response:
    `200 ok` or `404 not found`

12. **Display Welcome** `GET /app`
    Displays a simple web page with welcome message

## Contributing

Contributions to the Chirpy Backend Go API are welcome! If you find any issues or have improvements to suggest, please feel free to open an issue or a pull request.
