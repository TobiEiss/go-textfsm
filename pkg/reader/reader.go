package reader

import (
	"bufio"
	"os"
)

// ReadLineByLine reads a file line by line with a channel
func ReadLineByLine(path string, out chan string) {
	// open file
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	// read with scanner
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	defer file.Close()

	// read line by line
	for scanner.Scan() {
		out <- scanner.Text()
	}
	close(out)
}
