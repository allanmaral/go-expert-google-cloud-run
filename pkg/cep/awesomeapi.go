package cep

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type awesomeAPIResponse struct {
	Cep       string `json:"cep"`
	Street    string `json:"address"`
	District  string `json:"district"`
	City      string `json:"localidade"`
	State     string `json:"state"`
	Latitude  string `json:"lat"`
	Longitude string `json:"lng"`
}

type AwesomeAPILoader struct {
	client *http.Client
}

var _ Loader = &AwesomeAPILoader{}

func NewAwesomeAPILoader() *AwesomeAPILoader {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	return &AwesomeAPILoader{
		client: &http.Client{Transport: tr},
	}
}

func (l *AwesomeAPILoader) Load(ctx context.Context, cep string) (CEP, error) {
	if !Valid(cep) {
		return CEP{}, ErrInvalidCEP
	}

	url := fmt.Sprintf("https://cep.awesomeapi.com.br/json/%s", cep)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return CEP{}, err
	}

	req = req.WithContext(ctx)

	res, err := l.client.Do(req)
	if err != nil {
		return CEP{}, err
	}
	defer res.Body.Close()

	if res.StatusCode == 404 {
		return CEP{}, ErrCEPNotFound
	}

	if res.StatusCode != 200 {
		return CEP{}, ErrServiceUnavailable
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return CEP{}, err
	}

	var b awesomeAPIResponse
	err = json.Unmarshal(body, &b)
	if err != nil {
		return CEP{}, err
	}

	c := CEP{
		Cep:          b.Cep,
		Street:       b.Street,
		Neighborhood: b.District,
		City:         b.City,
		State:        b.State,
		Latitude:     b.Latitude,
		Longitude:    b.Longitude,
		Service:      "AwesomeAPI",
	}

	return c, nil
}
