http://quotes.money.163.com/service/chddata.html?code=0601512&start=20210316&end=20210316

wget -O 1000592.csv "http://quotes.money.163.com/service/chddata.html?code=1000592&start=20210324&end=20210324"

http://www.szse.cn/market/stock/list/index.html
http://www.sse.com.cn/assortment/stock/list/share/

awk '{print "wget -O "$1".csv \"http://quotes.money.163.com/service/chddata.html?code="$1"&start=20210324&end=20210324\""}' hs.txt  | sh 

LC_CTYPE=C&&cat *.csv | sort -u | sed 's/None/0/g'  > all210324.csv 
iconv -f gb2312 -t utf-8 all210324.csv > all210324_1.csv

load data infile '/Users/lihong/gitrep/joingp/rx210326.txt' into table dayline character set utf8  fields terminated by ' ' OPTIONALLY ENCLOSED BY '"'  lines terminated by '\n' (date,code,open,close, low,high,volume,money,factor, high_limit,low_limit,avg,pre_close,paused, m5,m10,m20,m60,m250, dif,dea,macd);


/Users/lihong/gitrep/src

mysql -uroot  stock -e "load data infile 'D:\\book\\vlog\\stock\\shen\\all.csv' into table rx character set gb2312  fields terminated by',' OPTIONALLY ENCLOSED BY '\"'  lines terminated by '\n' IGNORE 1 LINES (rq,dm,mc,sp,zg,zd,kp,qsp,zde,zdf,hs,cjl,cje,zsz,ltsz,cjb);"

CREATE TABLE `dayline` (
  `id` INT(10) UNSIGNED NOT NULL AUTO_INCREMENT,
  `date` DATE NOT NULL,
  `code` VARCHAR(16) COLLATE utf8mb4_unicode_ci NOT NULL,
  `open` DOUBLE UNSIGNED NOT NULL,
  `close` DOUBLE UNSIGNED NOT NULL,
  `low` DOUBLE UNSIGNED NOT NULL,
  `high` DOUBLE UNSIGNED NOT NULL,
  `volume` DOUBLE UNSIGNED NOT NULL,
  `money` DOUBLE NOT NULL,
  `factor` DOUBLE NOT NULL,
  `high_limit` DOUBLE NOT NULL,
  `low_limit` DOUBLE NOT NULL,
  `avg` DOUBLE NOT NULL,
  `pre_close` DOUBLE NOT NULL,
  `paused` VARCHAR(25) COLLATE utf8mb4_unicode_ci NOT NULL,
  `m5` DOUBLE DEFAULT NULL,
  `m10` DOUBLE UNSIGNED NOT NULL DEFAULT 0,
  `m20` DOUBLE UNSIGNED NOT NULL DEFAULT 0,
  `m60` DOUBLE UNSIGNED NOT NULL DEFAULT 0,
  `m250` DOUBLE UNSIGNED NOT NULL DEFAULT 0,
  `dif` DOUBLE DEFAULT 0,
  `dea` DOUBLE DEFAULT 0,
  `macd` DOUBLE DEFAULT 0,
  `r1` VARCHAR(25) DEFAULT NULL,
  `create_time` TIMESTAMP NULL DEFAULT CURRENT_TIMESTAMP(),
  PRIMARY KEY (`id`),
  UNIQUE KEY `date` (`date`,`code`)
) ENGINE=INNODB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

python gp.py 2021-01-22&&python gp.py 2021-01-21&&python gp.py 2021-01-20&&python gp.py 2021-01-19&&python gp.py 2021-01-18

市值在200亿以下的，50亿以上，最好是全流通的（大股东已经减持完毕）没有爆炒过的

每股收益大于0.5 ;近两年没翻过倍

《情深缘浅》《烟火里的尘埃》《可惜我是一个平凡的人》《南城》《没有什么不同》
《那个女孩》《不惑余生》《风度》《浪人情歌》《可乐》《微笑着告别了我的青春》《失忆》

# 技术面
## 趋势突破,
趋势线,均线,百分比,量能

## 形态突破,
W M 三角形,箱体,123法则

## 顺势回调,
三重滤网,百分比回调.

## 超跌反弹
2B破底翻

# 基本面
## 指数法 预期法(行业) 安全边际法(价值/估值高于价格,很难找),
行业/周期/政策
每股净资产高 
市盈率(PE)小 
每股公积金高 
体积中等,
通过分析流通股东了解庄家控筹-越高越好 
行业调查

# 如何判断假突破 
1. 量,向上一定要有量
2. 价格最好突破压力区
3. 判断真假没有意义,赚钱真,止损假. 坚定止损.


# 策略:滤网,趋势,位置,买入信号
## 滤网
* 如果做股票,需要基本面先进行筛选: 市场容量,行业利润,企业增长,行业龙头等
* 如果做商品期货,需要进行基差+库存+升贴水的研究,然后确定大方向
* 如果只是日内交易或超短线,就用开盘价+盘口作为滤网

## 趋势
* 高低点法: 只要高点和低点**不断抬高或降低**,就是趋势成立,坚定只按现在趋势做.
* 二分法: 当高点不再创新高,只要不跌破前一个波段的50%,仍然趋势成立
* 中轴法: 一个倒三角形,菱形,或者箱体的中轴可以作为趋势判断的标准,只要不跌破或者突破,现在趋势仍然成立

## 位置
* 顺势突破: 形态突破,N字突破,支撑压力
* 顺大逆小回调: 支撑压力关键区
* 短线逆势衰竭: 5浪行情尾盘+假动作
* 中长线吸筹完成: 长期或者巨大形态突破

## 买入信号
* 13个买入信号: 抛售高潮,2B,假诱空,头肩底,双突破,强势出现,跳跃小溪.回抽确认,N字突破.箱体弹簧,周期共振,初次急跌破50%波段位置 

## 盈利平仓(人性与理性的结合)
**最好的平仓方法是分批平仓,将资金分成三份**
1. 顺势平仓: 到达你的目标位,可以是前期的支撑与压力,可以是波浪理论的反推,可以是江恩的时间周期,可以是其他原因(30%)
2. 双跌破或者关键支撑压力位跌破.(30%)
3. 前一个波段50%被跌破

## 1. 抛售高潮
由于是逆势行情,一定要加基本面过滤,不能看到就买.
1. 顺大逆小
2. 巨量出现
3. 小5浪模式(可有可无)
4. 前期的关键支撑压力位

## 2. 2B
2B概率很高,但也是逆势行情,也一定要加基本面过滤,不能看到就买
1. 顺大逆小
2. 小5浪模式
3. 最好有一个小形态,可以是三角形,箱体或者棱形(与抛售高潮不同的是,他不需要巨量,必须阳线快速出现,必须第二根就要上去,这才是明确的猎杀止损,或者主力进入)

## 3 假诱空
在一轮下跌后,先向上突破形态,在快速下打,出巨量,然后慢慢V型反转回到愿区间
1. 最好顺大逆小,完全逆势也可以,但概率就不高了   
2. 小5浪模式
3. 一定有小形态,可以是三角形,箱体或者棱形的向上多,在快速向下空,然后又V型反转.(与前面波动相比,波动率明显加大)

## 4 顺势头肩底
头肩属于小级别逆势,所以最好顺大逆小  
趋势转折,市场情绪变弱
1. 左肩 头部 右肩 突破,每一步都配合量
  头肩底: 左肩到右肩不断缩量,突破头肩放量(理想状态) 
  头肩顶: 左肩头部巨量,右肩缩量,突破放量(理想状态) 
2. 顺大逆小是关键的关键(核心的核心)
3. 趋势线突破 (衰竭反转)
4. 箱体弹簧 (中继头肩)
5. 钻石型变异 (大底大顶)

## 5 双突破
### 双突破系统构建
1. 核心: 趋势线+2B/形态/量能积累/关键支撑压力点/筹码分布 (筹码分布和量能积累大概率相符)  
2. 计划  
  开: 止损:前期高低点/量能积累点-分价图,加仓:支撑与压力互换/形态突破  
  平: 分批:30%-50%趋势线跌破/盈亏比->沽压 30%; 50%双跌破 (永远不要想鱼头鱼尾)
3. 三重滤网-顺大逆小  
4. 四重滤网-指标背离/量价背离

### 市场规律
1. 形态周期越长,行情越大
2. 形态越完美,行情越好
3. 牛市中,小底部可能大涨,大顶部可能小跌
4. 长期牛市中,第一次双突破 80%被止损;长期熊市中,第一次双突破,70%回调
5. 对不同形态 顺势

## 6 强势出现
一波吸筹后面,出现明显缩量后的放量,可以用中轴法判断.
**中轴+放量=强势出现**

## 7 跳跃小溪
一波吸筹后期,放量突破颈线  
**突破颈线+放量=跳跃小溪**

## 8 回抽确认
跳跃小溪跳后马上开始回调,回调后在颈线出获得支撑,再次向上时,就是一个安全的买点  
**突破颈线+回抽+再次确认向上=回抽确认**

## 9 N字突破
趋势确认后,每一次的N字突破

## 10 支撑压力互换
趋势确认后,每一次支撑压力呼唤点位  
**确认趋势+支撑压力互换(事不过三原则)** 

## 11 箱体弹簧
趋势向上,走一个箱体,突然跌破,又快速拉起,做出箱体弹簧.
符合大趋势+快速跌破+快速拉起=箱体弹簧

## 12 形态突破 三角形,箱体或者其他形态顺势突破(多周期共振)

## 13 趋势急跌到50%
一波大趋势下跌,急跌50%(只能是第一次,而且必须急跌)

## 衰竭反转2B
趋势衰竭,吸筹急迫,猎杀止损.  
1. 确定周期(极为重要),安周期止损止盈,大周期很少,小周期较多
2. 第二根K线必须拉起,越快概率越大
3. 量价配合过滤  
4. 前期支撑压力过滤  
衰竭2B和威科夫理论中弹簧区别

## 支撑与压力测试
行为金融学和市场情绪,持仓成本
1. 前期高低点位
2. 前期成交,量能密集区(机构建仓区)
3. 趋势线与切线交叉
支撑与压力互换找开仓点

## 反转形态
1. 头肩
2. 扩散三角形
3. 菱形(钻石)
4. 圆弧(茶壶带柄)
5. 岛形反转

## 中继形态
1. 普通三角形
2. 箱体
3. 旗形
4. 通道(斜箱)
5. N字

## 极简双突破
1. 趋势线+形态
2. 趋势线+2B
3. 趋势线+关键压力支撑点互换
4. 趋势线+量能区
5. 趋势线+筹码区
6. 趋势线+巨量+孕线

# 策略
* 开: 双突破
* 加: 
  1. 回调+支撑与压力
  2. N字
  3. 形态突破
  4. 弹簧
* 平: 分批(理性首平+运气飘单)
  1. 强压区域
  2. 前高低点
  3. 前量能密集点
  4. 双突破/运气飘单
* 止损: 固定(总资产2%) 或者最近的高低点

# 威科夫
## 二元对立 
1. 机构关心供应/需求
2. 吸筹区 派发区
3. 对手成本
4. 支撑压力

比较PE 筛洗可比企业 区间日均总市值 净资产 固定资产 营收 地域 

去掉PE为负的企业  算平均 加权平均 权重=总市值/所有企业总市值 权重PE 中位数 去掉异常值

//==================

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
) ENGINE=InnoDB AUTO_INCREMENT=239304 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;


alter table valuation add yearday_avg_market_cap double NOT NULL DEFAULT 0; 

select a.yearday_avg_market_cap,a.pe_ratio,a.market_cap ,a.code,b.name ,b.zjw_code,b.zjw_name ,pe_ratio*a.yearday_avg_market_cap/( select sum(v.yearday_avg_market_cap) industry_market_cap  from valuation v,industry i  where v.code=i.code and zjw_code='C32' and  v.day='2021-04-16' and v.pe_ratio>0 and i.zjw_code = b.zjw_code )xx from  valuation a,industry b   where a.code=b.code and a.day='2021-04-16' and a.pe_ratio>0 and zjw_code='C32' order by 2 ;

select count(1) cnt,sum(v.yearday_avg_market_cap) industry_market_cap,i.zjw_code,i.zjw_name
  from valuation v,industry i
  where v.code=i.code and v.day='2021-04-16' and pe_ratio>0
  group by i.zjw_code,i.zjw_name ;

drop table industry1;
CREATE TABLE `industry1` (
`id` int(10) unsigned NOT NULL AUTO_INCREMENT,
`name` varchar(20) COLLATE utf8mb4_unicode_ci NOT NULL,
`code` varchar(16) COLLATE utf8mb4_unicode_ci NOT NULL,
`sw_l1_code` varchar(10) COLLATE utf8mb4_unicode_ci NOT NULL,
`sw_l1_name` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL,
`sw_l2_code` varchar(10) COLLATE utf8mb4_unicode_ci NOT NULL,
`sw_l2_name` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL,
`sw_l3_code` varchar(10) COLLATE utf8mb4_unicode_ci NOT NULL,
`sw_l3_name` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL,
`zjw_code` varchar(10) COLLATE utf8mb4_unicode_ci NOT NULL,
`zjw_name` varchar(128) COLLATE utf8mb4_unicode_ci NOT NULL,
`create_time` timestamp NULL DEFAULT current_timestamp(),
`cnt` int NOT NULL DEFAULT 0,

`inc_revenue_year` double NOT NULL DEFAULT 0,
`inc_revenue_year_avg` double NOT NULL DEFAULT 0,
`inc_revenue_year_rank` int NOT NULL DEFAULT 0,

`inc_revenue_annual` double NOT NULL DEFAULT 0,
`inc_revenue_annual_avg` double NOT NULL DEFAULT 0,
`inc_revenue_annual_rank` int NOT NULL DEFAULT 0,

`net_operate_cash_flow` double NOT NULL DEFAULT 0,
`operating_revenue` double NOT NULL DEFAULT 0,
`cfo_sales` double NOT NULL DEFAULT 0,
`cfo_sales_avg` double NOT NULL DEFAULT 0,
`cfo_sales_rank` int NOT NULL DEFAULT 0,

`total_assets` double NOT NULL DEFAULT 0,
`total_owner_equities` double NOT NULL DEFAULT 0,
`leverage_Ratio` double NOT NULL DEFAULT 0,
`leverage_Ratio_avg`  double NOT NULL DEFAULT 0,
`leverage_Ratio_rank` int NOT NULL DEFAULT 0,

`report_date` timestamp not NULL,
`public_date` timestamp not NULL,
PRIMARY KEY (`id`),
UNIQUE KEY `idx_code_day` (`code`)
) ENGINE=InnoDB AUTO_INCREMENT=1038851 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

inc_revenue_year_on_year 营收同比
inc_revenue_annual 营收环比
 
cfo_sales net_operate_cash_flow/operating_revenue 经营活动现金流量净额/营收

Leverage_Ratio balance.total_assets/balance.total_owner_equitie 总资产/股东权益 

select code,name , inc_revenue_year,
rank() over(ORDER BY inc_revenue_year desc) as 'inc_revenue_year_rank' ,
(select sum(inc_revenue_year)/count(1) from industry1 b where b.zjw_code=a.zjw_code) inc_revenue_year_avg , 

rank() over(ORDER BY inc_revenue_annual desc) as 'inc_revenue_annual_rank' ,
(select sum(inc_revenue_annual)/count(1) from industry1 b where b.zjw_code=a.zjw_code) inc_revenue_annual_avg

from industry1 a where zjw_code= 'I65' limit 5



select code,name , inc_revenue_year, report_date,zjw_name, (select count(1) from industry1 b where b.zjw_code=a.zjw_code) cnt
rank() over(partition by zjw_code ORDER BY inc_revenue_year desc)  inc_year_rank,
round((select sum(inc_revenue_year)/count(1) from industry1 b where b.zjw_code=a.zjw_code),4) inc_year_avg , 
inc_revenue_annual,
rank() over(partition by zjw_code ORDER BY inc_revenue_annual desc) inc_nnual_rank ,
round((select sum(inc_revenue_annual)/count(1) from industry1 b where b.zjw_code=a.zjw_code),4) inc_annual_avg,
round(net_operate_cash_flow/operating_revenue,4) as cfo_sales,

rank() over(partition by zjw_code ORDER BY net_operate_cash_flow/operating_revenue desc)  cfo_rank ,
round((select sum(net_operate_cash_flow/operating_revenue)/count(1) from industry1 b where b.zjw_code=a.zjw_code),4) cfo_avg,
round(total_assets/total_owner_equities,4) as leverage,
rank() over(partition by zjw_code ORDER BY total_assets/total_owner_equities)  leverage_rank ,
round((select sum(total_assets/total_owner_equities)/count(1) from industry1 b where b.zjw_code=a.zjw_code),4) leverage_avg
from industry1 a 