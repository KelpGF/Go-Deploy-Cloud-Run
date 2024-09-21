package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/KelpGF/Go-Deploy-Cloud-Run/internal/services"
)

type SearchCEPResult struct {
	Cep         string `json:"cep"`
	Logradouro  string `json:"logradouro"`
	Complemento string `json:"complemento"`
	Bairro      string `json:"bairro"`
	Localidade  string `json:"localidade"`
	Uf          string `json:"uf"`
	Ibge        string `json:"ibge"`
	Gia         string `json:"gia"`
	Ddd         string `json:"ddd"`
	Siafi       string `json:"siafi"`
	Error       string `json:"erro"`
}
type ErrorMessage struct {
	Message string `json:"message"`
}

type Result struct {
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

func WeatherByCepHandler(w http.ResponseWriter, r *http.Request) {
	zipCode := r.URL.Query().Get("zipCode")
	w.Header().Set("Content-Type", "application/json")

	if zipCode == "" {
		errorMessage := ErrorMessage{Message: "zipCode is required"}

		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(errorMessage)

		return
	}

	zipCodeData, err := services.GetZipCodeData(zipCode)

	if err != nil {
		errorString := err.Error()
		statusCode := http.StatusInternalServerError

		if err.Error() == "invalid zipcode" {
			statusCode = http.StatusUnprocessableEntity
		}

		if err.Error() == "can not find zipcode" {
			statusCode = http.StatusNotFound
		}

		errorMessage := ErrorMessage{Message: errorString}

		w.WriteHeader(statusCode)
		json.NewEncoder(w).Encode(errorMessage)

		return
	}

	weatherData, err := services.GetWeatherData(zipCodeData.Localidade)

	if err != nil {
		errorMessage := ErrorMessage{Message: "Error on weather request"}

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(errorMessage)

		return
	}

	result := Result{
		TempC: weatherData.Current.TempC,
		TempF: weatherData.Current.TempF,
		TempK: weatherData.Current.TempC + 273.15,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}
