<div align="center"> <a href="https://github.com/Harvey-Specter/plunge"> <img width="120" src="./img/logo_1024.png"> </a><br/><br/>

[![license](./img/mit_lic.svg)](LICENSE)
<h1>Plunge</h1>
</div>

<**中文** |  [日文](./README.md) >

## 介绍

Plunge是通过识别图表的形态量价来选股的工具,主要用来参考学习,不作为投资依据

## 准备

- [Python](https://www.python.org/) - 行情数据采集使用python开发
- [Golang](https://go.dev/) - 识别图表使用Go语言开发
- [MariaDB](https://mariadb.org/) - 数据库使用MariaDB（也可以使用MySQL的其他分支）
- [tushare.pro](https://tushare.pro/) - 行情数据源

## 安装和使用
- 获取代码  

      git clone https://github.com/Harvey-Specter/plunge.git
- 采集某天的数据并且计算那天的均线（计算均线需要历史数据，例如计算MA60需要60天的行情数据）

      cd plunge
      python savedate/tsdata.py 2022-07-01
- 编译&运行

      go build -o plunge
      chmod +x ./plunge && ./plunge dayline 2022-07-01

- 结果文件:程序会默认在工作目录中生成结果文件, 同时结果数据会保存到Plunge服务器, 也可以根据需要修改输出文件格式
      
## 识别的形态
下面是的符合形态,实际挑选出来会有一定误差,可以根据喜好修改程序来收紧/放宽筛选规则

1.	**头肩底放量突破颈线**    
   ![头肩底突破1](./img/tj1.png)  

<hr/>

2.	**头肩底/箱体中轴以上放量启动**  
   ![头肩底/箱体中轴以上放量启动](./img/tj2.png)
<hr/>

3.	**首板之后没跌下来并且出现吸筹形态**  
   ![吸筹](./img/xc.png)
<hr/>

4.	**上升趋势中连续出现高/低开并且收星** 
   ![收星](./img/star.png)
<hr/>

5.	**缺口附近出现支持** 
   ![缺口](./img/qk.png)

## 表结构 / desc dayline

      +-------------+------------------+------+-----+---------------------+----------------+
      | Field       | Type             | Null | Key | Default             | Extra          |

      +-------------+------------------+------+-----+---------------------+----------------+

      | id          | int(10) unsigned | NO   | PRI | NULL                | auto_increment |

      | date        | date             | NO   | MUL | NULL                |                |

      | code        | varchar(16)      | NO   | MUL | NULL                |                |

      | open        | double unsigned  | NO   |     | NULL                |                |

      | close       | double unsigned  | NO   |     | NULL                |                |

      | low         | double unsigned  | NO   |     | NULL                |                |

      | high        | double unsigned  | NO   |     | NULL                |                |

      | volume      | double unsigned  | NO   |     | NULL                |                |

      | pre_close   | double unsigned  | NO   |     | 0                   |                |

      | paused      | varchar(25)      | NO   |     | 0                   |                |

      | m5          | double unsigned  | NO   |     | 0                   |                |

      | m10         | double unsigned  | NO   |     | 0                   |                |

      | m20         | double unsigned  | NO   |     | 0                   |                |

      | m30         | double unsigned  | NO   |     | 0                   |                |

      | m60         | double unsigned  | NO   |     | 0                   |                |

      | create_time | timestamp        | YES  |     | current_timestamp() |                |

      +-------------+------------------+------+-----+---------------------+----------------+
## 许可
[MIT](./img/LICENSE)