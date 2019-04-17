package main

import (
//    "flag"
	"fmt"
	"strings"
	"net/http"
	"io/ioutil"
    "os"
    "github.com/gin-gonic/gin"
    "encoding/json"
)

import "github.com/buger/jsonparser"
const ABIOS_URL = "https://api.abiosgaming.com/v2"
const CLIENT_ID = "test-task"

type listOfJson struct {
    List []interface{}
}

var atoken string

func main() {
    //client_id := flag.String("id", "", "client id for Abios API")
    //client_key := flag.String("key", "", "client key for Abios API")
    //flag.Parse()

	server := gin.Default()

    /*var err error
    atoken, err = getAccessToken(*client_key, *client_id)

    if err != nil {
        fmt.Println("something is wrong")
    }*/

    // set routes to specific methods
	server.GET("/series/live", series)
	server.GET("/players/live", players)
	server.GET("/teams/live", teams)

	server.Run() // listen and serve on 0.0.0.0:8080
}

func series(c *gin.Context) {
    // download
    //statusCode, jsonData := getBody(ABIOS_URL + "/series")
    _, liveSeries := getLiveSeries()

    c.IndentedJSON(200, liveSeries)
}

func players(c *gin.Context) {
    // download
    //statusCode, jsonData := getBody(ABIOS_URL + "/series")
    _, data := getLiveSeries()

    // iterate over teams and return
    players := getNestedKeyFromArray(data, "rosters", "players")

    c.IndentedJSON(200, players)
}

func teams(c *gin.Context) {
    // download
    _, data := getLiveSeries()

    // iterate over teams and return
    teams := getNestedKeyFromArray(data, "rosters", "teams")

    c.IndentedJSON(200, teams)
}

func getAccessToken(clientToken string, clientId string) (string, error) {

	url := ABIOS_URL +  "/oauth/access_token"

	payload := strings.NewReader("grant_type=client_credentials&client_id=" + clientId + "&client_secret=" + clientToken)

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

    value, err := jsonparser.GetString(body, "access_token")
    if err != nil {
        fmt.Println("No access token found", string(body), res.StatusCode)
        os.Exit(1)
    }

    fmt.Println(value)
    return value, err
}


func getBody(url string) (int, []byte) {
    req, _ := http.NewRequest("GET", url + "?access_token=" + atoken, strings.NewReader(""))

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

    return res.StatusCode, body
}

func getLiveSeries() (error, []map[string]interface{}) {
    //var returnValue []byte
    var jn map[string]interface{}
    var returnValue []map[string]interface{}
    var err error

    //_, jsonData := getBody(ABIOS_URL + "/series")

    jsonData, err := ioutil.ReadFile("src/main/series.txt")
    json.Unmarshal(jsonData, &jn)

    // iterate over data
    series := getKeyAndReturnArray(jn, "data")

    for _, serie := range series {
        if serie["end"] == nil {
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

/*func jsonPrettyPrint(in []byte) string {
    var out bytes.Buffer
    err := json.Indent(&out, in, "", "\t")
    if err != nil {
        return string(in)
    }
    return out.String()
}*/
