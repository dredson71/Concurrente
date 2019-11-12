package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

func ReadFile2(file string) {
	f, err := os.Open(file)
	if err != nil {
		log.Fatalf("Cannot open '%s': %s\n", file, err.Error())
	}
	defer f.Close()
	n := 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		n++
		if n < 2 {
			continue
		}
		if n > 4 {
			break
		}
		res1 := strings.Split(scanner.Text(), ",")
		x1, _ := strconv.Atoi(strings.Replace(res1[0], " ", "", -1))
		x2, _ := strconv.Atoi(strings.Replace(res1[1], " ", "", -1))
		x3, _ := strconv.Atoi(strings.Replace(res1[2], " ", "", -1))
		x4, _ := strconv.Atoi(strings.Replace(res1[3], " ", "", -1))
		x5, _ := strconv.Atoi(strings.Replace(res1[4], " ", "", -1))
		println(x1)
		println(x2)
		println(x3)
		println(x4)
		println(x5)

	}
}

func ReadFile(file string) [][]string {
	f, err := os.Open(file)
	if err != nil {
		log.Fatalf("Cannot open '%s': %s\n", file, err.Error())
	}
	defer f.Close()
	r := csv.NewReader(f)
	r.Comma = ','
	rows, err := r.ReadAll()
	if err != nil {
		log.Fatalln("Cannot read CSV data:", err.Error())
	}
	return rows
}

func calculateDistance(rows [][]string, y1 int, y2 int, y3 int, y4 int, y5 int) {
	distance := make([]float64, 0)
	for i := range rows {
		if i != 0 {
			x1, _ := strconv.Atoi(strings.Replace(rows[i][0], " ", "", -1))
			x2, _ := strconv.Atoi(strings.Replace(rows[i][1], " ", "", -1))
			x3, _ := strconv.Atoi(strings.Replace(rows[i][2], " ", "", -1))
			x4, _ := strconv.Atoi(strings.Replace(rows[i][3], " ", "", -1))
			x5, _ := strconv.Atoi(strings.Replace(rows[i][4], " ", "", -1))
			distanciaP := distanciaEuclidiana(x1, x2, x3, x4, x5, y1, y2, y3, y4, y5)
			i++
			distance = append(distance, distanciaP)
		}
	}
	fmt.Println(distance[0])

}

func distanciaEuclidiana(x1 int, x2 int, x3 int, x4 int, x5 int, y1 int, y2 int, y3 int, y4 int, y5 int) float64 {
	sum_cuadrados := math.Pow(float64(y1)-float64(x1), 2) + math.Pow(float64(y2)-float64(x2), 2) +
		math.Pow(float64(y3)-float64(x3), 2) + math.Pow(float64(y4)-float64(x4), 2) + math.Pow(float64(y5)-float64(x5), 2)
	resultado := math.Sqrt(sum_cuadrados)
	return resultado
}

func main() {
	/*	rows := ReadFile("prueba.csv")
		calculateDistance(rows, 8, 2, 0, 0, 0)*/
	ReadFile2("prueba.csv")

}
