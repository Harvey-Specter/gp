#-*- coding: UTF-8 -*-
import sys
import pymysql
from jqdatasdk import *
import datetime
from jqdatasdk.technical_analysis import *
from sqlalchemy.sql import func
#from jqlib.technical_analysis 6214 6239 5000 1094 741 import * 

auth('15640293905', '1q2w3e4R?')
# 查询是否连接成功
is_auth = is_auth()
print(is_auth)

ss=get_query_count()
print(ss)

codelist=['600718.XSHG','002315.XSHE']

datalist=[]

df = get_fundamentals(query( # func.count('*'), 
        valuation.code, valuation.market_cap, valuation.circulating_market_cap, indicator.adjusted_profit,indicator.pubDate,indicator.statDate
    ).filter(
        # valuation.market_cap >= 50,
        valuation.market_cap <= 400, #市值
		# indicator.eps >=0.4,
		# valuation.market_cap * 0.5 <= valuation.circulating_market_cap
    ).order_by(
        # 按市值降序排列
        # valuation.market_cap.desc()
		valuation.market_cap.desc()
	).limit(
        # 最多返回100个
        2
    ), date='2021-03-26')

# print(df)
for index, row in df.iterrows():
    
    #if str(row['open'])!='nan' and is_number(row['open']):
	display_name=get_security_info(row['code']).display_name
	if 'st' in display_name or 'ST' in display_name:
		continue;

	# 600291.XSHG
	# income statement 
	qincome=query(# finance.STK_INCOME_STATEMENT.company_name,
        finance.STK_INCOME_STATEMENT.code,
        finance.STK_INCOME_STATEMENT.pub_date,
        finance.STK_INCOME_STATEMENT.start_date,
        finance.STK_INCOME_STATEMENT.end_date,
        finance.STK_INCOME_STATEMENT.operating_revenue,
		finance.STK_INCOME_STATEMENT.basic_eps,
		finance.STK_INCOME_STATEMENT.diluted_eps,
		finance.STK_INCOME_STATEMENT.report_type,
		finance.STK_INCOME_STATEMENT.report_date,
	finance.STK_INCOME_STATEMENT.np_parent_company_owners).filter(finance.STK_INCOME_STATEMENT.code==row['code'],finance.STK_INCOME_STATEMENT.pub_date>'2020-01-01',finance.STK_INCOME_STATEMENT.report_type==0
	).order_by(
        # 按市值降序排列
        # valuation.market_cap.desc()
		finance.STK_INCOME_STATEMENT.pub_date.desc()
	).limit(1)
	dfincome=finance.run_query(qincome)
	basic_eps=-1 #基础每股收益
	operating_revenue = -1
	# print(len(indf))
	pubincome = '-1'
	statincome = '-1'
	if(len(dfincome)>0):
		basic_eps = dfincome.iloc[0]['basic_eps']
		operating_revenue= dfincome.iloc[0]['operating_revenue']
		pubincome=dfincome.iloc[0]['pub_date']
		statincome=dfincome.iloc[0]['report_date']
		# print(indf.iloc[0]['basic_eps'])

	# STK_BALANCE_SHEET
	qbalance=query(
        finance.STK_BALANCE_SHEET.code,
        finance.STK_BALANCE_SHEET.pub_date,
        finance.STK_BALANCE_SHEET.start_date,
        finance.STK_BALANCE_SHEET.end_date,
        finance.STK_BALANCE_SHEET.advance_payment, # 预付
		finance.STK_BALANCE_SHEET.account_receivable, #应收
		finance.STK_BALANCE_SHEET.report_type
	).filter(finance.STK_BALANCE_SHEET.code==row['code'],finance.STK_BALANCE_SHEET.pub_date>'2020-01-01',finance.STK_BALANCE_SHEET.report_type==0
	).order_by(
		finance.STK_BALANCE_SHEET.pub_date.desc()
	).limit(1)
	dfbalance=finance.run_query(qbalance)
	account_receivable = -1 #应收帐款
	advance_payment = -1 #预付
	bsstart_date='-1'
	if(len(dfbalance)>0):
		account_receivable = dfbalance.iloc[0]['account_receivable']
		advance_payment = dfbalance.iloc[0]['advance_payment']
		bsstart_date = dfbalance.iloc[0]['start_date']

	# STK_CASHFLOW_STATEMENT
	qcash=query(
        finance.STK_CASHFLOW_STATEMENT.code,
        finance.STK_CASHFLOW_STATEMENT.pub_date,
        finance.STK_CASHFLOW_STATEMENT.start_date,
        finance.STK_CASHFLOW_STATEMENT.end_date,
        finance.STK_CASHFLOW_STATEMENT.net_operate_cash_flow, # 经营活动现金流量净额
		finance.STK_CASHFLOW_STATEMENT.net_invest_cash_flow, #投资活动现金流量净额
		finance.STK_CASHFLOW_STATEMENT.net_finance_cash_flow, #筹资活动现金流量净额

		finance.STK_CASHFLOW_STATEMENT.report_type
	).filter(finance.STK_CASHFLOW_STATEMENT.code==row['code'],finance.STK_CASHFLOW_STATEMENT.pub_date>'2020-01-01',finance.STK_CASHFLOW_STATEMENT.report_type==0
	).order_by(
		finance.STK_CASHFLOW_STATEMENT.pub_date.desc()
	).limit(1)
	dfcash=finance.run_query(qcash)
	net_operate_cash_flow = -1 #应收帐款
	net_invest_cash_flow = -1
	net_finance_cash_flow = -1
	if(len(dfcash)>0):
		net_operate_cash_flow = dfcash.iloc[0]['net_operate_cash_flow']
		net_invest_cash_flow = dfcash.iloc[0]['net_invest_cash_flow']
		net_finance_cash_flow = dfcash.iloc[0]['net_finance_cash_flow']

	print(row['code'][0:6],display_name,row['market_cap'], row['circulating_market_cap'],basic_eps,operating_revenue,row['adjusted_profit'],row['pubDate'],row['statDate'],account_receivable,advance_payment,bsstart_date,net_operate_cash_flow,net_invest_cash_flow,net_finance_cash_flow)

logout()
# news = finance.run_query(query(finance.CCTV_NEWS).filter(finance.CCTV_NEWS.day=='2021-03-01').limit(5))
