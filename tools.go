package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func minValue(values []float64) (int, float64) {
	min := values[0]
	idx := 0
	for i, v := range values {
		if v > 0 && v < min {
			min = v
			idx = i
		}
	}
	return idx, min
}
func maxValue(values []float64) (int, float64) {
	max := values[0]
	idx := 0
	for i, v := range values {
		if v > 0 && v > max {
			max = v
			idx = i
		}
	}
	return idx, max
}
func transCode(dmcode string) string {
	code := dmcode[0:6]
	if strings.Contains(dmcode, "XSHG") {
		code = "1" + code
	} else {
		code = "0" + code
	}
	return code
}

func saveEBK(cont string, fileName string) {
	var d1 = []byte(cont)
	err2 := ioutil.WriteFile(fileName, d1, 0666) //写入文件(字节数组)
	check(err2)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
func Decimal(value float64) string {

	value1 := fmt.Sprintf("%.2f", value)
	return value1
}

// 获取日期
func getRqsByDB(db *sql.DB, rq string, limit string) []string { //300859-西域; 600718-东软
	var rqs []string
	//dmstmt, err := db.Prepare(`select distinct code as dm from dayline where code = '688339.XSHG' `) // "'002037"`)
	dmstmt, err := db.Prepare(`select distinct date as rq from dayline where date<='` + rq + `' order by date desc limit ` + limit)
	if err != nil {
		fmt.Printf("query prepare err:%s\n", err.Error())
		return rqs
	}
	//defer dmstmt.Close()
	dmrows, err := dmstmt.Query()
	if err != nil {
		fmt.Printf("query err:%s\n", err.Error())
	}
	for dmrows.Next() {
		var rq string
		dmrows.Scan(&rq)
		rqs = append(rqs, rq)
	}
	defer Closedb(dmstmt, dmrows)
	return rqs
}

func getTpDmByDate(db *sql.DB, rq string) []map[string]string { //300859-西域; 600718-东软
	var dms []map[string]string
	sql := `select distinct a.code from tp a where a.date='` + rq + `' `

	dmstmt, err := db.Prepare(sql)
	if err != nil {
		fmt.Printf("query prepare err:%s\n", err.Error())
		return dms
	}
	//defer dmstmt.Close()
	dmrows, err := dmstmt.Query()
	if err != nil {
		fmt.Printf("query err:%s\n", err.Error())
	}
	for dmrows.Next() {
		dmMap := make(map[string]string)
		var code string
		dmrows.Scan(&code)
		dmMap["code"] = code
		dms = append(dms, dmMap)
	}
	defer Closedb(dmstmt, dmrows)
	return dms
}
func tvLog(funcName string, code string, day string) {

	if strings.Contains(code, ".JP") {
		code = code[0:4]
		fmt.Println(funcName, day, "https://www.tradingview.com/chart/CFSEAW1L/?symbol=TSE%3A"+code)
		//https://www.tradingview.com/chart/CFSEAW1L/?symbol=TSE%3A1375
	} else {
		//code = code[0:6]
		//fmt.Println(funcName, day, code)

	}
}
func getDm(db *sql.DB, rq string, tname string) []map[string]interface{} { //300859-西域; 600718-东软
	//var dms []string
	var dms []map[string]interface{}

	// 回踩缺口 605289 002932 600889 a.code like '002155%' and  002155%
	// (a.code like '002128%' or a.code like '002155%' or a.code like '603353%' ) and a.code like '000655%' and
	sql := `select distinct a.code from ` + tname + ` a where a.code not like '688%' and a.code not like '3%' and a.paused='0' and a.date='` + rq + `' `
	// fmt.Println(sql)
	//sql := `select distinct a.code from dayline a where a.code like '688%' and a.date='` + rq + `' `

	dmstmt, err := db.Prepare(sql)

	//dmstmt, err := db.Prepare(`select distinct a.code ,b.pe_ratio, b.turnover_ratio,c.zjw_name, c.name from dayline a ,valuation b ,industry c where a.code=b.code and b.pe_ratio>0 and a.paused='0' and b.code=c.code and c.name not like '%ST%' and a.date=b.day and a.date='` + rq + `' and a.code like '002413%'`)
	if err != nil {
		fmt.Printf("query prepare err:%s\n", err.Error())
		return dms
	}
	//defer dmstmt.Close()
	dmrows, err := dmstmt.Query()
	if err != nil {
		fmt.Printf("query err:%s\n", err.Error())
	}
	for dmrows.Next() {
		dmMap := make(map[string]interface{})
		var code string
		// var pe_ratio float64
		// var turnover_ratio float64
		// var zjw_name string
		// var name string
		// var cnt int
		// var inc_revenue_year_rank int
		// var inc_revenue_annual_rank int
		// var cfo_sales_rank int
		// var leverage_ratio_rank int
		//dmrows.Scan(&code, &pe_ratio, &turnover_ratio, &zjw_name, &name, &cnt, &inc_revenue_year_rank, &inc_revenue_annual_rank, &cfo_sales_rank, &leverage_ratio_rank)
		dmrows.Scan(&code)
		dmMap["code"] = code
		// dmMap["pe_ratio"] = pe_ratio
		// dmMap["turnover_ratio"] = turnover_ratio
		// dmMap["zjw_name"] = zjw_name
		// dmMap["name"] = name
		// dmMap["cnt"] = cnt
		// dmMap["inc_revenue_year_rank"] = inc_revenue_year_rank
		// dmMap["inc_revenue_annual_rank"] = inc_revenue_annual_rank
		// dmMap["cfo_sales_rank"] = cfo_sales_rank
		// dmMap["leverage_ratio_rank"] = leverage_ratio_rank

		//c.cnt,inc_revenue_year_rank,inc_revenue_annual_rank,cfo_sales_rank ,leverage_ratio_rank

		dms = append(dms, dmMap)
	}
	defer Closedb(dmstmt, dmrows)

	return dms
}

func reveSlice(s []string) []string {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	//fmt.Println(s)
	return s
}

func reveSliceF(s []float64) []float64 {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
	//fmt.Println(s)
	return s
}

func getRqs(buyDate string) []string {
	currentTime := time.Now()
	var rqs []string
	i := -2
	for {

		oldTime := currentTime.AddDate(0, 0, i)
		if int(oldTime.Weekday()) == 0 || int(oldTime.Weekday()) == 6 {
			i--
			continue
		}
		oldTimeStr := oldTime.Format("2006-01-02")
		rqs = append(rqs, oldTimeStr)
		//fmt.Println(oldTimeStr, oldTime.Weekday(), int(oldTime.Weekday()), buyDate)
		if oldTimeStr == buyDate || i == -100 {
			break
		}
		i--
	}
	return reveSlice(rqs)
	//return rqs
}
