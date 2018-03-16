package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strconv"
)

type Price struct {
	Index int
	Price float64
	Open  float64
	High  float64
	Low   float64
}

var priceList []Price
var constanta [3]float64
var constGrada [3]float64

func AutoRegression(hari int, constants [3]float64) float64 {

	ret := priceList[hari-1].Open*constants[0] +
		priceList[hari-2].Open*constants[1] +
		priceList[hari-3].Open*constants[2]

	// fmt.Println(ret)

	return ret

}

func FindGrad(constants [3]float64, h float64) [3]float64 {

	var constGrad [3]float64

	for i := 0; i < 3; i++ {
		temp := constants[i]
		constants[i] = temp + h
		MAE1 := MAE(constants)
		constants[i] = temp - h
		MAE2 := MAE(constants)
		constGrad[i] = (MAE1 - MAE2) / (2 * h)
	}

	return constGrad
}

func MAE(constants [3]float64) float64 {

	var ret float64
	ret = 0

	for i := 3; i < len(priceList); i++ {
		ret += math.Abs(priceList[i].Open - AutoRegression(i, constants))
	}

	return ret
}

func FindGradNorm(constGrad [3]float64) float64 {

	var temp float64

	for i := range constGrad {
		temp += (constGrad[i] * constGrad[i])
	}

	return math.Sqrt(temp)
}

func readFile() {

	csvFile, _ := os.Open("data")

	reader := csv.NewReader(bufio.NewReader(csvFile))

	for {
		var price Price
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			log.Fatal(error)
		}
		price.Price, _ = strconv.ParseFloat(line[1], 64)
		price.Open, _ = strconv.ParseFloat(line[2], 64)
		price.High, _ = strconv.ParseFloat(line[3], 64)
		price.Low, _ = strconv.ParseFloat(line[4], 64)
		priceList = append(priceList, price)
	}
}

func main() {

	var alpha float64
	var h float64

	readFile()

	constanta[0] = 1
	constanta[1] = 1
	constanta[2] = 1

	constGrada[0] = 1
	constGrada[1] = 1
	constGrada[2] = 1

	alpha = 1.e-9
	h = 1.e-7

	var iteration int64
	iteration = 0

	for FindGradNorm(constGrada) > 0.0001 {
		//for i := 0; i < 5; i++ {
		iteration++
		constGrada := FindGrad(constanta, h)
		for i := 0; i < 3; i++ {
			constanta[i] -= (constGrada[i] * alpha)
		}
		fmt.Println("-----------------------------------------")
		fmt.Println("iterasi ke-", iteration)
		fmt.Println("Gradien konstanta:", constGrada)
		fmt.Println("Constanta :", constanta)
		fmt.Println("MAE  :", MAE(constanta), "\n")
	}
}
