package main

import (
    "flag"
	"fmt"
	"strings"
	"net/http"
	"io/ioutil"
    "os"
    "encoding/json"
    "log"
    "golang.org/x/time/rate"
)

const ABIOS_URL = "https://api.abiosgaming.com/v2"

var atoken string
var limiter = rate.NewLimiter(2, 2)

func main() {
    // parse client_id and key
    client_id := flag.String("id", "", "client id for Abios API")
    client_key := flag.String("key", "", "client key for Abios API")
    flag.Parse()

    // start and set rate limit
    mux := http.NewServeMux()
    mux.HandleFunc("/series/live", series)
	mux.HandleFunc("/players/live", players)
	mux.HandleFunc("/teams/live", teams)


    // store token and let it global
    atoken = getAccessToken(*client_key, *client_id)
    log.Println("ACCESS_TOKEN:", atoken)

    // routers
    http.ListenAndServe(":8080", limit(mux))
}

func series(w http.ResponseWriter, r *http.Request) {
    _, liveSeries := getLiveSeries()

    writeJson(w, liveSeries)
}

func players(w http.ResponseWriter, r *http.Request)  {
    _, data := getLiveSeries()

    players := getNestedKeyFromArray(data, "rosters", "players")

    writeJson(w, players)
}

func teams(w http.ResponseWriter, r *http.Request)  {
    _, data := getLiveSeries()

    teams := getNestedKeyFromArray(data, "rosters", "teams")

    writeJson(w, teams)
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
        if serie["start"] != nil && serie["end"] == nil {
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

func writeJson(w http.ResponseWriter, js []map[string]interface{}) {
    indented, _ := json.MarshalIndent(js, "", " ")

    w.Header().Set("Content-Type", "application/json")
    w.Write(indented)
}


func limit(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if limiter.Allow() == false {
            http.Error(w, http.StatusText(429), http.StatusTooManyRequests)
            return
        }

        next.ServeHTTP(w, r)
    })
}
