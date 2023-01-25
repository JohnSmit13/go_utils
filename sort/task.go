
//Задача: реализовать утилиту Sort. Утилита должна поддерживать флаги '-n','-k','-u','-r'.
 
package main

import (
    "bufio"
    "fmt"
    "log"
    "os"
    "strconv"
    "time"
    "math/rand"
)

/******************
 *      Vars      *
 ******************/
var column int
var Sorter sorter
var files []string

func main() {

    if len(files) == 0 {
        fmt.Println("No file")
        return
    }

    for i := range files {
        f, err := os.Open(files[i])

        if err != nil {
            log.Fatal(err)
        }
        defer f.Close()
        scanner := bufio.NewScanner(f)

        // читаем все строки из файла в слайс строк
        var fileStrings []string
        for scanner.Scan() {
            fileStrings = append(fileStrings,scanner.Text())
        }

        // для хранения строк с количеством колонок < Sorter.col
        var shortStrings []columnString 
        // для хранения строк с количеством колонок >= Sorter.col
        var stringsToSort []columnString

        // если ищем не по первой колонке то для каждой строки
        // нужно найти символ с которого начинается column-ая колонка.
        // если в строке нет такого количества колонок,
        // "отправляем" эту строку ко всем коротким строкам
        // их отсортируем отдельно и выведем первыми
        if Sorter.col != 1 {
            for i := range fileStrings {
                if pos, ok := countColumns(fileStrings[i], Sorter.col); ok {
                    val := columnString{&fileStrings[i],pos}
                    stringsToSort = append(stringsToSort, val)
                } else {
                    val := columnString{&fileStrings[i],0}
                    shortStrings = append(shortStrings,val)
                }
            }
            Sorter.Sort(shortStrings)
            Sorter.Sort(stringsToSort)
            fname := "sorted_" + string(files[i])
            res,ok := os.Create(fname)
            if ok != nil {
                fmt.Fprintf(os.Stderr, "can't create file %s\n sorting result: \n\n",fname)
            } else {
                for i := range shortStrings {
                    val := shortStrings[i].str
                    fmt.Fprintln(res,*val)
                }
                for i := range stringsToSort {
                    val := stringsToSort[i].str
                    fmt.Fprintln(res,*val)
                }
            }
        } else { // column == 1
            for i:= range fileStrings {
                val := columnString{&fileStrings[i],0}
                stringsToSort = append(stringsToSort,val)
            }

            Sorter.Sort(stringsToSort)

            fname := "sorted_" + string(files[i])
            res,ok := os.Create(fname)

            // если не удалось создать файл, то результат выведем на экран
            if ok != nil { 
                fmt.Fprintf(os.Stderr, "can't create file %s\n sorting result: \n\n",fname)
                if !Sorter.uniq { // записываем повторения
                    for i := range stringsToSort {
                        val := stringsToSort[i].str
                        fmt.Printf("%v\n",*val)
                    }
                } else { // только уникальные строки
                    for i := 0; i < len(stringsToSort) -1; i++ {
                        if *stringsToSort[i].str == *stringsToSort[i+1].str {
                            // строка равна следующей, пропускаем ее
                        } else {
                            fmt.Printf("%v\n",*stringsToSort[i].str)
                        }
                    }
                    // записываем последнюю строку
                    fmt.Printf("%v\n",stringsToSort[len(stringsToSort)-1])
                }
            } else { // файл создан
                if !Sorter.uniq { // записываем повторения
                    for i := range stringsToSort {
                        val := stringsToSort[i].str
                        fmt.Fprintln(res,*val)
                    }
                } else { // только уникальные строки
                    for i := 0; i < len(stringsToSort) -1; i++ {
                        if *stringsToSort[i].str == *stringsToSort[i+1].str {
                            // строка равна следующей, пропускаем ее
                        } else {
                            fmt.Fprintln(res,*&stringsToSort[i].str)
                        }
                    }
                    // записываем последнюю строку
                    fmt.Fprintln(res,*&stringsToSort[len(stringsToSort)-1].str)
                }
            }
        }
    }// for range files
}

func countColumns(s string, c int) (int,bool) {
    if len(s) == 0 {return 0, false}
    if c == 1 {return 0, true}
    colCount := 1
    pos := 0

    sep := ' ' // разделитель колонок по умолчанию пробел
    for pos < len(s) && colCount != c {
        // проходим по строке и ищем пробел. 
        // как только нашли, пропускаем все рядом стоящие пробелы,
        // т.к. несколько пробелов разделяют строку также как и один.
        if s[pos] == byte(sep) {
            pos++ 
            colCount++
            if pos == len(s) {continue}
            for s[pos] == byte(sep) {
                pos++
            }
            continue
        }
        pos++
    }
    if colCount != c {
        return 0, false
    } else {
        return pos, true
    }
}

func init() {
    // был ли флаг -k, нужна для вывода сообщения, 
    // в случае некорректного номера колонки
    col := false 

    // проходим по всем аргументам, начиная со второго
    // и "включаем" соответствующие флаги
    for i := 1; i < len(os.Args); i++ {
        v := os.Args[i]
        if v == "-k" {
            col = true
            if i == len(os.Args)-1 {
                fmt.Fprint(os.Stderr,"missed column number")
                os.Exit(1)
            } else {
                if n,e := toInt(os.Args[i+1]); e {
                    Sorter.col = n
                    i++
                    continue
                } else {
                    fmt.Fprint(os.Stderr,"can't convert column number to int")
                    os.Exit(1)
                }
            }
        } else if v == "-n" {
            Sorter.num = true
            continue
        } else if v == "-r" {
            Sorter.rev = true
            continue
        } else if v == "-u" {
            Sorter.uniq = true
            continue
        } else {
            // все что не флаг - имя файла
            files = append(files, v)
            continue
        }
    }

    if Sorter.col < 1 {
        if col {
            fmt.Println("column can't be less then 1, it was set to 1")
        }
        Sorter.col = 1
    }

    if len(files) == 0 {
        fmt.Fprint(os.Stderr, "no file\n")
        os.Exit(1)
    }
}

// types 

// хранит строку и позицию с которой начинается слово 
// по которому будем сортировать,
// необходим для корректной обработки коротких строк 
// при больших значениях column
type columnString struct {
    str *string
    pos int
}

type sorter struct {
    num  bool
    uniq bool
    rev  bool
    col  int
}

func(s sorter) CmpString(x,y string) bool {
    if s.rev {
        return x>y
    } else {
        return x<y
    }
}

func(s sorter) CmpInt(x,y int) bool {
    if s.rev {
        return x>y
    } else {
        return x<y
    }
}

func(s sorter) Sort(strs []columnString) []columnString {

    if len(strs) < 2 {
        return strs
    }

    seed := int64(time.Now().UnixNano())
    rand.Seed(seed)
    left, right := 0, len(strs)-1

    pivot := rand.Int() % len(strs)

    strs[pivot], strs[right] = strs[right], strs[pivot]

    if !s.num {
        // нечисловая сортировка
        for i := range strs {
            p1 := strs[i].pos
            v1 := strs[i].str
            p2 := strs[right].pos
            v2 := strs[right].str
            if s.CmpString((*v1)[p1:], (*v2)[p2:]) {
                strs[left], strs[i] = strs[i], strs[left]
                left++
            }
        }
        strs[left], strs[right] = strs[right], strs[left]
        s.Sort(strs[:left])
        s.Sort(strs[left+1:])
    } else { // числовая сортировка
        for i := range strs {
            // первое число
            p1 := strs[i].pos
            v  := strs[i].str
            v1 := (*v)[p1:]
            valToCmp1, ok := toInt(v1)
            if !ok { 
                // v1 не число, пока ничего не делаем, отсортируем такие строки позже
            }
            // второе число
            p2 := strs[right].pos
            v   = strs[right].str
            v2 := (*v)[p2:]
            valToCmp2, ok := toInt(v2)
            if !ok { 
                // v2 не число, пока ничего не делаем, отсортируем такие строки позже
            }
            if s.CmpInt(valToCmp1,valToCmp2) {
                strs[left], strs[i] = strs[i], strs[left]
                left++
            }
        }

        strs[left], strs[right] = strs[right], strs[left]
        s.Sort(strs[:left])
        s.Sort(strs[left+1:])

        // строки которые не были переведены в число не отсортированы
        end := 0
        for end < len(strs) {
            v := strs[end].str
            p := strs[end].pos
            _, ok := toInt((*v)[p:])
            if !ok {
                end++
            } else { break }
        }
        s.num = false // чтобы пропустить попытку сконвертировать в число строки
                      // которые не конвертируются
        s.Sort(strs[:end])
    }

    return strs
}

func toInt(x string) (int, bool) {
    var res int
    var ok error
    
    end := 0 //первый символ который не является цифрой
    for end < len(x) && x[end] >= 48 && x[end] <= 57 {
        end++
    }

    res, ok = strconv.Atoi(x[:end])
    if ok != nil {
        return 0, false
    }

    return res, true
}
