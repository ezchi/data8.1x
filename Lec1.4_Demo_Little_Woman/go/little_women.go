package main

import (
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
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

	chapters := splitChapter(string(dat))[1:]

	fmt.Printf("There are %d chapters.\n", len(chapters))

	cntChristmas := countWord(chapters, "Christmas")
	fmt.Printf("Count \"Chistmas\": %v\n", cntChristmas)

	cntJo := countWord(chapters, "Jo")
	fmt.Printf("Count \"Jo\": %v\n", cntJo)

	cntMeg := countWord(chapters, "Meg")
	fmt.Printf("Count \"Meg\": %v\n", cntMeg)

	cntAmy := countWord(chapters, "Amy")
	fmt.Printf("Count \"Amy\": %v\n", cntAmy)

	cntBeth := countWord(chapters, "Beth")
	fmt.Printf("Count \"Beth\": %v\n", cntBeth)

	cntLaurie := countWord(chapters, "Laurie")
	fmt.Printf("Count \"Laurie\": %v\n", cntLaurie)
}

func countWord(chapters []string, word string) []int {
	counters := make([]int, len(chapters))

	for i, c := range chapters {
		counters[i] = strings.Count(c, word)
	}

	return counters
}

func splitChapter(text string) []string {
	return strings.Split(text, "CHAPTER")
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
