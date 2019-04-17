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
//    "strconv"
)

import "github.com/buger/jsonparser"
const ABIOS_URL = "https://api.abiosgaming.com/v2"
const CLIENT_ID = "test-task"

var atoken string

func main() {
    client_id := flag.String("id", "", "client id for Abios API")
    client_key := flag.String("key", "", "client key for Abios API")
    flag.Parse()

	server := gin.Default()

    var err error
    atoken, err = getAccessToken(*client_key, *client_id)

    if err != nil {
        fmt.Println("something is wrong")
    }

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

    var tj map[string]*interface{}
    json.Unmarshal(liveSeries, &tj)

    c.JSON(200, tj)
}

func players(c *gin.Context) {
    // download
    //statusCode, jsonData := getBody(ABIOS_URL + "/series")
    var allPlayers []byte
    _, data := getLiveSeries()

    // iterate over teams and return
    jsonparser.ArrayEach(data, func(player []byte, dataType jsonparser.ValueType, offset int, err error) {
        fmt.Println(string(player))
        allPlayers = append(allPlayers, player...)
    }, "rosters", "players")

    c.JSON(200, allPlayers)
}

func teams(c *gin.Context) {
    // download
    var allTeams []byte
    _, data := getLiveSeries()

    // iterate over teams and return
    jsonparser.ArrayEach(data, func(team []byte, dataType jsonparser.ValueType, offset int, err error) {
        fmt.Println(string(team))
        allTeams = append(allTeams, team...)
    }, "rosters", "teams")

    c.JSON(200, allTeams)
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

func getLiveSeries() (error, []byte) {
    var returnValue []byte
    var err error

    //_, jsonData := getBody(ABIOS_URL + "/series")

    jsonData, err := ioutil.ReadFile("src/main/series.txt")

    // iterate over data
    jsonparser.ArrayEach(jsonData, func(serie []byte, dataType jsonparser.ValueType, offset int, err error) {
        valu, _, _,_ := jsonparser.Get(serie, "end")

        // if end is null: still live
        if string(valu) == "null" {
            returnValue = append(returnValue, serie...)
        }
    }, "data")

    return err, returnValue
}

/*func jsonPrettyPrint(in []byte) string {
    var out bytes.Buffer
    err := json.Indent(&out, in, "", "\t")
    if err != nil {
        return string(in)
    }
    return out.String()
}*/
