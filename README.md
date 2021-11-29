<!--
 * @Author: your name
 * @Date: 2021-11-01 13:02:08
 * @LastEditTime: 2021-11-29 22:11:33
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /teccamp-envelop-rain/README.md
-->
# teccamp-envelop-rain
![](https://img.shields.io/badge/license-MIT-blue) ![](https://img.shields.io/badge/build-passing-blue)

## 1. Overview
项目实现三个接口，分别是抢红包接口(/snatch)，拆红包接口(/open)以及查看钱包列表接口(/get_wallet_list)。并且要求红包雨活动当中总金额、总个数、用户可抢最大个数、用户抢到的概率可配置，并且支持活动总金额的热更新。

## 2. Requirement
* mysql
* kafka
* linux
* redis

## 3. Quick Start
```bash
go build && ./envelop-rain
```

## 4. 性能优化
* 雪花算法生成分布式唯一红包ID；
* 在项目当中引入了Kafka对DB做异步写;
* 在local cache引入了布隆过滤器对抢满的用户、红包抢完、红包金额不足等情况做标记，尽可能减少redis网络请求，提高吞吐量；
* Redis lua脚本执行事务，保证高并发下数据的一致性和可靠性；
* Hystrix熔断器使用，我们在本次项目当中也加入了熔断器来应对突发的情况，在请求量超过一定阈值，并且错误率达到20%以上，熔断器会自动熔断对应的service，直接返回；ServiceStatusUnavaluable的响应，等待一定时间之后再去探测接口是否恢复；
* 令牌桶限流，这个在Hystrix-go当中好像就已经实现了令牌桶限流，当然我们也自己实现了一个令牌桶限流的算法。根据压测的实际结果，实际设置的流量控制分别为3w、2w、2w，超出流量限制的流量会立即返回。

## 5. Test
选择使用 **wrk** 工具对三个接口进行压力测试。通过读取写好的lua脚本对我们部署好的接口进行压力测试，我们配置线上共有10w个红包，压测时间为60s。

**Snatch接口压测结果：23487QPS**
![20211129220301](https://i.loli.net/2021/11/29/1HUE63oOiF2smuf.png)

**Open接口压测结果：5286QPS**
![20211129220338](https://i.loli.net/2021/11/29/rkxPXsdqBtwFEjm.png)

**Open接口压测结果：11918QPS**
![20211129220406](https://i.loli.net/2021/11/29/6Spf8rwVHUWienA.png)
