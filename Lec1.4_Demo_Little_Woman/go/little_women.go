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

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
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

	cntPeriod := countWord(chapters, ".")
	fmt.Printf("Count \".\": %v\n", cntPeriod)

	// Plot
	p, err := plot.New()

	if err != nil {
		log.Fatalf("can not create plot: %v", err)
	}

	p.Title.Text = "Plotutil example"
	p.X.Label.Text = "X"
	p.Y.Label.Text = "Y"

	err = plotutil.AddLinePoints(p,
		"Jo", accumulateLinePoints(cntJo),
		"Meg", accumulateLinePoints(cntMeg),
		"Amy", accumulateLinePoints(cntAmy),
		"Beth", accumulateLinePoints(cntBeth),
		"Laurie", accumulateLinePoints(cntLaurie))

	if err != nil {
		log.Fatalf("can not add line points: %v", err)
	}

	// Save the plot to a PNG file.
	if err := p.Save(10*vg.Inch, 10*vg.Inch, "charactor-name-plot.png"); err != nil {
		log.Fatalf("can not save plot: %v", err)
	}

	p, err = scatterPlot(lengthVSPeriods(chapters, cntPeriod), "Chapter length vs Num of periods", "num of periods", "chapter length", "chapter-length-vs-num-of-period.png")

	if err != nil {
		log.Fatalf("can not plot chapter length vs nub of periods: %v\n", err)
	}
}

func lengthVSPeriods(c []string, p []int) plotter.XYs {
	pts := make(plotter.XYs, len(c))

	for i := range pts {
		pts[i].X = float64(p[i])
		pts[i].Y = float64(len(c[i]))
	}
	return pts
}
func scatterPlot(xys plotter.XYer, title, labelX, labelY, path string) (*plot.Plot, error) {
	p, err := plot.New()

	if err != nil {
		return nil, err
	}

	p.Title.Text = title
	p.X.Label.Text = labelX
	p.Y.Label.Text = labelY

	p.Add(plotter.NewGrid())

	s, err := plotter.NewScatter(xys)
	if err != nil {
		return p, err
	}

	p.Add(s)
	// Save the plot to a PNG file.
	if err := p.Save(4*vg.Inch, 4*vg.Inch, path); err != nil {
		return p, err
	}
	return p, nil
}

func convertToXY(d []int) plotter.XYs {
	pts := make(plotter.XYs, len(d))

	for i := range pts {
		pts[i].X = float64(i)
		pts[i].Y = float64(d[i])
	}

	return pts
}

func accumulateLinePoints(d []int) plotter.XYs {
	pts := make(plotter.XYs, len(d))
	var sum float64

	for i := range pts {
		sum += float64(d[i])

		pts[i].X = float64(i)
		pts[i].Y = sum
	}

	return pts
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
