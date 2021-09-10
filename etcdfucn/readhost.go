package etcdfucn

import (
	"bufio"
	"log"
	"os"
)

func ReadHost(filePath string) ([]string, error)  {

	// use to read file
	// return data in from to []string

	var hostname []string
	file, err := os.Open(filePath)
	if err != nil {
		log.Println("open file faile ", file)
		return nil, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		hostname = append(hostname, scanner.Text())
	}

	return  hostname, scanner.Err()
}