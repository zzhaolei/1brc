# 1brc
1️⃣🐝🏎️ The One Billion Row Challenge

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

## 生成数据
克隆原始[仓库](https://github.com/gunnarmorling/1brc)：
```shell
git clone https://github.com/gunnarmorling/1brc

cd 1brc/src/main/python
python3 create_measurements.py 1000000000
```
生成的数据会在`1brc/measurements.txt`，大约有`15G`左右(平台)。

## 生成基础结果
复制生成的`1brc/measurements.txt`文件，或直接指定文件路径
```shell
go run baseline/main.go measurements.txt > base.txt
```
`base.txt`是基线结果，后续的优化可以和这个结果进行对比。
