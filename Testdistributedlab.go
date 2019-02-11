package main

import (
  "encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
  "strconv"
  "strings"
  "log"
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
  return Time{0, 0, 0}
} // Вычитает из minued subtrahend, предназначена для нахождения времени переезда, которое меньше 24 часов

func DifferenceAll(subtrahend Time, minuend Time) Time {
  diff := minuend.Hours*60*60 + minuend.Minutes*60 + minuend.Seconds - subtrahend.Hours*60*60 - subtrahend.Minutes*60 - subtrahend.Seconds
  return Time{diff/(60*60), (diff - (diff/(60*60))*3600)/60, (diff%3600)%60}
} // Тоже что и Difference, но предназначена для вычитание времен больше 24 часов

func DifferenceSec(subtrahend Time, minuend Time) int {
  diff := minuend.Hours*60*60 + minuend.Minutes*60 + minuend.Seconds - subtrahend.Hours*60*60 - subtrahend.Minutes*60 - subtrahend.Seconds
  return diff
} // Тоже что DifferenceAll, возвращает результат в секундах

func NewIter(min []ResultMatrix, visitedPoint []int) {
  for i, _ := range min {
    min[i].resPrice = maxPrice
    min[i].resTime.currentTime.Hours = 0
    min[i].resTime.currentTime.Minutes = 0
    min[i].resTime.currentTime.Seconds = 0
    min[i].resTime.sumTime.Hours = intconst
    min[i].resTime.sumTime.Minutes = 0
    min[i].resTime.sumTime.Seconds = 0
    visitedPoint[i] = 1
  }
} // Задает необходимые, для дальнейших вычислений, значение массивов минимальных путей и посещенных станций

func NewIterWithoutTime(min []ResultMatrix, visitedPoint []int) {
  for i, _ := range min {
    min[i].resPrice = maxPrice
    visitedPoint[i] = 1
  }
} // То же что и NewIter, но изменят только часть отвечающую за цену

func WayTime(DepartureStationId string, ArrivalStationId string, resultMatrix [][]ResultMatrix, initialMatrix [][]TrainLegs, uniqueStation []string, way []Way, infoTime [][]Time) {
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
  if (resultMatrix[start][finish].resTime.sumTime.Hours == intconst) {
    fmt.Println("No way")
    return
  }
  var timezero = Time{0, 0, 0}
  var str2 Time
  for finish != start {
    for i := 0; i < len(initialMatrix); i++ {
      str2 = Time{resultMatrix[start][i].resTime.sumTime.Hours%24, resultMatrix[start][i].resTime.sumTime.Minutes, resultMatrix[start][i].resTime.sumTime.Seconds}
      for j := 0; j < len(initialMatrix[i][finish].Trains); j++ {
        for z := 0; z < len(infoTime[start]); z++ {
          if (Difference(initialMatrix[i][finish].Trains[j].DepartureTime, initialMatrix[i][finish].Trains[j].ArrivalTime) != timezero) && SumTime(Difference(str2, initialMatrix[i][finish].Trains[j].DepartureTime), Difference(initialMatrix[i][finish].Trains[j].DepartureTime, initialMatrix[i][finish].Trains[j].ArrivalTime)) == DifferenceAll(resultMatrix[start][i].resTime.sumTime, resultMatrix[start][finish].resTime.sumTime) && DifferenceAll(resultMatrix[start][i].resTime.sumTime, resultMatrix[start][finish].resTime.sumTime) == infoTime[start][z] {
          infoTime[start][z] = timezero
          way[next].name = uniqueStation[i]
          way[next].id = initialMatrix[i][finish].Trains[j].TrainId
          next++
          finish = i
          break
          }
        }
      }
    }
  }
  showSlice(way)
  fmt.Printf("Time: Hours:%v, Minutes:%v, Seconds:%v\n", resultMatrix[start][print].resTime.sumTime.Hours, resultMatrix[start][print].resTime.sumTime.Minutes, resultMatrix[start][print].resTime.sumTime.Seconds)
  return
} // Восстанавливает кратчайший(по времени) путь с помощью матриц resultMatrix, initialMatrix и выводит его, а также время перезда

func AllShortestWayDijkstra(resultMatrix [][]ResultMatrix, initialMatrix [][]TrainLegs, min []ResultMatrix, visitedPoint []int, infoTime [][]Time, uniqueStation []string) {
  for i, _ := range min {
    ShortestWayDijkstraTimeInfo(initialMatrix, min, visitedPoint, infoTime, i, uniqueStation)
    NewIterWithoutTime(min, visitedPoint)
    ShortestWayDijkstraPrice(initialMatrix, min, visitedPoint, i)
    for j, _ := range min {
       resultMatrix[i][j] = min[j]
    }
    NewIter(min, visitedPoint)
  }
} // Находит наиболее выгодный(по цене и времени) проезд между каждой парой станций

func ShortestWayDijkstraTimeInfo(initialMatrix [][]TrainLegs, min []ResultMatrix, visitedPoint []int, infoTime [][]Time, stationindex int, uniqueStation []string) {
  min[stationindex].resTime.sumTime.Hours = 0
  var minIndex, per int
  var mins MinTime
  var temp, str1, str2 Time
  var minimum Time
  for minIndex < 10000 {
    minIndex = 10000
    minimum = Time{intconst, 0, 0}
    for i, _ := range min {
      if visitedPoint[i] == 1 && DifferenceSec(minimum, min[i].resTime.sumTime) < 0 {
        minimum = min[i].resTime.sumTime
        mins = min[i].resTime
        minIndex = i;
      }
    }
    str1 = Time{(mins.currentTime.Hours/24)*24, 0, 0}
    str2 = DifferenceAll(str1, mins.currentTime)
    if minIndex != 10000 {
      for i, _ := range initialMatrix[minIndex] {
        for j := 1; j < len(initialMatrix[minIndex][i].Trains); j++ {
          temp = SumTime(Difference(str2, initialMatrix[minIndex][i].Trains[j].DepartureTime), Difference(initialMatrix[minIndex][i].Trains[j].DepartureTime, initialMatrix[minIndex][i].Trains[j].ArrivalTime))
          str1 = SumTime(mins.currentTime, temp)
          per = DifferenceSec(str1, min[i].resTime.sumTime)
          if per > 0 {
            infoTime[stationindex] = append(infoTime[stationindex], temp)
            min[i].resTime.sumTime = str1
            min[i].resTime.currentTime = str1
          }
          str1 = Time{(mins.currentTime.Hours/24)*24, 0, 0}
        }
      }
      visitedPoint[minIndex] = 0
    }
  }} // Находит для станции под индексом stationindex кратчайший(по времени) путь до всех остальных станций исользуя алгоритм Дейкстры

func Initial(trainLegs TrainLegsXML) ([][]ResultMatrix, [][]TrainLegs, []ResultMatrix, []int, []string, [][]Time) {
  uniqueStation := make([]string, 0)
  uniqueStation = append(uniqueStation, trainLegs.TrainLegsSlice[0].DepartureStationId)
  uniqueStation = append(uniqueStation, trainLegs.TrainLegsSlice[0].ArrivalStationId)
  inside1 := true
  inside2 := true
  var num int
  var num64 float64
  var trainLeg TrainLeg
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
  initialMatrix := make([][]TrainLegs, len(uniqueStation))
  resultMatrix := make([][]ResultMatrix, len(uniqueStation))
  infoTime := make([][]Time, len(uniqueStation))
  for i, _ := range initialMatrix {
    initialMatrix[i] = make([]TrainLegs, len(uniqueStation))
    resultMatrix[i] = make([]ResultMatrix, len(uniqueStation))
    infoTime[i] = make([]Time, 0)
    for j, _ := range initialMatrix[i] {
      resultMatrix[i][j].resTime = MinTime{Time{0, 0, 0}, Time{0, 0, 0}}
    }
  }
  time1.Hours = 0
  time1.Minutes = 0
  time1.Seconds = 0
  time2.Hours = 0
  time2.Minutes = 0
  time2.Seconds = 0
  trainLeg = TrainLeg{ TrainId : 0, DepartureStationId : "", ArrivalStationId : "", Price : 0.0, DepartureTime : time1, ArrivalTime : time2 }
  for i, _ := range initialMatrix {
    for j, _ := range initialMatrix {
      initialMatrix[i][j].Trains = append(initialMatrix[i][j].Trains, trainLeg)
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
          num64, _ = strconv.ParseFloat(trainLegs.TrainLegsSlice[i].Price, 32)
          trainLeg = TrainLeg{ TrainId : int(num), DepartureStationId : trainLegs.TrainLegsSlice[i].DepartureStationId, ArrivalStationId : trainLegs.TrainLegsSlice[i].ArrivalStationId, Price : float64(num64), DepartureTime : time1, ArrivalTime : time2 }
          initialMatrix[j][z].Trains = append(initialMatrix[j][z].Trains, trainLeg)
        }
      }
    }
  }

  min := make([]ResultMatrix, len(uniqueStation))
  visitedPoint := make([]int, len(uniqueStation))
  NewIter(min, visitedPoint)
  return resultMatrix, initialMatrix, min, visitedPoint, uniqueStation, infoTime
} // На основе имеющейся информации создает и возвращает матрицы кратчайших путей между каждой парой станций(пока не заполнена), связей, мин. путей(одномерная), пройденных пунктов, уникальных имен пунктов и различных времен переезда между станциями

func WayPrice(DepartureStationId string, ArrivalStationId string, resultMatrix [][]ResultMatrix, initialMatrix [][]TrainLegs, uniqueStation []string, way []Way) {
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
  if (resultMatrix[start][finish].resPrice == maxPrice) {
    fmt.Println("No way")
    return
  }
  for finish != start {
    for i := 0; i < len(initialMatrix); i++ {
      for j, _ := range initialMatrix[i][finish].Trains {
        if (initialMatrix[i][finish].Trains[j].Price != 0.0) && (initialMatrix[i][finish].Trains[j].Price == (resultMatrix[start][finish].resPrice - resultMatrix[start][i].resPrice)) {
          way[next].name = uniqueStation[i]
          way[next].id = initialMatrix[i][finish].Trains[j].TrainId
          next++
          finish = i
          break
        }
      }
    }
  }
  showSlice(way)
  fmt.Println("Price: ", resultMatrix[start][print].resPrice, "\n")
  return
} //Восстанавливает кратчайший(по цене) путь с помощью матриц resultMatrix, initialMatrix и выводит его, а также цену переезда

func ShortestWayDijkstraPrice(initialMatrix [][]TrainLegs, min[]ResultMatrix, visitedPoint []int, stationindex int) {
  min[stationindex].resPrice = 0
  var minIndex int
  var mins float64
  var temp float64
  for minIndex < 10000 {
    minIndex = 10000
    mins = maxPrice
    for i, _ := range min {
      if visitedPoint[i] == 1 && (min[i].resPrice - mins) < 0 {
        mins = min[i].resPrice
        minIndex = i;
      }
    }
    if minIndex != 10000 {
      for i, _ := range initialMatrix {
        for j, _ := range initialMatrix[minIndex][i].Trains {
          if initialMatrix[minIndex][i].Trains[j].Price > 0.0 {
            temp = mins + initialMatrix[minIndex][i].Trains[j].Price
            if temp < min[i].resPrice {
              min[i].resPrice = temp
            }
          }
        }
      }
      visitedPoint[minIndex] = 0
    }
  }
} // Находит для станции под индексом stationindex кратчайший(по цене) путь до всех остальных станций исользуя алгоритм Дейкстры

func ParseXML(fileName string) TrainLegsXML {
  xmlFile, err := os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}
	defer xmlFile.Close()
  byteValue, err := ioutil.ReadAll(xmlFile)
  if err != nil {
		log.Fatal(err)
	}
  var trainLegs TrainLegsXML
  err = xml.Unmarshal(byteValue, &trainLegs)
  if err != nil {
    log.Fatal(err)
  }
  return trainLegs
} // Считывает название файла и возвращает структуру, в которой записана необходимая информация из XML файла

type ResultMatrix struct {
  resPrice float64
  resTime MinTime
} // Используется в матрице resultMatrix

type Way struct {
  name string
  id int
} // Служит для восстановления пути

type TrainLegs struct {
  Trains []TrainLeg
} // Используется в матрице initialMatrix хранит список различных прямых переездов между парой станций

type TrainLeg struct {
  TrainId int
  DepartureStationId string
  ArrivalStationId string
  Price float64
  DepartureTime Time
  ArrivalTime Time
} // Используется в матрице initialMatrix хранит 1 перезд между парой станций

type Time struct {
  Hours int
  Minutes int
  Seconds int
} // Хранит время в часах, минутах и секундах

type MinTime struct {
  currentTime Time
  sumTime Time
} // Используется в матрице resultMatrix хранит текущее время и суммарное время необходимое чтобы добраться от одной станции до другой

type TrainLegsXML struct {
	XMLName xml.Name `xml:"TrainLegs"`
	TrainLegsSlice   []TrainLegXML   `xml:"TrainLeg"`
} // Хранит елементы XML файла

type TrainLegXML struct {
  XMLName xml.Name `xml:"TrainLeg"`
	TrainId string `xml:"TrainId,attr"`
	DepartureStationId string `xml:"DepartureStationId,attr"`
	ArrivalStationId string `xml:"ArrivalStationId,attr"`
  Price string `xml:"Price,attr"`
  ArrivalTimeString string `xml:"ArrivalTimeString,attr"`
  DepartureTimeString string `xml:"DepartureTimeString,attr"`
} // Хранит элемент, в котором записана информация о переезде

func main() {
  var stationDep, stationArr string // Названия пунктов отправления и прибытия соответственно
  fmt.Print("Введите начальную станцию: ")
  fmt.Fscan(os.Stdin, &stationDep) // Считвание с консоли названия пункта отправления
  fmt.Print("Введите конечную станцию: ")
  fmt.Fscan(os.Stdin, &stationArr) // Считвание с консоли названия пункта прибытия
  var trainLegs TrainLegsXML // Структура, в которую считывается информация из XML файла
  var fileName string = "data.xml" // имя файла
  trainLegs = ParseXML(fileName)
  var uniqueStation []string // Хранит уникальные названия станций, в матрицах resultMatrix, initialMatrix, индексы совпадают с индексами в этой матрице
  var visitedPoint []int // Используется в некоторых функциях хранит информацию о том, была ли посещена станция(0) или нет(1)
  var initialMatrix [][]TrainLegs // Хранит информацию о перездах между каждой парой станций, элемент initialMatrix[i][j].Trains[z] хранит информацию о переезде из и i в j по маршруту z
  var resultMatrix [][]ResultMatrix // Хранит конечную информацию о перездах между каждой парой станций, элемент initialMatrix[i][j] хранит информацию о переезде из и i в j
  var min []ResultMatrix // Необходима для заполнения resultMatrix
  var infoTime [][]Time // Хранит уникальные времени перезда(время простоя + время пути) между парами станций, необходима для восстановления пути
  resultMatrix, initialMatrix, min, visitedPoint, uniqueStation, infoTime = Initial(trainLegs)
  AllShortestWayDijkstra(resultMatrix, initialMatrix, min, visitedPoint, infoTime, uniqueStation)
  wayTime := make([]Way, len(uniqueStation))
  wayPrice := make([]Way, len(uniqueStation))
  WayTime(stationDep, stationArr, resultMatrix, initialMatrix, uniqueStation, wayTime, infoTime)
  WayPrice(stationDep, stationArr, resultMatrix, initialMatrix, uniqueStation, wayPrice)
}
