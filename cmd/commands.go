package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/tecnologer/sudoku/clients/cli/game"
	"github.com/tecnologer/sudoku/clients/cli/utils"
	sudoku "github.com/tecnologer/sudoku/src"
)

var (
	commandsMap map[string]func(...string)
	commands    []*command
)

func init() {
	commands = []*command{
		{
			cmd:    "help",
			action: printHelp,
			about:  "Prints the list of commands available",
			alias:  []string{"?", "h"},
		},
		{
			cmd:    "exit",
			action: exit,
			about:  "Closes the game",
			alias:  []string{"close", "q"},
		},
		{
			cmd:    "new",
			action: newGame,
			about:  "Starts new game",
			alias:  []string{},
		},
		{
			cmd:    "set",
			action: setValue,
			about:  "Sets the value in the coordinate",
			alias:  []string{},
		},
		{
			cmd:    "display",
			action: printGame,
			about:  "Sets the value in the coordinate",
			alias:  []string{"print", "show", "p"},
		},
		{
			cmd:    "validate",
			action: validateGame,
			about:  "Validates if the game is correct",
			alias:  []string{"v"},
		},
		{
			cmd:    "solve",
			action: solveGame,
			about:  "Solves the current game",
			alias:  []string{},
		},
		{
			cmd:    "set-row",
			action: setRow,
			about: `Sets the values to the specified row.
		* Single line: set-row <row_number>,<data separated by comma>
			+ Example: row number 3
				>> set-row 3,0,0,1,0,3,0,0,0,6
		* Different line: <row_number> <data separated by comma>
			+ Example: row number 3
				>> set-row 
				>> 3 0,0,1,0,3,0,0,0,6`,
			alias: []string{"row"},
		},
	}

	commandsMap = make(map[string]func(...string))
	for _, cmd := range commands {
		commandsMap[cmd.cmd] = cmd.action
		for _, alias := range cmd.alias {
			if _, exists := commandsMap[alias]; exists {
				continue
			}

			commandsMap[alias] = cmd.action
		}
	}
}

func printHelp(args ...string) {
	if len(args) == 0 || (len(args) == 1 && args[0] == "") {
		fmt.Println("\nAvailable commands: ")
		for _, cmd := range commands {
			fmt.Println("\t", cmd)
		}
	} else {
		fmt.Printf("\nHelp for: %v\n", args)
		for _, arg := range args {
			for _, cmd := range commands {
				if cmd.isCommand(arg) {
					fmt.Println("\t", cmd)
				}
			}
		}
	}
	fmt.Println()
}

func exit(args ...string) {
	i := "y"
	// i := readInput("Are you sure? (Y/n): ")

	if strings.ToLower(i) == "y" {
		os.Exit(0)
	}
}

func newGame(args ...string) {
	var levelStr string
	if len(args) == 0 {
		fmt.Println("Difficulties availables: ")
		for _, c := range sudoku.GetComplexities() {
			fmt.Printf("\t- %s\n", c)
		}

		fmt.Println()
		levelStr, args = readInput("Choose you difficulty: ")
	} else {
		levelStr = strings.Trim(args[0], "")
	}

	level := sudoku.StringToComplexity(levelStr)
	if level == sudoku.InvalidLevel {
		fmt.Println("Invalid difficulty selected")
		newGame()
		return
	}

	game.Current = sudoku.NewGame(level)

	fmt.Printf("new %s game started\n", level)
	CallCmd("print", args...)
}

func setValue(args ...string) {
	var inputs []string
	input := ""
	for {

		if len(args) != 3 {
			input, args = readInput("Type the coordinate and the value separate by commas (x,y,z) [\"cancel\" to try other command]: ")
			if input == "cancel" {
				fmt.Println("Canceled")
				return
			}
			inputs = strings.Split(input, ",")
			if len(inputs) != 3 {
				fmt.Println("Invalid data. Please, try again")
				continue
			}
		} else {
			inputs = args
		}

		row, err := strconv.Atoi(inputs[0])
		if err != nil || row > 9 || row < 1 {
			fmt.Println("The value for the row is invalid. Should be a integer between 1 and 9")
			continue
		}

		col, err := strconv.Atoi(inputs[1])
		if err != nil || col > 9 || col < 1 {
			fmt.Println("The value for the column is invalid. Should be a integer between 1 and 9")
			continue
		}

		val, err := strconv.Atoi(inputs[2])
		if err != nil || val > 9 || val < 1 {
			fmt.Println("The value is invalid. Should be a integer between 1 and 9")
			continue
		}

		x, y := row-1, col-1
		if game.Current.IsCoordinateLockedXY(x, y) {
			fmt.Println("The coordinate is locked, please select other")
			setValue()
			return
		}

		game.Current.Set(x, y, val)
		break
	}

}

func printGame(args ...string) {
	if game.Current == nil {
		fmt.Println("there is not started game")
		return
	}
	for x := 0; x < 9; x++ {
		fmt.Print("|")
		for y := 0; y < 9; y++ {
			if !game.Current.IsEmpty(x, y) {
				fmt.Print(game.Current.Get(x, y))
			} else {
				fmt.Print("-")
			}
			fmt.Print(" ")

			if (y+1)%3 == 0 {
				fmt.Print("|")
			}
		}
		if (x+1)%3 == 0 && x < 8 {
			fmt.Println("\n|------|------|------|")
		} else {
			fmt.Println()
		}
	}
}

func validateGame(args ...string) {
	errs := game.Current.Validate()

	if errs.Count == 0 {
		fmt.Println("everything is ok")
		return
	}

	p := ""
	for t, es := range errs.Errs {
		if len(es) > 1 {
			p = "s"
		} else {
			p = ""
		}

		fmt.Printf("Error in the following %s%s:\n", t, p)
		for _, e := range es {
			fmt.Printf("\tX: %d, Y: %d\n", e.X+1, e.Y+1)
		}
	}
}

func solveGame(args ...string) {
	game.Current.Solve()
	CallCmd("p")
}

func setRow(args ...string) {
	var row int
	var data [9]int
	var err error
	if len(args) != 10 {
		var rowStr string
		rowStr, args = readInput("Type the values with the following format. <row> <data separated by comma>")

		args = append([]string{rowStr}, args...)
	}

	row, err = strconv.Atoi(args[0])
	if err != nil {
		fmt.Println("The value for the row is invalid. Should be a integer between 1 and 9")
		setRow()
		return
	}

	array, err := utils.ParseStringArrayToIntArray(args[1:])
	if err != nil {
		fmt.Println("The data to set in the error is not valid. Should be a integer between 1 and 9")
		setRow()
		return
	}

	copy(data[:], array[:9])

	game.Current.SetDataRow(row-1, data)
}
