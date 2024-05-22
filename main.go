package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type AddressBrazilApiDto struct {
	Cep          string `json:"cep"`
	State        string `json:"state"`
	City         string `json:"city"`
	Neighborhood string `json:"neighborhood"`
	Street       string `json:"street"`
	Service      string `json:"service"`
}

type AddressViaCepDto struct {
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
}

func main() {
	cpe := "70200001"

	fmt.Printf("Find postcode address %s\n", cpe)

	addressFromBrazilApi := make(chan AddressBrazilApiDto)
	addressFromViaCep := make(chan AddressViaCepDto)

	// Brazil Api
	go func() {
		dto, err := GetAddressByCepFromBrazilApiClient(cpe)
		if err != nil {
			fmt.Println("Error Get Address from BrazilApi", err)
		} else {
			addressFromBrazilApi <- dto
		}
	}()

	// ViaCep
	go func() {
		dto, err := GetAddressByCepFromViaCepClient(cpe)
		if err != nil {
			fmt.Println("Error Get Address from ViaCep", err)
		} else {
			addressFromViaCep <- dto
		}
	}()

	select {

	case address := <-addressFromBrazilApi: // BrazilApi
		fmt.Printf("Street from BrazilApi: %s\n", address.Street)

	case address := <-addressFromViaCep: // ViaCep
		fmt.Printf("Street from ViaCep: %s\n", address.Logradouro)

	case <-time.After(time.Second * 1):
		println("timeout")
	}
}

func GetAddressByCepFromBrazilApiClient(cpe string) (AddressBrazilApiDto, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	requestHttp, err := http.NewRequestWithContext(
		ctx,
		"GET",
		fmt.Sprintf("https://brasilapi.com.br/api/cep/v1/%s", cpe),
		nil,
	)
	if err != nil {
		return AddressBrazilApiDto{}, err
	}

	resultHttp, err := http.DefaultClient.Do(requestHttp)
	if err != nil {
		return AddressBrazilApiDto{}, err
	}
	defer resultHttp.Body.Close()

	readHttp, err := io.ReadAll(resultHttp.Body)
	if err != nil {
		return AddressBrazilApiDto{}, err
	}

	var data AddressBrazilApiDto
	err = json.Unmarshal(readHttp, &data)
	if err != nil {
		return AddressBrazilApiDto{}, err
	}
	return data, nil
}

func GetAddressByCepFromViaCepClient(cpe string) (AddressViaCepDto, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	requestHttp, err := http.NewRequestWithContext(
		ctx,
		"GET",
		fmt.Sprintf("http://viacep.com.br/ws/%s/json/", cpe),
		nil,
	)
	if err != nil {
		return AddressViaCepDto{}, err
	}

	resultHttp, err := http.DefaultClient.Do(requestHttp)
	if err != nil {
		return AddressViaCepDto{}, err
	}
	defer resultHttp.Body.Close()

	readHttp, err := io.ReadAll(resultHttp.Body)
	if err != nil {
		return AddressViaCepDto{}, err
	}

	var data AddressViaCepDto
	err = json.Unmarshal(readHttp, &data)
	if err != nil {
		return AddressViaCepDto{}, err
	}
	return data, nil
}
