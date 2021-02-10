package main

import (
	"flag"
	"runtime"
	"log"
	"io/ioutil"
	"os"
	"fmt"
	"time"
)

var cmdArgs flags
func init(){
	flag.StringVar(&cmdArgs.FilePath, "File", "", "Path to the json file(Number of requested will be defaulted to the number of objects in the file unless its 1)")
	flag.StringVar(&cmdArgs.URL, "URL", "http://example.com/", "URL")
	flag.StringVar(&cmdArgs.Method, "Method", "GET", "Method of the request")
	flag.IntVar(&cmdArgs.Limit, "Limit", runtime.NumCPU(), "Limit of CPU threads (default is the max)")
	flag.IntVar(&cmdArgs.Requests, "Requests", 1, "Number of requests")
	flag.BoolVar(&cmdArgs.Log, "Log", false, "True will log into log.txt")
	//If the file doesn't exist, create it or append to the file
	

}

func main(){
	flag.Parse()
	runtime.GOMAXPROCS(cmdArgs.Limit)
	var data []jsonBody
	if cmdArgs.Log == true{
		logger, err := os.OpenFile("logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatal(err)
		}

		log.SetOutput(logger)
	}
	if cmdArgs.FilePath != ""{
		file := openFile(cmdArgs.FilePath)
		data = decodeFile(file)
	}
	num := cmdArgs.Requests
	if (num > len(data) && len(data) != 1 && cmdArgs.FilePath != "") || (num == 1 && cmdArgs.FilePath != ""){
		num = len(data)
	}
	responses := parallelRequest(cmdArgs.Method, cmdArgs.URL, client(), num, data)

	var successfull int
	var responseTime time.Duration


	for _, i := range responses{
		if i.Error != nil {
			log.Println(i.Error.Error())
			continue
		}
		body, err := ioutil.ReadAll(i.Response.Body)
		if err != nil {
			log.Println("Error reading the response: ", err.Error())
		}
		successfull ++
		responseTime += i.Time
		log.Println("Response: ", string(body), "\nStatus code: ", i.Response.Status, "\nResponse Time: ", i.Time)
	}
	fmt.Println("Number of successfull requests: ", successfull)
	fmt.Println("Avarage response time: ", responseTime/time.Duration(successfull))
	fmt.Println("Logging	 to log.txt file")
}		