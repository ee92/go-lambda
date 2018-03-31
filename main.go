package main

import (
  "fmt"
  "strings"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/ses"
  "github.com/aws/aws-sdk-go/aws/awserr"
  "github.com/aws/aws-lambda-go/lambda"
)

var recievers = []Reciever{
	{
		email: "650egor@gmail.com",
		subbed: false,
	},
	{
		email: "eegor650@gmail.com",
		subbed: true,
	},
	{
		email: "broegorov@gmail.com",
		subbed: false,
	},
}

type Reciever struct {
	email string
	subbed bool
}

type Request struct {
	Commits []struct {
		Message string `json:"message"`
    Url string `json:"url"`
	} `json:"commits"`
}

type Response struct {
  Url string `json:"url"`
}

func Email(url string) {
  arr := []*string{}
	for _, r := range(recievers) {
		if r.subbed {
      s := aws.String(r.email)
			arr = append(arr, s)
		}
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

  result, err := svc.SendEmail(input)
  if err != nil {
    if aerr, ok := err.(awserr.Error); ok {
      switch aerr.Code() {
      case ses.ErrCodeMessageRejected:
          fmt.Println(ses.ErrCodeMessageRejected, aerr.Error())
      case ses.ErrCodeMailFromDomainNotVerifiedException:
          fmt.Println(ses.ErrCodeMailFromDomainNotVerifiedException, aerr.Error())
      case ses.ErrCodeConfigurationSetDoesNotExistException:
          fmt.Println(ses.ErrCodeConfigurationSetDoesNotExistException, aerr.Error())
      case ses.ErrCodeConfigurationSetSendingPausedException:
          fmt.Println(ses.ErrCodeConfigurationSetSendingPausedException, aerr.Error())
      case ses.ErrCodeAccountSendingPausedException:
          fmt.Println(ses.ErrCodeAccountSendingPausedException, aerr.Error())
      default:
          fmt.Println(aerr.Error())
      }
    } else {
      fmt.Println(err.Error())
    }
    return
  }
  fmt.Println(result)
}

func Handler(request Request) (Response, error) {
  if strings.Contains(request.Commits[0].Message, "jarvis") {
    Email(request.Commits[0].Url)
  }
  return Response {
    Url: request.Commits[0].Url,
  }, nil
}

func main() {
  lambda.Start(Handler)
}
