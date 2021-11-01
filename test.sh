###
 # @Author: your name
 # @Date: 2021-11-01 16:12:58
 # @LastEditTime: 2021-11-01 16:25:24
 # @LastEditors: Please set LastEditors
 # @Description: In User Settings Edit
 # @FilePath: /teccamp-envelop-rain/tmp.sh
### 

go test -c  common/util_test.go common/util.go && ./common.test

go test -c  db/db_test.go db/redis.go db/mysql.go db/model.go && ./db.test

go test -c redpacket/gen_red_packet_test.go redpacket/gen_red_packet.go && ./redpacket.test

rm *.test