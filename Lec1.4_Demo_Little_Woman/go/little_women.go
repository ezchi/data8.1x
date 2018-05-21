package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	fileName := "little_women.txt"

	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		err = getLittleWoman(fileName)

		if err != nil {
			log.Fatalf("can not get %s: %v", fileName, err)
		}
	}

	dat, err := ioutil.ReadFile(fileName)

	if err != nil {
		log.Fatalf("can not read file %s: %v", fileName, err)
	}

	fmt.Println(string(dat))
}

func getLittleWoman(path string) error {
	little_women_url := "https://www.inferentialthinking.com/chapters/01/3/little_women.txt"

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	// resp, err := client.Get("https://someurl:443/)
	res, err := client.Get(little_women_url)

	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	out, err := os.Create(path)

	if err != nil {
		return err
	}
	defer out.Close()
	io.Copy(out, res.Body)
	return nil
}
