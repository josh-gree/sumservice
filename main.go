package main

import (
	"fmt"
	"encoding/json"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"net/http"
	"bytes"
)

type Job struct{
	Data []float64 `json:"data"`
}

type Result struct{
	Out float64 `json:"out"`
	Service string `json:"service"`
}

func main(){
	Listen()
}

func Send(r Result) error {
	fmt.Println("Sending result: sum")
	data, err := json.Marshal(r)

	// How to deal with name issues?? service only needs to know public location
	publichost := "localhost:7000"
	_, err = http.Post(fmt.Sprintf("http://%s/result/",publichost),"application/json",bytes.NewBuffer(data))
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}


func Recivejob(c echo.Context) error {
	fmt.Println("Recieved Job: sum")
	j := Job{}
	err := c.Bind(&j)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Printf("Data : %f\n",j.Data)
	result := Sum(j)

	Send(result)

	return nil
}

func Sum(j Job) Result {
	fmt.Println("Doing computation: sum")
	sum := 0.0
	for _, v := range j.Data {
		sum += v
	}
	res := Result{Out:sum,Service:"sum"}
	return res
}

func Listen() {
	fmt.Println("Starting to Listen: sum")
	e := echo.New()
	logConfig := middleware.LoggerConfig{
		  Format: `[${time_rfc3339}] ${status} ${method} ${host}${path}` + "\n",
	}
	e.Use(middleware.LoggerWithConfig(logConfig))

	e.POST("/", Recivejob)
	e.Start(":8000")
}
