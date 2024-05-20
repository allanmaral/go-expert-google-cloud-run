package weather

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type weatherAPIResponse struct {
	Current struct {
		TempC float64 `json:"temp_c"`
	} `json:"current"`
}

type WeatherAPILoader struct {
	apikey string
}

var _ Loader = &WeatherAPILoader{}

func NewWeatherAPILoader(apikey string) *WeatherAPILoader {
	return &WeatherAPILoader{
		apikey: apikey,
	}
}

func (l *WeatherAPILoader) Load(ctx context.Context, lat, lng string) (Weather, error) {
	url := fmt.Sprintf("https://api.weatherapi.com/v1/current.json?key=%s&q=%s,%s&aqi=no", l.apikey, lat, lng)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Weather{}, err
	}

	req = req.WithContext(ctx)

	client := http.DefaultClient
	res, err := client.Do(req)
	if err != nil {
		return Weather{}, err
	}
	defer res.Body.Close()

	if res.StatusCode == 403 {
		return Weather{}, ErrUnauthorized
	}

	if res.StatusCode >= 400 && res.StatusCode < 500 {
		return Weather{}, ErrInvalidLocation
	}

	if res.StatusCode != 200 {
		return Weather{}, ErrServiceUnavailable
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return Weather{}, err
	}

	var b weatherAPIResponse
	err = json.Unmarshal(body, &b)
	if err != nil {
		return Weather{}, err
	}

	c := Weather{
		TempC:   b.Current.TempC,
		TempF:   CelsiusToFahrenheit(b.Current.TempC),
		TempK:   CelsiusToKelvin(b.Current.TempC),
		Service: "WeatherAPI",
	}

	return c, nil
}
