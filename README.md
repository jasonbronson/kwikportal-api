# Kwik Portal API

Kwik Portal API is a web application built with GO & Svelte that allows users to organize and manage their bookmarks more efficiently. It provides a user-friendly interface for storing and categorizing bookmarks, making it easier to find and access saved URLs.

## Features

- Create an account to start managing your bookmarks.
- Add bookmarks by providing the URL and assigning them to specific categories.
- Organize bookmarks into custom categories for easy navigation.
- Search for bookmarks based on keywords or tags.
- Edit and delete bookmarks as needed.
- User-friendly interface with a clean and intuitive design.

## Technologies Used

- Svelte: A lightweight JavaScript framework for building user interfaces.
- Go: Backend development using the Go programming language.
- SQLite: Database for storing user and bookmark data.
- GORM: Go library for interacting with the database.
- Gin: Web framework for building RESTful APIs in Go.
- JWT: Authentication and authorization using JSON Web Tokens.
- Bcrypt: Password hashing for secure storage.

## Installation

1. Clone the repository:

   ```shell
   git clone https://github.com/jasonbronson/kwikportal-api.git
   ```

2. Install dependencies

```cd kwikportal-api
go mod download
```

3. Setup env file

```cp .env.example .env
   #edit the env file adding JWT details etc..
```

4. Run project

```go run main.go

```

Contributing
Contributions are welcome! If you have any suggestions, bug reports, or feature requests, please open an issue or submit a pull request.

License
This project is licensed under the UnLicense.

Feel free to customize and expand upon this template according to your project's specific requirements and additional sections you want to include.
