package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/tealeg/xlsx"
	"os"
	"strings"
	"time"
)

func main() {
	conf := loadConfiguration()
	location,_ := time.LoadLocation("Europe/Rome")
	date := time.Now().In(location).Format("2006-01-02")
	respC := map[int]chan *htmlTable{}

	//Start fetching data
	for i := 0; i < len(conf.Sports); i++ {
		name:=conf.Sports[i]
		c := collect(conf.Url, name, date)
		respC[i] = make(chan *htmlTable)
		respI:=i

		//Populate respC map on response
		c.OnResponse(func(r *colly.Response) {
			print(fmt.Sprintf("Response received, status: %d", r.StatusCode))
			data := &htmlTable{}
			err := json.Unmarshal(r.Body, data)
			respC[respI] <- data

			if err != nil {
				onError(err)
			}
		})
	}

	//Wait each response and generate the xlsx object
	file := xlsx.NewFile()
	db:= createDb(conf.Storage + "/tmp.db", conf.Sports)
	parts:= strings.Split(conf.Filter, ";")



	for i:= 0; i< len(conf.Sports); i++{
		name:= conf.Sports[i]

		var filter string
		if len(parts) >= i{
			filter=parts[i]
		}else{
			filter=parts[0]
		}

		response := <- respC[i]
		seedTable(name, response.Html, db)
		saveSheet(name, name, response.Html, file, db, "")
		saveSheet(name + " - filtred",name, response.Html, file, db, filter)
		sendTelegramMessage(conf.Telegram.Token, conf.Telegram.Channel, generateMarkdown(response.Html, name, filter, db))
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
