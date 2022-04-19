
mysql -uroot -p123456 stock -e"select id,day,code,market_cap ,(select avg(market_cap) from valuation b where b.day>=date_sub(a.day,INTERVAL 365 DAY) and b.code=a.code ) yearday_avg_market_cap from valuation a where a.day>='2021-03-01' ;" > yearday_avg_market_cap.txt


update valuation a set yearday_avg_market_cap=(select avg(market_cap) from valuation b where b.day>=date_sub(a.day,INTERVAL 365 DAY) and b.code=a.code ) where a.day>='2021-03-01'

+------------+-------------+------------+--------------------+
| day        | code        | market_cap | avg_market_cap     |
+------------+-------------+------------+--------------------+
| 2021-04-16 | 600718.XSHG |   116.1616 |           136.8145 |
| 2021-04-15 | 600718.XSHG |   113.5526 | 136.83871959183674 |
| 2021-04-14 | 600718.XSHG |   113.6769 | 136.86425772357728 |
| 2021-04-13 | 600718.XSHG |   113.3042 |   136.880032388664 |
| 2021-04-12 | 600718.XSHG |   113.9254 |   136.880032388664 |
| 2021-04-09 | 600718.XSHG |   114.5465 | 136.96309156626504 |
| 2021-04-08 | 600718.XSHG |   114.5465 |        136.9991532 |
| 2021-04-07 | 600718.XSHG |   115.4162 | 137.03294780876493 |
| 2021-04-06 | 600718.XSHG |   115.7889 | 137.03294780876493 |
| 2021-04-02 | 600718.XSHG |    114.795 |  137.0805841897233 |
| 2021-04-01 | 600718.XSHG |   114.5465 | 137.09409409448818 |
| 2021-03-31 | 600718.XSHG |   114.9193 | 137.11821647058824 |
| 2021-03-30 | 600718.XSHG |    114.795 |    137.15137109375 |
| 2021-03-29 | 600718.XSHG |   116.7828 |    137.15137109375 |
| 2021-03-26 | 600718.XSHG |   118.0252 | 137.26843410852715 |
| 2021-03-25 | 600718.XSHG |   117.5282 | 137.32268996138998 |
| 2021-03-24 | 600718.XSHG |   116.5343 | 137.37127230769232 |
| 2021-03-23 | 600718.XSHG |   116.7828 | 137.40520229885058 |
| 2021-03-22 | 600718.XSHG |    115.292 | 137.40520229885058 |
| 2021-03-19 | 600718.XSHG |   114.2981 |  137.5025209125475 |
| 2021-03-18 | 600718.XSHG |   115.1677 | 137.54168598484847 |
| 2021-03-17 | 600718.XSHG |   116.4101 | 137.59133811320754 |
| 2021-03-16 | 600718.XSHG |   116.4101 |  137.6368804511278 |
| 2021-03-15 | 600718.XSHG |   115.9131 |  137.6368804511278 |
| 2021-03-12 | 600718.XSHG |   116.4101 | 137.80992499999996 |
| 2021-03-11 | 600718.XSHG |   117.7767 | 137.91187769516728 |
| 2021-03-10 | 600718.XSHG |   116.2859 |  138.0370025925926 |
| 2021-03-09 | 600718.XSHG |   118.6464 | 138.15203505535052 |
| 2021-03-08 | 600718.XSHG |   120.5099 | 138.15203505535052 |
| 2021-03-05 | 600718.XSHG |   121.8765 |  138.4419179487179 |
| 2021-03-04 | 600718.XSHG |   120.6342 |  138.5877664233576 |
| 2021-03-03 | 600718.XSHG |   121.6281 | 138.75288363636358 |
| 2021-03-02 | 600718.XSHG |   121.5038 | 138.90240036231876 |
+------------+-------------+------------+--------------------+

select * from (
	select left(code,6)code, 
	sum(if(date='2021-03-09', dif ,0)) as dif9 ,
	sum(if(date='2021-03-10', dif ,0)) as dif10,
	sum(if(date='2021-03-11', dif ,0)) as dif11 ,
	sum(if(date='2021-03-12', dif ,0)) as dif12 ,
	sum(if(date='2021-03-19', dif ,0)) as dif19,
	sum(if(date='2021-03-22', dif ,0)) as dif22,
	sum(if(date='2021-03-23', dif ,0)) as dif23,
	sum(if(date='2021-03-24', dif ,0)) as dif24,
	sum(if(date='2021-03-25', dif ,0)) as dif25,
	sum(if(date='2021-03-26', dif ,0)) as dif26,
	
	sum(if(date='2021-03-09', dea ,0)) as dea9,
	sum(if(date='2021-03-10', dea ,0)) as dea10,
	sum(if(date='2021-03-11', dea ,0)) as dea11,
	sum(if(date='2021-03-12', dea ,0)) as dea12,	
	sum(if(date='2021-03-19', dea ,0)) as dea19,
	sum(if(date='2021-03-22', dea ,0)) as dea22,
	sum(if(date='2021-03-23', dea ,0)) as dea23,
	sum(if(date='2021-03-24', dea ,0)) as dea24,
	sum(if(date='2021-03-25', dea ,0)) as dea25,
	sum(if(date='2021-03-26', dea ,0)) as dea26,
	
	sum(if(date='2021-03-26', close ,0)) as close, 
	sum(if(date='2021-03-26', m60 ,0)) as m60, 
	sum(if(date='2021-03-26', m20 ,0)) as m20
	
	 from dayline where 
	 date >= '2021-03-09'  group by code ) as x 
	 where  dif24 >0 -- dea12>0 and dea15>0 and dea16 >0
   and  dea26-dea24>=dea24-dea22 
	 -- and dea25-dea24>=dea24-dea23 
	 and dea23-dea22>=dea22-dea19 and dea12-dea11>=dea11-dea10 and dea10-dea9>=dea9  
	 and m20<close and m60<close 
	 order by dif9+dea9,dif10+dea10,dif11+dea11,dif12+dea12 limit 20