/*
 * @Author: your name
 * @Date: 2021-11-01 13:02:08
 * @LastEditTime: 2021-11-03 20:39:53
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /teccamp-envelop-rain/common/constant.go
 */
package router

const (
	SNATCH_SUCCESS           int = 0
	SNATCH_NOT_LUCKY         int = 1
	SNATCH_NO_RED_PACKET     int = 2
	SNATCH_EXCEED_MAX_AMOUNT int = 3
	SNATCH_EMPTY_UID         int = 4
	SNATCH_JSON_PARSE_ERROR  int = 5
)

const (
	SNATCH_SUCCESS_MESSAGE           string = "success snatch"
	SNATCH_NOT_LUCKY_MESSAGE         string = "you are not lucky, try again"
	SNATCH_NO_RED_PACKET_MESSAGE     string = "you are so slow, no red packet yet"
	SNATCH_EXCEED_MAX_AMOUNT_MESSAGE string = "you are exceed the max amount, no more red packet yet"
	SNATCH_EMPTY_UID_MESSAGE         string = "input a empty uid"
	SNATCH_JSON_PARSE_ERROR_MESSAGE  string = "json parsing error"
)

const (
	OPEN_SUCCESS          int = 0
	OPEN_INVALID_USER     int = 1
	OPEN_INVALID_PACKET   int = 2
	OPEN_REPEAT           int = 3
	OPEN_NOT_MATCH        int = 4
	OPEN_EMPTY_ID         int = 5
	OPEN_JSON_PARSE_ERROR int = 6
)

const (
	OPEN_SUCCESS_MESSAGE          string = "success open envelop"
	OPEN_INVALID_USER_MESSAGE     string = "invalid user id"
	OPEN_INVALID_PACKET_MESSAGE   string = "invalid packet id"
	OPEN_REPEAT_MESSAGE           string = "packet has been opened yet"
	OPEN_NOT_MATCH_MESSAGE        string = "user don't have this packet"
	OPEN_EMPTY_ID_MESSAGE         string = "input a empty uid or packet id"
	OPEN_JSON_PARSE_ERROR_MESSAGE string = "json parsing error"
)

const (
	WALLET_SUCCESS          int = 0
	WALLET_EMPTY_ID         int = 1
	WALLET_JSON_PARSE_ERROR int = 2
)

const (
	WALLET_SUCCESS_MESSAGE          string = "check wallet success"
	WALLET_EMPTY_ID_MESSAGE         string = "input a empty uid"
	WALLET_JSON_PARSE_ERROR_MESSAGE string = "json parsing error"
)
