package main

import (
	"database/sql"
	"fmt"
	"math"

	_ "github.com/go-sql-driver/mysql"
)

func xc(db *sql.DB, rqParam string, tname string) {
	fmt.Println("date==xc=" + rqParam)
	dms := getDm(db, rqParam, tname)

	fmt.Println("day", "industry", "code", "name", "turnover_ratio", "pe_ratio",
		"industry_cnt", "inc_revenue_year_rank", "inc_revenue_annual_rank", "cfo_sales_rank", "leverage_ratio_rank")
	enter := `
`
	rs := enter
	for _, dm := range dms {
		sps := []float64{}

		dmSql := "SELECT id, date rq,code dm, close sp, high zg, low zd, m5, volume as cjl ,pre_close as qsp FROM " + tname + " where code = '" + dm["code"].(string) + "' and date<='" + rqParam + "' ORDER BY date DESC limit 20"

		// fmt.Print(dmSql, "\n")
		//------------
		//if(stmt!=nil){}
		stmt, err := db.Prepare(dmSql)
		if err != nil {
			fmt.Printf("query prepare err:%s\n", err.Error())
			return
		}

		rows, err := stmt.Query()
		defer Closedb(stmt, rows)
		if err != nil {
			fmt.Printf("query err:%s\n", err.Error())
		}
		cnt := 0
		var cjl0 float64
		ok := false
		var qs0 = 0
		var qs1 = 0
		var qs1sum = 0
		var change0 float64
		var change1 float64
		//var minsp float64
		//var maxsp float64
		qs1sum0, days := 5, 10
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

			rows.Scan(&id, &rq, &dm, &sp, &zg, &zd, &m5, &cjl, &qsp)

			if rq > rqParam {
				sps = append(sps, sp)
				continue
			}

			// fmt.Println(ok, cnt, rq, sp, qsp, qs0, qs1, qs1sum, math.Abs((cjl-cjl0)/cjl), math.Abs((sp-qsp)/qsp))

			if cnt < days {
				if sp > qsp {
					qs1sum++
					if qs1 == 2 || math.Abs((cjl-cjl0)/cjl) >= 1.4 || math.Abs((sp-qsp)/qsp) >= 0.07 {
						//if qs1 == 2 || math.Abs((cjl-cjl0)/cjl) >= 2 || math.Abs((sp-qsp)/qsp) >= 0.07 {
						// fmt.Println(cnt, rq, sp, qsp, qs0, qs1, qs1sum, math.Abs((cjl-cjl0)/cjl), math.Abs((sp-qsp)/qsp))

						ok = false
						break
					} else {
						if change1 < math.Abs((sp-qsp)/qsp) {
							change1 = math.Abs((sp - qsp) / qsp)
						}
						qs1 += 1
						qs0 = 0
					}
				} else if sp <= qsp {
					if qs0 == 2 || math.Abs((cjl-cjl0)/cjl) >= 1.4 || math.Abs((sp-qsp)/qsp) >= 0.05 {
						// fmt.Println(cnt, rq, sp, qsp, qs0, qs1, qs1sum, math.Abs((cjl-cjl0)/cjl), math.Abs((sp-qsp)/qsp))
						ok = false
						break
					} else {
						if change0 < math.Abs((sp-qsp)/qsp) {
							change0 = math.Abs((sp - qsp) / qsp)
						}
						qs0 += 1
						qs1 = 0
					}
				}
				cjl0 = cjl
			} else if cnt >= days && cnt < 18 {
				if (sp-qsp)/qsp >= 0.098 {
					ok = true
					break
				} else {
					// ok = false
				}
			} else if cnt >= 18 {
				break
			}
			cnt++
		}
		// fmt.Println(cnt, ok, change0, change1)

		if ok && qs1sum0 >= 5 && change0 < change1 {
			//if ok && qs1sum >= 7 {
			tvLog("xc", dm["code"].(string), rqParam)

			fmt.Println("xc"+rqParam, dm["code"].(string)[0:6], reveSliceF(sps))
			// fmt.Println(rqParam, dm["zjw_name"].(string), dm["code"].(string)[0:6], dm["name"], dm["turnover_ratio"].(float64), dm["pe_ratio"].(float64), reveSliceF(sps), qs1sum)
			code := transCode(dm["code"].(string))
			rs += code + enter
		}
		Closedb(stmt, rows)
	}
	fileName := rqParam + "_xc.EBK"
	saveEBK(rs, fileName)
}
