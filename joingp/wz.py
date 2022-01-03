#-*- coding: UTF-8 -*-
import sys
import pymysql
from jqdatasdk import *
import datetime
import time
from jqdatasdk.technical_analysis import *
from sqlalchemy.sql import func
#from jqlib.technical_analysis 6214 6239 5000 1094 741 import * 

def baseInfo(stock_list,rq):
    df = get_fundamentals(query( # func.count('*'), 
            valuation.code, 
            valuation.market_cap, 
            valuation.circulating_market_cap, 
            valuation.pe_ratio,
            valuation.turnover_ratio,
            valuation.capitalization,  # 总股本
            valuation.circulating_cap, # 流通股本 
            indicator.adjusted_profit,
            indicator.pubDate,
            indicator.statDate,
            
        ).filter(
            # valuation.market_cap >= 10,
            # valuation.market_cap <= 500, #市值

            valuation.code.in_(stock_list)
            # indicator.eps >=0.4,
            # valuation.market_cap * 0.5 <= valuation.circulating_market_cap
        ).order_by(
            # 按市值降序排列
            # valuation.market_cap.desc()
            valuation.market_cap.desc()
        ).limit(
            # 最多返回100个
            5000
        ), date=rq)

    # print(df)
    print(
            'code',
            'name',
            '市值',#row['market_cap'], 
            '流通市值',#row['circulating_market_cap'],
            '市盈率',#row['pe_ratio'],
            '换手率',
            '每股收益','公积金',#basic_eps,
            '营收',#operating_revenue,
            '财务费用',#financial_expense, #财务费用
            #pubincome,
            #statincome,
            '扣非净利润',#row['adjusted_profit'],
            #row['pubDate'],    row['statDate'],
            '应收','预付',##account_receivable, advance_payment,
            #bspub_date,bsreport_date,
            '经营活动现金流量','投资活动现金流量','筹资活动现金流量',#net_operate_cash_flow,net_invest_cash_flow,net_finance_cash_flow,
            #cfpub_date,cfreport_date
            '行业'#inds[row['code']]['zjw']['industry_name']
            )
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
            finance.STK_INCOME_STATEMENT.financial_expense,
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
        financial_expense = -1
        # print(len(indf))
        pubincome = '-1'
        statincome = '-1'
        if(len(dfincome)>0):
            basic_eps = dfincome.iloc[0]['basic_eps']
            operating_revenue = dfincome.iloc[0]['operating_revenue']
            financial_expense = dfincome.iloc[0]['financial_expense']
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
            finance.STK_BALANCE_SHEET.report_type,
            finance.STK_BALANCE_SHEET.report_date,
            finance.STK_BALANCE_SHEET.capital_reserve_fund #资本公积金
        ).filter(finance.STK_BALANCE_SHEET.code==row['code'],finance.STK_BALANCE_SHEET.pub_date>'2020-01-01',finance.STK_BALANCE_SHEET.report_type==0
        ).order_by(
            finance.STK_BALANCE_SHEET.pub_date.desc()
        ).limit(1)
        dfbalance=finance.run_query(qbalance)
        account_receivable = -1 #应收帐款
        advance_payment = -1 #预付
        bspub_date='-1'
        bsreport_date = '-1'
        if(len(dfbalance)>0):
            account_receivable = dfbalance.iloc[0]['account_receivable']
            advance_payment = dfbalance.iloc[0]['advance_payment']
            bspub_date = dfbalance.iloc[0]['pub_date']
            bsreport_date = dfbalance.iloc[0]['report_date']
            capital_reserve_fund = dfbalance.iloc[0]['capital_reserve_fund']

        # STK_CASHFLOW_STATEMENT
        qcash=query(
            finance.STK_CASHFLOW_STATEMENT.code,
            finance.STK_CASHFLOW_STATEMENT.pub_date,
            finance.STK_CASHFLOW_STATEMENT.start_date,
            finance.STK_CASHFLOW_STATEMENT.end_date,
            finance.STK_CASHFLOW_STATEMENT.net_operate_cash_flow, # 经营活动现金流量净额
            finance.STK_CASHFLOW_STATEMENT.net_invest_cash_flow, #投资活动现金流量净额
            finance.STK_CASHFLOW_STATEMENT.net_finance_cash_flow, #筹资活动现金流量净额
            finance.STK_CASHFLOW_STATEMENT.report_date,
            finance.STK_CASHFLOW_STATEMENT.report_type
        ).filter(finance.STK_CASHFLOW_STATEMENT.code==row['code'],finance.STK_CASHFLOW_STATEMENT.pub_date>'2020-01-01',finance.STK_CASHFLOW_STATEMENT.report_type==0
        ).order_by(
            finance.STK_CASHFLOW_STATEMENT.pub_date.desc()
        ).limit(1)
        dfcash=finance.run_query(qcash)
        net_operate_cash_flow = -1 #应收帐款
        net_invest_cash_flow = -1
        net_finance_cash_flow = -1
        cfreport_date='-1'
        cfpub_date = '-1'
        if(len(dfcash)>0):
            net_operate_cash_flow = dfcash.iloc[0]['net_operate_cash_flow']
            net_invest_cash_flow = dfcash.iloc[0]['net_invest_cash_flow']
            net_finance_cash_flow = dfcash.iloc[0]['net_finance_cash_flow']
            cfreport_date=dfcash.iloc[0]['report_date']
            cfpub_date=dfcash.iloc[0]['pub_date']

        print(
            row['code'][0:6],
            display_name,
            row['market_cap'], 
            row['circulating_market_cap'],
            row['pe_ratio'],
            row['turnover_ratio'],
            basic_eps,capital_reserve_fund,
            operating_revenue,
            financial_expense, #财务费用
            #pubincome,
            #statincome,
            row['adjusted_profit'],
            #row['pubDate'],    row['statDate'],
            account_receivable, advance_payment,
            #bspub_date,bsreport_date,
            net_operate_cash_flow,net_invest_cash_flow,net_finance_cash_flow,
            #cfpub_date,cfreport_date
            inds[row['code']]['zjw']['industry_name']
            )

auth('15640293905', '1q2w3e4R?')
# 查询是否连接成功
is_auth = is_auth()
print(is_auth)

ss=get_query_count()
print(ss)

#stock_list=normalize_code(['000503','000566','000615','000620','000676','000709','000723','000778','000790','000825','000862','000886','000915','000959','001896','002053','002110','002161','002174','002285','002320','002329','002348','002378','002383','002435','002490','002512','002575','002596','002614','002639','002693','002730','002955','002962','003004','003008','003013','003025','003033','300342','300465','300622','300671','300687','300745','300782','300875','300881','300899','600022','600052','600079','600117','600200','600231','600403','600569','600581','600733','600818','600929','601005','601099','601113','601258','601969','603002','603377','603392','603398','603439','603555','603598','603893','603919','603986','603991','605018','605099','605118','605198','605333','688016','688068','688356','688357'])

codelist = ['310001','310002','310003','310004']

rq=time.strftime('%Y-%m-%d')

if len(sys.argv)==2:
    rq=sys.argv[1]
    #rq=time.strftime('%Y-%m-%d')
    
print(rq)

# df=finance.run_query(query(finance.STK_HK_HOLD_INFO).filter(finance.STK_HK_HOLD_INFO.link_id==310001).order_by(finance.STK_HK_HOLD_INFO.day.desc()).limit(2))
# print(df[:1])

for code in codelist:
    df=finance.run_query(query(
        finance.STK_HK_HOLD_INFO).
        filter(
            finance.STK_HK_HOLD_INFO.link_id==code,
            finance.STK_HK_HOLD_INFO.day==rq
        ).order_by(
            finance.STK_HK_HOLD_INFO.share_ratio.desc()).limit(20)
        )

    print(df)



logout()
# news = finance.run_query(query(finance.CCTV_NEWS).filter(finance.CCTV_NEWS.day=='2021-03-01').limit(5))