###
 # @Author: your name
 # @Date: 2021-11-01 16:12:58
 # @LastEditTime: 2021-11-01 16:25:24
 # @LastEditors: Please set LastEditors
 # @Description: In User Settings Edit
 # @FilePath: /teccamp-envelop-rain/tmp.sh
### 

go test -c  configs/config_test.go configs/config.go && ./configs.test

go test -c  db/db_test.go db/redis.go db/mysql.go  db/model.go && ./db.test

rm *.test