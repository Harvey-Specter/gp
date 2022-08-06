package main

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func breakwithvol(db *sql.DB, rqParam string, m int, tname string) []map[string]string {
	fmt.Println("date==dltp3l="+rqParam+"  ma_n=", m)
	voltimes := 2.5
	market := "1"
	if tname == "dayline_jp" {
		market = "2"
	}
	dms := getDm(db, rqParam, tname)
	enter := `
`
	rs := enter
	rs0 := enter
	var dataMapArray []map[string]string
	// var dataMapArray0 []map[string]string

	for _, dm := range dms {
		//fmt.Println(dm)
		dataMap := make(map[string]string)
		// dataMap0 := make(map[string]string)
		sps := []float64{}
		//dmSql := `SELECT id,rq,dm,sp,zg,zd,m5,cjl FROM rx where dm = "` + dm + `" and rq<'2021-03-10' ORDER BY dm,rq DESC`
		//dmSql := `SELECT id,rq,dm,sp,zg,zd,m5,cjl FROM rx where dm = "` +
		//  dm + `" and rq <='` + rqParam + `' ORDER BY dm,rq DESC`

		sqlMa := ""
		if m == 20 {
			sqlMa = "m20"
		} else if m == 30 {
			sqlMa = "m30"
		} else if m == 60 {
			sqlMa = "m60"
		}
		dmSql := "SELECT id, date rq,code dm, close sp, open kp, high zg, low zd, m5, volume as cjl,pre_close qsp, " + sqlMa + " , m60 FROM " + tname + " where code = '" + dm["code"].(string) + "' and date<='" + rqParam + "' ORDER BY date DESC"
		// fmt.Print(dmSql, "\n")

		stmt, err := db.Prepare(dmSql)
		if err != nil {
			fmt.Printf("query prepare err:%s\n", err.Error())
			return nil
		}
		rows, err := stmt.Query()
		defer Closedb(stmt, rows)
		if err != nil {
			fmt.Printf("query err:%s\n", err.Error())
			return nil
		}
		var lastsp float64
		cnt := 0
		var sp3 float64
		var sp2 float64
		var sp1 float64
		var sp0 float64
		var rq3 string
		var rq2 string
		var rq1 string
		var rq0 string
		var cjl0 float64
		var cjl1 float64
		var cjl2 float64
		var cjl3 float64
		var qs0sps []float64
		var qs1sps []float64
		var qs0rqs []string
		var qs1rqs []string
		cjlx := 0
		ok := false
		var qs int
		quekou := false
		pinbar := false
		for rows.Next() {
			//
			var id int
			var rq string
			var dm string
			var sp float64
			var kp float64

			var zg float64
			var zd float64
			var m5 float64
			var cjl float64
			var qsp float64
			var ma_n float64
			var ma_60 float64

			rows.Scan(&id, &rq, &dm, &sp, &kp, &zg, &zd, &m5, &cjl, &qsp, &ma_n, &ma_60)

			if rq > rqParam {
				sps = append(sps, sp)
				continue
			}
			//fmt.Println(rq, "---", sp, m5)
			if cnt == 0 {
				lowk := sp
				if sp >= kp {
					lowk = kp
				}
				if (lowk-zd)/(zg-zd) >= 0.66 {
					pinbar = true
				} else {
					pinbar = false
				}
				if (ma_n == 0 || sp <= ma_n || sp <= ma_60) && !pinbar { //|| sp < zg*0.95 {
					break
				} else {
					if zd > qsp*(1.01) {
						quekou = true
					}
					qs1sps = append(qs1sps, sp)
					qs1rqs = append(qs1rqs, rq)
					sps = append(sps, sp)
					cjl0 += cjl
					cjl2 += cjl
					lastsp = sp
					qs = 1
				}
			} else {
				// fmt.Println(sp, m5, cjl0, cjl)
				if sp*1.6 < lastsp {
					ok = false
					break
				}
				if cnt == 1 {
					// 暂时不考虑昨天的收盘价是否下跌1==2
					if sp <= m5 && 1 == 2 { //这块不能删 && sp < m5 { //cjl0 <= float64(cjl)*1.2 &&
						//if cnt == 2 && sp < m5 && cjl0 != 0 { //
						ok = false
						break
					} else {
						if float64(cjl)*voltimes <= cjl0 { // && cjl0 < cjl*2 { //这里做1, 1.5 和 2倍的切换
							cjlx = 1
						}
						cjl0 += cjl
						cjl2 += cjl
					}
				} else if cnt == 2 {
					cjl1 += cjl
					cjl2 += cjl
				} else if cnt == 3 {
					cjl1 += cjl
					cjl3 += cjl
				} else if cnt == 4 {
					if float64(cjl1)*voltimes <= cjl0 { // && cjl0 < cjl*2 { ////这里做1, 1.5 和 2倍的切换
						cjlx = 2
					}
					cjl3 += cjl
				} else if cnt == 5 {
					cjl3 += cjl
					if float64(cjl3)*voltimes <= cjl2 { // && cjl0 < cjl*2 { ////这里做1, 1.5 和 2倍的切换
						cjlx = 3
					}
				}
				// fmt.Println(rq, sp, ma_n)
				if sp < ma_n {
					qs0sps = append(qs0sps, sp)
					qs0rqs = append(qs0rqs, rq)
					if qs == 1 {
						qs = 0
					}
					if len(qs1sps) > 0 {
						if sp3 == 0 {
							idx3 := -1
							idx3, sp3 = maxValue(qs1sps)
							rq3 = qs1rqs[idx3]
						} else if sp1 == 0 {
							idx1 := -1
							idx1, sp1 = maxValue(qs1sps)
							rq1 = qs1rqs[idx1]
						}
						qs1sps = []float64{}
						qs1rqs = []string{}
					}

				} else {
					qs1sps = append(qs1sps, sp)
					qs1rqs = append(qs1rqs, rq)
					if qs == 0 {
						qs = 1
					}
					if len(qs0sps) > 0 {
						if sp2 == 0 {
							idx2 := -1
							idx2, sp2 = minValue(qs0sps)
							rq2 = qs0rqs[idx2]
						} else if sp0 == 0 {
							idx0 := -1
							idx0, sp0 = minValue(qs0sps)
							rq0 = qs0rqs[idx0]
						}
						qs0sps = []float64{}
						qs0rqs = []string{}
					}
				}
			}
			if sp0 > 0 && sp1 > 0 && sp2 > 0 && sp3 > 0 {
				ok = true
				break
			}
			//fmt.Println(dm, reveSliceF(sps), max2, min2, rq2, max1, min1, rq1, min0, rq0, cjlx, lastsp)
			cnt++
		}
		// 头肩放宽成箱体
		if ok && sp3 <= sp1*1.14 && sp2 >= sp0 {

			//if sp3*(1+0.04) >= sp1 && (quekou || cjlx > 0) {
			// 颈线附近放量修改成中轴放量或者跳空
			if sp3 >= sp0+(sp1-sp0)/2 && (quekou || cjlx > 0) {
				code := transCode(dm["code"].(string))

				if market == "2" {
					if sp3 >= sp1 && (sp3-sp0)/sp0 >= 0.15 {
						dataMap = setDataMap(rqParam, strings.Split(dm["code"].(string), ".")[0], "1", market)
						dataMapArray = append(dataMapArray, dataMap)
					}
				}
				fmt.Println("dltp"+rqParam, strings.Split(dm["code"].(string), ".")[0], sp0, rq0, sp1, rq1, sp2, rq2, sp3, rq3, cjlx, quekou)
				rs += code + enter

			} else if pinbar {
				code := transCode(dm["code"].(string))
				//fmt.Println("dltp00-"+rqParam, dm["code"].(string)[0:6], sp0, rq0, sp1, rq1, sp2, rq2, sp3, rq3, cjlx, quekou)
				rs0 += code + enter
			}
		}
		Closedb(stmt, rows)
	}
	if market == "1" {
		fileName := rqParam + "_" + strconv.Itoa(m) + "_breakwithvol.EBK"
		fileName0 := rqParam + "_" + strconv.Itoa(m) + "_breakwithvol0.EBK"
		saveEBK(rs, fileName)
		saveEBK(rs0, fileName0)
	}
	return dataMapArray

}
