package main

import (
        "encoding/csv"
        "strings"
        "io/ioutil"
        "log"
        "strconv"
        "gopkg.in/mgo.v2"
)

type Building struct{
        Bid    int
        Yob    int
        Ctp    string
        Height float64
        Btp    int
}

func main(){
        session, err := mgo.Dial("localhost:27017")
        if err != nil {
                log.Fatal(err)
        }
        session.SetMode(mgo.Monotonic, true)
        c := session.DB("Building_data").C("NYC_building")
        defer session.Close()
        fileName := "building.csv"
        cntb,err := ioutil.ReadFile(fileName)
        if err != nil {
                log.Fatal(err)
        }
        r2 := csv.NewReader(strings.NewReader(string(cntb)))
        for i:=0;i<10000;i++{
                ss,_ := r2.Read()
                var bd1 Building
                bd1.Bid,_ = strconv.Atoi(ss[1])
                bd1.Yob,_ = strconv.Atoi(ss[2])
                bd1.Ctp = ss[5]
                bd1.Height,_ = strconv.ParseFloat(ss[7],64)
                bd1.Btp,_ = strconv.Atoi(ss[8])
                err = c.Insert(bd1)
        }
}