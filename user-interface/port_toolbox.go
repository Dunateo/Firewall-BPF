package main

import (
	"bufio"
	"log"
	"os"
)


/**
Delete a port in the file
 */
func delete_port( fileName string , port string ) []string {
	// open the file
	file, err := os.Open(fileName)
	check(err)

	//update file
	DeleteFile(fileName)
	CreateFile(fileName)

	//line by line scanner
	fileScanner := bufio.NewScanner(file)

	var count int = 0
	var contents []string
	// read line by line
	for fileScanner.Scan() {
		if fileScanner.Text() != port {
			contents = append(contents,fileScanner.Text())
			count++
		}

	}
	// handle first encountered error while reading
	if err := fileScanner.Err(); err != nil {
		log.Fatalf("Error while reading file: %s", err)
	}
	file.Close()

	//write in the file without port
	WriteFile(fileName, contents)
	return contents
}

/**
add a port to the file
 */
func AddPort(fileName string,port string)  {
	WriteLine(fileName, port)
}