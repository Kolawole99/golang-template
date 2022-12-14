
# Diary Application

This helps you handle entries. It lets you authenticate and input an entry into the application.

## Environment Variables

To run this project, you will need to add the following environment variables to your .env file

`cp .env.local > .env`

This gives you the default keys and values. You can then customize the values to your variables of choosing.

## Hot Reloading
We use Air package to get hotrealing while developing.

To install the package use `go install github.com/cosmtrek/air@latest`

This project should have a generated `.air.toml` file but if it does not, use this command to generate the file and make updates to it as required
`air init`

## Generating Documentation before Push
To generate documentation before pushing the changes either to create a new documentation and to update the changes to the documentation

`swag fmt` - This formats the Swagger documentation annotations

`swag init` - This generates the latest version of the documentation

## Authors

- [@Kolawole99](https://github.com/Kolawole99)

# This formats the codebase
go fmt **/*.go
# This runs and compiles all files in the directory
go run ./cmd/api/
go run .
