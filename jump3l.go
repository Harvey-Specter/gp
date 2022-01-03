package main

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

func jump3l(db *sql.DB, rqParam string) {
	fmt.Println("date==jump3l=" + rqParam)
	voltimes := 0.8 //1.8
	dms := getDm(db, rqParam)
	for _, dm := range dms {
		sps := []float64{}

		dmSql := "SELECT id, date rq,code dm, close sp, high zg, low zd, m5, volume as cjl,pre_close qsp FROM dayline where code = '" + dm["code"].(string) + "'  ORDER BY code,date DESC"

		//dmSql := "SELECT id, date rq,code dm, close sp, high zg, low zd, m5, volume as cjl,pre_close qsp FROM dayline where code like '600540%'  ORDER BY code,date DESC"

		//fmt.Print(dmSql, "\n")
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
		var lastsp float64
		cnt := 0
		var sp2 float64
		//var max2 float64
		//var min2 float64
		//var rq2 string
		var sp1 float64
		//var max1 float64
		//var min1 float64
		//var rq1 string
		var min0 float64
		//var rq0 string
		var cjl0 float64
		var cjl1 float64
		var cjl2 float64
		var cjl3 float64
		cjlx := 0
		ok := false
		var qs int
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
			//fmt.Println(cnt, rq, sp, qsp, sps, rq, rqParam)

			if rq > rqParam {

				sps = append(sps, sp)
				continue
			}

			if cnt == 0 {
				if m5 == 0 || sp < m5 || sp <= qsp { //|| sp < zg*0.95 {
					ok = false
					break
				} else {
					sps = append(sps, sp)
					lastsp = sp
				}
			}
			if cnt == 1 {
				//fmt.Println(sp < zg*0.99, zg*0.99)
				if m5 == 0 || sp > m5 || sp >= qsp { //|| sp < zg*0.95 {
					ok = false
					break
				} else {
					cjl0 += cjl
					cjl2 += cjl
					qs = 1
				}
			} else {
				// fmt.Println(sp, m5, cjl0, cjl)
				if sp*1.6 < lastsp {
					ok = false
					break
				}
				if cnt == 2 {
					// 暂时不考虑昨天的收盘价是否下跌1==2
					if sp > m5 || sp >= qsp { //这块不能删 && sp < m5 { //cjl0 <= float64(cjl)*1.2 &&
						//if cnt == 2 && sp < m5 && cjl0 != 0 { //
						ok = false
						break
					} else {
						if float64(cjl)*voltimes >= cjl0 { // && cjl0 < cjl*2 { //这里做1, 1.5 和 2倍的切换
							cjlx = 1
						}
						cjl0 += cjl
						cjl2 += cjl
					}
				} else if cnt == 3 {
					cjl1 += cjl
					cjl2 += cjl
				} else if cnt == 4 {
					cjl1 += cjl
					cjl3 += cjl
				} else if cnt == 5 {
					if float64(cjl1)*voltimes >= cjl0 { // && cjl0 < cjl*2 { ////这里做1, 1.5 和 2倍的切换
						cjlx = 2
					}
					cjl3 += cjl
				} else if cnt == 6 {
					cjl3 += cjl
					if float64(cjl3)*voltimes >= cjl2 { // && cjl0 < cjl*2 { ////这里做1, 1.5 和 2倍的切换
						cjlx = 3
					}
				}
				if sp < m5 { // 最后一次下跌过m5的收盘价
					//fmt.Println("111", sp, m5, qs, rq)
					if qs == 1 {
						if sp2 == 0.0 {
							//max2 = zg
							//min2 = zd
							sp2 = sp
							//rq2 = rq
						} else if sp1 == 0.0 {
							//max1 = zg
							//min1 = zd
							sp1 = sp
							//rq1 = rq
						} else if min0 == 0.0 {
							//fmt.Println("22222222", lastsp, sp)
							if lastsp <= sp*1.05 {
								ok = false
								break
							}
							min0 = zd
							//rq0 = rq
						} else if sp != 0 && sp < min0 {
							//fmt.Println("33333333")
							min0 = sp
							//ok = false
							//break
						} else if min0 > 0 && sp > min0 || sp > lastsp {
							//fmt.Println("444444")
							ok = true
							break
						}
					}
					qs = 0
				} else {
					qs = 1
					//ok = true
				}
			}
			//fmt.Println(dm, reveSliceF(sps), max2, min2, rq2, max1, min1, rq1, min0, rq0, cjlx, lastsp)
			cnt++
		}
		//fmt.Println(ok, cjlx, min0)
		//fmt.Println(dm, reveSliceF(sps), max2, min2, rq2, max1, min1, rq1, min0, rq0, cjlx)
		if ok && min0 > 0 { // && cjlx > 1
			//fmt.Println(sps)
			//fmt.Println(rqParam, dm["zjw_name"].(string), dm["code"].(string)[0:6], dm["name"], dm["turnover_ratio"].(float64), dm["pe_ratio"].(float64), reveSliceF(sps), max2, min2, rq2, max1, min1, rq1, min0, rq0, cjlx)
			fmt.Println(rqParam, dm["zjw_name"].(string), dm["code"].(string)[0:6], dm["name"], dm["turnover_ratio"].(float64), dm["pe_ratio"].(float64), reveSliceF(sps), cjlx)
		}
		Closedb(stmt, rows)
	}
}
