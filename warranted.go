package warranted

import (
	"net/http"

	"github.com/warranted-io/warranted-go/client"
	"github.com/warranted-io/warranted-go/decisions"
	"github.com/warranted-io/warranted-go/lawEnforcementRequests"
	"github.com/warranted-io/warranted-go/me"
	"github.com/warranted-io/warranted-go/schema"
)

type Client struct {
	client                 *client.Client
	requestHandler         *client.RequestHandler
	Decisions              *decisions.DecisionsService
	Me                     *me.MeService
	Schema                 *schema.SchemaService
	LawEnforcementRequests *lawEnforcementRequests.LawEnforcementRequestsService
}

// NewClient creates a new warranted client instance
func NewClient(accountID, authToken string) (*Client, error) {
	c, err := client.NewClient(accountID, authToken)
	if err != nil {
		return nil, err
	}
	rh := client.NewRequestHandler(*c)
	ds := decisions.NewDecisionsService(rh)
	ms := me.NewMeService(rh)
	ss := schema.NewSchemaService(rh)
	ls := lawEnforcementRequests.NewLawEnforcementRequestsService(rh)
	return &Client{
		client:                 c,
		requestHandler:         rh,
		Decisions:              ds,
		Me:                     ms,
		Schema:                 ss,
		LawEnforcementRequests: ls,
	}, nil
}

func (c *Client) SetHeader(headers http.Header) {
	c.client.SetHeader(headers)
}

func (c *Client) SetHost(host string) {
	c.client.SetHost(host)
}
