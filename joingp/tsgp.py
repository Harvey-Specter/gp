#-*- coding: UTF-8 -*-
import sys
import pymysql
import tushare 
import time 
from sqlalchemy.sql import func
import pandas 

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
        return id,ma
    finally:
        closeConn(conn,cs)

def getSQL(tName):
    sql = ''
    if tName=='dayline':
        sql = "INSERT INTO dayline(date,code,open,close,low,high,volume,money,factor, high_limit,low_limit,avg,pre_close,paused, m5,m10,m20,m60,m30) \
            VALUES (%s,%s,%s,%s,%s, %s,%s,%s,%s,%s, %s,%s,%s,%s,%s,  %s,%s,%s,%s )"
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

def updateBatchM510203060(vals):
    print('======updateBatchM5203060=======')
    conn,cs=getConn()
    try:
        sql = 'UPDATE dayline SET m5 = (%s), m10=(%s), m20=(%s), m30 = (%s), m60=(%s) WHERE id = (%s) '
        cs.executemany(sql, vals)
        #print('sql==',sql)
    except Exception as e:
        print("updateBatchM510203060 出现如下异常%s"%e)
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

    rq1=rq.replace('-' , '')
    pro = tushare.pro_api('f84256115f164f99aba02574320a9279fc584d9b13674bd79ed0d9cc')
    df = pro.daily(trade_date=rq1)
    
    datalist=[]

    for index, row in df.iterrows():
        
        if str(row['open'])!='nan' and is_number(row['open']):
            rowlist = [row['trade_date'],row['ts_code'].replace('SH','XSHG').replace('SZ','XSHE'), row['open'],row['close'],row['low'], row['high'],row['vol'],row['amount'],0,0,0,0,row['pre_close'],'0', 
            0,0,0,0,0]

            datalist.append(rowlist)

    if len(datalist)>0:
        saveBatch(datalist,'dayline')
    else:
        print("empty datalist")

def updateMA(rq):
    codes = getCodes(rq)
    updateList=[]
    for code in codes:

        id,ma5 =genMA(code[0],rq,5)
        id,ma10=genMA(code[0],rq,10)
        id,ma20=genMA(code[0],rq,20)
        id,ma30=genMA(code[0],rq,30)
        id,ma60=genMA(code[0],rq,60)
        updateList.append([ma5,ma10,ma20,ma30,ma60,id])

    updateBatchM510203060(updateList)

rq='2021-01-01'

if len(sys.argv)!=2:
    print ('err args cnt')

print(sys.argv)
rq=sys.argv[1]

start=time.time()

savePrice(rq) #每天执行
updateMA(rq) #每天执行

end=time.time()

print('Running time: %s Seconds'%(end-start))



