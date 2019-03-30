package main

import (
    "encoding/json"
    "log"
    "fmt"
    "net/http"
    "strconv"
    "goji.io"
    "goji.io/pat"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
)

func ErrorWithJSON(w http.ResponseWriter, message string, code int) {
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    w.WriteHeader(code)
    fmt.Fprintf(w, "{message: %q}", message)
}

func ResponseWithJSON(w http.ResponseWriter, json []byte, code int) {
    w.Header().Set("Content-Type", "application/json; charset=utf-8")
    w.WriteHeader(code)
    w.Write(json)
}

type Building struct{
        Bid    int
        Yob    int
        Ctp    string
        Height float64
        Btp    int
}

func main() {
    session, err := mgo.Dial("localhost:27017")
    if err != nil {
        log.Fatal(err)
    }
    defer session.Close()

    session.SetMode(mgo.Monotonic, true)
    ensureIndex(session)

    mux := goji.NewMux()
    mux.HandleFunc(pat.Get("/buildings"), allBuildings(session))
    mux.HandleFunc(pat.Get("/buildings/:bid"), buildingByID(session))
    mux.HandleFunc(pat.Get("/buildings/height/:hgt"), buildingHeight(session))
    mux.HandleFunc(pat.Get("/buildings/groupbytype/"), buildingGroupByType(session))
    http.ListenAndServe("localhost:8080", mux)
}

func ensureIndex(s *mgo.Session) {
    session := s.Copy()
    defer session.Close()

    c := session.DB("Building_data").C("NYC_building")

    index := mgo.Index{
        Key:        []string{"bid"},
        Unique:     true,
        DropDups:   true,
        Background: true,
        Sparse:     true,
    }
    err := c.EnsureIndex(index)
    if err != nil {
        log.Fatal(err)
    }
}

func allBuildings(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        session := s.Copy()
        defer session.Close()

        c := session.DB("Building_data").C("NYC_building")

        var bd1 []Building
        err := c.Find(bson.M{}).All(&bd1)
        if err != nil {
            ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
            log.Println("Failed getting all buildings: ", err)
            return
        }

        respBody, err := json.MarshalIndent(bd1, "", "  ")
        if err != nil {
            log.Fatal(err)
        }

        ResponseWithJSON(w, respBody, http.StatusOK)
    }
}

func buildingByID(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        session := s.Copy()
        defer session.Close()

        bid := pat.Param(r, "bid")
        bid1,_ := strconv.Atoi(bid)

        c := session.DB("Building_data").C("NYC_building")

        var bd1 Building
        err := c.Find(bson.M{"bid": bid1}).One(&bd1)
        if err != nil {
            ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
            log.Println("Failed finding building: ", err)
            return
        }

        if bd1.Bid == 0 {
            ErrorWithJSON(w, "Building not found", http.StatusNotFound)
            return
        }

        respBody, err := json.MarshalIndent(bd1, "", "  ")
        if err != nil {
            log.Fatal(err)
        }

        ResponseWithJSON(w, respBody, http.StatusOK)
    }
}

func buildingHeight(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        session := s.Copy()
        defer session.Close()
        hgt := pat.Param(r, "hgt")
        hgt1,_ := strconv.ParseFloat(hgt,64)

        c := session.DB("Building_data").C("NYC_building")
        var bd1 []Building
        err := c.Find(bson.M{"height": bson.M{"$gte": hgt1}}).All(&bd1)
        if err != nil {
            ErrorWithJSON(w, "Database error", http.StatusInternalServerError)
            log.Println("Failed getting all buildings: ", err)
            return
        }

        respBody, err := json.MarshalIndent(bd1, "", "  ")
        if err != nil {
            log.Fatal(err)
        }

        ResponseWithJSON(w, respBody, http.StatusOK)
    }
}

func buildingGroupByType(s *mgo.Session) func(w http.ResponseWriter, r *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        session := s.Copy()
        defer session.Close()
        c := session.DB("Building_data").C("NYC_building")
        pipeline := []bson.M{
            {
                "$group": bson.M{
                    "_id": "$btp",
                    "sum": bson.M{"$sum": 1},
                },
            },
        }
        pipe := c.Pipe(pipeline)
        resp := []bson.M{}
        err := pipe.All(&resp)
        if err != nil {
            log.Println("Errored: ", err)
        }
        respBody, err := json.MarshalIndent(resp, "", "  ")
        if err != nil {
            log.Fatal(err)
        }
        ResponseWithJSON(w, respBody, http.StatusOK)
    }
}