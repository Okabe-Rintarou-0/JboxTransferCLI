package main

import (
	_ "jtrans/db"
	"jtrans/transfer"
)

func main() {
	transfer.Execute()
}
