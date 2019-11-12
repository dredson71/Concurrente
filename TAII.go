package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
)

func ReadFile2(file string, begin int, end int, y1 int, y2 int, y3 int, y4 int, y5 int) ([]float64, []int) {
	distance := make([]float64, 0)
	team := make([]int, 0)
	f, err := os.Open(file)
	if err != nil {
		log.Fatalf("Cannot open '%s': %s\n", file, err.Error())
	}
	defer f.Close()
	n := 0
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		n++
		if n < begin {
			continue
		}
		if n > end {
			break
		}
		res1 := strings.Split(scanner.Text(), ",")
		x1, _ := strconv.Atoi(strings.Replace(res1[0], "?", "0", -1))
		x2, _ := strconv.Atoi(strings.Replace(res1[1], "?", "0", -1))
		x3, _ := strconv.Atoi(strings.Replace(res1[2], "?", "0", -1))
		x4, _ := strconv.Atoi(strings.Replace(res1[3], "?", "0", -1))
		x5, _ := strconv.Atoi(strings.Replace(res1[4], "?", "0", -1))
		x6, _ := strconv.Atoi(strings.Replace(res1[5], "?", "0", -1))
		distanciaP := distanciaEuclidiana(x1, x2, x3, x4, x5, y1, y2, y3, y4, y5)
		distance = append(distance, distanciaP)
		team = append(team, x6)
	}
	fmt.Println(distance)
	fmt.Println(team)
	return (distance), team
}

func quicksort(a []float64, team []int) ([]float64, []int) {
	if len(a) < 2 {
		return a, team
	}

	left, right := 0, len(a)-1

	pivot := rand.Int() % len(a)

	a[pivot], a[right] = a[right], a[pivot]
	team[pivot], team[right] = team[right], team[pivot]

	for i, _ := range a {
		if a[i] < a[right] {
			a[left], a[i] = a[i], a[left]
			team[left], team[i] = team[i], team[left]
			left++
		}
	}

	a[left], a[right] = a[right], a[left]
	team[left], team[right] = team[right], team[left]

	quicksort(a[:left], team[:left])
	quicksort(a[left+1:], team[left+1:])

	return a, team
}

func distanciaEuclidiana(x1 int, x2 int, x3 int, x4 int, x5 int, y1 int, y2 int, y3 int, y4 int, y5 int) float64 {
	sum_cuadrados := math.Pow(float64(y1)-float64(x1), 2) + math.Pow(float64(y2)-float64(x2), 2) +
		math.Pow(float64(y3)-float64(x3), 2) + math.Pow(float64(y4)-float64(x4), 2) + math.Pow(float64(y5)-float64(x5), 2)
	resultado := math.Sqrt(sum_cuadrados)
	return resultado
}

func main() {
	slice, team := ReadFile2("mammographic_masses.csv", 1, 8, 8, 2, 0, 0, 0)
	quicksort(slice, team)
	fmt.Println(slice)
	fmt.Println(team)

}
