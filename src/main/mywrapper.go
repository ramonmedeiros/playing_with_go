package main

import (
    "flag"
	"fmt"
	"strings"
	"net/http"
	"io/ioutil"
    "os"
    "github.com/gin-gonic/gin"
    "encoding/json"
    "time"
    "log"
)

const ABIOS_URL = "https://api.abiosgaming.com/v2"

var atoken string

func main() {
    client_id := flag.String("id", "", "client id for Abios API")
    client_key := flag.String("key", "", "client key for Abios API")
    flag.Parse()

	server := gin.Default()

    atoken = getAccessToken(*client_key, *client_id)

    log.Println("ACCESS_TOKEN:", atoken)

	server.GET("/series/live", series)
	server.GET("/players/live", players)
	server.GET("/teams/live", teams)

	server.Run()
}

func series(c *gin.Context) {
    _, liveSeries := getLiveSeries()

    c.IndentedJSON(200, liveSeries)
}

func players(c *gin.Context) {
    _, data := getLiveSeries()

    players := getNestedKeyFromArray(data, "rosters", "players")

    c.IndentedJSON(200, players)
}

func teams(c *gin.Context) {
    _, data := getLiveSeries()

    teams := getNestedKeyFromArray(data, "rosters", "teams")

    c.IndentedJSON(200, teams)
}

func getAccessToken(clientToken string, clientId string) (string) {

	url := ABIOS_URL +  "/oauth/access_token"

	payload := strings.NewReader("grant_type=client_credentials&client_id=" + clientId + "&client_secret=" + clientToken)

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

    var at map[string]string
    json.Unmarshal(body, &at)

    if res.StatusCode != 200 {
        fmt.Println("No token found", at)
        os.Exit(1)
    }

    return at["access_token"]
}


func getBody(url string) (int, []byte) {
    limiter := time.Tick(500 * time.Millisecond)
    <-limiter

    req, _ := http.NewRequest("GET", url + "?access_token=" + atoken, strings.NewReader(""))

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
    return res.StatusCode, body
}

func getLiveSeries() (error, []map[string]interface{}) {
    var jn map[string]interface{}
    var returnValue []map[string]interface{}
    var err error

    _, jsonData := getBody(ABIOS_URL + "/series")
    json.Unmarshal(jsonData, &jn)

    series := getKeyAndReturnArray(jn, "data")

    for _, serie := range series {
        if serie["end"] == nil {
            log.Println("Id:", serie["id"], " starts", serie["start"])
            returnValue = append(returnValue, serie)
        }
    }
    return err, returnValue
}

func getKeyAndReturnArray(array map[string]interface{}, key string) ([]map[string]interface{}) {
   var ret []map[string]interface{}

    mid := array[key]
    byt, _ := json.Marshal(mid)
    json.Unmarshal(byt, &ret)
    return ret
}


func getNestedKeyFromArray(array []map[string]interface{}, key1 string, key2 string) ([]map[string]interface{}) {

    var returnValue []map[string]interface{}

    for _, elem1 := range array {
        key1Elem := getKeyAndReturnArray(elem1, key1)

        for _, elem2 := range key1Elem {
            key2Elem := getKeyAndReturnArray(elem2, key2)
            returnValue = append(returnValue, key2Elem...)
        }
    }

    return returnValue
}

