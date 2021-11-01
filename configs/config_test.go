/*
 * @Author: your name
 * @Date: 2021-11-01 13:02:08
 * @LastEditTime: 2021-11-01 16:16:44
 * @LastEditors: Please set LastEditors
 * @Description: In User Settings Edit
 * @FilePath: /teccamp-envelop-rain/configs/config_test.go
 */
package configs

import (
	"testing"
)

func TestGenerateConfigFromFile(t *testing.T) {
	config := GenerateConfigFromFile()
	t.Log(config)
}
