package main

import (
	"fmt"
	"io/ioutil"
    "encoding/json"
    "os"
)

type series struct {
      Id int
      Title string
      Start string
      End string
      Tier int
      BestOf int
      Tournament_id int
      Substage_id int
      Deleted_at int
      Pbp_status string
      Postponed_from string
      Scores {}
      chain string
      forfeit {}
      streamed bool
      seeding {}
      game {}
      bracket_pos": {...},
      rosters": [Roster],
      tournament": {...}, //optional
      matches": [Match], //optional
      casters": [Caster], //optional
      sportsbook_odds": {...} //optional
    

}

func main() {
    dat, err := ioutil.ReadFile("src/main/series.txt")
    fmt.Println(err)

	var data map[string]interface{}
	jsonErr := json.Unmarshal(dat, &data)

    if jsonErr != nil {
        fmt.Println("nao funciona o json")
        os.Exit(1)
    }

    fmt.Println(data["data"][0])

}

