# sugoku
A case-study for a sudoku solver in go

## Solving Algorithm
The solving algorithm we are using is the constraint satisfaction combined with backtracking.

## Case Study
The purpose of this case study is to see how using goroutines even is simple apps can improve the performance of an app.

## Running the app
Simply add your sudoku matrix in a file such as matrix.txt and then run:
```
go run main.go -matrix-file matrix.txt 
```

Optionally add set `memory-dump` to true to see the memory usage of the app while solving the sudoku:
```
go run main.go -matrix-file matrix.txt -memory-dump true
```