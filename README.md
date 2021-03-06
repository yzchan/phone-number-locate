国内手机号归属地查询
-----

## 关于

手机号归属地查询是一个十分常见的需求，一般采用三方接口获取。在实际项目中如果查询量比较多，一般会对查询结果进行缓存。其原理就是缓存手机号段（手机号前7位）和查询结果。

缓存结果一般可以长期使用，开放携号转网之后，会出现运营商信息不准的情况，但是归属地信息依然是有效的。当然也可以根据需要自行制定缓存更新策略

下述手机号说明来源于网络：

```
我国使用的号码为11位，其中各段有不同的编码方向：前3位—网络识别号；第4-7位—地区编码；第8-11位—用户号码。
MDN号码的结构如下：CC + MAC + H0 H1 H2 H3 + ABCD其中：
【CC】：国家码，中国使用86。
【MAC】：移动接入码，本网采用网号方案，为133。
【H0H1H2H3】：HLR识别码，由运营商统一分配。
【ABCD】：移动用户号，由各HLR自行分配。
```

## 运营商号段

关于手机号中前3位运营商号段，具体详细可以[查询这里](https://zh.wikipedia.org/wiki/%E4%B8%AD%E5%9B%BD%E5%86%85%E5%9C%B0%E7%A7%BB%E5%8A%A8%E7%BB%88%E7%AB%AF%E9%80%9A%E8%AE%AF%E5%8F%B7%E6%AE%B5)

截止到2021年4月为止：中国移动共计26个号段、中国联通共计13个号段、中国电信共计15个号段，共计54个运营商号段。

去除掉虚拟号段、卫星号段和物联网号段后约为40多个平时手机号常用的号段。如果4位HLR识别码全部用满的话，7位号段总共有40多万种组合。

# 存储策略

用mysql或者mongodb存储所有可能的结果，用redis做缓存提升查询性能。查询未命中（系统中不存在）时去三方查询结果并保存。

借助redis的hash类型进行缓存。几十万数据存一个hash会导致一个该键太大，影响查询性能。可以按照前3位MAC码进行分片，每个分片里面最多存储1w个键值对。

如果能购买或者下载到一些初始数据就更好了。

- [buybook归属地查询](http://www.buybook.com.cn/) 数据比较规范。就是有点老了。而且号段不全
- [ip138手机归属地查询](https://ip138.com/sj/) 可以免费查询，同时也开放了api查询
- [360手机归属地查询](http://cx.shouji.360.cn/phonearea.php?number=1360150) 可以免费抓取
- [淘宝手机归属地查询](https://tcc.taobao.com/cc/json/mobile_tel_segment.htm?tel=13333333333) 淘宝免费归属地查询，只返回省，无法定位到市
- [易源数据](https://www.showapi.com/apiGateway/view/6) 提供api接口查询
- [挖数据](https://www.wapi.cn/source/3.html) 有接口查询，并且可以购买归属地数据源
- [聚合数据](https://www.juhe.cn/docs/api/id/11) 接口查询
- [阿里云市场](https://www.aliyun.com/ss/?k=%E6%89%8B%E6%9C%BA%E5%8F%B7%E5%BD%92%E5%B1%9E%E5%9C%B0) 整合了很多三方接口

# 新增46w+号段数据

新增phone_sections.csv.zip文件为csv格式的zip压缩包。包含46w+的手机号段数据。数据主要来源于 [360手机归属地查询](http://cx.shouji.360.cn/phonearea.php?number=1360150)
、[buybook归属地查询](http://www.buybook.com.cn/)、以及CSDN下载的部分数据，号段数据可能并不十分准确。
