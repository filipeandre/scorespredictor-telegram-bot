package main

import (
	"database/sql"
	"encoding/json"
	"github.com/PuerkitoBio/goquery"
	"github.com/tealeg/xlsx"
	"strconv"
	"strings"
)


//Get rowspan or colspan index
func getHeaderAttrIndex(attr string, tableHeading *goquery.Selection) int{
	val:=0
	span, exists := tableHeading.Attr(attr)
	if exists{
		i, err :=strconv.Atoi(span)
		if err == nil {
			val = i -1
		}
	}
	return val
}

//Append headers to the xlsx
func appendHeaders(tableHtml *goquery.Selection, sheet *xlsx.Sheet, addLeague bool){

	var (
		row *xlsx.Row
		cell *xlsx.Cell
		rowI int
		colI int
	)

	subIndex:= 0

	tableHtml.Find("thead").Find("tr").Each(func(indexTr int, rowHtml *goquery.Selection) {
		row= sheet.AddRow()

		for n:=1; n<subIndex; n++{
			row.AddCell().Value=""
		}

		rowHtml.Find("td").Each(func(indexTh int, tableHeading *goquery.Selection) {

			//Add league after first
			if addLeague && indexTr == 0 && indexTh == 1 {
				cell = row.AddCell()
				cell.Value = "League"
				cell.Merge(0, 1)
				subIndex = subIndex +1
			}

			rowI = getHeaderAttrIndex("rowspan",tableHeading)
			colI = getHeaderAttrIndex("colspan",tableHeading)

			if indexTr == 0 && rowI > 0{
				subIndex = subIndex +1
			}

			cell = row.AddCell()

			//Exclude final score
			if "Final Score" != tableHeading.Text(){
				if colI + rowI > 0 {
					cell.Merge(colI, rowI)
				}
				cell.Value = tableHeading.Text()
			}

			for n:=0; n<colI; n++{
				row.AddCell().Value=""
			}

		})
	})
}


//Save a new sheet
func saveSheet(name string, table string, htmlStr string, file *xlsx.File, db *sql.DB, where string){
	sheet, err := file.AddSheet(name)
	if err != nil {
		stderr(err)
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader("<html><body>" + htmlStr + "</body></html>"))
	if err != nil {
		stderr(err)
	}

	var (
		row *xlsx.Row
		cell *xlsx.Cell
	)

	// Select the first table and write headers
	tableHtml:= doc.Find("table").First()
	appendHeaders(tableHtml, sheet, table == "HOCKEY")

	rows, _ := db.Query("SELECT line FROM " + table + " " + where)
	var line string
	for rows.Next() {
		row= sheet.AddRow()
		err:= rows.Scan(&line)
		if err != nil {
			stderr(err)
		}
		var cells []string
		err =json.Unmarshal([]byte(line), &cells)
		if err != nil {
			stderr(err)
		}

		for i:=0;i< len(cells); i++{
			cell = row.AddCell()
			cell.Value= cells[i]
		}
	}

	row= sheet.AddRow()
}