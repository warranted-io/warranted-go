package schema

import (
	"bytes"
	"encoding/json"
	"io"

	warranted "github.com/warranted-io/warranted-go/client"
)

type SchemaService struct {
	requestHandler *warranted.RequestHandler
}

type SchemaResponse struct {
	Schema Schema `'json:"schema"`
}

type Schema struct {
	Id        string  `json:"id"`
	CreatedAt int     `json:"createdAt"`
	UpdatedAt int     `json:"updatedAt"`
	Fields    []Field `json:"fields"`
}

type Field struct {
	Key         string `json:"key"`
	Description string `json:"description"`
}

func NewSchemaService(requestHandler *warranted.RequestHandler) *SchemaService {
	return &SchemaService{
		requestHandler: requestHandler,
	}
}

// Get Schema
func (ss *SchemaService) Get() (Schema, error) {
	resp, err := ss.requestHandler.Get("/api/v1/schema", nil, nil)
	if err != nil {
		return Schema{}, err
	}
	defer resp.Body.Close()
	var buffer bytes.Buffer
	_, err = io.Copy(&buffer, resp.Body)
	if err != nil {
		return Schema{}, err
	}
	body := buffer.Bytes()

	var schemaResponse SchemaResponse
	err = json.Unmarshal(body, &schemaResponse)
	return schemaResponse.Schema, err
}

// Update Schema
func (ss *SchemaService) Update(schema Schema) (Schema, error) {
	body, err := json.Marshal(schema)
	if err != nil {
		return Schema{}, err
	}
	resp, err := ss.requestHandler.Put("/api/v1/schema", nil, nil, body...)
	if err != nil {
		return Schema{}, err
	}
	defer resp.Body.Close()
	var buffer bytes.Buffer
	_, err = io.Copy(&buffer, resp.Body)
	if err != nil {
		return Schema{}, err
	}
	body = buffer.Bytes()

	var schemaResponse SchemaResponse
	err = json.Unmarshal(body, &schemaResponse)
	return schemaResponse.Schema, err
}
