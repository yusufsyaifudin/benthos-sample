package main

import (
	"github.com/mitchellh/cli"
	"github.com/yusufsyaifudin/benthos-sample/golang/calc"
	"github.com/yusufsyaifudin/benthos-sample/golang/wallet"
	"log"
	"os"
)

func main() {
	c := cli.NewCLI("app", "1.0.0")
	c.Args = os.Args[1:]
	c.Commands = map[string]cli.CommandFactory{
		"generate-file-number": func() (cli.Command, error) {
			return &calc.CMDGenerateFileNumber{}, nil
		},
		"calculate-file-number": func() (cli.Command, error) {
			return &calc.CMDCalculateFileNumber{}, nil
		},
		"calculate-hardcoded": func() (cli.Command, error) {
			return &calc.CMDEvalHardcoded{}, nil
		},
		"calculate-server": func() (cli.Command, error) {
			return &calc.CMDCalculateServer{}, nil
		},
		"wallet-server": func() (cli.Command, error) {
			return &wallet.CMDServer{}, nil
		},
		"wallet-generate": func() (cli.Command, error) {
			return &wallet.CMDGenerate{}, nil
		},
	}

	exitStatus, err := c.Run()
	if err != nil {
		log.Println(err)
	}

	os.Exit(exitStatus)
}
