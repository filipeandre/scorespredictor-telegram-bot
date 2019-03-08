package main

import (
	"database/sql"
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	tb "gopkg.in/tucnak/telebot.v2"
	"strings"
)

type to struct {
	Channel string
}

func (obj to) Recipient() string {
	return  obj.Channel
}

var bot *tb.Bot

func initBot(token string){
	var err error
	bot, err = tb.NewBot(tb.Settings{
		Token:  token,
		URL: "https://api.telegram.org",
	})
	if err != nil {
		stderr(err)
	}
}

//Generate telegram message
func generateMarkdown(htmlStr string, tableName string, where string, db *sql.DB, ) string{

	doc, err := goquery.NewDocumentFromReader(strings.NewReader("<html><body>" + htmlStr + "</body></html>"))
	if err != nil {
		stderr(err)
	}

	// Select the first table and write headers
	tableHtml:= doc.Find("table").First()
	first:=tableHtml.Find("thead").Find("tr").First()
	sec:=first.Next()

	var (
		headers   []string
		rowI      int
		colI      int
		secondary *goquery.Selection
	)

	//Single line table head
	first.Find("td").Each(func(indexTh int, tableHeading *goquery.Selection) {

		//Exclude final score
		if "Final Score" == tableHeading.Text(){
			return
		}

		//Add league header
		if indexTh == 1 {
			headers = append(headers, "League")
		}


		rowI = getHeaderAttrIndex("rowspan",tableHeading)
		colI = getHeaderAttrIndex("colspan",tableHeading)

		if rowI == 0{

			for i:= 0;i<= colI;i++{
				if secondary == nil {
					secondary = sec.Find("td").First()
				}else{
					secondary = secondary.Next()
				}

				headers = append(headers, secondary.Text())
			}

		}else{
			headers = append(headers, tableHeading.Text())
		}
	})

	//Will add all the info into each row / cell
	var data []string
	rows, _ := db.Query("SELECT line FROM " + tableName + " " + where)
	var line string
	for rows.Next() {
		err:= rows.Scan(&line)
		if err != nil {
			stderr(err)
		}
		var cells []string
		err =json.Unmarshal([]byte(line), &cells)
		if err != nil {
			stderr(err)
		}
		compact := ""
		for i:=0; i< len(cells);i++{
			compact = compact + "*" +headers[i] + "*: " + cells[i] + " \n"
		}

		data = append(data, compact)
	}

	var tableData string

	if len(data) > 0{
		tableData = strings.Join(data, "\n")
	}else {
		tableData = "No Matches found today."
	}

	return "*" +tableName + " " + where + "*\n" + tableData
}

//Send telegram message to a specific chanel or chat id
func sendTelegramMessage(token string, channel string, msg string){

	if channel == ""{
		return
	}

	if bot == nil{
		initBot(token)
	}

	_, err := bot.Send( &to{Channel: channel}, msg, &tb.SendOptions{
		ParseMode: tb.ModeMarkdown,
	})
	if err != nil {
		stderr(err)
	}
}

func sendTelegramFile(token string, channel string, path string){

	if channel == ""{
		return
	}

	if bot == nil{
		initBot(token)
	}

	f := &tb.Document{File: tb.FromDisk(path)}
	_, err := bot.Send( &to{Channel: channel}, f)
	if err != nil {
		stderr(err)
	}

}