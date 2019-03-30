package main

import (
    "fmt"
    "net/http"
    "io/ioutil"
)

func getHeight(hgt float64 ){
	resp, err := http.Get("http://localhost:8080/buildings/height/"+fmt.Sprintf("%f", hgt))
	if err != nil {
        panic(err)
    
    }
    defer resp.Body.Close()
    s,err:=ioutil.ReadAll(resp.Body)
    filename := "height_test.txt"
    ioutil.WriteFile(filename,s,0644)
}

func getAggregate(){
	resp, err := http.Get("http://localhost:8080/buildings/groupbytype/")
	if err != nil {
        panic(err)
    
    }
    defer resp.Body.Close()
    s,err:=ioutil.ReadAll(resp.Body)
    filename := "agg.txt"
    ioutil.WriteFile(filename,s,0644)
}

func getAll(){
	resp, err := http.Get("http://localhost:8080/buildings")
	if err != nil {
        panic(err)
    
    }
    defer resp.Body.Close()
    s,err:=ioutil.ReadAll(resp.Body)
    filename := "all_building.txt"
    ioutil.WriteFile(filename,s,0644)

}

func main(){
	getHeight(24.5)
	getAggregate()
	getAll()
}