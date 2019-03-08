package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/tealeg/xlsx"
	"os"
	"time"
)

func main() {
	conf := loadConfiguration()
	location,_ := time.LoadLocation("Europe/Rome")
	date := time.Now().In(location).Format("2006-01-02")
	collectors := map[string]*colly.Collector{}
	responses := map[string]*htmlTable{}

	//Start fetching data
	for i := 0; i < len(conf.Sports); i++ {
		name:=conf.Sports[i]
		c := collect(conf.Url, name, date)

		//Populate responses map on response
		c.OnResponse(func(r *colly.Response) {
			print(fmt.Sprintf("Response received, status: %d", r.StatusCode))
			data := &htmlTable{}
			err := json.Unmarshal(r.Body, data)
			responses[name]  = data

			if err != nil {
				onError(err)
			}
		})
		collectors[name]  = c
	}

	//Wait each response and generate the xlsx object
	file := xlsx.NewFile()
	db:= createDb(conf.Storage + "/tmp.db", conf.Sports)
	filter:= conf.Filter
	for name, c := range collectors {
		c.Wait()
		seedTable(name, responses[name].Html, db)
		saveSheet(name, name, responses[name].Html, file, db, "")
		saveSheet(name + " - filtred",name, responses[name].Html, file, db, filter)
		sendTelegramMessage(conf.Telegram.Token, conf.Telegram.Channel, generateMarkdown(responses[name].Html, name, filter, db))
	}

	//Write file to disk
	fileName:= conf.Storage + "/" + date + "_scorespredictor.xlsx"
	err := file.Save(fileName)
	if err != nil {
		onError(err)
	}

	print(fmt.Sprintf("Saved: %s", fileName))
	sendTelegramFile(conf.Telegram.Token, conf.Telegram.Channel2, fileName)

	//Remove file at end
	if _, err := os.Stat(fileName); err == nil {
		err := os.Remove(fileName)
		if err != nil {
			onError(err)
		}
	}
}
