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
  maxPrice = 10000.0 // Цена обозначающая невозможность прямого проезда
  intconst = 10000 // Аналогичная величина для времени, записывается в часы
)

func showSlice(ourSlice []Way) { // Выводит путь записанный в way в нужном порядке
  var next1, next2 int
  next1 = 0
  next2 = 0
  for i := len(ourSlice) - 1; i >= 0; i-- {
    if ourSlice[i].name != "" && next1 == 0 && i != 0 {
      next1 = i
      continue
    }
    if ourSlice[i].name != "" && next2 == 0 {
      next2 = i
    }
    if next1 != 0 && ourSlice[next2].name != "" {
      fmt.Printf("From %s to %s, id = %v, ", ourSlice[next1].name, ourSlice[next2].name, ourSlice[next1].id)
      next1 = next2
      next2 = 0
    }
  }
}

func NewIterTime(minTime []MinTime, visitedPoint []int) {
  for i, _ := range minTime {
    minTime[i].currentTime.Hours = 0
    minTime[i].currentTime.Minutes = 0
    minTime[i].currentTime.Seconds = 0
    minTime[i].sumTime.Hours = intconst
    minTime[i].sumTime.Minutes = 0
    minTime[i].sumTime.Seconds = 0
    visitedPoint[i] = 1
  }
} // Задает необходимые, для дальнейших вычислений, значение массива мин. времен и посещенных станций

func SumTime(add1 Time, add2 Time) Time {
  sum := add1.Hours*60*60 + add1.Minutes*60 + add1.Seconds + add2.Hours*60*60 + add2.Minutes*60 + add2.Seconds
  return Time{sum/(60*60), (sum - (sum/(60*60))*3600)/60, (sum%3600)%60}
} // Суммирует 2 времени

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
  return Time{0, 0, 0}} // Вычитает из minued subtrahend, предназначена для нахождения времени переезда, которое меньше 24 часов

func DifferenceAll(subtrahend Time, minuend Time) Time {
  diff := minuend.Hours*60*60 + minuend.Minutes*60 + minuend.Seconds - subtrahend.Hours*60*60 - subtrahend.Minutes*60 - subtrahend.Seconds
  return Time{diff/(60*60), (diff - (diff/(60*60))*3600)/60, (diff%3600)%60}
} // Тоже что и Difference, но предназначена для вычитание времен больше 24 часов

func DifferenceSec(subtrahend Time, minuend Time) int {
  diff := minuend.Hours*60*60 + minuend.Minutes*60 + minuend.Seconds - subtrahend.Hours*60*60 - subtrahend.Minutes*60 - subtrahend.Seconds
  return diff
} // Тоже что DifferenceAll, возвращает результат в секундах

func NewIterPrice(minPrice []float64, visitedPoint []int) {
  for i, _ := range minPrice {
    minPrice[i] = maxPrice
    visitedPoint[i] = 1
  }
} // Задает необходимые, для дальнейших вычислений, значение массива мин. цен и посещенных станций

func WayTime(DepartureStationId string, ArrivalStationId string, resultInitialMatrixTime [][]MinTime, initialMatrixTime [][]TrainLegsTime, uniqueStation []string, way []Way, infoTime [][]Time) {
  if DepartureStationId == ArrivalStationId {
    fmt.Printf("You are already in %s\n", ArrivalStationId)
    return
  }
  var start, finish, print, next int
  next = 0
  for i, _ := range uniqueStation {
    if uniqueStation[i] == DepartureStationId {
      start = i
    } else if uniqueStation[i] == ArrivalStationId {
      finish = i
      print = finish
      way[next].name = uniqueStation[i]
      next++
    }
  }
  if (resultInitialMatrixTime[start][finish].sumTime.Hours == intconst) {
    fmt.Println("No way")
    return
  }
  var timezero = Time{0, 0, 0}
  var str2 Time
  for finish != start {
    for i := 0; i < len(initialMatrixTime); i++ {
      str2 = Time{resultInitialMatrixTime[start][i].sumTime.Hours%24, resultInitialMatrixTime[start][i].sumTime.Minutes, resultInitialMatrixTime[start][i].sumTime.Seconds}
      for j := 0; j < len(initialMatrixTime[i][finish].TrainTime); j++ {
        for z := 0; z < len(infoTime[start]); z++ {
          if (Difference(initialMatrixTime[i][finish].TrainTime[j].DepartureTime, initialMatrixTime[i][finish].TrainTime[j].ArrivalTime) != timezero) && SumTime(Difference(str2, initialMatrixTime[i][finish].TrainTime[j].DepartureTime), Difference(initialMatrixTime[i][finish].TrainTime[j].DepartureTime, initialMatrixTime[i][finish].TrainTime[j].ArrivalTime)) == DifferenceAll(resultInitialMatrixTime[start][i].sumTime, resultInitialMatrixTime[start][finish].sumTime) && DifferenceAll(resultInitialMatrixTime[start][i].sumTime, resultInitialMatrixTime[start][finish].sumTime) == infoTime[start][z] {
          infoTime[start][z] = timezero
          way[next].name = uniqueStation[i]
          way[next].id = initialMatrixTime[i][finish].TrainTime[j].TrainId
          next++
          finish = i
          break
          }
        }
      }
    }
  }
  showSlice(way)
  fmt.Printf("Time: Hours:%v Minutes:%v, Seconds:%v ", resultInitialMatrixTime[start][print].sumTime.Hours, resultInitialMatrixTime[start][print].sumTime.Minutes, resultInitialMatrixTime[start][print].sumTime.Seconds)
  return
} // Восстанавливает кратчайший(по времени) путь с помощью матриц resultInitialMatrixTime, initialMatrixTime и выводит его, а также время перезда

func AllShortestWayDijkstraTime(resultInitialMatrixTime [][]MinTime, initialMatrixTime [][]TrainLegsTime, minTime []MinTime, visitedPoint []int, infoTime [][]Time, uniqueStation []string) {
  for i, _ := range minTime {
    ShortestWayDijkstraTimeInfo(initialMatrixTime, minTime, visitedPoint, infoTime, i, uniqueStation)
    for j, _ := range minTime {
       resultInitialMatrixTime[i][j] = minTime[j]
    }
    NewIterTime(minTime, visitedPoint)
  }
} // Находит наиболее короткий(по времени) проезд между каждой парой станций

func ShortestWayDijkstraTimeInfo(initialMatrixTime [][]TrainLegsTime, minTime []MinTime, visitedPoint []int, infoTime [][]Time, stationindex int, uniqueStation []string) { // Находит для станции под индексом stationindex кратчайший(по времени) путь до всех остальных станций исользуя алгоритм Дейкстры
  minTime[stationindex].sumTime.Hours = 0
  var minIndex int
  var min MinTime
  var temp, per, str1, str2 Time
  var minimum Time
  for minIndex < 10000 {
    minIndex = 10000
    minimum = Time{intconst, 0, 0}
    for i, _ := range minTime {
      if visitedPoint[i] == 1 && DifferenceSec(minimum, minTime[i].sumTime) < 0 {
        minimum = minTime[i].sumTime
        min = minTime[i]
        minIndex = i;
      }
    }
    str1 = Time{(min.currentTime.Hours/24)*24, 0, 0}
    str2 = DifferenceAll(str1, min.currentTime)
    if minIndex != 10000 {
      for i, _ := range initialMatrixTime[minIndex] {
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

func InitialTime(trainLegs TrainLegs) ([][]MinTime, [][]TrainLegsTime, []MinTime, []int, []string, [][]Time) {
  uniqueStation := make([]string, 0)
  uniqueStation = append(uniqueStation, trainLegs.TrainLegsSlice[0].DepartureStationId)
  uniqueStation = append(uniqueStation, trainLegs.TrainLegsSlice[0].ArrivalStationId)
  inside1 := true
  inside2 := true
  var num int
  var trainLegTime TrainLegTime
  var time1, time2 Time
  time := make([]string, 3)
  timeint := make([]int, 3)
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
          timeint[0], _ = strconv.Atoi(time[0])
          timeint[1], _ = strconv.Atoi(time[1])
          timeint[2], _ = strconv.Atoi(time[2])
          time1 = Time{ Hours : int(timeint[0]), Minutes : int(timeint[1]), Seconds : int(timeint[2])}
          time = strings.Split(trainLegs.TrainLegsSlice[i].ArrivalTimeString, ":")
          timeint[0], _ = strconv.Atoi(time[0])
          timeint[1], _ = strconv.Atoi(time[1])
          timeint[2], _ = strconv.Atoi(time[2])
          time2 = Time{ Hours : int(timeint[0]), Minutes : int(timeint[1]), Seconds : int(timeint[2])}
          num, _ = strconv.Atoi(trainLegs.TrainLegsSlice[i].TrainId)
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
} // На основе имеющейся информации создает и возвращает матрицы кратчайших(по времени) путей между каждой парой станций(пока не заполнена), связей, мин. времен, пройденных пунктов, уникальных имен пунктов и различных времен переезда между станциями

func WayPrice(DepartureStationId string, ArrivalStationId string, resultInitialMatrixPrice [][]float64, initialMatrixPrice [][]TrainLegPrice, uniqueStation []string, way []Way) {
  if DepartureStationId == ArrivalStationId {
    fmt.Printf("You are already in %s\n", ArrivalStationId)
    return
  }
  var start, finish, print, next int
  next = 0
  for i, _ := range uniqueStation {
    if uniqueStation[i] == DepartureStationId {
      start = i
    } else if uniqueStation[i] == ArrivalStationId {
      finish = i
      print = finish
      way[next].name = uniqueStation[i]
      next++
    }
  }
  if (resultInitialMatrixPrice[start][finish] == maxPrice) {
    fmt.Println("No way")
    return
  }
  for finish != start {
    for i := 0; i < len(initialMatrixPrice); i++ {
      if (initialMatrixPrice[i][finish].Price != 0.0) && (initialMatrixPrice[i][finish].Price == (resultInitialMatrixPrice[start][finish] - resultInitialMatrixPrice[start][i])) {
        way[next].name = uniqueStation[i]
        way[next].id = initialMatrixPrice[i][finish].TrainId
        next++
        finish = i
      }
    }
  }
  showSlice(way)
  fmt.Println("Price ", resultInitialMatrixPrice[start][print])
  return
} //Восстанавливает кратчайший(по цене) путь с помощью матриц resultInitialMatrixPrice, initialMatrixTime и выводит его, а также цену переезда

func AllShortestWayDijkstraPrice(resultInitialMatrixPrice [][]float64, initialMatrixPrice [][]TrainLegPrice, minPrice []float64, visitedPoint []int) {
  for i, _ := range minPrice {
    ShortestWayDijkstraPrice(initialMatrixPrice, minPrice, visitedPoint, i)
    for j, _ := range minPrice {
       resultInitialMatrixPrice[i][j] = minPrice[j]
    }
    NewIterPrice(minPrice, visitedPoint)
  }
} // Находит наиболее дешевый проезд между каждой парой станций

func ShortestWayDijkstraPrice(initialMatrixPrice [][]TrainLegPrice, minPrice []float64, visitedPoint []int, stationindex int) {
  minPrice[stationindex] = 0
  var minIndex int
  var min float64
  var temp float64
  var minimum float64
  for minIndex < 10000 {
    minIndex = 10000
    min = maxPrice
    minimum = maxPrice
    for i, _ := range minPrice {
      if visitedPoint[i] == 1 && (minPrice[i] - minimum) < 0 {
        minimum = minPrice[i]
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
} // Находит для станции под индексом stationindex кратчайший(по цене) путь до всех остальных станций исользуя алгоритм Дейкстры

func InitialPrice(trainLegs TrainLegs) ([][]float64, [][]TrainLegPrice, []float64, []int, []string) {
  uniqueStation := make([]string, 0)
  uniqueStation = append(uniqueStation, trainLegs.TrainLegsSlice[0].DepartureStationId)
  uniqueStation = append(uniqueStation, trainLegs.TrainLegsSlice[0].ArrivalStationId)
  inside1 := true
  inside2 := true
  var num float64
  var id int
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
          id, _ = strconv.Atoi(trainLegs.TrainLegsSlice[i].TrainId)
          initialMatrixPrice[j][z].TrainId = id
          //initialMatrixPrice[i][j].Ways[0].station = uniqueStation[i]
        }
      }
    }
  }
  minPrice := make([]float64, len(uniqueStation))
  visitedPoint := make([]int, len(uniqueStation))
  NewIterPrice(minPrice, visitedPoint)
  return resultInitialMatrixPrice, initialMatrixPrice, minPrice, visitedPoint, uniqueStation
} // На основе имеющейся информации создает и возвращает матрицы кратчайших(по цене) путей между каждой парой станций(пока не заполнена), связей, мин. цен, пройденных пунктов, уникальных имен пунктов

func ParseXML(fileName string) TrainLegs {
  xmlFile, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)
	}
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
} // Считывает название файла и возвращает структуру, в которой записана необходимая информация из XML файла

type Way struct {
  name string
  id int
} // Служит для восстановления пути

type TrainLegsTime struct {
  TrainTime []TrainLegTime
} // Используется в матрице initialMatrixTime хранит список различных прямых переездов между парой станций

type TrainLegTime struct {
  TrainId int
  DepartureStationId string
  ArrivalStationId string
  DepartureTime Time
  ArrivalTime Time
} // Используется в матрице initialMatrixTime хранит 1 перезд между парой станций

type Time struct {
  Hours int
  Minutes int
  Seconds int
} // Хранит время в часах, минутах и секундах

type MinTime struct {
  currentTime Time
  sumTime Time
} // Используется в матрице resultInitialMatrixTime хранит текущее время и суммарное время необходимое чтобы добраться от одной станции до другой

type TrainLegPrice struct {
  TrainId int
  DepartureStationId string
  ArrivalStationId string
  Price float64
} // Используется в матрице initialMatrixPrice хранит информацию о переезде между парой станций, хранит всего 1 переезд. Так как время не интересует, то выбираем только самый дешевый переезд

type TrainLegs struct {
	XMLName xml.Name `xml:"TrainLegs"`
	TrainLegsSlice   []TrainLeg   `xml:"TrainLeg"`
} // Хранит елементы XML файла

type TrainLeg struct {
  XMLName xml.Name `xml:"TrainLeg"`
	TrainId string `xml:"TrainId,attr"`
	DepartureStationId string `xml:"DepartureStationId,attr"`
	ArrivalStationId string `xml:"ArrivalStationId,attr"`
  Price string `xml:"Price,attr"`
  ArrivalTimeString string `xml:"ArrivalTimeString,attr"`
  DepartureTimeString string `xml:"DepartureTimeString,attr"`
} // Хранит элементы, в которых записана информация о переездах

func main() {
  var station1, station2 string
  fmt.Print("Введите начальную станцию: ")
  fmt.Fscan(os.Stdin, &station1)
  fmt.Print("Введите конечную станцию: ")
  fmt.Fscan(os.Stdin, &station2)
  var trainLegs TrainLegs // Структура, в которую считывается информация из XML файла
  var fileName string = "data.xml" // имя файла
  trainLegs = ParseXML(fileName)
  var uniqueStation []string // Хранит уникальные названия станций, в матрицах resultInitialMatrixTime, initialMatrixTime, resultInitialMatriPrice, initialMatrixPrice индексы совпадают с индексами в этой матрице
  var visitedPoint []int // Используется в некоторых функциях хранит информацию о том, была ли посещена станция(0) или нет(1)
  // Test Price
  var initialMatrixPrice [][]TrainLegPrice // Хранит информацию о перездах между каждой парой станций, элемент initialMatrixPrice[i][j] хранит информацию о переезде из и i в j(для цен)
  var resultInitialMatrixPrice [][]float64 // Хранит конечную информацию о перездах между каждой парой станций, элемент initialMatrixPrice[i][j] хранит информацию о переезде из и i в j(для цен)
  var minPrice []float64 // Необходима для заполнения resultInitialMatrixPrice
  resultInitialMatrixPrice, initialMatrixPrice, minPrice, visitedPoint, uniqueStation = InitialPrice(trainLegs)
  AllShortestWayDijkstraPrice(resultInitialMatrixPrice, initialMatrixPrice, minPrice, visitedPoint)
  way1 := make([]Way, len(uniqueStation)) // Хранит инверсированный маршрут
  WayPrice(station1, station2, resultInitialMatrixPrice, initialMatrixPrice, uniqueStation, way1)

  //Test Time
  var initialMatrixTime [][]TrainLegsTime // Хранит информацию о перездах между каждой парой станций, элемент initialMatrixPrice[i][j] хранит информацию о переезде из и i в j(для времени)
  var resultInitialMatrixTime [][]MinTime // Хранит конечную информацию о перездах между каждой парой станций, элемент initialMatrixPrice[i][j] хранит информацию о переезде из и i в j(для времени)
  var minTime []MinTime // Необходима для заполнения resultInitialMatrixTime
  var infoTime [][]Time // Хранит уникальные времени перезда(время простоя + время пути) между парами станций, необходима для восстановления пути
  resultInitialMatrixTime, initialMatrixTime, minTime, visitedPoint, uniqueStation, infoTime = InitialTime(trainLegs)
  AllShortestWayDijkstraTime(resultInitialMatrixTime, initialMatrixTime, minTime, visitedPoint, infoTime, uniqueStation)
  way2 := make([]Way, len(uniqueStation)) // аналогично
  WayTime(station1, station2, resultInitialMatrixTime, initialMatrixTime, uniqueStation, way2, infoTime)
}
