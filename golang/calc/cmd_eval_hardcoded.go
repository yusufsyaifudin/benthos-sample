package calc

import (
	"github.com/mitchellh/cli"
	"log"
)

type CMDEvalHardcoded struct{}

var _ cli.Command = (*CMDEvalHardcoded)(nil)

func (c *CMDEvalHardcoded) Help() string {
	return "Hardcoded evaluation of value"
}

func (c *CMDEvalHardcoded) Run(args []string) int {
	valueToCalc := +100 + 28 - 14*3 + 27 - 19 + 21*1 - 8 - 9 + 13 - 18 + 8 + 17*1 + 9 + 21 + 7*2*3 + 3*3/1 - 11/2 + 15/1/2 - 3 + 28/2/2*2*3*2/1/2*2*1 + 39/1*1/1/1 + 8/1 - 12*2 - 7 + 44*3 - 16*1*2 - 3 + 39 + 5 - 4*1/2*3*2*3*2/1/1/1/1/1 + 45/2 - 1 + 27*3 - 3 - 3*3 + 35*3/2 - 8 - 3*1*2 + 5 + 15/2/2/1*3 + 11/2/2*3*3/1*1/1*3/1/1*1 + 26*3 - 1 - 19 - 16 - 7/2 - 6 + 39 + 10 + 41/2*1/1 - 18*1 + 21/1*3/1 - 15 - 9/1/1 + 30 + 40 + 15 + 47*1 - 10/2/1/2/1/2 + 29/2 + 25 + 36 + 1 - 8 - 1*3/1/2*3/1/1 + 36/2/2*2 + 36*3/1 + 6 + 22*2 - 9*3 - 17 + 20/2 - 8 - 17 - 15 + 30 + 10 + 26 - 9*2/1 + 12*3 - 11 - 7 + 16*3/2*3 - 3 - 14 - 16*3/1 - 14*3 + 40 - 1 + 35*3*3 + 45 - 16 + 27/2 - 5 - 4/1*1 - 17 - 12*2*3/1 - 3*1 + 36 - 19 - 1*3/2*2/1/2 + 48 - 9/1 - 12 + 28 - 18 + 17 + 37*3/1/2/2 + 4 - 14 + 2 + 17/1*3/1*1/1 + 49/2*3*2 - 6*2 - 6*2/2 + 1/1/2 + 41/2*2 + 16*3 - 18/2 - 4 + 34/2 + 39 + 16/2 - 13 - 13 - 7/2 + 11 + 26 - 13 + 37*3*2/1 - 14 - 19*3*1 - 6/1 - 18*1*2 - 14 - 17/2*3 - 13 - 7*1/1*1 + 46/2*2 + 11*3/2 + 5 - 1/1/2 + 37/2 + 30 - 18/1 - 13 - 19/1/1*3 - 4*2 + 49*3/2/2 - 13 + 18 - 18 + 15/2 - 9*1 - 9/1*2*3*2 - 12 + 24 - 1*1*2 + 29 + 41 + 10 - 2 - 19/1 - 12*3 - 9/1 - 10 - 19/2*2*3*1 - 3 - 10*3/2 + 10 - 15*2 - 1 + 3 + 32*1/2*2/1 + 37 - 12 - 2 - 3*2 + 6 + 11 - 3/1 - 15/1 + 32 - 14 + 4/2 - 3*2*3 - 19*3 + 20 + 16/1/2*3*1 + 30/1 + 38 + 29/1 + 45*1/2 - 3*2 + 33 - 13 - 13 + 38*1 - 4 - 1*3 + 36 - 6 + 37/1 + 25 - 16/1 + 30 + 22*1 - 2*1/2/2/1/2 - 17/1 + 43 + 9*1*1*3/2 + 7/2/2*2/1/1 - 10/2 - 4/1 - 11 - 10*3 + 47 + 19 + 16/2/2 + 40*1*2 + 1 + 9*3 + 10 + 9 + 45 + 39 + 2/1 + 2 + 12 - 5/1 + 30 - 18 + 21*1*1 + 33 + 32 - 12 + 44 - 5 + 4*3 - 8 - 16*2 + 42*1 - 16 + 39 - 16/1 + 4/2/1 + 3 + 41 + 36 - 14 + 30 - 19 - 19*1*2 + 22*2/2*3 + 35*2*1 + 5*1 - 17 + 25/1 + 26 + 29/1 - 9 - 7 - 10*2/2 - 15 - 17 - 16*1*2 + 48*3/1 + 44*3 - 14 + 20/1/2 - 6/1 + 4*2*3*1 + 45*2 + 21 - 2/1 + 19/1 + 47*1/2 + 40/2*1/2/1*1/2/1 + 37/2 - 16/1*1 + 46 + 1/2 - 19 - 4*3/2*3*1 + 22 + 31 + 42*1/1 - 6/1/1*1/2 + 20*1*3 + 20 + 24*2 + 28/2 + 42 - 9/2 - 7 + 30/2 + 44/2*3 - 4/2/2/1 + 34 + 24/2*2*3 - 11/2 - 15 + 46 + 2/1/2/1 - 16*1 - 16 - 15 + 19 + 40/2 - 12*1 - 7*3 - 18 - 13/1 + 9/2 - 15 + 47 + 15 - 19 - 6*2 - 8 - 1*2 + 16*2 + 2 - 4 + 2 + 36/2*2 + 29*3*3/1/1/1*1*3*3/2 + 13*1*2 + 1 + 17 + 45/2/1/1*1/2/2 - 4 - 7 - 11 - 9 + 13*1 + 27*3/1 + 15 + 30/1 - 19 + 39 - 1/1 - 17/1*1/2/1*1 + 24 + 8/1 + 5*3/2 - 6/1 - 15/2/1 - 9*3*2 - 14*1/2/2 - 10/2/1/2 + 43 + 40 - 17*2/1*1 + 10 - 14 + 30 + 46*3 - 2*3 - 19/1 - 16 - 12 + 18 + 41/1*3/1/2/1 - 6 + 15 + 16*2 + 24/2*1/2/2 - 8*1*1*1/1/1*3/2 + 10 - 8 - 14/2*2*3 + 48/2 + 23/1 - 17*3 + 36 - 13/2 + 14 + 13 - 13 + 45 + 34*1/1*2 + 4 - 7 - 19 + 4/1 + 21 + 15*3*2*1/1 + 6*3/2/2*2*2 - 18 - 8*3 + 28 + 40 - 17*3 - 8 - 19 - 3*1*2 - 14*2 - 3/2 - 14 - 13*2*3 - 2 - 9/2 + 8*1 - 16 + 48 + 44*2 + 5*1 - 6 - 19/2 - 11/1/2*3 + 42/1*2*1 + 38/1 - 14*2 - 16 - 19*3 + 27 + 18/1 + 11/2*2/1 - 3 + 49*3*1/2 - 11 - 19*2/1*1*1*2/2*1 - 2 - 9 - 10/2 - 4 + 1 - 15*3/1*2*1*1 - 4/1/1*3 - 6*3/2 + 44*1*2 + 47 + 24/1 - 10/1*3/2*1/1 + 49/2*3 + 22 - 19*2 - 16*2*1 + 48 + 21/2/2/1 - 18 - 1 - 9 + 29 + 23 - 18/1 + 24 - 13 - 4*1 - 4*3 + 10*2 - 14/2*1 - 9/1 - 13 + 10*2/1*2 + 32 + 34/1*1 - 4 + 46/1 - 18*2 + 26/2 + 30 + 11*2 + 13/2*3/2 - 9 - 10 - 11*2 - 11 - 7 - 9 + 4 + 27*2*1*3 - 4 + 12*1 + 9 - 6/2*1/2/1 - 7 + 4/2 - 5/2 - 11 - 8 + 45 - 15 + 34 - 19/1 + 44 + 15 - 2 + 36*1*1 + 13*2/1 + 27/2 + 40 - 16

	log.Printf("calculated value: %d\n", valueToCalc)
	return 0
}

func (c *CMDEvalHardcoded) Synopsis() string {
	return "Hardcoded evaluation of value"
}