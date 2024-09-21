package integration_test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/KelpGF/Go-Deploy-Cloud-Run/internal/handlers"
	"github.com/KelpGF/Go-Deploy-Cloud-Run/internal/server"
	"github.com/KelpGF/Go-Deploy-Cloud-Run/internal/services"
	"github.com/stretchr/testify/suite"
)

type IntegrationTestSuite struct {
	suite.Suite

	baseUrl string
	server  server.ServerHttp
}

func (s *IntegrationTestSuite) SetupSuite() {
	s.baseUrl = "http://localhost:8080/zip-code/weather?zipCode="
	s.server = server.ServerHttp{}
	go s.server.Run()
}

func (s *IntegrationTestSuite) TestSuccess() {
	url := s.baseUrl + "01311000"

	output, err := services.HttpRequest(url)

	s.Suite.NoError(err)
	s.Suite.Equal(http.StatusOK, output.Code)

	var data handlers.Result
	s.Suite.NoError(json.Unmarshal(output.Data, &data))

	s.Suite.NotEmpty(data.TempC)
	s.Suite.NotEmpty(data.TempF)
	s.Suite.NotEmpty(data.TempK)
}

func (s *IntegrationTestSuite) TestMissingZipCode() {
	url := s.baseUrl

	output, err := services.HttpRequest(url)

	s.Suite.NoError(err)
	s.Suite.Equal(http.StatusBadRequest, output.Code)

	var errorMessage handlers.ErrorMessage
	s.Suite.NoError(json.Unmarshal(output.Data, &errorMessage))

	s.Suite.Equal("zipCode is required", errorMessage.Message)
}

func (s *IntegrationTestSuite) TestInvalidZipCode() {
	url := s.baseUrl + "000000aa"

	output, err := services.HttpRequest(url)

	s.Suite.NoError(err)
	s.Suite.Equal(http.StatusUnprocessableEntity, output.Code)

	var errorMessage handlers.ErrorMessage
	s.Suite.NoError(json.Unmarshal(output.Data, &errorMessage))

	s.Suite.Equal("invalid zipcode", errorMessage.Message)
}

func (s *IntegrationTestSuite) TestNotFoundZipCode() {
	url := s.baseUrl + "00000000"

	output, err := services.HttpRequest(url)

	s.Suite.NoError(err)
	s.Suite.Equal(http.StatusNotFound, output.Code)

	var errorMessage handlers.ErrorMessage
	s.Suite.NoError(json.Unmarshal(output.Data, &errorMessage))

	s.Suite.Equal("can not find zipcode", errorMessage.Message)
}

func TestSuiteIntegration(t *testing.T) {
	suite.Run(t, new(IntegrationTestSuite))
}
