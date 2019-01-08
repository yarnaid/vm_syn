package main

import (
	"log"

	"github.com/pkg/errors"

	"github.com/spf13/cobra"
	vm "github.com/yarnaid/vm_syn/vminstance"
)

var rootCmd = &cobra.Command{
	Use:   "vm",
	Short: "Execute code in vm",
	Run: func(cmd *cobra.Command, args []string) {
		vm := vm.New()
		vm.LoadFile(args[0])
		go vm.Run()
		startGui(vm)
	},
	Args: cobra.ExactValidArgs(1),
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {}

func main() {
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(errors.Wrap(err, "main execute error"))
	}
}
