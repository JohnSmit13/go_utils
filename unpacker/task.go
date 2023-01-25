// Задача: Реализовать простую утилиту для "распаковки" строки, например:
// "a5" -> "aaaaa"
// "b2" -> "bb"
// "c0" -> ""
// Дополнительно реализовать поддержку обратного слэша, например:
// "\5" -> "5"
// "\55"-> "55555"
// "\\" -> "\"
// "\\\\" -> "\\"

package main

import (
	"fmt"
	"os"
    "errors"
)

type dype struct {
    d int
}

func main() {
    r1,ok1 := unpack(`add0r4d\3c\3\\4`) // adrrrrd3c3\\\\
    r2,ok2 := unpack(`a10`)             // BadString
    r3,ok3 := unpack(`a1`)              // a
    r4,ok4 := unpack(`\\`)              // \
    r5,ok5 := unpack(`ad0e`)            // ae
    r6,ok6 := unpack(`\\\\\\\e`)        // \\\e
    if ok1 == nil {
        fmt.Printf("add0r4d\\3c\\3\\\\4 -> %s\n", r1)
    } else {
        fmt.Fprintf(os.Stderr, "%s", ok1.Error())
    }
    if ok2 == nil {
        fmt.Printf("%s\n", r2)
    } else {
        fmt.Fprintf(os.Stderr, "a10 -> %s\n", ok2.Error())
    }
    if ok3 == nil {
        fmt.Printf("a1 -> %s\n", r3)
    } else {
        fmt.Fprintf(os.Stderr, "%s", ok3.Error())
    }
    if ok4 == nil {
        fmt.Printf("\\\\ -> %s\n", r4)
    } else {
        fmt.Fprintf(os.Stderr, "%s", ok4.Error())
    }
    if ok5 == nil {
        fmt.Printf("ad0e -> %s\n", r5)
    } else {
        fmt.Fprintf(os.Stderr, "%s", ok5.Error())
    }
    if ok6 == nil {
        fmt.Printf("\\\\\\\\\\\\\\e -> %s\n", r6)
    } else {
        fmt.Fprintf(os.Stderr, "%s", ok6.Error())
    }
}

func unpack(s string) (string, error){
    var r = []rune(s)
    var res []rune
    var haveLast bool = false // есть ли что добавить в результат?
    var last rune             // руна для добавления результата

    for i := 0; i < len(r); i++ {
        if r[i] == '\\' { // backslash
            if haveLast { // если есть символ добавляем его, чтобы не потерять
                res = append(res,last)
                haveLast = false
            }
            // если 2ой символ после бэкслэша число
            // следующий символ нужно добавить несколько раз
            // для этого "ставим в очередь" этот символ,
            // и увеличиваем индекс на 1 чтобы пропустить этот символ
            // и добавить его сразу несколько раз на следующей итерации
            if i < len(r) - 2 && isNumber(r[i+2]) {
                last = r[i+1]
                haveLast = true
                i++
                continue
            // если 2ой символ после бэкслэша не число
            // то ничего не делаем. \we превратится в we
            } else if i < len(r) -1 {
                last = r[i+1]
                haveLast = true
                i++
                continue
            } else {// последний символ в строке был '\', это некорректная строка
                return "", errors.New("BadString")
            }
        } else if isNumber(r[i]) { // если символ число
            if haveLast { // если есть что добавить, добавляем count раз
                count := r[i] - '0'
                appendSlice := make([]rune,count)
                for j := range appendSlice {
                    appendSlice[j] = last
                }
                res = append(res,appendSlice...)
                haveLast = false
                continue
            } else { // если нечего добавлять то это некорректная строка
                return "", errors.New("BadString")
            }
        } else { // все остальное это руна
            if haveLast { // если есть что добавить - добавляем
                res = append(res,last)
                last = r[i] // "ставим в очередь" текущий символ
                haveLast = true // можно опустить, ничего не поменяется
            } else {
                last = r[i]
                haveLast = true
            }
        }
    } // for
    if haveLast { // эта проверка нужна чтобы вывести последний символ,
                  // если предпоследний символ был '\'
        res = append(res,last)
    }
    return string(res),nil
}

func isNumber(r rune) bool {
    n := r-'0'
    return n >= 0 && n <= 9
}
