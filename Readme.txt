# Project setup:

step 1: Install Go from https://go.dev, to verify enter the following command in terminal :-  "go version"

step 2: Install Go extension in vscode 

step 3: Provide your mysql db connection string in the .env file in root folder or insert your user and password in the following line and paste it in .env file
"DSN=<username>:<password>@tcp(localhost:3306)/<schemaName>?charset=utf8mb4&parseTime=True&loc=Local"


# Run the application:

To run the application, open this folder in integrated terminal by double clicking on the main.go file in root directory and enter command, "go run main.go"

Use command "go build" to build the application

Note: Do install the dependency packages that vscode suggests you to download like "gopls", if it does not suggests run the following command in integrated terminal:

"go install -v golang.org/x/tools/gopls@latest"

# Run tests:

Go to _test.go file in respective folders and click on run test button appearing on top of test functions.



