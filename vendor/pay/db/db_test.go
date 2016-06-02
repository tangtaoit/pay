package db

import (
	"testing"
	"fmt"
)



func TestTrade_Update(t *testing.T) {
	sess := NewSession()
	sess.Update("trades").Where("trade_no=?","1160527553280001").Set("status",1).Exec()
}

func TestTrade_UpdateTX(t *testing.T) {
	sess := NewSession()
	tx, _ :=sess.Begin()
	tx.Update("trades").Where("trade_no=?","1160527553280001").Set("status",1).Exec()
	tx.Commit()
}





