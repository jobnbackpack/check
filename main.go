/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"jobnbackpack/check/cmd"
	"jobnbackpack/check/db"
)

func main() {
	db.ConnectDB()
	cmd.Execute()
}
