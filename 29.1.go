package main

import (
	"fmt"
	"math"
	"strconv"
)

func main() {

	for {
		var enter string
		fmt.Println("Введите число:")
		_, err := fmt.Scan(&enter)
		if err != nil {
			return
		}
		if enter == "стоп" {
			break
		}
		enterNumber, err := strconv.Atoi(enter)
		if err != nil {
			fmt.Println("Введите целое число или слово стоп для завершения программы")
			continue
		}
		fc := square(enterNumber)
		sc := multiplication(fc)

		fmt.Println("Произведение:", <-sc)
	}

}

func square(number int) chan float64{
	firstChan := make(chan float64)
	squareNumber := math.Pow(float64(number), 2)
	fmt.Println("Ввод:", number)
	fmt.Println("Квадрат:", squareNumber)
	go func() {
		firstChan <- squareNumber
	}()
	return firstChan
}

func multiplication(firstChan chan float64) chan float64{
	secondChan := make(chan float64)
	multipliedNumber := 2 * (<- firstChan)
	go func() {
		secondChan <- multipliedNumber
	}()
	return secondChan
}

