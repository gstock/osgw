package client

import (
"encoding/json"
"log"
"errors"
"net/http"
"github.com/ddliu/go-httpclient"
"net/url"
)


type WWOResponse struct {

    Data struct {
        Weathers []Weather `json:"weather,omitempty"`
    } `json:"data,omitempty"`
}

type Weather struct {
    Date        string   `json:"date,omitempty"`
    MaxTempC    string   `json:"maxtempC,omitempty"`
    MinTempC    string   `json:"mintempC,omitempty"`
}

type WeatherClient interface { 
    GetWeather(string, string) (Weather, error)
} 

type WWOClient struct{} 

func (c WWOClient) GetWeather(location string, date string) (Weather, error) {
    
    var response WWOResponse


    if date == "" {        
        return Weather{}, errors.New("Date cannot be empty")
    }

    // TODO: don't concatenate the url
    res, err := httpclient.Get("http://api.worldweatheronline.com/premium/v1/past-weather.ashx?key=68b656ec93704a38847141640180806&q=" + url.QueryEscape(location) + "&format=json&date=" + url.QueryEscape(date))

    if err != nil || res.StatusCode != http.StatusOK {
        log.Println("[WWOClient > HTTP ERROR] ", err)        
        return Weather{}, errors.New("Error obtaining weather from WWO")
    }
    
    jsonErr := json.NewDecoder(res.Body).Decode(&response)

    if jsonErr != nil {
        log.Println("[WWOClient > JSON ERROR] ", jsonErr)
        return Weather{}, errors.New("Error parsing weather from WWO")
    }

    if len(response.Data.Weathers) > 0 {
        return response.Data.Weathers[0], nil
    } else {
        return Weather{}, errors.New("Could not find weather data")
    }

    
}

