###
 # @Author: your name
 # @Date: 2021-11-01 16:12:58
 # @LastEditTime: 2021-11-02 21:05:40
 # @LastEditors: Please set LastEditors
 # @Description: In User Settings Edit
 # @FilePath: /teccamp-envelop-rain/tmp.sh
### 

go test -c  common/util_test.go common/util.go common/snowflake.go common/snowflake_test.go && ./common.test

go test -c  repository/mysql_test.go repository/redis_test.go repository/redis.go repository/mysql.go repository/model.go && ./repository.test

rm *.test