package main

import (
  "encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
  "strconv"
)

const (
  maxPrice = 10000.0
)

type Show interface {
  showSlice()
}

type TrainLegPriceSlice [][]TrainLegPrice

type float32Slice [][]float32

func (ourSlice TrainLegPriceSlice) showSlice() { // Более наглядно выводит матрицу начальных цен
  for i, _ := range ourSlice {
    for j, _ := range ourSlice {
      fmt.Printf("%v ", ourSlice[i][j].Price)
    }
    fmt.Println("\n")
  }
}

func (ourSlice float32Slice) showSlice() { // Более наглядно выводит матрицу конечных цен
  for i, _ := range ourSlice {
    for j, _ := range ourSlice {
      fmt.Printf("%v ", ourSlice[i][j])
    }
    fmt.Println("\n")
  }
}

func NewIter(minPrice []float32, visitedPoint []int) { // Задает необходимые, для дальнейших вычислений, значение массива мин. цен и посещенных станций
  for i, _ := range minPrice {
    minPrice[i] = maxPrice
    visitedPoint[i] = 1
  }
}

func AllShortestWayDijkstraPrice(resultInitialMatrixPrice [][]float32, initialMatrixPrice [][]TrainLegPrice, minPrice []float32, visitedPoint []int) { // Находит наиболее дешевый проезд между каждой парой станций
  for i, _ := range minPrice {
    ShortestWayDijkstraPrice(initialMatrixPrice, minPrice, visitedPoint, i)
    for j, _ := range minPrice {
       resultInitialMatrixPrice[i][j] = minPrice[j]
    }
    NewIter(minPrice, visitedPoint)
  }
}

func ShortestWayDijkstraPrice(initialMatrixPrice [][]TrainLegPrice, minPrice []float32, visitedPoint []int, stationindex int) { // Находит для станции под индексом stationindex кратчайший путь до всех остальных станций исользуя алгоритм Дейкстры
  minPrice[stationindex] = 0
  var minIndex int
  var min float32
  var temp float32
  for minIndex < 10000 {
    minIndex = 10000
    min = maxPrice
    for i, _ := range minPrice {
      if visitedPoint[i] == 1 && minPrice[i] < maxPrice {
        min = minPrice[i]
        minIndex = i;
      }
    }
    if minIndex != 10000 {
      for i, _ := range initialMatrixPrice {
        if initialMatrixPrice[minIndex][i].Price > 0.0 {
          temp = min + initialMatrixPrice[minIndex][i].Price
          if temp < minPrice[i] {
            minPrice[i] = temp
          }
        }
      }
      visitedPoint[minIndex] = 0
    }
  }
}

func InitialPrice(trainLegs TrainLegs) ([][]float32, [][]TrainLegPrice, []float32, []int, []string) { // На основе имеющейся информации создает и возвращает матрицы связей, мин. цены, пройденных пунктов, и уникальных имен пунктов
  uniqueStation := make([]string, 0)
  uniqueStation = append(uniqueStation, trainLegs.TrainLegsSlice[0].DepartureStationId)
  uniqueStation = append(uniqueStation, trainLegs.TrainLegsSlice[0].ArrivalStationId)
  inside1 := true
  inside2 := true
  var num float64
  var id int64
  for i := 1; i < len(trainLegs.TrainLegsSlice); i++ {
    inside1 = false
    inside2 = false
    for j := 0; j < len(uniqueStation); j++ {
      if trainLegs.TrainLegsSlice[i].DepartureStationId == uniqueStation[j] {
        inside1 = true
      }
      if trainLegs.TrainLegsSlice[i].ArrivalStationId == uniqueStation[j] {
        inside2 = true
      }
    }
    if inside1 == false {
      uniqueStation = append(uniqueStation, trainLegs.TrainLegsSlice[i].DepartureStationId)
    }
    if inside2 == false {
      uniqueStation = append(uniqueStation, trainLegs.TrainLegsSlice[i].ArrivalStationId)
    }
  }
  initialMatrixPrice := make([][]TrainLegPrice, len(uniqueStation))
  resultInitialMatrixPrice := make([][]float32, len(uniqueStation))
  for i, _ := range initialMatrixPrice {
    initialMatrixPrice[i] = make([]TrainLegPrice, len(uniqueStation))
    resultInitialMatrixPrice[i] = make([]float32, len(uniqueStation))
    for j, _ := range initialMatrixPrice[i] {
      initialMatrixPrice[i][j].Price = maxPrice
      resultInitialMatrixPrice[i][j] = 0.0
      if i == j { initialMatrixPrice[i][j].Price = 0.0 }
    }
  }
  for i := 0; i < len(trainLegs.TrainLegsSlice); i++ {
    for j, _ := range uniqueStation {
      for z, _ := range uniqueStation {
        num, _ = strconv.ParseFloat(trainLegs.TrainLegsSlice[i].Price, 32)
        if trainLegs.TrainLegsSlice[i].DepartureStationId == uniqueStation[j] && trainLegs.TrainLegsSlice[i].ArrivalStationId == uniqueStation[z] && float32(num) < initialMatrixPrice[j][z].Price {
          initialMatrixPrice[j][z].Price = float32(num)
          initialMatrixPrice[j][z].DepartureStationId = trainLegs.TrainLegsSlice[i].DepartureStationId
          initialMatrixPrice[j][z].ArrivalStationId = trainLegs.TrainLegsSlice[i].ArrivalStationId
          id, _ = strconv.ParseInt(trainLegs.TrainLegsSlice[i].TrainId, 16, 32)
          initialMatrixPrice[j][z].TrainId = int(id)
        }
      }
    }
  }
  minPrice := make([]float32, len(uniqueStation))
  visitedPoint := make([]int, len(uniqueStation))
  NewIter(minPrice, visitedPoint)
  return resultInitialMatrixPrice, initialMatrixPrice, minPrice, visitedPoint, uniqueStation
}

func ParseXML(fileName string) TrainLegs { // Считывает название файла и возвращает структуру, в которой записана необходимая информация
  xmlFile, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("Successfully Opened %s\n", fileName)
	defer xmlFile.Close()
  byteValue, err := ioutil.ReadAll(xmlFile)
  if err != nil {
		fmt.Println(err)
	}
  var trainLegs TrainLegs
  err = xml.Unmarshal(byteValue, &trainLegs)
  if err != nil {
    fmt.Printf("error: %trainLegs", err)
  }
  return trainLegs
}

type TrainLegPrice struct {
  TrainId int
  DepartureStationId string
  ArrivalStationId string
  Price float32
}

type TrainLegs struct {
	XMLName xml.Name `xml:"TrainLegs"`
	TrainLegsSlice   []TrainLeg   `xml:"TrainLeg"`
}

type TrainLeg struct {
  XMLName xml.Name `xml:"TrainLeg"`
	TrainId string `xml:"TrainId,attr"`
	DepartureStationId string `xml:"DepartureStationId,attr"`
	ArrivalStationId string `xml:"ArrivalStationId,attr"`
  Price string `xml:"Price,attr"`
  ArrivalTimeString string `xml:"ArrivalTimeString,attr"`
  DepartureTimeString string `xml:"DepartureTimeString,attr"`
}

func main() {
  var trainLegs TrainLegs
  var fileName string = "user.xml"
  trainLegs = ParseXML(fileName)
  var uniqueStation []string
  var initialMatrixPrice TrainLegPriceSlice
  var resultInitialMatrixPrice float32Slice
  var minPrice []float32
  var visitedPoint []int
  resultInitialMatrixPrice, initialMatrixPrice, minPrice, visitedPoint, uniqueStation = InitialPrice(trainLegs)
  AllShortestWayDijkstraPrice(resultInitialMatrixPrice, initialMatrixPrice, minPrice, visitedPoint)
  var s Show = initialMatrixPrice
  s.showSlice()
  s = resultInitialMatrixPrice
  s.showSlice()
  fmt.Println(minPrice)
  fmt.Println(visitedPoint)
  fmt.Println(uniqueStation)
}
