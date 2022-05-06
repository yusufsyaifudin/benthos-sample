package wallet

import (
	"fmt"
	"github.com/jessevdk/go-flags"
	"github.com/mitchellh/cli"
	"log"
	"os"
)

type CMDGenerate struct{}

var _ cli.Command = (*CMDGenerate)(nil)

func (c *CMDGenerate) Help() string {
	return "Generate wallet pattern data"
}

func (c *CMDGenerate) Run(args []string) int {
	type Argument struct {
		Output    string `short:"o" long:"out" description:"Output file"`
		WalletID  string `long:"wallet-id" description:"Wallet ID"`
		KafkaMeta bool   `long:"kafka-meta" description:"Kafka meta data included or not"`
	}

	var argsVal Argument
	args, err := flags.ParseArgs(&argsVal, args)
	if err != nil {
		err = fmt.Errorf("failed parsing flag: %w", err)
		log.Println(err)
		return -1
	}

	err = c.numberGenerator(argsVal.Output, argsVal.WalletID, argsVal.KafkaMeta, 1001)
	if err != nil {
		log.Println(err)
		return -1
	}

	log.Printf("success generate file on %s\n", argsVal.Output)
	return 0
}

func (c *CMDGenerate) Synopsis() string {
	return "Generate wallet pattern data"
}

func (c *CMDGenerate) numberGenerator(dir string, walletID string, withKafkaMeta bool, size int) (err error) {

	// create if not exist
	file, err := os.Create(dir)
	defer func() {
		if file == nil {
			return
		}

		if _err := file.Close(); _err != nil {
			_err = fmt.Errorf("closing file '%s' error: %w", dir, err)
			log.Println(err)
			return
		}
	}()

	if err != nil {
		err = fmt.Errorf("error create file '%s': %w", dir, err)
		return
	}

	amt := 0
	for i := 0; i < size; i++ {
		amt += 10

		op := "topup"
		if amt == 110 {
			amt = 11 // spend amount always same
			op = "spend"
		}

		data := fmt.Sprintf(`{"walletID":"%s", "operation": "%s", "amount": %d}`,
			walletID,
			op,
			amt,
		)
		if withKafkaMeta {
			data = fmt.Sprintf(`{"kafkaPartition": "%s", "kafkaOffset": "%d", "walletID":"%s", "operation": "%s", "amount": %d}`,
				"1",
				i+1,
				walletID,
				op,
				amt,
			)
		}

		if op == "spend" {
			amt = 0 // reset to topup amount
		}

		data = fmt.Sprintf("%s\n", data) // add new line
		_, err = file.WriteString(data)
		if err != nil {
			err = fmt.Errorf("error write file '%s': %w", dir, err)
			return
		}
	}

	return
}
