# warranted-go
A helper library for using the Warranted.io API.

## Installation
The recommended way to install `warranted-go` is by using [Go modules](https://go.dev/ref/mod#go-get).

You can run the command below from your terminal in the project directory to install the library:

```
go get github.com/warranted-io/warranted-go
```

### Test your installation
To make sure the installation was successful, try hitting the `/api/v1/me` API, like this:
```go
package main

import (
  "fmt"
  "os"

  "github.com/warranted-io/warranted-go"
)

func main() {
  // Get your Account Id and Auth Token from https://app.warranted.io/settings/webhook
  accountId := os.Getenv("WARRANTED_ACCOUNT_ID")
  authToken := os.Getenv("WARRANTED_AUTH_TOKEN")

  // Initialize the client
  warrantedClient := warranted.NewClient(accountId, authToken)

  // Fetch and print the response object
  response, err := warrantedClient.Me.Get()
  if err != nil {
    fmt.Println(err)
  } else {
    fmt.Println(response)
  }
}
```

## Usage
Check out [our docs](https://app.warranted.io/docs) for more details.