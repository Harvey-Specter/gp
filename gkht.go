package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func gkht(db *sql.DB, rqParam string) {
	fmt.Println("date==gkht=" + rqParam)
	dms := getDm(db, rqParam)

	//fmt.Println(len(dms))

	for _, dm := range dms {
		sps := []float64{}

		dmSql := "SELECT id, date rq,code dm, close sp, high zg, low zd, m5, volume as cjl ,pre_close as qsp, open as kp FROM dayline where code = '" + dm["code"].(string) + "'  ORDER BY code,date DESC limit 10"

		//fmt.Println(dmSql)
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
		dataType := 0
		//ok := false

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
			qType := "x"

			rows.Scan(&id, &rq, &dm, &sp, &zg, &zd, &m5, &cjl, &qsp, &kp)

			if rq > rqParam {
				sps = append(sps, sp)
				continue
			}
			/* 倒数 dataType 1234
			   0回  1回  2上  3高
			   0回  1上  2高
			   0回  1回  2高
			   0回  1高
			*/
			if cnt > 3 {
				break
			}
			if cnt == 0 {
				//fmt.Println(rq, cnt, sp, qsp)
				if sp > qsp { // 第0天 上退出
					break
				}
				sps = append(sps, sp)
				//cnt++
			} else if cnt == 1 {
				//fmt.Println(rq, cnt, kp, sp, qsp)
				if sp > kp {
					if kp > qsp*(1.02) { // 第二天s
						dataType = 4
						break
					} else {
						qType = "s"
					}
				} else if sp <= kp { //  第二天回
					qType = "h"
				}
				//cnt++
			} else if cnt == 2 {
				if sp > kp {
					//fmt.Println(rq, cnt, kp, sp, qsp)
					if kp > qsp*(1.02) { // 第3天s
						if qType == "s" {
							dataType = 2
							break
						} else if qType == "h" {
							dataType = 3
							break
						}

						break
					} else {
						qType = "s"
					}
				} else if sp <= kp { //  第3天回
					break
				}
				//cnt++
			} else if cnt == 3 {
				//fmt.Println(rq, cnt, kp, sp, qsp)
				if sp > kp && kp > qsp*(1.02) {
					dataType = 1
					break
				} else {
					break
				}
				//cnt++
			}
			cnt++
		}
		if dataType > 0 {

			//fmt.Println(rqParam, dm["zjw_name"].(string), dm["code"].(string)[0:6], dm["name"], dm["turnover_ratio"].(float64), dm["pe_ratio"].(float64),
			//	dm["cnt"], dm["inc_revenue_year_rank"], dm["inc_revenue_annual_rank"], dm["cfo_sales_rank"], dm["leverage_ratio_rank"],
			//	reveSliceF(sps))
			fmt.Println("gkht"+rqParam, dm["code"].(string)[0:6], reveSliceF(sps), dataType)
			// fmt.Println(rqParam, dm["zjw_name"].(string), dm["code"].(string)[0:6], dm["name"], dm["turnover_ratio"].(float64), dm["pe_ratio"].(float64), reveSliceF(sps), qs1sum)
		}
		Closedb(stmt, rows)
	}
}
