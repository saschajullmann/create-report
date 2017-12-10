package main

import (
    "net/http"
    "encoding/json"
    "github.com/go-redis/redis"
    "os"
    "log"
    "strings"
    "io/ioutil"
)

type Document struct {
    Data map[string]interface{}
    Template string
}

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func getTemplateFiles() []string {
    file, err := os.Open("./templates")
    if err != nil {
        log.Fatalf("failed opening directory: %s", err)
    }
    defer file.Close()

    list,_ := file.Readdirnames(0) // 0 to read all files and folders
    return list
}

func loadTemplates(client *redis.Client) {
    tempList := getTemplateFiles()

    for _, name := range tempList {
        temp_string, err := ioutil.ReadFile("./templates/" + name)
        check(err)

        tempStringSep := strings.Split(string(name), ".")
        var template_string interface{}
        template_string = string(temp_string)

        client.Set(tempStringSep[0], template_string, 0)
    }
    return
}

func index_handler(w http.ResponseWriter, r *http.Request, client *redis.Client) {
    decoder := json.NewDecoder(r.Body)

    var d Document
    err:= decoder.Decode(&d)
    check(err)
    defer r.Body.Close()
    s, err := json.Marshal(d)
    check(err)

    myerr := client.LPush(os.Getenv("REDIS_QUEUE"), s).Err()
    check(myerr)

    return
}

func main() {
    client := redis.NewClient(&redis.Options{
        Addr:     os.Getenv("REDIS_HOST"),
        Password: os.Getenv("REDIS_PASSWORD"),
        DB:       0,
    })

    defer client.Close()
    loadTemplates(client)
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        index_handler(w, r, client)
    })
    http.ListenAndServe(":8000", nil)
}
