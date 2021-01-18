package file_management

import (
	"fmt"
	"io/ioutil"
	"os"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

/**
GET all the file in a string
 */
func read_file(fileName string) string{
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
	for _,content := range contents {
		_, err = file.WriteString(content)
		check(err)
	}

	// Save file changes.
	err = file.Sync()
	check(err)

	fmt.Println("File Updated.")
}

func writeLine(fileName string, content string) {
	// Open file using READ & WRITE permission.
	var file, err = os.OpenFile(fileName, os.O_RDWR, 0644)
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





