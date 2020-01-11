package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var matrixFile = flag.String("matrix-file", "matrix.txt", "path to a custom matrix file")

var initialMatrix [9][9]int
var possibleValuesMatrix [10][9][9]int
var solvedMatrix [9][9]int

func init() {
	flag.Parse()
	err := loadMatrix(matrixFile)
	if err != nil {
		fmt.Printf("Could not load a valid matrix: %s", err)
		os.Exit(1)
	}

	if !isValidMatrix(&initialMatrix) {
		fmt.Printf("Loaded matrix is invalid")
		os.Exit(1)
	}

}

func main() {

	fmt.Println("Initial Sudoku")
	printMatrix(initialMatrix)

	useConstrantSatisfaction := false
	useBacktracking := true
	allowPrintMemStats := false

	if allowPrintMemStats {
		go printMemStats()
	}

	initializePossibleValuesMatrix(initialMatrix)
	iterations := 0
	sudokuIsSolved := false
	start := time.Now()
	if useConstrantSatisfaction {
		//while no new value is eliminated check for constraints
		for eliminatePossibleValues() {
			iterations++
			checkForSinglePossibleValues()
			if sudokuIsSolved = isSudokuSolved(possibleValuesMatrix[0]); sudokuIsSolved {
				break
			}
		}
	}

	//start backtracking
	if !sudokuIsSolved && useBacktracking {
		sudokuIsSolved = backtrack(&possibleValuesMatrix[0])
	}
	duration := time.Since(start)

	solvedMatrix = possibleValuesMatrix[0]
	sudokuStatus := "not solved"
	if sudokuIsSolved {
		sudokuStatus = "solved"
	}
	fmt.Printf("\nSudoku %s with %d Iterations & %d backtracks within %d ns\n", sudokuStatus, iterations, totalBacktracks, duration.Nanoseconds())
	printMatrix(solvedMatrix)
}

func printMemStats() {
	for {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		fmt.Printf("Alloc = %v KiB", bToKb(m.Alloc))
		fmt.Printf("\tTotalAlloc = %v KiB", bToKb(m.TotalAlloc))
		fmt.Printf("\tSys = %v KiB", bToKb(m.Sys))
		fmt.Printf("\tNumGC = %v\n", m.NumGC)
		time.Sleep(time.Second)
	}
}

func bToKb(b uint64) uint64 {
	return b / 1024
}

var totalBacktracks int

func backtrack(matrix *[9][9]int) bool {
	totalBacktracks++
	if !hasEmptyCell(matrix) {
		return true
	}
	for col := 0; col < 9; col++ {
		for row := 0; row < 9; row++ {
			if matrix[col][row] == 0 {
				for val := 9; val >= 1; val-- {
					matrix[col][row] = val
					if hasAcceptableValue(col, row, *matrix) {
						if backtrack(matrix) {
							return true
						}
						matrix[col][row] = 0
					} else {
						matrix[col][row] = 0
					}
				}
				return false
			}
		}
	}
	return false
}

func checkForSinglePossibleValues() {
	for col := 0; col < 9; col++ {
		for row := 0; row < 9; row++ {
			eligibleValues := 0
			eligibleValue := 0
			for val := 1; val <= 9; val++ {
				if possibleValuesMatrix[val][row][col] != 0 {
					eligibleValues++
					eligibleValue = val
				}
			}
			if eligibleValues == 1 {
				possibleValuesMatrix[0][row][col] = eligibleValue
			}
		}
	}
}

func eliminatePossibleValues() (eliminated bool) {
	eliminated = false
	for col := 0; col < 9; col++ {
		for row := 0; row < 9; row++ {
			if possibleValuesMatrix[0][row][col] == 0 {
				for possibleVal := 1; possibleVal <= 9; possibleVal++ {
					if possibleValuesMatrix[possibleVal][row][col] != 0 {
						possibleValuesMatrix[0][row][col] = possibleVal
						if !isValidMatrix(&possibleValuesMatrix[0]) {
							eliminated = true
							possibleValuesMatrix[possibleVal][row][col] = 0
						}
						possibleValuesMatrix[0][row][col] = 0
					}
				}
			}
		}
	}

	return
}

func isSudokuSolved(matrix [9][9]int) bool {
	for col := 0; col < 9; col++ {
		for row := 0; row < 9; row++ {
			if matrix[row][col] == 0 {
				return false
			}
		}
	}
	return true
}

func loadMatrix(filename *string) error {
	filebuffer, err := ioutil.ReadFile(*filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	inputdata := string(filebuffer)
	data := bufio.NewScanner(strings.NewReader(inputdata))
	data.Split(bufio.ScanRunes)

	insRow, insCol := 0, 0
	for data.Scan() {
		if data.Text() == "\n" || data.Text() == " " {
			continue
		}
		if insRow == 9 || insCol == 9 {
			break
		}

		elem, err := strconv.Atoi(data.Text())
		if err != nil {
			return err
		}
		initialMatrix[insRow][insCol] = elem
		if insCol == 8 {
			insCol = 0
			insRow++
		} else {
			insCol++
		}
	}
	return nil
}

func isValidMatrix(matrix *[9][9]int) bool {
	for col := 0; col < 9; col++ {
		for row := 0; row < 9; row++ {
			if !hasAcceptableValue(row, col, *matrix) {
				return false
			}
		}
	}
	return true
}

func hasAcceptableValue(row int, col int, matrix [9][9]int) bool {
	if matrix[row][col] == 0 {
		return true
	}
	// Row search
	for searchRow := 0; searchRow < 9; searchRow++ {
		if searchRow != row && matrix[searchRow][col] == matrix[row][col] && matrix[searchRow][col] != 0 {
			return false
		}
	}
	// Col search
	for searchCol := 0; searchCol < 9; searchCol++ {
		if searchCol != col && matrix[row][searchCol] == matrix[row][col] && matrix[row][searchCol] != 0 {
			return false
		}
	}

	// Neighbor search
	rowNeighborMin, rowNeighborMax := getMinMaxNeighbor(row)
	colNeighborMin, colNeighborMax := getMinMaxNeighbor(col)

	for searchCol := colNeighborMin; searchCol <= colNeighborMax; searchCol++ {
		for searchRow := rowNeighborMin; searchRow <= rowNeighborMax; searchRow++ {
			if searchCol != col && searchRow != row && matrix[searchRow][searchCol] == matrix[row][col] && matrix[searchRow][searchCol] != 0 {
				return false
			}
		}
	}

	return true
}

func getMinMaxNeighbor(idx int) (int, int) {
	if idx <= 2 {
		return 0, 2
	} else if idx <= 5 {
		return 3, 5
	} else {
		return 6, 8
	}
}

func initializePossibleValuesMatrix(matrix [9][9]int) {
	for col := 0; col < 9; col++ {
		for row := 0; row < 9; row++ {
			possibleValuesMatrix[0][row][col] = matrix[row][col]
			if matrix[row][col] == 0 {
				for val := 0; val < 9; val++ {
					possibleValuesMatrix[val+1][row][col] = val + 1
				}
			} else {
				for val := 0; val < 9; val++ {
					possibleValuesMatrix[val+1][row][col] = 0
				}
				possibleValuesMatrix[matrix[row][col]][row][col] = matrix[row][col]
			}
		}
	}
}

func hasEmptyCell(matrix *[9][9]int) bool {
	for col := 0; col < 9; col++ {
		for row := 0; row < 9; row++ {
			if matrix[row][col] == 0 {
				return true
			}
		}
	}
	return false
}

func printMatrix(matrix [9][9]int) {
	fmt.Println("+-------+-------+-------+")
	for row := 0; row < 9; row++ {
		fmt.Print("| ")
		for col := 0; col < 9; col++ {
			if col == 3 || col == 6 {
				fmt.Print("| ")
			}
			if matrix[row][col] == 0 {
				fmt.Printf("_ ")
			} else {
				fmt.Printf("%d ", matrix[row][col])
			}
			if col == 8 {
				fmt.Print("|")
			}
		}
		if row == 2 || row == 5 || row == 8 {
			fmt.Println("\n+-------+-------+-------+")
		} else {
			fmt.Println()
		}
	}
}
