
VirtualBoxVM --startvm 18b71cd9-52e6-43b7-8c3d-0a4ba40ed26b

STK_SHAREHOLDER_FLOATING_TOP10

CREATE TABLE `stk_shareholder_floating_top10` (
	`id` int(10) unsigned NOT NULL AUTO_INCREMENT,
	`date` date NOT NULL,
	`code` varchar(16) COLLATE utf8mb4_unicode_ci NOT NULL,
	`open` double unsigned NOT NULL,
	`close` double unsigned NOT NULL,
	`low` double unsigned NOT NULL,
	`high` double unsigned NOT NULL,
	`volume` double unsigned NOT NULL,
	`money` double NOT NULL,
	`factor` double NOT NULL,
	`high_limit` double NOT NULL,
	`low_limit` double NOT NULL,
	`avg` double NOT NULL,
	`pre_close` double NOT NULL,
	`paused` varchar(25) COLLATE utf8mb4_unicode_ci NOT NULL,
	`m5` double DEFAULT NULL,
	`m10` double unsigned NOT NULL DEFAULT 0,
	`m20` double unsigned NOT NULL DEFAULT 0,
	`m60` double unsigned NOT NULL DEFAULT 0,
	`m250` double unsigned NOT NULL DEFAULT 0, 高级趋势技术分析
	`dif` double DEFAULT 0,
	`dea` double DEFAULT 0,
	`macd` double DEFAULT 0,
	`r1` varchar(25) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
	`create_time` timestamp NULL DEFAULT current_timestamp(),
	PRIMARY KEY (`id`),
	UNIQUE KEY `date` (`date`,`code`)
  ) ENGINE=InnoDB AUTO_INCREMENT=751761 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci

  drop table stk_shareholder_floating_top10;
 CREATE TABLE `stk_shareholder_floating_top10` (
  `id` int(11) NOT NULL,
  `company_id` int(11) NOT NULL,
  `company_name` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL,
  `code` varchar(16) COLLATE utf8mb4_unicode_ci NOT NULL,
  `end_date` date NOT NULL,
  `pub_date` date NOT NULL,
  `change_reason_id` int(11) NOT NULL,
  `change_reason` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL,
  `shareholder_rank` int(11) NOT NULL,
  `shareholder_id` int(11) NOT NULL,
  `shareholder_name` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL,
  `shareholder_name_en` varchar(1000) COLLATE utf8mb4_unicode_ci ,
  `shareholder_class_id` int(11) NOT NULL,
  `shareholder_class` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL,
  `share_number` bigint(20) DEFAULT NULL,
  `share_ratio` double unsigned NOT NULL DEFAULT 0,
  `sharesnature_id` int(11) NOT NULL,
  `sharesnature` varchar(200) COLLATE utf8mb4_unicode_ci NOT NULL,
  `create_time` timestamp NULL DEFAULT current_timestamp(),
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_top10` (`code`,`end_date`,`shareholder_rank`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci



select shareholder_class_id,shareholder_name,shareholder_class,round(sum(share_number*close)) from dayline a ,stk_shareholder_floating_top10 b 
  where a.code = b.code and a.date='2021-09-30' and shareholder_class_id not in ('307099','307001','307024','307023','307006','307011','307012','307019') and shareholder_name like '%易方达%'
  group by shareholder_id,shareholder_name,shareholder_class_id,shareholder_class 

  select sum(share_ratio) ,code from stk_shareholder_floating_top10 where code not like '300%' and code not like '68%' and  sharesnature='流通A股' and shareholder_class not in ('自然人','其他机构','上市公司','保险公司','地方国资委','保险公司和上市公司','国有资产经营公司','资产管理公司资产管理计划','风险投资','银行') group by code order by 1 desc limit 10 ;

  select sum(share_ratio) ,code from stk_shareholder_floating_top10 where code not like '300%' and code not like '68%' and  sharesnature='流通A股' and shareholder_class like '%基金%' group by code order by 1 desc limit 10 ;

select code,end_date,pub_date,shareholder_rank,shareholder_name,shareholder_class,share_number,share_ratio,sharesnature from stk_shareholder_floating_top10 where code like '600841.XSHG' ;

select o.code,
round(al/(p.circulating_market_cap/market_cap),2) alc,
round(jg/(p.circulating_market_cap/market_cap),2) jgc,
round(jj/(p.circulating_market_cap/market_cap),2) jjc,
al,jg,jj , end_date from (
  select a.code, round(sum(share_ratio),2) as al ,
  round((select sum(share_ratio) from stk_shareholder_floating_top10 x where x.code = a.code and  sharesnature='流通A股' and shareholder_class not in ('自然人','其他机构','上市公司','保险公司','地方国资委','保险公司和上市公司','国有资产经营公司','资产管理公司资产管理计划','风险投资') ),2) jg ,
  round(ifnull( 
	  (select sum(share_ratio) from stk_shareholder_floating_top10 x where x.code = a.code and sharesnature='流通A股' and ( shareholder_class like '%基金%' or shareholder_class like '%QFII%') ) ,0),2) jj ,end_date
   from stk_shareholder_floating_top10 a  where 
   code not like '68%' and code not like '300%' and 
   sharesnature='流通A股' group by a.code ,a.end_date)
as o ,valuation p where o.code=p.code and p.day='2021-11-01'
   order by 4 desc ,2 desc ,3 desc limit 400 ;

drop table valuation;
   CREATE TABLE `valuation` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `day` date NOT NULL,
  `code` varchar(16) COLLATE utf8mb4_unicode_ci NOT NULL,
  `capitalization` double unsigned NOT NULL,
  `circulating_cap` double unsigned NOT NULL,
  `market_cap` double unsigned NOT NULL,
  `circulating_market_cap` double unsigned NOT NULL,
  `turnover_ratio` double unsigned NOT NULL,
  `pe_ratio` double NOT NULL,
  `pe_ratio_lyr` double NOT NULL,
  `pb_ratio` double NOT NULL,
  `ps_ratio` double NOT NULL,
  `pcf_ratio` double NOT NULL,
  `create_time` timestamp NULL DEFAULT current_timestamp(),
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_code_day` (`day`,`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci