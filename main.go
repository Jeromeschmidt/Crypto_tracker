package main

import (
	"log"
    // "strings"
    "fmt"
    "time"
	"encoding/json"
	"os"
	"github.com/gocolly/colly"
    // twilio "github.com/kevinburke/twilio-go"
)

// var thCsrf string
//
// type Assignments struct {
//     Users []Assignment `json:"assignments"`
// }
//
// type Assignment struct {
// 	Title	string
// 	Link	string
// 	Class string
//     DueDate string
//     Reminded bool
// }

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
        time.Sleep(10 * time.Second)
    }
}

func scrape() {
		type Coins struct {
			Coins []Coin `json:"Coins"`
		}

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

			fmt.Println(time.Now().String()[1])


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
	// fmt.Println(string(JSONcoin))

	f, err := os.OpenFile("coins.json", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0660)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(string(JSONcoin)); err != nil {
		panic(err)
	}
}

// func main() {
// 	// create a new collector
// 	c := colly.NewCollector(
//         // colly.AllowedDomains("gradescope.com", "https://www.gradescope.com"),
//     )
//     //
//     // client := twilio.NewClient("ACaf356b30e339a6e9d0dba1f51aa4d989", "284305ff72993b5ab89fb22732ed3f1d", nil)
//     //
//     // msg, _ := client.Messages.SendMessage("+18787898352", "+14028407963", "Sent via go :) ✓", nil)
//     // fmt.Println(msg.Sid, msg.FriendlyPrice())
//
//
//     // authenticity_token
//     // csrf-token
//     // extract TH_CSRF token for the session
// 	// c.OnHTML("form[role=form] input[type=hidden][name=csrf-token]", func(e *colly.HTMLElement) {
// 	// 	thCsrf := e.Attr("value")
// 	// 	log.Println(thCsrf)
// 	// 	return
//     // })
//     c.OnHTML("[name=csrf-token]", func(e *colly.HTMLElement) {
// 		thCsrf := e.Attr("content")
// 		log.Println(thCsrf)
// 		return
//     })
//
// 	// authenticate
// 	err := c.Post("https://www.gradescope.com/login", map[string]string{"session_email": "jerome.schmidt@students.makeschool.com", "session_password": "Jms1014*neb123", "csrf-token": thCsrf})
// 	if err != nil {
// 		log.Fatal(err)
// 	}
//
//     // On every a element which has href attribute call callback
// 	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
//         log.Println(e.Text)
//         link := e.Attr("href")
//         if !strings.Contains(link, "/courses/") {
//             return
//         }
//
//
//         // if !strings.HasPrefix(link, "/courses/") {
//         //     return
//         // }
//         // if e.Attr("class") == "courseBox" {
// 		// /	return
//     // })
//     //
//         // e.Request.Visit(link)
//         log.Println(link)
//         c.Visit(e.Request.AbsoluteURL(link))
//     //     // link := e.Attr("href")
// 	// 	// // If link start with browse or includes either signup or login return from callback
// 	// 	// if !strings.HasPrefix(link, "/browse") || strings.Index(link, "=signup") > -1 || strings.Index(link, "=login") > -1 {
// 	// 	// 	return
// 	// 	// }
// 	// 	// // start scaping the page under the link found
// 	// 	// e.Request.Visit(link)
//     //
//     })
//
//     // Before making a request print "Visiting ..."
// 	c.OnRequest(func(r *colly.Request) {
// 		log.Println("visiting", r.URL.String())
// 	})
//
// 	// attach callbacks after login
// 	c.OnResponse(func(r *colly.Response) {
// 		log.Println("response received", r.StatusCode)
//         c.Visit("https://www.gradescope.com/courses/233222")
// 	})
//
// 	// start scraping
// 	c.Visit("https://www.gradescope.com/login")
// }
