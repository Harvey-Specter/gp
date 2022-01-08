#-*- coding: UTF-8 -*-
import sys
import pymysql
from jqdatasdk import *
import datetime
import time 
from jqdatasdk.technical_analysis import *
from sqlalchemy.sql import func
import pandas 
#from jqlib.technical_analysis 6214 6239 5000 1094 741 import * 
def updateYeardayCap(rq):
    conn,cs=getConn()
    try:
        idsql = " SELECT id FROM valuation WHERE day = '%s' " % (rq)
        #idsql = "SELECT id FROM valuation WHERE id = 241392 " % (rq)
        #print(idsql)
        cs.execute(idsql)
        results = cs.fetchall()
        idList = []
        for row in results:
            idList.append(row[0])
        saveBatch(idList,'yearday_avg_market_cap')
        #print (results)
    except Exception as e:
        print("出现如下异常%s"%e)
        return 
    closeConn(conn,cs)

def replaceFloatingTop10(rq):

    codes = getCodes(rq)

    for code in codes:
        floatingTop10=[]
        q=query(finance.STK_SHAREHOLDER_FLOATING_TOP10).filter(finance.STK_SHAREHOLDER_FLOATING_TOP10.code==code[0],finance.STK_SHAREHOLDER_FLOATING_TOP10.end_date>'2021-09-01').limit(10)
        df=finance.run_query(q)

        for index, row in df.iterrows():
            
            code = row['code']
                
            rowlist=[
                row['id'],row['company_id'],row['company_name'],row['code'],row['end_date'],row['pub_date'],
                row['change_reason_id'],row['change_reason'],row['shareholder_rank'],row['shareholder_id'],
                row['shareholder_name'],row['shareholder_name_en'],row['shareholder_class_id'],row['shareholder_class'],
                row['share_number'],row['share_ratio'],row['sharesnature_id'],row['sharesnature'],
            
            ]
            floatingTop10.append(rowlist)

        #print(code+':'+a[key])
        saveBatch(floatingTop10,'stk_shareholder_floating_top10')


def replaceIndustry(rq):

    codes = getCodes(rq)
    codelist=[]
    for code in codes:
        codelist.append(code[0])

    industryList=[]
    inds=get_industry(security=codelist, date=rq)

    df = get_fundamentals(query(
            balance.code,
            indicator.inc_revenue_year_on_year, 
            indicator.inc_revenue_annual, 
            cash_flow.net_operate_cash_flow,
            income.operating_revenue, #cfo_sales net_operate_cash_flow/operating_revenue
            income.basic_eps, #每股收益
            balance.total_assets,
            balance.total_owner_equities,
            balance.pubDate,
            balance.statDate,
            indicator.ocf_to_revenue, 
            indicator.adjusted_profit
        ).filter(
          # 这里不能使用 in 操作, 要使用in_()函数
          balance.code.in_(codelist)
        ).order_by(
            balance.code
        ).limit(
            5000
        ), 
            date=rq
            #statDate='2020'
        )

    for index, row in df.iterrows():
        
        code = row['code']

        if  code in inds.keys():
            sInfo = get_security_info(code)
            display_name = sInfo.display_name
            end_date = sInfo.end_date
            #print(code,display_name,end_date)
            if str(end_date)=='2200-01-01' and  not 'ST' in display_name:
                # print(code,display_name,'----------',end_date)
                rowlist=[code,display_name,
                inds[code]['sw_l1']['industry_code'],inds[code]['sw_l1']['industry_name'],
                inds[code]['sw_l2']['industry_code'],inds[code]['sw_l2']['industry_name'],
                inds[code]['sw_l3']['industry_code'],inds[code]['sw_l3']['industry_name'],
                #'','','','','','',
                inds[code]['zjw']['industry_code'],  inds[code]['zjw']['industry_name'],

                row['inc_revenue_year_on_year'] if str(row['inc_revenue_year_on_year'])!='nan' else 0,
                row['inc_revenue_annual'] if str(row['inc_revenue_annual'])!='nan' else 0,

                row['net_operate_cash_flow'] if str(row['net_operate_cash_flow'])!='nan' else 0,
                row['operating_revenue'] if str(row['operating_revenue'])!='nan' else 0,
                row['total_assets'] if str(row['total_assets'])!='nan' else 0,
                row['total_owner_equities'] if str(row['total_owner_equities'])!='nan' else 0,

                row['statDate'] if str(row['statDate'])!='nan' else 0,
                row['pubDate'] if str(row['pubDate'])!='nan' else 0,
                row['ocf_to_revenue'] if str(row['ocf_to_revenue'])!='nan' else 0,
                row['adjusted_profit'] if str(row['adjusted_profit'])!='nan' else 0,
                rq
                ]
                industryList.append(rowlist)

        #print(code+':'+a[key])
    saveBatch(industryList,'industry1')

    ranks = getRank(rq)
    
    ranklist=[]
    for rank in ranks:
        ranklist.append([rank[0],rank[1],rank[2],rank[3],rank[4],rank[5],rank[6],rank[7],rank[8],rank[9],rank[10],rank[11],rank[12],rank[13],rank[14],rank[15],rank[16]])
    saveBatch(ranklist,'updateIndustry')

def valuationDay(rq,tName):
    valuationlist=[]
    codelist=[]
    df = get_fundamentals( query( 
        valuation.code , 
        valuation.day,
        valuation.capitalization, 
        valuation.circulating_cap, 
        valuation.market_cap,
		valuation.circulating_market_cap,  
		valuation.turnover_ratio, 
        valuation.pe_ratio, 
        valuation.pe_ratio_lyr, 
        valuation.pb_ratio, 
        valuation.ps_ratio, 
        valuation.pcf_ratio		
    ).limit(
        5000
    ), date=rq)
    for index, row in df.iterrows():
        
        if str(row['pe_ratio'])!='nan' and is_number(row['pe_ratio']) and row['pe_ratio']>0  and row['pe_ratio']>0 and row['pb_ratio']>0 and row['ps_ratio']>0:
            #a if a>1 else b
            rowlist = [row['code'],row['day'], 
            row['capitalization'] if str(row['capitalization'])!='nan' else 0,
            row['circulating_cap'] if str(row['circulating_cap'])!='nan' else 0,
            row['market_cap'] if str(row['market_cap'])!='nan' else 0,
            row['circulating_market_cap'] if str(row['circulating_market_cap'])!='nan' else 0,
            row['turnover_ratio'] if str(row['turnover_ratio'])!='nan' else 0,
            row['pe_ratio'] if str(row['pe_ratio'])!='nan' else 0,
            row['pe_ratio_lyr'] if str(row['pe_ratio_lyr'])!='nan' else 0,
            row['pb_ratio'] if str(row['pb_ratio'])!='nan' else 0,
            row['ps_ratio'] if str(row['ps_ratio'])!='nan' else 0,
            row['pcf_ratio'] if str(row['pcf_ratio'])!='nan' else 0
            ]
            valuationlist.append(rowlist)
            #codelist.append(row['code'])
    if len(valuationlist)>0:
        saveBatch(valuationlist,tName)
    else:
        print("empty valuationlist")

    # valuationRanks = getValuationRanks(rq)
    
    # updateValuationlist=[]
    # for rank in valuationRanks:
    #     updateValuationlist.append([rank[0],rank[1],rank[2],rank[3],rank[4],rank[5],rank[6]])
    # saveBatch(updateValuationlist,'updateValuation')
    #return codelist

def is_number(s):
    try:
        float(s)
        return True
    except ValueError:
        pass
 
    try:
        import unicodedata
        unicodedata.numeric(s)
        return True
    except (TypeError, ValueError):
        pass
 
    return False

def getCodes(rq):
    conn,cs=getConn()
    try:
        #sql = "SELECT distinct code FROM industry1 "
        #sql = "SELECT distinct code FROM dayline WHERE date = '2021-03-01' " 
        sql = "SELECT distinct code FROM dayline WHERE date = '%s' " % (rq)
        #sql = "SELECT distinct code FROM dayline WHERE date = '%s' and code in('600126.XSHG','600718.XSHG') " % (rq)
        #print(sql)
        cs.execute(sql)
        results = cs.fetchall()
        return results
    except Exception as e:
        print("出现如下异常%s"%e)
        return
    closeConn(conn,cs)

def getRank(rq):
    conn,cs=getConn()
    try:
        sql = "select (select count(1) from industry1 b where b.sw_l1_code=a.sw_l1_code and b.day=a.day) cnt, \
rank() over(partition by sw_l1_code ORDER BY inc_revenue_year desc)  inc_year_rank, \
round((select sum(inc_revenue_year)/count(1) from industry1 b where b.sw_l1_code=a.sw_l1_code and b.day=a.day),4) inc_year_avg , \
rank() over(partition by sw_l1_code ORDER BY inc_revenue_annual desc) inc_annual_rank , \
round((select sum(inc_revenue_annual)/count(1) from industry1 b where b.sw_l1_code=a.sw_l1_code and b.day=a.day),4) inc_annual_avg, \
round(net_operate_cash_flow/operating_revenue,4) as cfo_sales, \
rank() over(partition by sw_l1_code ORDER BY net_operate_cash_flow/operating_revenue desc)  cfo_sales_rank , \
round((select sum(net_operate_cash_flow/operating_revenue)/count(1) from industry1 b where b.sw_l1_code=a.sw_l1_code and b.day=a.day),4) cfo_sales_avg, \
round(total_assets/total_owner_equities,4) as leverage_ratio, \
rank() over(partition by sw_l1_code ORDER BY total_assets/total_owner_equities)  leverage_ratio_rank , \
round((select sum(total_assets/total_owner_equities)/count(1) from industry1 b where b.sw_l1_code=a.sw_l1_code and b.day=a.day),4) leverage_ratio_avg, \
rank() over(partition by sw_l1_code ORDER BY ocf_to_revenue desc)  ocf_to_revenue_rank , \
round((select sum(ocf_to_revenue)/count(1) from industry1 b where b.sw_l1_code=a.sw_l1_code and b.day=a.day),4) ocf_to_revenue_avg, \
round(adjusted_profit/operating_revenue,4) as adjusted_profit_revenue, \
rank() over(partition by sw_l1_code ORDER BY adjusted_profit/operating_revenue desc)  adjusted_profit_revenue_rank , \
round((select sum(adjusted_profit/operating_revenue)/count(1) from industry1 b where b.sw_l1_code=a.sw_l1_code and b.day=a.day),4) adjusted_profit_revenue_avg, id \
from industry1 a where day='"+rq+"'" 
        #sql = "SELECT distinct code FROM dayline WHERE date = '%s' and code in('600126.XSHG','600718.XSHG') " % (rq)
        #print(sql)
        cs.execute(sql)
        results = cs.fetchall()
        return results
    except Exception as e:
        print("出现如下异常%s"%e)
        return
    closeConn(conn,cs)

def getValuationRanks(rq):
    conn,cs=getConn()
    try:
        sql = "select \
rank() over(partition by sw_l1_code ORDER BY pe_ratio)  pe_ratio_rank , \
round((select sum(pe_ratio)/count(1) from valuation a ,industry1 b \
where a.day='"+rq+"' and a.pe_ratio>0 and a.pb_ratio>0 and a.ps_ratio>0 and a.code =b.code and b.sw_l1_code=y.sw_l1_code),4) pe_ratio_avg , \
rank() over(partition by sw_l1_code ORDER BY pb_ratio)  pb_ratio_rank , \
round((select sum(pb_ratio)/count(1) from valuation a ,industry1 b  \
where a.day='"+rq+"' and a.pe_ratio>0 and a.pb_ratio>0 and a.ps_ratio>0 and a.code =b.code and b.sw_l1_code=y.sw_l1_code),4) pb_ratio_avg , \
rank() over(partition by sw_l1_code ORDER BY ps_ratio)  ps_ratio_rank , \
round((select sum(ps_ratio)/count(1) from valuation a ,industry1 b \
where a.day='"+rq+"' and a.pe_ratio>0 and a.pb_ratio>0 and a.ps_ratio>0 and a.code =b.code and b.sw_l1_code=y.sw_l1_code),4) ps_ratio_avg ,x.id \
from valuation x ,industry1 y \
where x.code=y.code and x.day='"+rq+"' and pe_ratio>0 and pb_ratio>0 and ps_ratio>0 " 
        #sql = "SELECT distinct code FROM dayline WHERE date = '%s' and code in('600126.XSHG','600718.XSHG') " % (rq)
        #print(sql)
        cs.execute(sql)
        results = cs.fetchall()
        return results
    except Exception as e:
        print("getValuationRanks出现如下异常%s"%e)
        return
    closeConn(conn,cs)

def genMA(code,rq,n):
    #print(rq+'======genM5======='+code)
    conn,cs=getConn()
    try:
        sql = "SELECT id,close FROM dayline WHERE code = '%s' and date<='%s' and paused='0' order by date desc limit %d" % (code,rq ,n)

        cs.execute(sql)

        results = cs.fetchall()
        ma = 0 
        id = -1
        if len(results)==n:
          for row in results:
            if id==-1:
                id=row[0]
            ma+=row[1]

          ma=round(ma/n,2) 
        else:
          ma = -1
        return id,ma
    except Exception as e:
        print("出现如下异常%s"%e)
        return ma
    finally:
        closeConn(conn,cs)

def getSQL(tName):
    sql = ''
    if tName=='dayline':
        sql = "INSERT INTO dayline(date,code,open,close,low,high,volume,money,factor, high_limit,low_limit,avg,pre_close,paused, m5,m10,m20,m60,m250) \
            VALUES (%s,%s,%s,%s,%s, %s,%s,%s,%s,%s, %s,%s,%s,%s,%s,  %s,%s,%s,%s )"
    elif tName=='valuation':
        sql = "INSERT INTO valuation(code,day,capitalization,circulating_cap,market_cap,circulating_market_cap,turnover_ratio,pe_ratio,pe_ratio_lyr, pb_ratio,ps_ratio,pcf_ratio) \
            VALUES (%s,%s,%s,%s,%s, %s,%s,%s,%s,%s, %s,%s )"
    elif tName=='industry':
        sql = "REPLACE INTO industry(code,name,sw_l1_code,sw_l1_name,sw_l2_code,sw_l2_name,sw_l3_code,sw_l3_name, \
            zjw_code,zjw_name ) \
            VALUES (%s,%s, %s,%s,%s,%s, %s,%s,%s,%s  )"
    elif tName=='industry1':
        sql = "REPLACE INTO industry1(code,name,sw_l1_code,sw_l1_name,sw_l2_code,sw_l2_name,sw_l3_code,sw_l3_name, \
            zjw_code,zjw_name ,inc_revenue_year, inc_revenue_annual, net_operate_cash_flow, operating_revenue, \
                total_assets, total_owner_equities ,report_date , public_date,ocf_to_revenue,adjusted_profit ,day ) \
            VALUES (%s,%s,%s,%s,%s, %s,%s,%s,%s,%s ,%s,%s,%s, %s,%s,%s,%s,%s ,%s,%s, %s )"
    elif tName=='updateIndustry':
        # cnt inc_year_rank inc_year_avg inc_nnual_rank inc_annual_avg cfo_sales cfo_sales_rank cfo_sales_avg leverage_ratio leverage_ratio_rank leverage_ratio_avg id
        sql = "update industry1 set cnt=%s ,inc_revenue_year_rank=%s , inc_revenue_year_avg=%s , inc_revenue_annual_rank=%s , inc_revenue_annual_avg=%s , cfo_sales=%s , cfo_sales_rank =%s ,cfo_sales_avg=%s , leverage_ratio=%s , leverage_ratio_rank=%s , leverage_ratio_avg=%s , ocf_to_revenue_rank=%s , ocf_to_revenue_avg=%s , adjusted_profit_revenue=%s , adjusted_profit_revenue_rank=%s , adjusted_profit_revenue_avg=%s  where id =%s "
    
    elif tName=='updateValuation':
        # cnt inc_year_rank inc_year_avg inc_nnual_rank inc_annual_avg cfo_sales cfo_sales_rank cfo_sales_avg leverage_ratio leverage_ratio_rank leverage_ratio_avg id
        sql = "update valuation set pe_ratio_rank=%s ,pe_ratio_avg=%s , pb_ratio_rank=%s , pb_ratio_avg=%s , ps_ratio_rank=%s , ps_ratio_avg=%s where id =%s "

    elif tName=='yearday_avg_market_cap':
        sql="update valuation a set yearday_avg_market_cap=(select avg(market_cap) from valuation b where b.day>=date_sub(a.day,INTERVAL 365 DAY) and b.code=a.code )  where  a.id =  %s"

    elif tName=='stk_shareholder_floating_top10':

        sql= "INSERT INTO stk_shareholder_floating_top10(id,company_id,company_name,code,end_date,pub_date,change_reason_id,change_reason,shareholder_rank, shareholder_id,shareholder_name,shareholder_name_en,shareholder_class_id,shareholder_class, share_number,share_ratio,sharesnature_id,sharesnature) \
            VALUES (%s,%s,%s,%s,%s, %s,%s,%s,%s,%s, %s,%s,%s,%s,%s,  %s,%s,%s )"

    return sql

def saveBatch(vals,tName):
    # print('======saveBatch======='+tName)
    conn,cs=getConn()
    try:
        sql = getSQL(tName)

        # sql = "INSERT INTO dayline(date,code,open,close,low,high,volume,money,factor, high_limit,low_limit,avg,pre_close,paused, m5,m10,m20,m60,m250, dif,dea,macd) \
        #     VALUES (%s,%s,%s,%s,%s, %s,%s,%s,%s,%s, %s,%s,%s,%s,%s,  %s,%s,%s,%s,%s, %s,%s )"
        
        cs.executemany(sql, vals)
        #print('sql==',sql)
    except Exception as e:
        print("出现如下异常%s"%e)
        return
    conn.commit()
    closeConn(conn,cs)

def updateBatchM5(vals):
    print('======updateBatchM5=======')
    conn,cs=getConn()
    try:
        sql = 'UPDATE dayline SET m5 = (%s),m20=(%s) WHERE id = (%s) '
        cs.executemany(sql, vals)
        #print('sql==',sql)
    except Exception as e:
        print("updateBatchM5 出现如下异常%s"%e)
        return
    conn.commit()
    closeConn(conn,cs)

def closeConn(cursor,conn):
    cursor.close()
    conn.close()

def getConn():
    user = 'root' # = input('username:')
    pwd ='123456' # = input('password:')

    conn = pymysql.connect(host='localhost',user='root',password='123456',database='stock')
    cursor = conn.cursor()   
    return conn, cursor

def savePrice(rq):
    # codelist=['000979.XSHE'] #600718.XSHG','002315.XSHE']
    codelist = list(get_all_securities(['stock']).index)
    #print(stocks)

    # MA5 = MA(codelist, check_date=rq, timeperiod=5)
    # MA10 = MA(codelist, check_date=rq, timeperiod=10)
    # MA20 = MA(codelist, check_date=rq, timeperiod=20)
    # MA60 = MA(codelist, check_date=rq, timeperiod=60)
    # MA250 = MA(codelist, check_date=rq, timeperiod=250)
    # macd_dif, macd_dea, macd_macd = MACD(codelist,check_date=rq, SHORT = 12, LONG = 26, MID = 9)

    cols=['open','close','low','high','volume','money','factor','high_limit','low_limit','avg','pre_close','paused']
    df = get_price(codelist, start_date=rq, end_date=rq, frequency='daily', panel=False ,fields=cols)

    #date,code,open,close,low, high,volume,money,factor,high_limit,  low_limit,avg,pre_close,paused,m5,  m10,m20,m60,m250, dif,dea,macd)

    datalist=[]

    for index, row in df.iterrows():
        
        if str(row['open'])!='nan' and is_number(row['open']):
            rowlist = [row['time'],row['code'], row['open'],row['close'],row['low'], row['high'],row['volume'],row['money'],row['factor'],row['high_limit'],  row['low_limit'],row['avg'],row['pre_close'],row['paused'], 
            0,0,0,0,0] # ,0,0,0]
            # round(MA5[row['code']],2) , 
            # round(MA10[row['code']],2) ,  
            # round(MA20[row['code']],2),  
            # round(MA60[row['code']],2) , 
            # round(MA250[row['code']],2) ,
            # macd_dif[row['code']],
            # macd_dea[row['code']],
            # macd_macd[row['code']]]

            datalist.append(rowlist)
    
    #print(row['time'],row['code'],row['open'],row['close'],row['low'],   row['high'],row['volume'],row['money'],row['factor'],row['high_limit'],  row['low_limit'],row['avg'],row['pre_close'],row['paused'],MA5[row['code']],  MA10[row['code']],MA20[row['code']],MA60[row['code']],MA250[row['code']] ,macd_dif[row['code']],macd_dea[row['code']],macd_macd[row['code']] ) 

    if len(datalist)>0:
        saveBatch(datalist,'dayline')
    else:
        print("empty datalist")

def updateMA(rq):
    codes = getCodes(rq)
    updateList=[]
    for code in codes:

        id,ma5=genMA(code[0],rq,5)
        id,ma20=genMA(code[0],rq,20)
        updateList.append([ma5,ma20,id])

    updateBatchM5(updateList)

def floatingTop10():
    codes = getCodes(rq)
    pandas.set_option('display.max_columns', None)
    pandas.set_option('display.max_rows',None)
    for code in codes:

        q=query(finance.STK_SHAREHOLDER_FLOATING_TOP10).filter(finance.STK_SHAREHOLDER_FLOATING_TOP10.code==code[0],finance.STK_SHAREHOLDER_FLOATING_TOP10.pub_date>'2021-06-01').limit(10)
        df=finance.run_query(q)
        
        print(df)

rq='2021-01-01'

if len(sys.argv)!=4:
    print ('err args cnt')

rq=sys.argv[1]
userName=sys.argv[2]
password=sys.argv[3]

# logout()
auth(userName, password)

# 查询是否连接成功
is_auth = is_auth()
print(is_auth)

ss=get_query_count()
print(ss,rq)
start=time.time()

savePrice(rq) #每天执行
updateMA(rq) #每天执行

# valuationDay(rq,'valuation') #每天执行 暂时只用来统计基金持仓  平时就不用跑了 
# replaceFloatingTop10(rq)

# tdays = get_trade_days(start_date="2021-01-01",end_date="2021-12-17")

# print(len(tdays))

# for  d in tdays:
#     print (d)
#     savePrice(d)
#     updateMA(d)





# replaceIndustry(rq) #上次执行日期 20210624

# tdays = get_trade_days(start_date="2021-03-20",end_date="2021-04-16")
# for  d in tdays:
#     print (d)
#     start1=time.time()
#     updateYeardayCap(d)
#     end1=time.time()
#     print(str(d)+'-==== Running time: %s Seconds' %(end1-start1))

# updateYeardayCap(rq) #每天执行 暂停

end=time.time()

print('Running time: %s Seconds'%(end-start))
#exit(0)

ss=get_query_count()
print(ss,rq)

logout()

# news = finance.run_query(query(finance.CCTV_NEWS).filter(finance.CCTV_NEWS.day=='2021-03-01').limit(5))

# print(news)
