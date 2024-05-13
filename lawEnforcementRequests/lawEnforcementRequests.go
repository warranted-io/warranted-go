package lawEnforcementRequests

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"os"
	"path"
	"strconv"

	warranted "github.com/warranted-io/warranted-go/client"
)

const maxLimit = 500

type LawEnforcementRequestsService struct {
	requestHandler *warranted.RequestHandler
}

type Options struct {
	StartAt int
	Limit   int
}

type LawEnforcementRequestsResponse struct {
	LawEnforcementRequests []LawEnforcementRequest `json:"lawEnforcementRequests"`
}

type LawEnforcementRequestResponse struct {
	LawEnforcementRequest LawEnforcementRequest `json:"lawEnforcementRequest"`
}

type LawEnforcementRequest struct {
	Id         string   `json:"id"`
	CreatedAt  int      `json:"createdAt"`
	UpdatedAt  int      `json:"updatedAt"`
	Name       string   `json:"name"`
	Type       string   `json:"type"`
	DecisionId string   `json:"decisionId"`
	Statuses   []Status `json:"statuses"`
}

type Status struct {
	Status string `json:"status"`
	Reason string `json:"reason"`
}

type DeleteResponse struct {
	Message string `json:"message"`
}

func NewLawEnforcementRequestsService(requestHandler *warranted.RequestHandler) *LawEnforcementRequestsService {
	return &LawEnforcementRequestsService{
		requestHandler: requestHandler,
	}
}

// Get LawEnforcementRequests
func (ls *LawEnforcementRequestsService) Get(options Options) ([]LawEnforcementRequest, error) {
	baseURL := "/api/v1/lawEnforcementRequests"

	limit := options.Limit
	if limit < 1 {
		// Default non-sensical limits to maxLimit
		limit = maxLimit
	}
	baseURL += "?startAt=" + strconv.Itoa(options.StartAt) + "&limit=" + strconv.Itoa(limit)
	resp, err := ls.requestHandler.Get(baseURL, nil, nil)
	if err != nil {
		return []LawEnforcementRequest{}, err
	}
	defer resp.Body.Close()
	var buffer bytes.Buffer
	_, err = io.Copy(&buffer, resp.Body)
	if err != nil {
		return []LawEnforcementRequest{}, err
	}
	body := buffer.Bytes()

	var lawEnforcementRequestsResponse LawEnforcementRequestsResponse
	err = json.Unmarshal(body, &lawEnforcementRequestsResponse)
	return lawEnforcementRequestsResponse.LawEnforcementRequests, err
}

func (ls *LawEnforcementRequestsService) GetById(id string) (LawEnforcementRequest, error) {
	baseURL := "/api/v1/lawEnforcementRequests"
	baseURL = path.Join(baseURL, id)

	resp, err := ls.requestHandler.Get(baseURL, nil, nil)
	if err != nil {
		return LawEnforcementRequest{}, err
	}
	defer resp.Body.Close()
	var buffer bytes.Buffer
	_, err = io.Copy(&buffer, resp.Body)
	if err != nil {
		return LawEnforcementRequest{}, err
	}
	body := buffer.Bytes()

	var lawEnforcementRequestResponse LawEnforcementRequestResponse
	err = json.Unmarshal(body, &lawEnforcementRequestResponse)
	return lawEnforcementRequestResponse.LawEnforcementRequest, err
}

// New LawEnforcementRequest
func (ls *LawEnforcementRequestsService) Add(file *os.File) (LawEnforcementRequest, error) {
	fileData, err := ioutil.ReadAll(file)
	if err != nil {
		return LawEnforcementRequest{}, err
	}
	fi, err := file.Stat()
	if err != nil {
		return LawEnforcementRequest{}, err
	}
	var body bytes.Buffer
	writer := multipart.NewWriter(&body)
	filePart, err := writer.CreateFormFile("lawEnforcementRequest", fi.Name())
	if err != nil {
		return LawEnforcementRequest{}, err
	}
	_, err = filePart.Write(fileData)
	if err != nil {
		return LawEnforcementRequest{}, err
	}

	headers := map[string]interface{}{
		"Content-Type": writer.Boundary(),
	}

	err = writer.Close()
	if err != nil {
		return LawEnforcementRequest{}, err
	}

	resp, err := ls.requestHandler.Post("/api/v1/lawEnforcementRequest/new", nil, headers, body.Bytes()...)
	if err != nil {
		return LawEnforcementRequest{}, err
	}
	defer resp.Body.Close()
	var buffer bytes.Buffer
	_, err = io.Copy(&buffer, resp.Body)
	if err != nil {
		return LawEnforcementRequest{}, err
	}
	responseBody := buffer.Bytes()

	var lawEnforcementRequestResponse LawEnforcementRequestResponse
	err = json.Unmarshal(responseBody, &lawEnforcementRequestResponse)
	return lawEnforcementRequestResponse.LawEnforcementRequest, err
}

// Update LawEnforcementRequest
func (ss *LawEnforcementRequestsService) Update(lawEnforcementRequest LawEnforcementRequest) (LawEnforcementRequest, error) {
	body, err := json.Marshal(lawEnforcementRequest)
	if err != nil {
		return LawEnforcementRequest{}, err
	}
	resp, err := ss.requestHandler.Put("/api/v1/lawEnforcementRequest", nil, nil, body...)
	if err != nil {
		return LawEnforcementRequest{}, err
	}
	defer resp.Body.Close()
	var buffer bytes.Buffer
	_, err = io.Copy(&buffer, resp.Body)
	if err != nil {
		return LawEnforcementRequest{}, err
	}
	body = buffer.Bytes()

	var lawEnforcementRequestResponse LawEnforcementRequestResponse
	err = json.Unmarshal(body, &lawEnforcementRequestResponse)
	return lawEnforcementRequestResponse.LawEnforcementRequest, err
}

// Delete LawEnforcementRequest
func (ls *LawEnforcementRequestsService) Delete(id string) (DeleteResponse, error) {
	baseURL := "/api/v1/lawEnforcementRequests"
	baseURL = path.Join(baseURL, id)

	resp, err := ls.requestHandler.Delete(baseURL, nil, nil)
	if err != nil {
		return DeleteResponse{}, err
	}
	defer resp.Body.Close()
	var buffer bytes.Buffer
	_, err = io.Copy(&buffer, resp.Body)
	if err != nil {
		return DeleteResponse{}, err
	}
	body := buffer.Bytes()

	var dr DeleteResponse
	err = json.Unmarshal(body, &dr)
	return dr, err
}
