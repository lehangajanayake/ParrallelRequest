package main

import(
	"log"
	"encoding/json"
	"io/ioutil"
)

func openFile(path string)[]byte{
	file, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal("Error opening the file, ", err.Error())
	}
	return file
}

func decodeFile(file []byte)([]jsonBody){
	var result []jsonBody
	err := json.Unmarshal(file, &result)
	if err != nil {
		log.Fatal("Error decoding the file, ", err.Error())
	}
	return result
}