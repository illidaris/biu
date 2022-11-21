/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
)

// propertyCmd represents the property command
var propertyCmd = &cobra.Command{
	Use:   "property",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		absPath, err := filepath.Abs(filePath)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("[%s] will generate property \n", absPath)
		// property.Input(filePath)
	},
}
