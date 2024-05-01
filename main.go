package main

import (
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
	"strconv"
	"strings"
)

// DO NOT MODIFY
const version = "1.0.0"

type tempService struct {
	logger *slog.Logger
	client *http.Client
}

type Climate struct {
	Temperature float64 `json:"temperature"`
	Humidity    float64 `json:"humidity"`
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{}))
	logger.Info("starting temp-api", "version", version)

	service := tempService{
		logger: logger,
		client: &http.Client{},
	}

	service.run()
}

func (s *tempService) getMetrics(url string) ([]byte, error) {
	result, err := s.client.Get(url)
	if err != nil {
		s.logger.Error("error getting url", err)
		return nil, err
	}
	defer result.Body.Close()

	body, err := io.ReadAll(result.Body)
	if err != nil {
		s.logger.Error("error reading body", err)
		return nil, err
	}

	return body, nil
}

func (s *tempService) parseMetrics(metrics []byte) (Climate, error) {
	metricsString := string(metrics)
	metricsSlice := strings.Split(metricsString, "\n")

	tempSlice := strings.Split(metricsSlice[0], " ")
	humiSlice := strings.Split(metricsSlice[1], " ")

	temp, err := strconv.ParseFloat(tempSlice[1], 64)
	if err != nil {
		s.logger.Error("error parsing temperature float", err)
		return Climate{}, err
	}

	humi, err := strconv.ParseFloat(humiSlice[1], 64)
	if err != nil {
		s.logger.Error("error parsing humidity float", err)
		return Climate{}, err
	}

	c := Climate{
		Temperature: temp,
		Humidity:    humi,
	}

	return c, nil
}

func (s *tempService) handler(w http.ResponseWriter, r *http.Request) {
	endpoint, ok := os.LookupEnv("METRICS_ENDPOINT")
	if !ok {
		s.logger.Error("required environment variable missing", "variable", "METRICS_ENDPOINT")
		http.Error(w, "Server is missing configuration", http.StatusInternalServerError)
		return
	}

	rawMetrics, err := s.getMetrics(endpoint + "/metrics")
	if err != nil {
		s.logger.Error("error getting metrics", err)
		http.Error(w, "Failed to retrieve metrics", http.StatusInternalServerError)
		return
	}

	climate, err := s.parseMetrics(rawMetrics)
	if err != nil {
		s.logger.Error("error parsing metrics")
		http.Error(w, "Failed to parse metrics", http.StatusInternalServerError)
		return
	}

	s.logger.Info("current climate",
		"temperature", climate.Temperature,
		"humidity", climate.Humidity,
	)

	climateJSON, err := json.Marshal(climate)
	if err != nil {
		s.logger.Error("error marshalling climate json", err)
		http.Error(w, "Error marshaling JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, writeErr := w.Write(climateJSON)
	if writeErr != nil {
		s.logger.Error("error writing response", writeErr)
		return
	}
}

func (s *tempService) run() {
	http.HandleFunc("/", s.handler)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		s.logger.Error("error starting http server", err)
	}
}
