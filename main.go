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
	"github.com/joho/godotenv"
)

type Coin struct {
	Name	string
	Price	string
	DayChange string
	timeAdded string
}

func main() {
    run()
}

func goDotEnvVariable(key string) string {

  // load .env file
  err := godotenv.Load(".env")

  if err != nil {
    log.Fatalf("Error loading .env file")
  }

  return os.Getenv(key)
}

func run(){
	// Update Frequency in a 24hr period
	updateFreq := 1

	client := twilio.NewClient(goDotEnvVariable("twilioSid"), goDotEnvVariable("twilioKey"), nil)

    for true {
        scrape()

		// msg, _ := client.Messages.SendMessage("+18787898352", "+14028407963", name + " is up " + dayChange + " in the past 24 hours", nil)
		msg, _ := client.Messages.SendMessage(goDotEnvVariable("twilioNumber"), goDotEnvVariable("userNumber"), "Daily Crypto prices have been updated", nil)
		// fmt.Println(msg.Sid, msg.FriendlyPrice())
		fmt.Println(msg.Body)

		// 86400
        time.Sleep(time.Duration(86400 / updateFreq) * time.Second)
    }
}

func scrape() {

		c := colly.NewCollector()

        c.OnHTML("tr", func(e *colly.HTMLElement) {
    		name, price, dayChange := "", "", ""

    		e.ForEach("td", func(i int, e2 *colly.HTMLElement) {

                if i == 2 {

                    e2.ForEach("p", func(j int, e3 *colly.HTMLElement) {

                        if j == 0 {
                            name = e3.Text
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
