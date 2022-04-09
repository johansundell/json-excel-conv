package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/xuri/excelize/v2"
)

func init() {
	routes = append(routes, Route{"simpleHandler", "POST", "/sudde", simpleHandler})
}

type fmData struct {
	Data []struct {
		Pos   string `json:"pos"`
		Type  int    `json:"type"`
		Value string `json:"value"`
	} `json:"data"`
}

func simpleHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Making excel file")
	data := fmData{}
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		logger.Error(err)
		http.Error(w, err.Error(), 500)
		return
	}

	f := excelize.NewFile()
	defer f.Close()
	for _, row := range data.Data {
		if row.Type == 0 {
			f.SetCellValue("Sheet1", row.Pos, row.Value)
		} else if row.Type == 1 {
			n, err := strconv.Atoi(row.Value)
			if err == nil {
				f.SetCellValue("Sheet1", row.Pos, n)
			}
		} else if row.Type == 2 {
			fmt.Println("Calc", row.Value[1:])
			if err := f.SetCellFormula("Sheet1", row.Pos, row.Value[1:]); err != nil {
				fmt.Println(err)
			}
			test, _ := f.CalcCellValue("Sheet1", row.Pos)
			fmt.Println(test)
		} else if row.Type == 3 {
			n, err := strconv.ParseFloat(row.Value, 64)
			//fmt.Println(n, err)
			if err == nil {
				//fmt.Println("here")
				err := f.SetCellValue("Sheet1", row.Pos, n)
				//fmt.Println(err)
				_ = err
			}
		}
	}
	//f.UpdateLinkedValue()

	/*if err := f.SaveAs("Book1.xlsx"); err != nil {
		fmt.Println(err)
	}*/
	w.Header().Set("Content-Disposition", "attachment; filename=Book1.xlsx")
	// application/octet-stream
	// application/vnd.ms-excel
	// application/vnd.openxmlformats-officedocument.spreadsheetml.sheet
	w.Header().Set("Content-Type", "application/octet-stream")
	if _, err := f.WriteTo(w); err != nil {
		fmt.Fprintf(w, err.Error())
		fmt.Println(err)
		return
	}
}
