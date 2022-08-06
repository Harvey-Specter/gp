package main

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func stars(db *sql.DB, rqParam string, tname string) []map[string]string {
	fmt.Println("date==stars=" + rqParam)
	dms := getDm(db, rqParam, tname)
	market := "1"
	if tname == "dayline_jp" {
		market = "2"
	}
	//fmt.Println("day", "industry", "code", "name", "turnover_ratio", "pe_ratio",
	//	"industry_cnt", "inc_revenue_year_rank", "inc_revenue_annual_rank", "cfo_sales_rank", "leverage_ratio_rank")
	enter := `
`
	rs := enter
	var dataMapArray []map[string]string
	for _, dm := range dms {
		sps := []float64{}
		dataMap := make(map[string]string)
		dmSql := "SELECT id, date rq,code dm, close sp, high zg, low zd, m5, volume as cjl ,pre_close as qsp,open as kp,round(abs(open-close)/open,4) as ch ,round(abs(open-pre_close)/pre_close,4) as ch1 FROM " + tname + " where code = '" + dm["code"].(string) + "' and date<='" + rqParam + "' ORDER BY date DESC limit 8"

		// fmt.Print(dmSql, "\n")
		//------------
		//if(stmt!=nil){}
		stmt, err := db.Prepare(dmSql)
		if err != nil {
			fmt.Printf("query prepare err:%s\n", err.Error())
			return nil
		}

		rows, err := stmt.Query()
		defer Closedb(stmt, rows)
		if err != nil {
			fmt.Printf("query err:%s\n", err.Error())
		}
		cnt := 0
		//var cjl0 float64
		ok := false
		//var qs0 = 0
		//var qs1 = 0
		//var qs1sum = 0
		//var change0 float64
		//var change1 float64
		//var minsp float64
		//var maxsp float64
		//qs1sum0,
		gks := 0
		days := 4
		for rows.Next() {
			//
			var id int
			var rq string
			var dm string
			var sp float64

			var zg float64
			var zd float64
			var m5 float64
			var cjl float64
			var qsp float64
			var kp float64
			var ch float64
			var ch1 float64

			rows.Scan(&id, &rq, &dm, &sp, &zg, &zd, &m5, &cjl, &qsp, &kp, &ch, &ch1)

			if rq > rqParam {
				sps = append(sps, sp)
				continue
			}

			// fmt.Println(ok, cnt, rq, sp, qsp, qs0, qs1, qs1sum, math.Abs((cjl-cjl0)/cjl), math.Abs((sp-qsp)/qsp))
			// fmt.Println(rq, ch, ch1)
			if cnt <= days {
				if ch > 0.027 {
					ok = false
					break
				} else {
					ok = true
				}
				if ch1 > 0.018 {
					gks += 1
				}
			} else if cnt > days {
				break
			}

			cnt++
			//fmt.Println(rq, ch, ch1, ok, gks)
		}
		// fmt.Println("zuih", cnt, ok, gks)

		if ok && gks >= 2 {

			tvLog("stars", dm["code"].(string), rqParam)

			//if ok && qs1sum >= 7 {
			fmt.Println("stars_"+rqParam, dm["code"].(string), reveSliceF(sps))
			// fmt.Println(rqParam, dm["zjw_name"].(string), dm["code"].(string)[0:6], dm["name"], dm["turnover_ratio"].(float64), dm["pe_ratio"].(float64), reveSliceF(sps), qs1sum)
			code := transCode(dm["code"].(string))
			dataMap = setDataMap(rqParam, strings.Split(dm["code"].(string), ".")[0], "3", market)
			dataMapArray = append(dataMapArray, dataMap)
			rs += code + enter
		}
		Closedb(stmt, rows)
	}
	if market == "1" {
		fileName := rqParam + "_stars.EBK"
		saveEBK(rs, fileName)
	}

	return dataMapArray
}
