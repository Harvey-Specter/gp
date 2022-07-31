package main

import (
	"database/sql"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func closerValue(zd float64, values []float64) float64 {
	min := values[0]
	diff := 0.0
	for _, v := range values {
		if v > 0 && v < zd && diff == 0.0 {
			min = v
			diff = zd - v
		} else if v > 0 && v < zd && diff > 0.0 && diff > zd-v {
			min = v
			diff = zd - v
		}
	}
	return min
}

func qkht(db *sql.DB, rqParam string, tname string) []map[string]string {
	market := "1"
	if tname == "dayline_jp" {
		market = "2"
	}
	fmt.Println("date==qkht=" + rqParam)
	dms := getDm(db, rqParam, tname)

	fmt.Println("day", "industry", "code", "name", "turnover_ratio", "pe_ratio",
		"industry_cnt", "inc_revenue_year_rank", "inc_revenue_annual_rank", "cfo_sales_rank", "leverage_ratio_rank")

	enter := `
`
	rs := enter
	var dataMapArray []map[string]string
	for _, dm := range dms {
		//fmt.Println(dm)
		sps := []float64{}
		dataMap := make(map[string]string)
		// dmSql := "select * from (select id, date rq,code dm, close sp, high zg, low zd, m5 ,pre_close qsp,open kp from dayline where code = '" +
		// 	dm["code"].(string) +
		// 	//	"000063.XSHE" +
		// 	"'  and date<='" + rqParam + "' order by date desc  limit 60 ) as x order by rq  "
		dmSql := "SELECT id, date rq,code dm, close sp, high zg, low zd, m5 ,pre_close qsp,open kp FROM " + tname + " where code = '" + dm["code"].(string) + "' and date<='" + rqParam + "'  ORDER BY date "
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

		var gkqzg []float64
		var gkdate []string

		gkqzg = append(gkqzg, -1.0)
		gkdate = append(gkdate, "0000-00-00")

		ok := false
		qzg := 0.0
		todayhigh := 0.10
		qksize := 1.003
		//var sp float64
		for rows.Next() {
			var id int
			var rq string
			var dm string
			var sp float64
			var zg float64
			var zd float64
			var m5 float64
			var qsp float64
			var kp float64
			rows.Scan(&id, &rq, &dm, &sp, &zg, &zd, &m5, &qsp, &kp)

			if rq >= rqParam {
				sps = append(sps, sp)

				//fmt.Println(sp, qsp, (sp-qsp)/qsp)
				if (sp-qsp)/qsp > todayhigh {
					//fmt.Println(sp, qsp, (sp-qsp)/qsp, "break")
					ok = false
					break
				}
				//continue
			}
			//fmt.Println("qzg ", rq, qzg, zd, sp)
			if qzg == 0.0 {
				qzg = zg
			} else {

				if qzg*qksize < zd {
					ok = true
					gkqzg = append(gkqzg, qzg)
					gkdate = append(gkdate, rq)

				} else {
					// fmt.Println(len(gkqzg), gkqzg, gkdate)
					if len(gkqzg) > 1 {
						for i := len(gkqzg) - 1; i >= 0; i-- {
							// fmt.Println(rq, i, gkqzg[i])
							if zd <= gkqzg[i]*qksize {

								gkqzg = gkqzg[:i]
								gkdate = gkdate[:i]

								//fmt.Println("-----", rq, i, gkqzg)
							} else {
								break
							}
						}
					}
				}
				if len(gkqzg) > 0 {
					if zd > closerValue(zd, gkqzg)*(1+todayhigh) { //0.1
						ok = false
					} else {
						lowk := sp
						if sp >= kp {
							lowk = kp
						}
						if (lowk-zd)/(zg-zd) >= 0.66 || (sp-qsp)/qsp >= 0.02 {
							ok = true
						} else {
							ok = false
						}
						// ok = true
					}
				}

				qzg = zg
				// fmt.Println("temp ", gkqzg, gkdate, rq, qzg, zd, sp, ok)
			}
			cnt++
		}
		if ok && len(gkqzg) > 1 {

			tvLog("qkht", dm["code"].(string), rqParam)

			fmt.Println("qkht"+rqParam, dm["code"].(string)[0:6], gkqzg, gkdate, reveSliceF(sps))
			code := transCode(dm["code"].(string))
			rs += code + enter

			dataMap = setDataMap(rqParam, strings.Split(dm["code"].(string), ".")[0], "4", market)
			dataMapArray = append(dataMapArray, dataMap)
		}
		Closedb(stmt, rows)
	}
	fileName := rqParam + "_qk.EBK"
	saveEBK(rs, fileName)
	return dataMapArray
}
