package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"weather/structs"
)

func GetCoordinates(city string) (float64, float64, error) {
	url := fmt.Sprintf(
		"https://geocoding-api.open-meteo.com/v1/search?name=%s&count=1",
		city,
	)

	resp, err := http.Get(url)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()
	var response structs.GeocodingResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return 0, 0, err
	}
	if len(response.Results) == 0 {
		return 0, 0, errors.New("city not found")

	}
	return response.Results[0].Latitude,
		response.Results[0].Longitude, nil

}

func GetTemperature(lat, lon float64) (float64, error) {

	url := fmt.Sprintf(
		"https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&current=temperature_2m",
		lat,
		lon,
	)
	resp, err := http.Get(url)

	if err != nil {
		fmt.Printf("Failed to fatch the weather: %s", err.Error())
		os.Exit(1)
	}

	var weather structs.WeatherResponse
	err = json.NewDecoder(resp.Body).Decode(&weather)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return 0, err
	}
	return weather.Current.Temperature, nil
}

func worker(city string, wg *sync.WaitGroup) {
	defer wg.Done()
	lat, lot, err := GetCoordinates(city)
	if err != nil {
		return
	}
	temp, err := GetTemperature(lat, lot)
	if err != nil {
		return
	}
	fmt.Println(city, temp)
}

func main() {

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Write City name please:")
	if !scanner.Scan() {
		fmt.Println("Something went wrong")
		os.Exit(1)
	}
	cities := strings.Fields(scanner.Text())
	wg := &sync.WaitGroup{}

	for _, city := range cities {
		wg.Add(1)
		go worker(city, wg)
	}
	wg.Wait()
	os.Exit(0)

}
