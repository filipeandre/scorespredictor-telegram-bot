package main

import (
	"encoding/json"
	"fmt"
	"github.com/gocolly/colly"
	"github.com/tealeg/xlsx"
	"strings"
	"time"
)

var (
	dbFileName =  "/tmp5711903412.db"
	xlsxFileName = "scorespredictor.xlsx"
	homeTimeZone =  "Europe/Rome"
	websiteTimeZone = "America/New_York"
)

func main() {
	conf := loadConfiguration()
	location,_ := time.LoadLocation(homeTimeZone)
	date := time.Now().In(location).Format("2006-01-02")
	respC := map[int]chan *htmlTable{}

	//Start fetching data
	for i := 0; i < len(conf.Sports); i++ {
		name:=conf.Sports[i]
		c := collect(conf.Url, name, date)
		respC[i] = make(chan *htmlTable)
		respI:=i

		//Populate chanel store on response received
		c.OnResponse(func(r *colly.Response) {
			stdout(fmt.Sprintf("Response received, status: %d", r.StatusCode))
			data := &htmlTable{}
			err := json.Unmarshal(r.Body, data)
			respC[respI] <- data

			if err != nil {
				stderr(err)
			}
		})
	}

	//Wait each response and generate the xlsx object
	file := xlsx.NewFile()
	dbFilePath := conf.Temp + dbFileName
	removeFile(dbFilePath)
	db:= createDb(dbFilePath, conf.Sports)
	parts:= strings.Split(conf.Filter, ";")

	for i:= 0; i< len(conf.Sports); i++{
		name:= conf.Sports[i]

		var filter string
		if len(parts) > i{
			filter=parts[i]
		}else{
			filter=parts[0]
		}

		response := <- respC[i]
		hasLeague := name == "HOCKEY"
		seedTable(name, response.Html, db, hasLeague)
		saveSheet(name, name, response.Html, file, db, "", hasLeague)
		saveSheet(name + " - filtred",name, response.Html, file, db, filter, hasLeague)
		sendTelegramMessage(conf.Telegram.Token, conf.Telegram.Channel, generateMarkdown(response.Html, name, filter, db, hasLeague))
	}

	//Close the database and remove it
	err:= db.Close()
	if err != nil {
		stderr(err)
	}
	removeFile(dbFilePath)

	//Remove any existent xlsx and write the new xlsx to disk, them remove it
	fileName:= conf.Temp + "/" + date + "_" + xlsxFileName
	removeFile(fileName)
	err = file.Save(fileName)
	if err != nil {
		stderr(err)
	}
	stdout(fmt.Sprintf("Saved: %s", fileName))
	sendTelegramFile(conf.Telegram.Token, conf.Telegram.Channel2, fileName)
	removeFile(fileName)
}
