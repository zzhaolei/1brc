# 1brc
十亿行挑战（1️⃣🐝🏎️ The One Billion Row Challenge）

## 目录
- [目标](#目标)
  - [生成十亿行挑战所需的数据](#生成十亿行挑战所需的数据)
- [限制](#限制)
- [更多细节](#更多细节)

## 目标
文本文件包含了一系列气象站的温度值。 每行是一个测量值，格式为`<string: station name>;<double: measurement>`，其中测量值精确到一位小数。以下是一些示例行：
```txt
Hamburg;12.0
Bulawayo;8.9
Palembang;38.8
St. John's;15.2
Cracow;12.6
Bridgetown;26.9
Istanbul;6.2
Roseau;34.4
Conakry;31.2
Istanbul;23.0
```

任务是编写一个程序，该程序读取文本文件，计算每个气象站的最低、平均和最高温度值，并将结果输出到`stdout`，
格式如下（按气象站名称字母顺序排序，并且每个气象站的结果值格式为`<min>/<mean>/<max>`，保留一位小数点）：
```
{Abha=-23.0/18.0/59.2, Abidjan=-16.2/26.0/67.3, Abéché=-10.0/29.4/69.0, Accra=-10.1/26.4/66.4, Addis Ababa=-23.7/16.0/67.0, Adelaide=-27.8/17.3/58.5, ...}
```

## 限制
只能使用标准库实现。

### 生成十亿行挑战所需的数据
克隆原始[仓库](https://github.com/gunnarmorling/1brc)：
```shell
git clone https://github.com/gunnarmorling/1brc

cd 1brc/src/main/python
python3 create_measurements.py 1000000000
```
生成的数据会在`1brc/measurements.txt`，约为`15Gi`的大小。

## 更多细节
关于优化的更多细节，请看这篇[文章](https://zzhaolei.github.io/posts/%E5%8D%81%E4%BA%BF%E8%A1%8C%E6%8C%91%E6%88%98/)
