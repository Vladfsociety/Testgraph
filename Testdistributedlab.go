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

func ShortestWayDijkstraPrice(initialMatrixPrice [][]float32, minPrice []float32, visitedPoint []int, uniqueStation []string) {

}

func InitialPrice(trainLegs TrainLegs) ([][]float32, []float32, []int, []string) {
  uniqueStation := make([]string, 0)
  uniqueStation = append(uniqueStation, trainLegs.TrainLegsSlice[0].DepartureStationId)
  uniqueStation = append(uniqueStation, trainLegs.TrainLegsSlice[0].ArrivalStationId)
  inside1 := true
  inside2 := true
  var num float64
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
  initialMatrixPrice := make([][]float32, len(uniqueStation))
  for i, _ := range initialMatrixPrice {
    initialMatrixPrice[i] = make([]float32, len(uniqueStation))
    for j, _ := range initialMatrixPrice[i] {
      initialMatrixPrice[i][j] = maxPrice
    }
  }
  for i := 0; i < len(trainLegs.TrainLegsSlice); i++ {
    for j, _ := range uniqueStation {
      for z, _ := range uniqueStation {
        num, _ = strconv.ParseFloat(trainLegs.TrainLegsSlice[i].Price, 32)
        if trainLegs.TrainLegsSlice[i].DepartureStationId == uniqueStation[j] && trainLegs.TrainLegsSlice[i].ArrivalStationId == uniqueStation[z] && float32(num) < initialMatrixPrice[j][z] {
          initialMatrixPrice[j][z] = float32(num)
        }
      }
    }
  }
  for i, _ := range initialMatrixPrice {
    for j, _ := range initialMatrixPrice {
      fmt.Printf("%v ", initialMatrixPrice[i][j])
    }
    fmt.Println("\n")
  }
  minPrice := make([]float32, len(uniqueStation))
  visitedPoint := make([]int, len(uniqueStation))
  for i, _ := range uniqueStation {
    minPrice[i] = maxPrice
    visitedPoint[i] = 1
  }
  return initialMatrixPrice, minPrice, visitedPoint, uniqueStation
}

func ParseXML(fileName string) TrainLegs {
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
  trainLegs = ParseXML(fileName) // Считывает название файла и возвращает структуру, в которой записана необходимая информация
  var uniqueStation []string
  var initialMatrixPrice [][]float32
  var minPrice []float32
  var visitedPoint []int
  initialMatrixPrice, minPrice, visitedPoint, uniqueStation = InitialPrice(trainLegs) // На основе имеющейся информации создает и возвращает матрицы связей, мин. цены, пройденных пунктов, и уникальных имен пунктов
  /*uniqueStation = append(uniqueStation, trainLegs.TrainLegsSlice[0].DepartureStationId)
  uniqueStation = append(uniqueStation, trainLegs.TrainLegsSlice[0].ArrivalStationId)
  inside1 := true
  inside2 := true
  for i := 0; i < len(trainLegs.TrainLegsSlice); i++ {
    fmt.Println("TrainLeg TrainId: " + trainLegs.TrainLegsSlice[i].TrainId)
    fmt.Println("TrainLeg DepartureStationId: " + trainLegs.TrainLegsSlice[i].DepartureStationId)
    fmt.Println("TrainLeg ArrivalStationId: " + trainLegs.TrainLegsSlice[i].ArrivalStationId)
    fmt.Println("TrainLeg Price: " + trainLegs.TrainLegsSlice[i].Price)
    fmt.Println("TrainLeg ArrivalTimeString: " + trainLegs.TrainLegsSlice[i].ArrivalTimeString)
    fmt.Println("TrainLeg DepartureTimeString: " + trainLegs.TrainLegsSlice[i].DepartureTimeString)
  }

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
  initialMatrixPrice := make([][]float32, len(uniqueStation))
  for i, _ := range initialMatrixPrice {
    initialMatrixPrice[i] = make([]float32, len(uniqueStation))
    for j, _ := range initialMatrixPrice[i] {
      initialMatrixPrice[i][j] = maxPrice
    }
  }
  var num float64
  for i := 0; i < len(trainLegs.TrainLegsSlice); i++ {
    for j, _ := range uniqueStation {
      for z, _ := range uniqueStation {
        num, _ = strconv.ParseFloat(trainLegs.TrainLegsSlice[i].Price, 32)
        if trainLegs.TrainLegsSlice[i].DepartureStationId == uniqueStation[j] && trainLegs.TrainLegsSlice[i].ArrivalStationId == uniqueStation[z] && float32(num) < initialMatrixPrice[j][z] {
          initialMatrixPrice[j][z] = float32(num)
        }
      }
    }
  }
  for i, _ := range initialMatrixPrice {
    for j, _ := range initialMatrixPrice {
      fmt.Printf("%v ", initialMatrixPrice[i][j])
    }
    fmt.Println("\n")
  }
  for i, _ := range uniqueStation {
    minPrice[i] = maxPrice
    visitedPoint[i] = 1
  }*/
  fmt.Println(initialMatrixPrice)
  fmt.Println(minPrice)
  fmt.Println(visitedPoint)
  fmt.Println(uniqueStation)
}
