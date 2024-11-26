/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package main

import (
	"log"

	"github.com/inoth/toybox/cmd/toybox/internal/project"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "toybox",
	Short:   "Toybox: An simple toolkit for Go services.",
	Long:    ``,
	Version: version,
}

func init() {
	rootCmd.AddCommand(project.CmdNew)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
