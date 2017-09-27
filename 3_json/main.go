package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

const citiBikeURL = "https://gbfs.citibikenyc.com/gbfs/en/station_status.json"

type stationData struct {
	LastUpdated int `json:"last_updated"`
	TTL         int `json:"ttl"`
	Data        struct {
		Stations []station `json:"stations"`
	} `json:"data"`
}

type station struct {
	ID                string `json:"station_id"`
	NumBikesAvailable int    `json:"num_bikes_available"`
	NumBikesDisabled  int    `json:"num_bike_disabled"`
	NumDocksAvailable int    `json:"num_docks_available"`
	NumDocksDisabled  int    `json:"num_docks_disabled"`
	IsInstalled       int    `json:"is_installed"`
	IsRenting         int    `json:"is_renting"`
	IsReturning       int    `json:"is_returning"`
	LastReported      int    `json:"last_reported"`
	HasAvailableKeys  bool   `json:"eightd_has_available_keys"`
}

func main() {
	response, err := http.Get(citiBikeURL)
	if err != nil {
		log.Printf("Error fetching data %s\n", err.Error())
		return
	}
	defer response.Body.Close()

	// Read response into byte array
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Printf("Error reading data %s\n", err.Error())
		return
	}

	// Unmarshal station data
	var sd stationData
	if err := json.Unmarshal(body, &sd); err != nil {
		log.Printf("Error unmarshaling station data %s\n", err.Error())
		return
	}

	// Print first station
	fmt.Printf("%+v\n\n", sd.Data.Stations[0])

	// Save data to file
	output, err := json.Marshal(sd)
	if err != nil {
		log.Printf("Error marhaling data %s\n", err.Error())
		return
	}
	if err := ioutil.WriteFile("citibike.json", output, 0644); err != nil {
		fmt.Printf("Error writing json to file %s\n", err.Error())
		return
	}

	fmt.Println("Done")
}
