package main

import (
	"testing"
	"github.com/tangtaoit/util"
	"os"
	"pay/config"
)

func TestNewAccountRecord(t *testing.T) {
	os.Setenv("GO_ENV", "tests")
	config.GetSetting()

	err := AccountRecharge("1160527760470002",int64(10))
	util.CheckErr(err)
}