package main

import (
    "flag"
	"fmt"
	"strings"
	"net/http"
	"io/ioutil"
    "github.com/gin-gonic/gin"
    "encoding/json"
)

const ABIOS_URL = "https://api.abiosgaming.com/v2"
const CLIENT_ID = "test-task"


func main() {
    client_id := flag.String("id", "", "client id for Abios API")
    client_key := flag.String("key", "", "client key for Abios API")
    flag.Parse()

    var token map[string]interface{}
    var err error
    token, err = getAccessToken(*client_key, *client_id)
    fmt.Println(token["access_token"])
    fmt.Println(err)

	r := gin.Default()

	r.GET("/series/live", series)
	r.GET("/players/live", players)
	r.GET("/teams/live", teams)

	r.Run() // listen and serve on 0.0.0.0:8080
}

func series(c *gin.Context) {
    c.JSON(200, gin.H{
			"message": "pong",
		})
}

func players(c *gin.Context) {
    c.JSON(200, gin.H{
			"message": "pong",
		})
}

func teams(c *gin.Context) {
    c.JSON(200, gin.H{
			"message": "pong",
		})
}

func getAccessToken(clientToken string, clientId string) (map[string]interface{}, error) {

	url := ABIOS_URL +  "/oauth/access_token"

	payload := strings.NewReader("grant_type=client_credentials&client_id=" + clientId + "&client_secret=" + clientToken)

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("content-type", "application/x-www-form-urlencoded")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var token map[string]interface{}
	jsonErr := json.Unmarshal(body, &token)

	return token, jsonErr
}



