package main

import (
  "time"
  "strings"
  // "net/http"
  // "io/ioutil"
  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
  "github.com/aws/aws-sdk-go/service/dynamodb"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/ses"
  "github.com/aws/aws-lambda-go/lambda"
)

// payload from github
type Request struct {
	Push struct {
		Changes []struct {
			Forced bool `json:"forced"`
			Old    struct {
				Type  string `json:"type"`
				Name  string `json:"name"`
				Links struct {
					Commits struct {
						Href string `json:"href"`
					} `json:"commits"`
					Self struct {
						Href string `json:"href"`
					} `json:"self"`
					HTML struct {
						Href string `json:"href"`
					} `json:"html"`
				} `json:"links"`
				Target struct {
					Hash  string `json:"hash"`
					Links struct {
						Self struct {
							Href string `json:"href"`
						} `json:"self"`
						HTML struct {
							Href string `json:"href"`
						} `json:"html"`
					} `json:"links"`
					Author struct {
						Raw  string `json:"raw"`
						Type string `json:"type"`
					} `json:"author"`
					Summary struct {
						Raw    string `json:"raw"`
						Markup string `json:"markup"`
						HTML   string `json:"html"`
						Type   string `json:"type"`
					} `json:"summary"`
					Parents []struct {
						Type  string `json:"type"`
						Hash  string `json:"hash"`
						Links struct {
							Self struct {
								Href string `json:"href"`
							} `json:"self"`
							HTML struct {
								Href string `json:"href"`
							} `json:"html"`
						} `json:"links"`
					} `json:"parents"`
					Date    time.Time `json:"date"`
					Message string    `json:"message"`
					Type    string    `json:"type"`
				} `json:"target"`
			} `json:"old"`
			Links struct {
				Commits struct {
					Href string `json:"href"`
				} `json:"commits"`
				HTML struct {
					Href string `json:"href"`
				} `json:"html"`
				Diff struct {
					Href string `json:"href"`
				} `json:"diff"`
			} `json:"links"`
			Truncated bool `json:"truncated"`
			Commits   []struct {
				Hash  string `json:"hash"`
				Links struct {
					Self struct {
						Href string `json:"href"`
					} `json:"self"`
					Comments struct {
						Href string `json:"href"`
					} `json:"comments"`
					Patch struct {
						Href string `json:"href"`
					} `json:"patch"`
					HTML struct {
						Href string `json:"href"`
					} `json:"html"`
					Diff struct {
						Href string `json:"href"`
					} `json:"diff"`
					Approve struct {
						Href string `json:"href"`
					} `json:"approve"`
					Statuses struct {
						Href string `json:"href"`
					} `json:"statuses"`
				} `json:"links"`
				Author struct {
					Raw  string `json:"raw"`
					Type string `json:"type"`
				} `json:"author"`
				Summary struct {
					Raw    string `json:"raw"`
					Markup string `json:"markup"`
					HTML   string `json:"html"`
					Type   string `json:"type"`
				} `json:"summary"`
				Parents []struct {
					Type  string `json:"type"`
					Hash  string `json:"hash"`
					Links struct {
						Self struct {
							Href string `json:"href"`
						} `json:"self"`
						HTML struct {
							Href string `json:"href"`
						} `json:"html"`
					} `json:"links"`
				} `json:"parents"`
				Date    time.Time `json:"date"`
				Message string    `json:"message"`
				Type    string    `json:"type"`
			} `json:"commits"`
			Created bool `json:"created"`
			Closed  bool `json:"closed"`
			New     struct {
				Type  string `json:"type"`
				Name  string `json:"name"`
				Links struct {
					Commits struct {
						Href string `json:"href"`
					} `json:"commits"`
					Self struct {
						Href string `json:"href"`
					} `json:"self"`
					HTML struct {
						Href string `json:"href"`
					} `json:"html"`
				} `json:"links"`
				Target struct {
					Hash  string `json:"hash"`
					Links struct {
						Self struct {
							Href string `json:"href"`
						} `json:"self"`
						HTML struct {
							Href string `json:"href"`
						} `json:"html"`
					} `json:"links"`
					Author struct {
						Raw  string `json:"raw"`
						Type string `json:"type"`
					} `json:"author"`
					Summary struct {
						Raw    string `json:"raw"`
						Markup string `json:"markup"`
						HTML   string `json:"html"`
						Type   string `json:"type"`
					} `json:"summary"`
					Parents []struct {
						Type  string `json:"type"`
						Hash  string `json:"hash"`
						Links struct {
							Self struct {
								Href string `json:"href"`
							} `json:"self"`
							HTML struct {
								Href string `json:"href"`
							} `json:"html"`
						} `json:"links"`
					} `json:"parents"`
					Date    time.Time `json:"date"`
					Message string    `json:"message"`
					Type    string    `json:"type"`
				} `json:"target"`
			} `json:"new"`
		} `json:"changes"`
	} `json:"push"`
	Repository struct {
		Scm     string `json:"scm"`
		Website string `json:"website"`
		Name    string `json:"name"`
		Links   struct {
			Self struct {
				Href string `json:"href"`
			} `json:"self"`
			HTML struct {
				Href string `json:"href"`
			} `json:"html"`
			Avatar struct {
				Href string `json:"href"`
			} `json:"avatar"`
		} `json:"links"`
		FullName string `json:"full_name"`
		Owner    struct {
			Username    string `json:"username"`
			Type        string `json:"type"`
			DisplayName string `json:"display_name"`
			UUID        string `json:"uuid"`
			Links       struct {
				Self struct {
					Href string `json:"href"`
				} `json:"self"`
				HTML struct {
					Href string `json:"href"`
				} `json:"html"`
				Avatar struct {
					Href string `json:"href"`
				} `json:"avatar"`
			} `json:"links"`
		} `json:"owner"`
		Type      string `json:"type"`
		IsPrivate bool   `json:"is_private"`
		UUID      string `json:"uuid"`
	} `json:"repository"`
	Actor struct {
		Username    string `json:"username"`
		Type        string `json:"type"`
		DisplayName string `json:"display_name"`
		UUID        string `json:"uuid"`
		Links       struct {
			Self struct {
				Href string `json:"href"`
			} `json:"self"`
			HTML struct {
				Href string `json:"href"`
			} `json:"html"`
			Avatar struct {
				Href string `json:"href"`
			} `json:"avatar"`
		} `json:"links"`
	} `json:"actor"`
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
  diff := " +++ -- jkfwuf84h4821938ybd992es"
  // res, _ := http.Get(request.Push.Changes[0].Target.Html.Href + "/raw")
  // diff, _ := ioutil.ReadAll(res.Body)
  // new SES session
  svc := ses.New(session.New())
  // email body
  body := "SOMTHING CHANGED: \n" + request.Push.Changes[0].New.Target.Links.HTML.Href + "\n" +
      "MESSAGE: \n" + request.Push.Changes[0].New.Target.Summary.Raw + "\n" +
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
        Data:    aws.String(strings.Split(request.Push.Changes[0].New.Target.Summary.Raw, "\n")[0]),
      },
    },
    Source: aws.String("650egor@gmail.com"),
  }
  // send email out
  svc.SendEmail(input)
}

// get subs and email them if condition met
func Handler(request Request) {
  if strings.Contains(request.Push.Changes[0].New.Target.Summary.Raw, "egor") {
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
