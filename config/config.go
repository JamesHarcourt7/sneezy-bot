package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

var (
	Token     string
	BotPrefix string

	config *configStruct
)

type configStruct struct {
	Token     string `json:"Token"`
	BotPrefix string `json:"BotPrefix"`
}

func ReadConfig() error {
	fmt.Println("Reading config file...")
	file, err := ioutil.ReadFile("./config.json")

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	err = json.Unmarshal(file, &config)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	Token = config.Token
	BotPrefix = config.BotPrefix

	return nil
}

func ReadData() (map[string]int, error) {
	fmt.Println("Reading data file...")
	file, err := ioutil.ReadFile("./data.json")

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	var result map[string]int
	err = json.Unmarshal(file, &result)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}
	return result, nil
}

func WriteData(data map[string]int) error {
	fmt.Println("Writing to data file...")

	str, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile("./data.json", str, 0755)
	if err != nil {
		return err
	}
	fmt.Println("Updated data file.")

	return nil
}
