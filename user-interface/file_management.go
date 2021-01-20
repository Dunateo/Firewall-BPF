package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

/**
GET all the file in a string
*/
func read_file(fileName string) string {
	//check if the file exist and handle an error
	data, err := ioutil.ReadFile(fileName)
	check(err)

	fmt.Print(string(data))

	return string(data)
}

/**
CREATE an File
*/
func CreateFile(fileName string) {
	// check if file exists
	var _, err = os.Stat(fileName)

	// create file if not exists
	if os.IsNotExist(err) {
		var file, err = os.Create(fileName)
		check(err)
		defer file.Close()
	}

	fmt.Println("File Created ", fileName)
}

/**
WRITE in a file
filePath STRING
Content []String
*/
func WriteFile(fileName string, contents []string) {
	// Open file using READ & WRITE permission.
	var file, err = os.OpenFile(fileName, os.O_RDWR, 0644)
	check(err)
	//at the end of operation it will close
	defer file.Close()

	//all the content is writed in the file
	for _, content := range contents {
		_, err = file.WriteString(content + "\n")
		check(err)
	}

	// Save file changes.
	err = file.Sync()
	check(err)

	fmt.Println("File Updated.")
}

func WriteLine(fileName string, content string) {
	// Open file using READ & WRITE permission.
	var file, err = os.OpenFile(fileName, os.O_APPEND|os.O_WRONLY, 0644)
	check(err)
	//at the end of operation it will close
	defer file.Close()

	//content is added in the File
	_, err = file.WriteString(content)
	check(err)

	// Save file changes.
	err = file.Sync()
	check(err)

	fmt.Println("File Line Updated.")
}

/**
DELETE an file
*/
func DeleteFile(fileName string) {
	// delete file
	var err = os.Remove(fileName)
	check(err)
	fmt.Println("File Deleted")
}

/**
Search for duplicata
*/
func doublonPort(fileName string, numport string) bool {
	//open file
	file, err := os.Open(fileName)
	check(err)
	//at the end of operation it will close
	defer file.Close()

	// scan and stock the file in a buffer
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var txtlines []string
	for scanner.Scan() {
		txtlines = append(txtlines, scanner.Text())
	}

	//iterate on each line
	for _, eachline := range txtlines {

		//fmt.Println(eachline)
		if strings.Compare(eachline, numport) == 0 {
			//fmt.Println(eachline)
			//return true when there is a duplicata
			return true
		}

	}
	//return false when there is no duplicata
	return false
}

func fileToSlice(fileName string) []string {
	//open file
	file, err := os.Open(fileName)
	check(err)
	//at the end of operation it will close
	defer file.Close()

	// scan and stock the file in a buffer
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var txtlines []string
	for scanner.Scan() {
		txtlines = append(txtlines, scanner.Text())
	}
	return txtlines
}
