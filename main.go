package main

import (
	"log"
    // "strings"
    "fmt"
    "time"
	"encoding/json"
	"os"
	"github.com/gocolly/colly"
    twilio "github.com/kevinburke/twilio-go"
)

// type Coins struct {
// 	Coins []Coin `json:"Coins"`
// }

type Coin struct {
	Name	string
	Price	string
	DayChange string
	timeAdded string
}

func main() {
    run()
}


func run(){
	client := twilio.NewClient("ACaf356b30e339a6e9d0dba1f51aa4d989", "284305ff72993b5ab89fb22732ed3f1d", nil)

    for true {
        scrape()

		// msg, _ := client.Messages.SendMessage("+18787898352", "+14028407963", name + " is up " + dayChange + " in the past 24 hours", nil)
		msg, _ := client.Messages.SendMessage("+18787898352", "+14028407963", "Daily Crypto prices have been updated", nil)
		fmt.Println(msg.Sid, msg.FriendlyPrice())
		// 86400
        time.Sleep(28800 * time.Second)
    }
}

func scrape() {

		c := colly.NewCollector()

        c.OnHTML("tr", func(e *colly.HTMLElement) {
    		name, price, dayChange := "", "", ""


            // class="sc-1v2ivon-0 gClTFY"
    		e.ForEach("td", func(i int, e2 *colly.HTMLElement) {

                if i == 2 {

                    e2.ForEach("p", func(j int, e3 *colly.HTMLElement) {

                        if j == 0 {
                            name = e3.Text
                            // fmt.Println(e3.Text)
                        }
                    })

                } else if i == 3 {
                    price = e2.Text
                } else if i == 4 {
                    dayChange = e2.Text
                }
            })


            coin := Coin{
    			Name:       name,
    			Price:      price,
    			DayChange:	dayChange,
				timeAdded:  time.Now().String(),
    		}

            fmt.Printf("%+v\n", coin)
			writeFile(coin)

        })

    	c.OnRequest(func(r *colly.Request) {
    		log.Println("visiting", r.URL.String())
            log.Println()
    	})

        c.Visit("https://coinmarketcap.com/")
}

func writeFile(coin Coin){
	JSONcoin, _ := json.MarshalIndent(coin, "", " ")

	f, err := os.OpenFile("coins.json", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0660)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(string(JSONcoin)); err != nil {
		panic(err)
	}
}
