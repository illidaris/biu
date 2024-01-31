/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"biu/internal/excel2model"
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
)

// excel2modelCmd represents the excel2model command
var excel2modelCmd = &cobra.Command{
	Use:   "excel2model",
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
		excel2model.Invoke(filePath)
	},
}

func init() {
	rootCmd.AddCommand(excel2modelCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// excel2modelCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// excel2modelCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
