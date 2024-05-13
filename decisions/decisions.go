package decisions

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"

	warranted "github.com/warranted-io/warranted-go/client"
)

type DecisionsService struct {
	requestHandler *warranted.RequestHandler
}

type DecisionResponse struct {
	Decision Decision `json:"decision"`
}

type Decision struct {
	Id                      string `json:"id"`
	CreatedAt               int64  `json:"createdAt"`
	UpdatedAt               int64  `json:"updatedAt"`
	LawEnforcementRequestId string `json:"lawEnforcementRequestId"`
	Owner                   string `json:"owner"`
	Result                  Result `json:"result"`
}

type Result struct {
	Subject   string                `json:"subject"`
	Requester string                `json:"requester"`
	StartTime string                `json:"startTime"`
	EndTime   string                `json:"endTime"`
	Schema    map[string]ResultData `json:"schema"`
}

type ResultData struct {
	IsRequested bool   `json:"isRequested"`
	IsError     bool   `json:"isError"`
	Reason      string `json:"reason"`
}

func NewDecisionsService(requestHandler *warranted.RequestHandler) *DecisionsService {
	return &DecisionsService{
		requestHandler: requestHandler,
	}
}

// Get Decisions
func (ds *DecisionsService) Get(id string) (Decision, error) {
	path := fmt.Sprintf("/api/v1/decisions/%s", id)
	resp, err := ds.requestHandler.Get(path, nil, nil)
	if err != nil {
		return Decision{}, err
	}
	defer resp.Body.Close()
	var buffer bytes.Buffer
	_, err = io.Copy(&buffer, resp.Body)
	if err != nil {
		return Decision{}, err
	}
	body := buffer.Bytes()

	var decisionResponse DecisionResponse
	err = json.Unmarshal(body, &decisionResponse)
	return decisionResponse.Decision, err
}
