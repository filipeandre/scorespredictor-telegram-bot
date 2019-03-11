package main

import (
	"database/sql"
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	_ "github.com/mattn/go-sqlite3"
	"math"
	"strconv"
	"strings"
	"time"
)

//Create a simple temp database
func createDb(path string, tables []string) *sql.DB{

	database, err := sql.Open("sqlite3", path)
	if err != nil {
		stderr(err)
	}

	for i := 0; i < len(tables); i++ {
		statement, err := database.Prepare("CREATE TABLE "+ tables[i] +" (line TEXT, goalsNum NUMBER, confidenceNum NUMBER)")
		if err != nil {
			stderr(err)
		}
		_, err = statement.Exec()
		if err != nil{
			stderr(err)
		}
	}

	return database
}

//Seed a table
func seedTable(name string, htmlStr string, db *sql.DB, hasLeague bool){

	doc, err := goquery.NewDocumentFromReader(strings.NewReader("<html><body>" + htmlStr + "</body></html>"))
	if err != nil {
		stderr(err)
	}

	// Select the first table
	tableHtml:= doc.Find("table").First()

	//HOCKEY have more on column
	extraOffset:=0
	if name == "HOCKEY"{
		extraOffset=1
	}

	tableHtml.Find("tbody").Find("tr").Each(func(indexTr int, rowHtml *goquery.Selection) {

		var cells []string
		rowHtml.Find("td").Each(func(indexTh int, tableCell *goquery.Selection) {
			var value string

			if tableCell.HasClass("Date"){
				value, _ = tableCell.Find("[name=gdate]").Attr("value")
				location1,_ := time.LoadLocation(websiteTimeZone)
				t1, err :=  time.ParseInLocation("2006-01-02 15:04:05", value, location1 )

				if err != nil {
					stderr(err)
				}
				location2,_ := time.LoadLocation(homeTimeZone)
				value = t1.In(location2).Format("2006-01-02 15:04:05")

			}else if tableCell.Children().HasClass("Name"){
				value = tableCell.Find(".team").First().Contents().First().Text()

			}else{
				value = tableCell.Text()
			}

			cells = append(cells,value)

			//Add the league value
			if tableCell.HasClass("Date") && extraOffset == 1 && hasLeague {
				value2, _ := tableCell.Find("[name=c_league]").Attr("value")
				cells = append(cells,value2)
			}
		})

		if len(cells) > 2 {
			//Ignore last cell (cells[:len(cells) -1]) (Empty result)
			b, err := json.Marshal(cells[:len(cells) -1])
			if err != nil {
				stderr(err)
			}

			parts := strings.Split(cells[3+extraOffset], ":")
			if len(parts) == 2{
				p1, _ := strconv.Atoi(parts[0])
				p2, _ := strconv.Atoi(parts[1])
				confidenceNum, _ := strconv.ParseFloat(strings.TrimSuffix(cells[4+extraOffset], "%"), 64)
				insertIntoTable(db, name, string(b), int(math.Abs(float64(p2-p1))), confidenceNum)
			}
		}
	})
}

//Insert a line into a table
func insertIntoTable(database *sql.DB, table string, line string, goalsNum int, confidenceNum float64){
	statement, err := database.Prepare("INSERT INTO "+ table +" (line, goalsNum, confidenceNum) VALUES (?, ?, ?)")
	if err != nil {
		stderr(err)
	}
	_, err =statement.Exec(line, goalsNum, confidenceNum)
	if err != nil{
		stderr(err)
	}
}