package main

import (
	"fmt"
	"log"
	"os"

	"github.com/go-resty/resty/v2"
	"github.com/gocolly/colly"
	"github.com/joho/godotenv"
)

func main() {
	lineNotifyToken := ""

	err := godotenv.Load(".env")
	if err != nil {
		if len(os.Args) == 1 {
			lineNotifyToken = os.Args[0]
		} else {
			log.Fatalln("Please provide at least one argument. or .env file")
			return
		}
	} else {
		lineNotifyToken = os.Getenv("LINE_NOTIFY_TOKEN")
	}

	client := resty.New()

	// Create a new collector
	c := colly.NewCollector()

	c.OnHTML("#content", func(e *colly.HTMLElement) {
		updatedDate := e.ChildText("#rightCol > div.divgta.goldshopf > table > tbody > tr:nth-child(4) > td.span.bg-span.txtd.al-r")
		updatedTime := e.ChildText("#rightCol > div.divgta.goldshopf > table > tbody > tr:nth-child(4) > td.em.bg-span.txtd.al-r")

		// ทองแท่ง
		goldBarBuyPrice := e.ChildText("#rightCol > div.divgta.goldshopf > table > tbody > tr:nth-child(1) > td:nth-child(3)")
		goldBarSellPrice := e.ChildText("#rightCol > div.divgta.goldshopf > table > tbody > tr:nth-child(1) > td:nth-child(2)")

		// ทองรูปพรรณ
		goldBuyPrice := e.ChildText("#rightCol > div.divgta.goldshopf > table > tbody > tr:nth-child(2) > td:nth-child(3)")
		goldSellPrice := e.ChildText("#rightCol > div.divgta.goldshopf > table > tbody > tr:nth-child(2) > td:nth-child(2)")

		goldM := fmt.Sprintf("ทองคำรูปพรรณ 96.5 \n ขาย %s \n ซื้อ %s", goldSellPrice, goldBuyPrice)
		goldB := fmt.Sprintf("ทองคำแท่ง 96.5 \n ขาย %s \n ซื้อ %s", goldBarSellPrice, goldBarBuyPrice)

		message := fmt.Sprintf("%s %s \n\n %s \n\n %s ", updatedDate, updatedTime, goldM, goldB)

		fmt.Println(message)

		// Send notification to LINE
		_, err = client.R().
			SetHeader("Authorization", fmt.Sprintf("Bearer %s", lineNotifyToken)).
			SetFormData(map[string]string{
				"message": message,
			}).
			Post("https://notify-api.line.me/api/notify")

		if err != nil {
			log.Fatalf("Error sending LINE notification: %v", err)
		}

	})

	c.OnError(func(_ *colly.Response, err error) {
		log.Println("Error occurred:", err)
	})

	// Start the crawling
	err = c.Visit("https://xn--42cah7d0cxcvbbb9x.com/")

	if err != nil {
		log.Fatal(err)
	}
}
