#-*- coding: UTF-8 -*-
import sys
import pymysql
import pandas_datareader.data as web
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

def getCoJP(rq):
    conn,cs=getConn()
    try:
        sql = "SELECT distinct code FROM co_jp WHERE market != 'ETF・ETN' and code>='2761' and date = %s " % (rq)
        #print(sql)
        cs.execute(sql)
        results = cs.fetchall()
        return results
    except Exception as e:
        print("出现如下异常%s"%e)
        return
    closeConn(conn,cs)

def getCodes(rq):
    conn,cs=getConn()
    try:
        sql = "SELECT distinct code FROM dayline_jp WHERE date = '%s' " % (rq)
        #print(sql)
        cs.execute(sql)
        results = cs.fetchall()
        return results
    except Exception as e:
        print("出现如下异常%s"%e)
    closeConn(conn,cs)

def genPerClose(code,rq):
    #print(rq+'======genM5======='+code)
    conn,cs=getConn()
    try:
        sql = "SELECT id,close FROM dayline_jp WHERE code = '%s' and date<='%s' order by date desc limit 2" % (code,rq)

        cs.execute(sql)

        results = cs.fetchall()
        preClose = 0 
        id = -1
        if len(results)==2:
          for row in results:
            if id==-1:
                id=row[0]
            else:
                preClose=row[1]
        else:
          preClose = -1
        return id,preClose
    except Exception as e:
        print("genPerClose 出现如下异常%s"%e)
        return id,preClose
    finally:
        closeConn(conn,cs)

def genMA(code,rq,n):
    #print(rq+'======genM5======='+code)
    conn,cs=getConn()
    try:
        sql = "SELECT id,close FROM dayline_jp WHERE code = '%s' and date<='%s' order by date desc limit %d" % (code,rq ,n)

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
    if tName=='dayline_jp':
        sql = "INSERT INTO dayline_jp(date,code,open,close,low,high,volume,money,factor, high_limit,low_limit,avg,pre_close,paused, m5,m10,m20,m60,m30) \
            VALUES (%s,%s,%s,%s,%s, %s,%s,%s,%s,%s, %s,%s,%s,%s,%s,  %s,%s,%s,%s )"
    return sql

def saveBatch(vals,tName):
    # print('======saveBatch======='+tName)
    conn,cs=getConn()
    try:
        sql = getSQL(tName)
        cs.executemany(sql, vals)
        #print('sql==',sql)
    except Exception as e:
        print("出现如下异常%s"%e)
        return
    conn.commit()
    closeConn(conn,cs)

def updateBatchM510203060(vals):
    print('======updateBatchM5203060pre_close=======')
    conn,cs=getConn()
    try:
        sql = 'UPDATE dayline_jp SET pre_close=(%s), m5 = (%s), m10=(%s), m20=(%s), m30 = (%s), m60=(%s) WHERE id = (%s) '
        cs.executemany(sql, vals)
        print('sql==',sql)
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

    conn = pymysql.connect(host='localhost',user='root',password='123456',database='stock',port=3307 )
    cursor = conn.cursor()   
    return conn, cursor

def savePrice(rq):

    coJPs = getCoJP(20220428)

    print(coJPs)

    for coJP in coJPs:

        print(coJP[0])
        # df = web.DataReader(coJP[0]+'.JP', 'stooq', start=rq, end=rq)
        df = web.DataReader(coJP[0]+'.JP', 'stooq', start='2022-01-01', end='2022-05-20')

        datalist=[]
#                 0                                1              2           3          4             5
# Pandas(Index=Timestamp('2022-05-19 00:00:00'), Open=3220.0, High=3255.0, Low=3200.0, Close=3235.0, Volume=17500)

        for row in df.itertuples():
            print(row)
            # print(row[0], row[1])

            if str(row[1])!='nan' and is_number(row[1]):

                rowlist = [row[0],coJP[0], row[1],row[4],row[3], row[2],row[5],
                0,0,0,0,0,0,
                '0', #paused 
                0,0,0,0,0]

                datalist.append(rowlist)

        if len(datalist)>0:
            saveBatch(datalist,'dayline_jp')
        else:
            print("empty datalist")

def updateMA(rq):
    codes = getCodes(rq)
    updateList=[]
    for code in codes:

        id,preClose= genPerClose(code[0],rq)
        id,ma5 =genMA(code[0],rq,5)
        id,ma10=genMA(code[0],rq,10)
        id,ma20=genMA(code[0],rq,20)
        id,ma30=genMA(code[0],rq,30)
        id,ma60=genMA(code[0],rq,60)
        updateList.append([preClose,ma5,ma10,ma20,ma30,ma60,id])

    updateBatchM510203060(updateList)

rq='2021-01-01'

if len(sys.argv)!=2:
    print ('err args cnt')

print(sys.argv)
rq=sys.argv[1]

start=time.time()

# savePrice(rq) #每天执行
updateMA(rq) #每天执行

end=time.time()

print('Running time: %s Seconds'%(end-start))



