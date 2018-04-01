package main

import (
  "fmt"
  "strings"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
  "github.com/aws/aws-sdk-go/service/dynamodb"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/ses"
  "github.com/aws/aws-lambda-go/lambda"
)

type Request struct {
	Commits []struct {
		Message string
    Url string
	}
  Repository struct {
    Name string
  }
}

type Response struct {
  Url string
}

type Subscribers struct {
  Email string
}


func GetSubs(repo string) *dynamodb.ScanOutput {
  svc := dynamodb.New(session.New())
  input := &dynamodb.ScanInput{
    ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
      ":r": {
        S: aws.String(repo),
      },
    },
    FilterExpression:     aws.String("contains (repos, :r)"),
    ProjectionExpression: aws.String("email"),
    TableName:            aws.String("webhook-subs"),
  }

  result, _ := svc.Scan(input)
  return result
}


func Email(url string, subs []Subscribers) {
  arr := []*string{}
	for _, r := range(subs) {
		arr = append(arr, aws.String(r.Email))
	}
  svc := ses.New(session.New())
  input := &ses.SendEmailInput{
    Destination: &ses.Destination{
      ToAddresses: arr,
    },
    Message: &ses.Message{
      Body: &ses.Body{
        Html: &ses.Content{
          Charset: aws.String("UTF-8"),
          Data:    aws.String(fmt.Sprintf("something changed here: %s", url)),
        },
        Text: &ses.Content{
          Charset: aws.String("UTF-8"),
          Data:    aws.String(fmt.Sprintf("something changed here: %s", url)),
        },
      },
      Subject: &ses.Content{
        Charset: aws.String("UTF-8"),
        Data:    aws.String("Test email"),
      },
    },
    Source: aws.String("650egor@gmail.com"),
  }
  svc.SendEmail(input)
}


func Handler(request Request) (Response, error) {
  if strings.Contains(request.Commits[0].Message, "jarvis") {
    subs := GetSubs(request.Repository.Name)
    items := []Subscribers{}
    dynamodbattribute.UnmarshalListOfMaps(subs.Items, &items)
    Email(request.Commits[0].Url, items)
  }
  return Response {
    Url: request.Commits[0].Url,
  }, nil
}


func main() {
  lambda.Start(Handler)
}
