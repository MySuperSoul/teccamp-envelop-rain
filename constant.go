package main

const (
	SNATCH_SUCCESS           int = 0
	SNATCH_NOT_LUCKY         int = 1
	SNATCH_NO_RED_PACKET     int = 2
	SNATCH_EXCEED_MAX_AMOUNT int = 3
)

const (
	SNATCH_SUCCESS_MESSAGE           string = "success snatch"
	SNATCH_NOT_LUCKY_MESSAGE         string = "you are not lucky, try again"
	SNATCH_NO_RED_PACKET_MESSAGE     string = "you are so slow, no red packet yet"
	SNATCH_EXCEED_MAX_AMOUNT_MESSAGE string = "you are exceed the max amount, no more red packet yet"
)

const (
	OPEN_SUCCESS        int = 0
	OPEN_INVALID_USER   int = 1
	OPEN_INVALID_PACKET int = 2
	OPEN_REPEAT         int = 3
	OPEN_NOT_MATCH      int = 4
)

const (
	OPEN_SUCCESS_MESSAGE        string = "success open envelop"
	OPEN_INVALID_USER_MESSAGE   string = "invalid user id"
	OPEN_INVALID_PACKET_MESSAGE string = "invalid packet id"
	OPEN_REPEAT_MESSAGE         string = "packet has been opened yet"
	OPEN_NOT_MATCH_MESSAGE      string = "user don't have this packet"
)

const (
	WALLET_SUCCESS int = 0
)

const (
	WALLET_SUCCESS_MESSAGE string = "check wallet success"
)

// system config path
const (
	CONFIG_PATH string = "configs/config.json"
)
