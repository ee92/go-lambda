package main

import (
  "strings"
  "net/http"
  "io/ioutil"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
  "github.com/aws/aws-sdk-go/service/dynamodb"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/ses"
  "github.com/aws/aws-lambda-go/lambda"
)

// payload from github
type Request struct {
	Commits []struct {
		Message string
    Url string
	}
  Repository struct {
    Name string
  }
}
// attributes for each email recipient
type Subscriber struct {
  Email string
}

// given repo name, return all subscribers
func GetSubs(repo string) *dynamodb.ScanOutput {
  // new dynamoDB session
  svc := dynamodb.New(session.New())
  // scan for all emails who have meet conditions
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

// send email to all subscribers of a repo
func Email(request Request, subs []Subscriber) {
  // convert subs to list of aws strings
  arr := []*string{}
	for _, r := range(subs) {
		arr = append(arr, aws.String(r.Email))
	}
  // get diff
  res, _ := http.Get(request.Commits[0].Url + ".diff")
  diff, _ := ioutil.ReadAll(res.Body)
  // new SES session
  svc := ses.New(session.New())
  // email body
  body := "SOMTHING CHANGED: \n" + request.Commits[0].Url + "\n" +
      "MESSAGE: \n" + request.Commits[0].Message + "\n" +
      "DIFF: \n" + string(diff) + "\n"
  // fill out email inputs
  input := &ses.SendEmailInput{
    Destination: &ses.Destination{
      ToAddresses: arr,
    },
    Message: &ses.Message{
      Body: &ses.Body{
        Text: &ses.Content{
          Charset: aws.String("UTF-8"),
          Data:    aws.String(body),
        },
      },
      Subject: &ses.Content{
        Charset: aws.String("UTF-8"),
        Data:    aws.String(strings.Split(request.Commits[0].Message, "\n")[0]),
      },
    },
    Source: aws.String("650egor@gmail.com"),
  }
  // send email out
  svc.SendEmail(input)
}

// get subs and email them if condition met
func Handler(request Request) {
  if strings.Contains(request.Commits[0].Message, "jarvis") {
    subs := GetSubs(request.Repository.Name)
    items := []Subscriber{}
    dynamodbattribute.UnmarshalListOfMaps(subs.Items, &items)
    Email(request, items)
  }
}

// call aws lambda function
func main() {
  lambda.Start(Handler)
}
