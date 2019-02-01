package main

import (
  "encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
  "strconv"
  "strings"
)

const (
  //maxPrice = 10000.0
  intconst = 10000
)

type Show interface {
  showSlice()
}

//type TrainLegPriceSlice [][]TrainLegPrice

type float64Slice [][]float64

type stringSlice []string

/*func (ourSlice TrainLegPriceSlice) showSlice() { // Более наглядно выводит матрицу начальных цен
  for i, _ := range ourSlice {
    for j, _ := range ourSlice {
      fmt.Printf("%v ", ourSlice[i][j].Price)
    }
    fmt.Println("\n")
  }
}*/

func (ourSlice stringSlice) showSlice() { // Более наглядно выводит матрицу начальных цен
  for i := len(ourSlice) - 1; ; i-- {
    if i == 0 {
      fmt.Printf("%s.", ourSlice[i])
      return
    }
    if ourSlice[i] == "" {
      continue
    }
    fmt.Printf("%s, ", ourSlice[i])
  }
}

func (ourSlice float64Slice) showSlice() { // Более наглядно выводит матрицу конечных цен
  for i, _ := range ourSlice {
    for j, _ := range ourSlice {
      fmt.Printf("%v ", float32(ourSlice[i][j]))
    }
    fmt.Println("\n")
  }
}

func NewIterTime(minTime []MinTime, visitedPoint []int) { // Задает необходимые, для дальнейших вычислений, значение массива мин. цен и посещенных станций
  for i, _ := range minTime {
    minTime[i].currentTime.Hours = 0
    minTime[i].currentTime.Minutes = 0
    minTime[i].currentTime.Seconds = 0
    minTime[i].sumTime.Hours = intconst
    minTime[i].sumTime.Minutes = 0
    minTime[i].sumTime.Seconds = 0
    visitedPoint[i] = 1
  }
}

func SumTime(add1 Time, add2 Time) Time {
  sum := add1.Hours*60*60 + add1.Minutes*60 + add1.Seconds + add2.Hours*60*60 + add2.Minutes*60 + add2.Seconds
  return Time{sum/(60*60), (sum - (sum/(60*60))*3600)/60, (sum%3600)%60}
}

func Difference(subtrahend Time, minuend Time) Time {
  if subtrahend.Hours > 23 && minuend.Hours < 24 {
    sub := Time{subtrahend.Hours%24, subtrahend.Minutes, subtrahend.Seconds}
    return Difference(sub, minuend)
  }
  diff := minuend.Hours*60*60 + minuend.Minutes*60 + minuend.Seconds - subtrahend.Hours*60*60 - subtrahend.Minutes*60 - subtrahend.Seconds
  if diff > 0 {
    return Time{diff/(60*60), (diff - (diff/(60*60))*3600)/60, (diff%3600)%60}
  }
  if diff < 0 {
    Time1, Time2 := Time{24, 0, 0}, Time{0, 0, 0}
    return SumTime(Difference(subtrahend, Time1), Difference(Time2, minuend))
  }
  return Time{0, 0, 0}
}

func DifferenceAll(subtrahend Time, minuend Time) Time {
  diff := minuend.Hours*60*60 + minuend.Minutes*60 + minuend.Seconds - subtrahend.Hours*60*60 - subtrahend.Minutes*60 - subtrahend.Seconds
  return Time{diff/(60*60), (diff - (diff/(60*60))*3600)/60, (diff%3600)%60}
}


/*func NewIterPrice(minPrice []float64, visitedPoint []int) { // Задает необходимые, для дальнейших вычислений, значение массива мин. цен и посещенных станций
  for i, _ := range minPrice {
    minPrice[i] = maxPrice
    visitedPoint[i] = 1
  }
}*/

func WayTime(DepartureStationId string, ArrivalStationId string, resultInitialMatrixTime [][]MinTime, initialMatrixTime [][]TrainLegsTime, uniqueStation []string, way stringSlice, infoTime [][]Time) {
  if DepartureStationId == ArrivalStationId {
    fmt.Printf("You are already in %s", ArrivalStationId)
    return
  }
  var start, finish, next int
  next = 0
  for i, _ := range uniqueStation {
    if uniqueStation[i] == DepartureStationId {
      start = i
    } else if uniqueStation[i] == ArrivalStationId {
      finish = i
      way[next] = uniqueStation[i]
      next++
    }
  }
  if (resultInitialMatrixTime[start][finish].sumTime.Hours == intconst) {
    fmt.Println("No way")
    return
  }
  var timezero = Time{0, 0, 0}
  fmt.Println(way)
  for finish != start {
    for i := 0; i < len(initialMatrixTime); i++ {
      for j := 0; j < len(initialMatrixTime[i]); j++ {
        for z := 0; z < len(infoTime); z++ {
          if (Difference(initialMatrixTime[i][finish].TrainTime[j].DepartureTime, initialMatrixTime[i][finish].TrainTime[j].ArrivalTime) != timezero) && Difference(initialMatrixTime[i][finish].TrainTime[j].DepartureTime, initialMatrixTime[i][finish].TrainTime[j].ArrivalTime) == Difference(resultInitialMatrixTime[start][i].sumTime, resultInitialMatrixTime[start][finish].sumTime) && Difference(resultInitialMatrixTime[start][i].sumTime, resultInitialMatrixTime[start][finish].sumTime) == infoTime[i][z] {
          infoTime[i][z] = timezero
          way[next] = uniqueStation[i]
          next++
          finish = i
          }
        }
      }
    }
  }
  var s Show
  s = way
  s.showSlice()
  return
}

func AllShortestWayDijkstraTime(resultInitialMatrixTime [][]MinTime, initialMatrixTime [][]TrainLegsTime, minTime []MinTime, visitedPoint []int, infoTime [][]Time) { // Находит наиболее дешевый проезд между каждой парой станций
  for i, _ := range minTime {
    ShortestWayDijkstraTimeInfo(initialMatrixTime, minTime, visitedPoint, infoTime, i)
    for j, _ := range minTime {
       resultInitialMatrixTime[i][j] = minTime[j]
    }
    NewIterTime(minTime, visitedPoint)
  }
}

func ShortestWayDijkstraTimeInfo(initialMatrixTime [][]TrainLegsTime, minTime []MinTime, visitedPoint []int, infoTime [][]Time, stationindex int) { // Находит для станции под индексом stationindex кратчайший путь до всех остальных станций исользуя алгоритм Дейкстры
  minTime[stationindex].sumTime.Hours = 0
  var minIndex int
  var min MinTime
  var temp, per, str1, str2 Time
  for minIndex < 10000 {
    minIndex = 10000
    for i, _ := range minTime {
      if visitedPoint[i] == 1 && minTime[i].sumTime.Hours < intconst {
        min = minTime[i]
        minIndex = i;
      }
    }
    str1 = Time{(min.currentTime.Hours/24)*24, 0, 0}
    str2 = DifferenceAll(str1, min.currentTime)
    if minIndex != 10000 {
      for i, _ := range initialMatrixTime[minIndex] {
        if i == minIndex {continue}
        for j := 1; j < len(initialMatrixTime[minIndex][i].TrainTime); j++ {
          temp = SumTime(Difference(str2, initialMatrixTime[minIndex][i].TrainTime[j].DepartureTime), Difference(initialMatrixTime[minIndex][i].TrainTime[j].DepartureTime, initialMatrixTime[minIndex][i].TrainTime[j].ArrivalTime))
          str1 = SumTime(min.currentTime, temp)
          per = DifferenceAll(str1, minTime[i].sumTime)
          if per.Hours > 0 || per.Minutes > 0 || per. Seconds > 0 {
            infoTime[stationindex] = append(infoTime[stationindex], temp)
            minTime[i].sumTime = str1
            minTime[i].currentTime = str1
          }
          str1 = Time{(min.currentTime.Hours/24)*24, 0, 0}
        }
      }
      visitedPoint[minIndex] = 0
    }
  }
}

func ShortestWayDijkstraTime(initialMatrixTime [][]TrainLegsTime, minTime []MinTime, visitedPoint []int, stationindex int) { // Находит для станции под индексом stationindex кратчайший путь до всех остальных станций исользуя алгоритм Дейкстры
  minTime[stationindex].sumTime.Hours = 0
  var minIndex int
  var min MinTime
  var temp, per, str1, str2 Time
  fmt.Println(minTime)
  for minIndex < 10000 {
    minIndex = 10000
    for i, _ := range minTime {
      if visitedPoint[i] == 1 && minTime[i].sumTime.Hours < intconst {
        min = minTime[i]
        minIndex = i;
      }
    }
    str1 = Time{(min.currentTime.Hours/24)*24, 0, 0}
    str2 = DifferenceAll(str1, min.currentTime)
    if minIndex != 10000 {
      for i, _ := range initialMatrixTime[minIndex] {
        if i == minIndex {continue}
        for j := 1; j < len(initialMatrixTime[minIndex][i].TrainTime); j++ {
          fmt.Println("\n", str1, str2, minIndex, i, j)
          fmt.Println(minTime)
          fmt.Println(initialMatrixTime[minIndex][i].TrainTime[j].DepartureTime, initialMatrixTime[minIndex][i].TrainTime[j].ArrivalTime)
          temp = SumTime(Difference(str2, initialMatrixTime[minIndex][i].TrainTime[j].DepartureTime), Difference(initialMatrixTime[minIndex][i].TrainTime[j].DepartureTime, initialMatrixTime[minIndex][i].TrainTime[j].ArrivalTime))
          str1 = SumTime(min.currentTime, temp)
          fmt.Println(minTime[i].sumTime)
          per = DifferenceAll(str1, minTime[i].sumTime)
          fmt.Println(str1, temp, per, minTime[i].sumTime)
          if per.Hours > 0 || per.Minutes > 0 || per. Seconds > 0 {
            minTime[i].sumTime = str1
            minTime[i].currentTime = str1
            fmt.Println(minTime[i].sumTime, minTime[i].currentTime)
          }
          str1 = Time{(min.currentTime.Hours/24)*24, 0, 0}
        }
      }
      visitedPoint[minIndex] = 0
    }
  }
}

func InitialTime(trainLegs TrainLegs) ([][]MinTime, [][]TrainLegsTime, []MinTime, []int, []string, [][]Time) { // На основе имеющейся информации создает и возвращает матрицы связей, мин. цены, пройденных пунктов, и уникальных имен пунктов
  uniqueStation := make([]string, 0)
  uniqueStation = append(uniqueStation, trainLegs.TrainLegsSlice[0].DepartureStationId)
  uniqueStation = append(uniqueStation, trainLegs.TrainLegsSlice[0].ArrivalStationId)
  inside1 := true
  inside2 := true
  var num int64
  var trainLegTime TrainLegTime
  var time1, time2 Time
  time := make([]string, 3)
  timeint := make([]int64, 3)
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
  initialMatrixTime := make([][]TrainLegsTime, len(uniqueStation))
  resultInitialMatrixTime := make([][]MinTime, len(uniqueStation))
  infoTime := make([][]Time, len(uniqueStation))
  for i, _ := range initialMatrixTime {
    initialMatrixTime[i] = make([]TrainLegsTime, len(uniqueStation))
    resultInitialMatrixTime[i] = make([]MinTime, len(uniqueStation))
    infoTime[i] = make([]Time, 0)
    for j, _ := range initialMatrixTime[i] {
      resultInitialMatrixTime[i][j] = MinTime{Time{0, 0, 0}, Time{0, 0, 0}}
    }
  }
  time1.Hours = 0
  time1.Minutes = 0
  time1.Seconds = 0
  time2.Hours = 0
  time2.Minutes = 0
  time2.Seconds = 0
  trainLegTime = TrainLegTime{ TrainId : 0, DepartureStationId : "", ArrivalStationId : "", DepartureTime : time1, ArrivalTime : time2 }
  for i, _ := range initialMatrixTime {
    for j, _ := range initialMatrixTime {
      initialMatrixTime[i][j].TrainTime = append(initialMatrixTime[i][j].TrainTime, trainLegTime)
    }
  }
  for i := 0; i < len(trainLegs.TrainLegsSlice); i++ {
    for j, _ := range uniqueStation {
      for z, _ := range uniqueStation {
        if trainLegs.TrainLegsSlice[i].DepartureStationId == uniqueStation[j] && trainLegs.TrainLegsSlice[i].ArrivalStationId == uniqueStation[z] {
          time = strings.Split(trainLegs.TrainLegsSlice[i].DepartureTimeString, ":")
          timeint[0], _ = strconv.ParseInt(time[0], 0, 32)
          timeint[1], _ = strconv.ParseInt(time[1], 0, 32)
          timeint[2], _ = strconv.ParseInt(time[2], 0, 32)
          time1 = Time{ Hours : int(timeint[0]), Minutes : int(timeint[1]), Seconds : int(timeint[2])}
          time = strings.Split(trainLegs.TrainLegsSlice[i].ArrivalTimeString, ":")
          timeint[0], _ = strconv.ParseInt(time[0], 0, 32)
          timeint[1], _ = strconv.ParseInt(time[1], 0, 32)
          timeint[2], _ = strconv.ParseInt(time[2], 0, 32)
          time2 = Time{ Hours : int(timeint[0]), Minutes : int(timeint[1]), Seconds : int(timeint[2])}
          num, _ = strconv.ParseInt(trainLegs.TrainLegsSlice[i].TrainId, 16, 32)
          trainLegTime = TrainLegTime{ TrainId : int(num), DepartureStationId : trainLegs.TrainLegsSlice[i].DepartureStationId, ArrivalStationId : trainLegs.TrainLegsSlice[i].ArrivalStationId, DepartureTime : time1, ArrivalTime : time2 }
          initialMatrixTime[j][z].TrainTime = append(initialMatrixTime[j][z].TrainTime, trainLegTime)
        }
      }
    }
  }

  minTime := make([]MinTime, len(uniqueStation))
  visitedPoint := make([]int, len(uniqueStation))
  NewIterTime(minTime, visitedPoint)
  return resultInitialMatrixTime, initialMatrixTime, minTime, visitedPoint, uniqueStation, infoTime
}

/*func WayPrice(DepartureStationId string, ArrivalStationId string, resultInitialMatrixPrice float64Slice, initialMatrixPrice TrainLegPriceSlice, uniqueStation []string, way stringSlice) {
  if DepartureStationId == ArrivalStationId {
    fmt.Printf("You are already in %s", ArrivalStationId)
    return
  }
  var start, finish, next int
  next = 0
  for i, _ := range uniqueStation {
    if uniqueStation[i] == DepartureStationId {
      start = i
    } else if uniqueStation[i] == ArrivalStationId {
      finish = i
      way[next] = uniqueStation[i]
      next++
    }
  }
  if (resultInitialMatrixPrice[start][finish] == maxPrice) {
    fmt.Println("No way")
    return
  }
  fmt.Println(way)
  for finish != start {
    for i := 0; i < len(initialMatrixPrice); i++ {
      fmt.Println("\n", way, i)
      if (initialMatrixPrice[i][finish].Price != 0.0) && (initialMatrixPrice[i][finish].Price == (resultInitialMatrixPrice[start][finish] - resultInitialMatrixPrice[start][i])) {
        way[next] = uniqueStation[i]
        next++
        finish = i
        fmt.Println("\n", way)
      }
    }
  }
  var s Show
  s = way
  s.showSlice()
  return
}

func AllShortestWayDijkstraPrice(resultInitialMatrixPrice [][]float64, initialMatrixPrice [][]TrainLegPrice, minPrice []float64, visitedPoint []int) { // Находит наиболее дешевый проезд между каждой парой станций
  for i, _ := range minPrice {
    ShortestWayDijkstraPrice(initialMatrixPrice, minPrice, visitedPoint, i)
    for j, _ := range minPrice {
       resultInitialMatrixPrice[i][j] = minPrice[j]
    }
    NewIterPrice(minPrice, visitedPoint)
  }
}

func ShortestWayDijkstraPrice(initialMatrixPrice [][]TrainLegPrice, minPrice []float64, visitedPoint []int, stationindex int) { // Находит для станции под индексом stationindex кратчайший путь до всех остальных станций исользуя алгоритм Дейкстры
  minPrice[stationindex] = 0
  var minIndex int
  var min float64
  var temp float64
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

func InitialPrice(trainLegs TrainLegs) ([][]float64, [][]TrainLegPrice, []float64, []int, []string) { // На основе имеющейся информации создает и возвращает матрицы связей, мин. цены, пройденных пунктов, и уникальных имен пунктов
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
  resultInitialMatrixPrice := make([][]float64, len(uniqueStation))
  for i, _ := range initialMatrixPrice {
    initialMatrixPrice[i] = make([]TrainLegPrice, len(uniqueStation))
    resultInitialMatrixPrice[i] = make([]float64, len(uniqueStation))
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
        if trainLegs.TrainLegsSlice[i].DepartureStationId == uniqueStation[j] && trainLegs.TrainLegsSlice[i].ArrivalStationId == uniqueStation[z] && num < initialMatrixPrice[j][z].Price {
          initialMatrixPrice[j][z].Price = num
          initialMatrixPrice[j][z].DepartureStationId = trainLegs.TrainLegsSlice[i].DepartureStationId
          initialMatrixPrice[j][z].ArrivalStationId = trainLegs.TrainLegsSlice[i].ArrivalStationId
          id, _ = strconv.ParseInt(trainLegs.TrainLegsSlice[i].TrainId, 0, 32)
          initialMatrixPrice[j][z].TrainId = int(id)
          //initialMatrixPrice[i][j].Ways[0].station = uniqueStation[i]
        }
      }
    }
  }
  minPrice := make([]float64, len(uniqueStation))
  visitedPoint := make([]int, len(uniqueStation))
  NewIterPrice(minPrice, visitedPoint)
  return resultInitialMatrixPrice, initialMatrixPrice, minPrice, visitedPoint, uniqueStation
}*/

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

type TrainLegsTime struct {
  TrainTime []TrainLegTime
}

type TrainLegTime struct {
  TrainId int
  DepartureStationId string
  ArrivalStationId string
  DepartureTime Time
  ArrivalTime Time
}

type Time struct {
  Hours int
  Minutes int
  Seconds int
}

type MinTime struct {
  currentTime Time
  sumTime Time
}

/*type TrainLegPrice struct {
  TrainId int
  DepartureStationId string
  ArrivalStationId string
  Price float64
}*/

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
  var fileName string = "users.xml"
  trainLegs = ParseXML(fileName)
  var uniqueStation []string
  var visitedPoint []int





  // Test Price
  /*var initialMatrixPrice TrainLegPriceSlice
  var resultInitialMatrixPrice float64Slice
  var minPrice []float64
  resultInitialMatrixPrice, initialMatrixPrice, minPrice, visitedPoint, uniqueStation = InitialPrice(trainLegs)
  AllShortestWayDijkstraPrice(resultInitialMatrixPrice, initialMatrixPrice, minPrice, visitedPoint)
  var s Show = initialMatrixPrice
  s.showSlice()
  s = resultInitialMatrixPrice
  s.showSlice()
  fmt.Println(minPrice)
  fmt.Println(visitedPoint)
  fmt.Println(uniqueStation)
  way := make(stringSlice, len(uniqueStation))
  //allway := make([][]string, len(uniqueStation))
  WayPrice("1929", "1981", resultInitialMatrixPrice, initialMatrixPrice, uniqueStation, way)*/





  //Test Time
  var initialMatrixTime [][]TrainLegsTime
  var resultInitialMatrixTime [][]MinTime
  var minTime []MinTime
  var infoTime [][]Time
  resultInitialMatrixTime, initialMatrixTime, minTime, visitedPoint, uniqueStation, infoTime = InitialTime(trainLegs)
  for i, _ := range initialMatrixTime {
    for j, _ := range initialMatrixTime {
        fmt.Println(" ", initialMatrixTime[i][j])
    }
    fmt.Println("\n")
  }
  fmt.Println(uniqueStation, minTime)
  ShortestWayDijkstraTimeInfo(initialMatrixTime, minTime, visitedPoint, infoTime, 0)
  fmt.Println(minTime)
  fmt.Println(infoTime)
  AllShortestWayDijkstraTime(resultInitialMatrixTime, initialMatrixTime, minTime, visitedPoint, infoTime)
  for i, _ := range resultInitialMatrixTime {
    for j, _ := range resultInitialMatrixTime {
      fmt.Println(resultInitialMatrixTime[i][j], i, j)
    }
    fmt.Println("\n")
  }
  way := make(stringSlice, len(uniqueStation))
  WayTime("Kyiv", "Melitopol", resultInitialMatrixTime, initialMatrixTime, uniqueStation, way, infoTime)

}
