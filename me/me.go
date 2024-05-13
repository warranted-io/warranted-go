package me

import (
	"bytes"
	"encoding/json"
	"io"

	warranted "github.com/warranted-io/warranted-go/client"
)

type MeService struct {
	requestHandler *warranted.RequestHandler
}

type Me struct {
	Id   string `json:"id"`
	Role string `json:"role"`
}

func NewMeService(requestHandler *warranted.RequestHandler) *MeService {
	return &MeService{
		requestHandler: requestHandler,
	}
}

// Get Me
func (ms *MeService) Get() (Me, error) {
	resp, err := ms.requestHandler.Get("/api/v1/me", nil, nil)
	if err != nil {
		return Me{}, err
	}
	defer resp.Body.Close()
	var buffer bytes.Buffer
	_, err = io.Copy(&buffer, resp.Body)
	if err != nil {
		return Me{}, err
	}
	body := buffer.Bytes()

	var me Me
	err = json.Unmarshal(body, &me)
	return me, err
}
