package calc

import (
	"bytes"
	"fmt"
	"github.com/mitchellh/cli"
	"log"
	"math/rand"
	"os"
	"time"
)

type CMDGenerateFileNumber struct{}

var _ cli.Command = (*CMDGenerateFileNumber)(nil)

func (c *CMDGenerateFileNumber) Help() string {
	return "Generate mathematical operator number on file"
}

func (c *CMDGenerateFileNumber) Run(args []string) int {
	if len(args) != 1 {
		log.Printf("must followed one exact argument contain file name, number of args: %d\n", len(args))
		return -1
	}

	fileName := args[0]
	err := c.numberGenerator(fileName, 1000)
	if err != nil {
		log.Println(err)
		return -1
	}

	log.Printf("success generate file on %s\n", fileName)
	return 0
}

func (c *CMDGenerateFileNumber) Synopsis() string {
	return "Generate mathematical operator number on file"
}

func (c *CMDGenerateFileNumber) numberGenerator(dir string, size int) (err error) {
	rand.Seed(time.Now().UnixNano())

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

	for i := 0; i < size; i++ {
		if i == 0 {
			_, err = file.WriteString(fmt.Sprintf("%s%d\n", "+", 100))
			if err != nil {
				err = fmt.Errorf("error write file '%s': %w", dir, err)
				return
			}
			continue
		}

		operator := c.randomMathOperator()
		value := c.randomValueNonZero(operator)
		_, err = file.WriteString(fmt.Sprintf("%s%d\n", operator, value))
		if err != nil {
			err = fmt.Errorf("error write file '%s': %w", dir, err)
			return
		}
	}

	return
}

func (c *CMDGenerateFileNumber) randomMathOperator() string {
	const mathOp = "*/+-"
	var buffer bytes.Buffer
	for i := 0; i < 1; i++ {
		buffer.WriteString(string(mathOp[rand.Intn(len(mathOp))]))
	}
	return buffer.String()
}

func (c *CMDGenerateFileNumber) randomValueNonZero(op string) int {
	value := rand.Intn(100)
	switch op {
	case "*":
		value = rand.Intn(4) // maximum 4 times
	case "/":
		value = rand.Intn(3) // if operator is divide, maximum divide by 3
	case "+":
		value = rand.Intn(50)
	case "-":
		value = rand.Intn(20)
	}

	// never return 0
	if value == 0 {
		return c.randomValueNonZero(op)
	}

	return value
}
