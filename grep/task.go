
//Задача: реализовать утилиту Grep. Утилита должно поддерживать флаги '-A','-B','-C','-c','-i','-v','-F','-n'.

package main

import (
	"bufio"
	"fmt"
	"os" 
	"regexp"
	"strconv"
	"strings"
)

var Grep greper

func main() {
    Grep.find()
}

/****************
 *     Init     *
 ****************/
func init() {
    // сначала обрабатываем все опции
    i := 1
    for i < len(os.Args) {

        if os.Args[i] == "-A" {
            fmt.Printf("have -A\n")
            if i == len(os.Args) -1 {
                fmt.Fprint(os.Stderr,"missing number of lines \"after\"\n")
                os.Exit(1)
            } else {
                num, ok := strconv.Atoi(os.Args[i+1])
                if ok != nil {
                    fmt.Fprint(os.Stderr,"can't convert number\n")
                    os.Exit(1)
                } else {
                    Grep.after = num
                }
                i += 2
                continue
            }
        } else if os.Args[i] == "-B" {
            fmt.Printf("have -B\n")
            if i == len(os.Args) -1 {
                fmt.Fprint(os.Stderr,"missing number of lines \"before\"\n")
                os.Exit(1)
            } else {
                num, ok := strconv.Atoi(os.Args[i+1])
                if ok != nil {
                    fmt.Fprint(os.Stderr,"can't convert number\n")
                    os.Exit(1)
                } else {
                    Grep.before = num
                }
                i += 2
                continue
            }
        } else if os.Args[i] == "-C" {
            fmt.Printf("have -C\n")
            if i == len(os.Args) -1 {
                fmt.Fprint(os.Stderr,"missing number of lines of context\n")
                os.Exit(1)
            } else {
                num, ok := strconv.Atoi(os.Args[i+1])
                if ok != nil {
                    fmt.Fprint(os.Stderr,"can't convert number")
                    os.Exit(1)
                } else {
                    Grep.context = num
                }
                i += 2
                continue
            }
        } else if os.Args[i] == "-c" {
            fmt.Printf("have -c\n")
            Grep.count = true;
            i++
            continue
        } else if os.Args[i] == "-i" {
            fmt.Printf("have -i\n")
            Grep.ignore = true;
            i++
            continue
        } else if os.Args[i] == "-v" {
            fmt.Printf("have -v\n")
            Grep.invert = true;
            i++
            continue
        } else if os.Args[i] == "-F" {
            fmt.Printf("have -F\n")
            Grep.fixed = true;
            i++
            continue
        } else if os.Args[i] == "-n" {
            fmt.Printf("have -n\n")
            Grep.lineNum = true;
            i++
            continue
        } else {
            // все опции должны быть до паттерна
            // если не распознали опцию, значит это паттерн
            Grep.pattern = os.Args[i]
            i++
            break
        }
    } // for
    // далее все аргументы это файлы
    if i == len(os.Args) {
        fmt.Printf("go to interact\n")
        // нет файла, переходим в интерактивный режим
        Grep.inter = true
    } else {
        for i < len(os.Args) {
            Grep.files = append(Grep.files,os.Args[i])
            i++
        }
    }
    // если был флаг -i приведем паттерн к нижнему регистру
    if Grep.ignore {
        Grep.pattern = strings.ToLower(Grep.pattern)
    }
    fmt.Printf("g.count: %t\n", Grep.count)
}

/*****************
 *    Types      *
 *****************/
// greper хранит все флаги и выполняет основную работу
type greper struct {
    after   int
    before  int
    context int
    count   bool
    ignore  bool
    invert  bool
    fixed   bool
    lineNum bool

    // если нет файла, переходим в интерактивный режим
    // в этом режиме принимаем строки из os.Stdin 
    // и ищем паттерн среди этих строк
    inter   bool

    pattern string
    files   []string
    // указывает какой файл сейчас проходит обработку, нужен для вывода
    // имени файла перед каждой строкой если файлов было несколько
    //index   int

}
/*****************
 *    Methods    *
 *****************/

// данный метод проверяет строку на совпадение с паттерном
// если включен флаг -F проверка идет с помощью strings.Contains
func(g greper) match(s string) (bool, error) {
    //fmt.Printf("g.match %t\n", g.fixed)
    if g.ignore {
        s = strings.ToLower(s)
    }
    if g.fixed {
        return strings.Contains(s, g.pattern), nil
    } else {
        v,ok := regexp.Match(g.pattern, []byte(s))
        /*fmt.Printf("pattern:%s\n",g.pattern)
        fmt.Printf("string :%s\n",s)
        fmt.Printf("regexp.Match(s,[]byte(g.pattern)) %t\n",v)*/
        return v,ok
    }
}

// данный метод выводит результат поиска
// принимает мапу с индексами строк которые нужно вывести,
// и слайс всех строк файла. проходя по всем строкам
// проверяем по мапе нужно ли выводить эту строку.
func(g greper) showResult(m map[int]bool, s []string, f int) {
    fmt.Printf("showing result\n")
    // выводим совпадения
    if !g.invert {
        for i := 0; i < len(s); i++ {
            // в мапу мы записываем только true,
            // а значит нет смысла проверять значение
            // если нет даже записи о нем в мапе.
            if _,ok := m[i]; ok {
                if g.lineNum {
                    fmt.Println(g.getFileName(f),i,s[i])
                } else {
                    fmt.Println(g.getFileName(f),s[i])
                }
            }
        }
    // выводим все кроме совпадений
    } else {
        for i := 0; i < len(s); i++ {
            // выводим строку только если ее индекс не в мапе
            if _,ok := m[i]; !ok {
                if g.lineNum {
                    fmt.Println(g.getFileName(f),i,s[i])
                } else {
                    fmt.Println(g.getFileName(f),s[i])
                }
            }
        }
    }
}

// данный метод возвращает имя файла с двоеточием в конце
// если файл один, то возвращает пустую строку, 
// т.к. нет необходимости указывать в каком файле было найдено совпадение
func(g greper) getFileName(i int) string {
    if len(g.files) > 1 {
        res := g.files[i] + ": "
        return res
    } else { return "" }
}

// данный метод выполняет поиск паттерна по всем файлам
// и выводит результат поиска в os.stdout
// для вывода используется метод showResult
func(g greper) find() {
    fmt.Printf("find\n")
    // в интерактивном режиме читаем из os.Stdin и выводим если совпадает
    if g.inter {
        fmt.Printf("interact\n")
        scan := bufio.NewScanner(os.Stdin)
        for scan.Scan() {
            if ok,_ := g.match(scan.Text()); ok {
                fmt.Println(scan.Text())
            }
        }
    }
    for f := range g.files {
        fmt.Printf("g.file: %s\n", g.files[f])
        // хранятся адреса тех строк которые нужно вывести в качестве результата
        // булевая переменная в значении показывает нужно ли выводить данную строку
        res := make(map[int]bool)
        // слайс всех строк в файле
        var s []string

        file,ok := os.Open(g.files[f])
        if ok != nil {
            fmt.Fprintf(os.Stderr,"%s",ok.Error())
        } else {
            scan := bufio.NewScanner(file)
            for scan.Scan() {
                s = append(s,scan.Text())
            }
        }

        // если нужно посчитать строки то вне зависимости от флагов -B, -A, -C
        // считаем только те строки которые подходят по паттерну.
        if g.count {
            fmt.Printf("count\n")
            var count int
            for i := range s {
                fmt.Printf("s[i]: %s\n",s[i])
                if r,ok := g.match(s[i]); ok == nil {
                    if r {
                        count++
                    }
                } else {
                    fmt.Fprintln(os.Stderr, ok.Error())
                }
            }
            fmt.Printf("%s%d\n",g.getFileName(f),count)
            continue
        }
        for i := 0; i < len(s); i++ {
            if r,ok := g.match(s[i]); ok == nil {
                if r { // строка подходит по паттерну
                    // добавляем ее в результат
                    res[i] = true
                    // проверяем флаги влияющие на количество выводимых строк
                    if g.before != 0 {
                        // добавляем в результат g.before строк до строки i
                        // начиная с позиции start
                        var start int
                        if (i - g.before) < 0 {
                            start = 0
                        } else {
                            start = i - g.before
                        }
                        for start < i {
                            res[start] = true
                            start++
                        }
                    }
                    if g.after != 0 {
                        // добавляем в результат g.after строк после строки i
                        // позиция последней строки для добавления
                        var end int
                        if g.after + i > len(s)-1 {
                            end = len(s) -1
                        } else {
                            end = g.after + i
                        }

                        for end != i {
                            res[end] = true
                            end--
                        }
                    }

                    if g.context != 0 {
                        // добавляем в результат g.context строк до строки i
                        // начиная с позиции start
                        var start int
                        if (i - g.context) < 0 {
                            start = 0
                        } else {
                            start = i - g.context
                        }
                        for start < i {
                            res[start] = true
                            start++
                        }
                        // добавляем в результат g.context строк после строки i
                        // позиция последней строки для добавления
                        var end int
                        if g.context + i > len(s)-1 {
                            end = len(s) -1
                        } else {
                            end = g.context + i
                        }
                        for end != i {
                            res[end] = true
                            end--
                        }
                    }
                } // if r
            } else { //g.match вернул ошибку
                fmt.Fprintln(os.Stderr, ok.Error())
            }
        }
        g.showResult(res,s,f)
    }
}
