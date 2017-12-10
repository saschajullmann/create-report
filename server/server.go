package main

import (
    "fmt"
    "github.com/go-redis/redis"
    "strings"
    "io/ioutil"
    "regexp"
    "path/filepath"
    "log"
    "os/exec"
    "time"
    "os"
    "encoding/json"
)

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func createDateString() string {
    t := time.Now()
    time_string := fmt.Sprintf("%d-%02d-%02d_%02d:%02d:%02d:%02d",
    t.Year(), t.Month(), t.Day(),
    t.Hour(), t.Minute(), t.Second(), t.Nanosecond())
    return time_string
}

func processMessage(message string, client *redis.Client) {
    var dat map[string]interface{}
    if err := json.Unmarshal([]byte(message), &dat); err != nil {
        panic(err)
    }

    // Get the template name
    template_name := dat["Template"].(string)
    // Store the contents of the template in a string to replace variables
    template_string, err := client.Get(template_name).Result()
    if err != nil {
        fmt.Println("Template not found. Cannot process message")
        return
    }

    // Loop over the keys of the json in order to replace the variables
    var data map[string]interface{}
    data = dat["Data"].(map[string]interface{})
    for k, v := range data {
        lower := strings.ToLower(k)
        regex_string := `\\newcommand{\\` + lower + `}{}`
        r, err := regexp.Compile(regex_string)
        check(err)
        v_string := v.(string)
        new_string := "\\newcommand{\\" + lower + "}" + "{" + v_string + "}"
        template_string = r.ReplaceAllString(template_string, new_string)
    }

    //write the new string back to file in tmp dir
    dir, err := ioutil.TempDir("", "doc")
    if err != nil {
        log.Fatal(err)
    }

    defer os.RemoveAll(dir) //clean the temp dir when we are done

    name := createDateString() + "_" + template_name
    tmpfn := filepath.Join(dir, name + ".tex")
    new_byte := []byte(template_string)
    if err := ioutil.WriteFile(tmpfn, new_byte, 0666); err != nil {
        log.Fatal(err)
    }

    // create pdf
    exec.Command("pdflatex", tmpfn).Output()

    // remove aux und log files
    exec.Command("rm", name + ".aux").Run()
    exec.Command("rm", name + ".log").Run()
}

func main() {
    client := redis.NewClient(&redis.Options{
        Addr:     os.Getenv("REDIS_HOST"),
        Password: os.Getenv("REDIS_PASSWORD"),
        DB:       0,
    })

    defer client.Close()

    for {
        msg, err := client.BLPop(0, os.Getenv("REDIS_QUEUE")).Result()
        check(err)
        for index, message := range msg {
            if index == 1 {
                go processMessage(message, client)
            }
        }
    }
}
