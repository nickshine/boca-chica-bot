package param

import (
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
)

// Client is an alias for ssm.SSM with some convenience receiver functions added.
type Client ssm.SSM

// GetClient returns a new Client.
func GetClient() *Client {
	sess := session.Must(session.NewSession())
	svc := ssm.New(sess)

	return (*Client)(svc)
}

// GetParams returns a map of key/value secrets from AWS SSM Parameter Store.
func (c *Client) GetParams(path string) (map[string]string, error) {
	input := &ssm.GetParametersByPathInput{
		Path:           aws.String(path),
		Recursive:      aws.Bool(true),
		WithDecryption: aws.Bool(true),
	}

	res, err := (*ssm.SSM)(c).GetParametersByPath(input)
	if err != nil {
		return nil, err
	}

	m := make(map[string]string)
	for _, p := range res.Parameters {
		name := strings.TrimPrefix(aws.StringValue(p.Name), path)
		m[name] = aws.StringValue(p.Value)
	}

	return m, nil
}
