/*
Copyright © 2025 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"github.com/end1essrage/tempro/cmd"
	"github.com/end1essrage/tempro/factory"
)

func main() {
	factory.Init()
	cmd.Execute()
}
