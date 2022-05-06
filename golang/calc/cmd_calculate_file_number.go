package calc

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/mitchellh/cli"
	"io"
	"log"
	"math"
	"math/big"
	"os"
	"strings"
)

type CMDCalculateFileNumber struct{}

var _ cli.Command = (*CMDCalculateFileNumber)(nil)

func (c *CMDCalculateFileNumber) Help() string {
	return "Calculate line by line file content in sequential mode."
}

func (c *CMDCalculateFileNumber) Run(args []string) int {
	if len(args) != 1 {
		log.Printf("must followed one exact argument contain file name, number of args: %d\n", len(args))
		return -1
	}

	fileName := args[0]
	finalCalculated, oneLineContent, err := c.calculateLineByLine(fileName)
	if err != nil {
		log.Println(err)
		return -1
	}

	log.Printf("one line content: %s\n", oneLineContent)
	log.Printf("calculated value inside file are: %s\n", finalCalculated)
	return 0
}

func (c *CMDCalculateFileNumber) Synopsis() string {
	return "Calculate line by line file content in sequential mode."
}

// calculateLineByLine return final calculated value, one line string of all mathematical operation and optional error.
func (c *CMDCalculateFileNumber) calculateLineByLine(fileName string) (string, string, error) {
	f, err := os.OpenFile(fileName, os.O_RDONLY, os.ModePerm)
	defer func() {
		if f == nil {
			return
		}

		if _err := f.Close(); _err != nil {
			_err = fmt.Errorf("error read line by line file %s: %w", fileName, _err)
			log.Println(_err)
		}
	}()

	if err != nil {
		err = fmt.Errorf("open file error: %w", err)
		return "", "", err
	}

	finalValBigInt := big.NewInt(0)
	rd := bufio.NewReader(f)

	oneline := &bytes.Buffer{}
	lineNum := 0
	exceedMaxInt := false
	for {
		lineNum++
		line, err := rd.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}

			err = fmt.Errorf("read file line error: %w", err)
			return "", "", err
		}

		// parse line string
		line = strings.TrimSpace(line)
		op := line[0]
		val := line[1:]

		currValBigInt := new(big.Int)
		_, err = fmt.Sscan(val, currValBigInt)
		if err != nil {
			err = fmt.Errorf("cannot parse string to int %s line number %d: %w", line[1:], lineNum, err)
			return "", "", err
		}

		if finalValBigInt.Cmp(big.NewInt(math.MaxInt)) >= 0 && !exceedMaxInt {
			exceedMaxInt = true
			log.Printf("current value larger than int after operation on line %d, this is why we need bigInt (%s vs %d)\n",
				lineNum, finalValBigInt.String(), math.MaxInt,
			)
		}

		switch op {
		case '*':
			finalValBigInt = finalValBigInt.Mul(finalValBigInt, currValBigInt)
		case '/':
			finalValBigInt = finalValBigInt.Div(finalValBigInt, currValBigInt)
		case '+':
			finalValBigInt = finalValBigInt.Add(finalValBigInt, currValBigInt)
		case '-':
			finalValBigInt = finalValBigInt.Sub(finalValBigInt, currValBigInt)
		default:
			return "", "", fmt.Errorf("unknown operator %s on line number %d\n", string(op), lineNum)
		}

		oneline.WriteString(line)
	}

	return finalValBigInt.String(), oneline.String(), nil
}
