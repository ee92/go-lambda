package main

import (
  "github.com/aws/aws-lambda-go/lambda"
  "fmt"
  "log"
)

type Request struct {
  Before string `json:"before"`
  After string `json:"after"`
}

type Response struct {
  Message string `json:"message"`
}

func Handler(request Request) (Response, error) {
  log.Printf("before: %s after: %s", request.Before, request.After)
  return Response {
    Message: fmt.Sprintf("before: %s after: %s", request.Before, request.After),
  }, nil
}

func main() {
  lambda.Start(Handler)
}
