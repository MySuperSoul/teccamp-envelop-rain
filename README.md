<!--
 * @Author: your name
 * @Date: 2021-11-01 13:02:08
 * @LastEditTime: 2021-11-01 16:42:05
 * @LastEditors: your name
 * @Description: In User Settings Edit
 * @FilePath: /teccamp-envelop-rain/README.md
-->
# teccamp-envelop-rain

### 1. 红包达到上限后还会返回如下响应：
{
    "code": 1,
    "data": {},
    "msg": "you are not lucky, try again"
}
### 修改方法：将对概率的判断移到Amount的判断之后
### 另：是否当用户的红包达到上限后，将用户id存到redis。这样每次访问数据库前多加一个redis的访问，性能对比？

### 2. post请求不输入uid，默认是0， 能否拒绝这个请求？还是由前端来处理这个问题。
### 修改方法：在ConvertString函数中进行特判，如果传入空字符串，直接报错退出

### 3. 对数据库链接的测试
### db_test.go中添加了对数据库的测试

### 4. todo: 配置文件是否在程序启动前，利用命令行参数可以配置。提供默认配置。
### 利用flag包，不知道需不需要，目前没有完成

### ......