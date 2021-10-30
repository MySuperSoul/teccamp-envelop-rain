package main

const (
	SNATCH_SUCCESS       int = 0
	SNATCH_NOT_LUCKY     int = 1
	SNATCH_NO_RED_PACKET int = 2
)

const (
	SNATCH_SUCCESS_MESSAGE       string = "success snatch"
	SNATCH_NOT_LUCKY_MESSAGE     string = "you are not lucky, try again"
	SNATCH_NO_RED_PACKET_MESSAGE string = "you are so slow, no red packet yet"
)

// system config path
const (
	CONFIG_PATH string = "configs/config.json"
)
