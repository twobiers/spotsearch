package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/twobiers/spotsearch/internal/pkg/client"
)

type authCmdConfig struct {
	test bool
}

var (
	authCmd = &cobra.Command{
		Use:   "auth",
		Run: func(cmd *cobra.Command, args []string) {
	  		if Test {
				test()
				return
			}

			client.Authenticate()
		},
  	}
	Test bool
)

func init() {
	authCmd.Flags().BoolVarP(&Test, "test", "t", false, "")
	rootCmd.AddCommand(authCmd)
}

func test() {
	result := client.TestAuth()
	if !result {
		fmt.Println("Authentication is invalid")
		os.Exit(1)
	}

	fmt.Println("Authentication works as intended")
}