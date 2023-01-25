
// Задача: Написать программу которая из массива слов будет находить все анаграммы
// группировать их и выдавать результат в виде
// $(первое_слово_из_группы_анаграмм):$(все_слова_группы_анаграмм).

package main

import (
	"fmt"
	"sort"
	"os"
)

/****************
 *    Types     *
 ****************/

// для хранения массива слов
type Words []string

// общий вид для слов, сколько и каких символов есть в этом слове,
// строковое представление данного типа будет ключем в промежуточной мапе
type form []symbolCount

// для подсчета количества символов в слове
type symbolCount struct {
    r     rune
    count uint8
}

// является value в промежуточной мапе
type record struct {
    // храним первое слово, чтобы использовать его как ключ в результате.
    firstWord string 
    // слайс всех слов, включая первое
    words Words
}

// для хранения промежуточных результатов
type book map[string] *record

/****************
 *     Vars     *
 ****************/

var words = Words {
    "пятак", "пятка", "тяпка",
    // слово "слиток" не случайно встречается дважды
    "столик", "листок", "слиток","слиток",
    "дура", "руда", "удар",
    "рост", "сорт", "торс", "трос",
    "автор", "отвар", "рвота", "товар",
    "аскет", "секта", "сетка", "тесак",
    "казан", "казна", "наказ",
    "карат", "карта", "катар",
    "клоун", "колун", "кулон", "уклон",
    "лапша", "палаш", "шпала",
    "носик", "носки", "оникс",
    "отвес", "отсев", "совет",
    "монета", "немота", "отмена",
    "колосок", "осколок", "соколок",
    "пасечник", "песчаник", "песчинка",
    "воспрещение", "всепрощение", "просвещение",
}

/****************
 *     Main     *
 ****************/
func main() {
    answer := getResult(words)
    for k,v := range answer {
        fmt.Printf("%s : %s\n",k,v)
    }
}

/****************
 *     Init     *
 ****************/
func init() {
    sort.Sort(words)
}

/****************
 *  Functions   *
 ****************/
func getResult(w Words) map[string][]string {
    var tmpRes = make(book)
    var Res = make(map[string][]string) // ответ
    for i := range w { 
        // находим "форму" слова
        tmp := toForm(w[i])
        haveWord := false
        // если такая форма уже есть в промежуточной мапе, 
        // и этого слова еще нет в слайсе, добавим слово в запись
        v,ok := tmpRes[tmp] 
        if ok {
            for j := range tmpRes[tmp].words {
                // если в слайсе нашли такое слово, не добавляем его
                if tmpRes[tmp].words[j] == w[i] {
                    haveWord = true
                    break
                }
            }
            if !haveWord {
                v.words = append(v.words,w[i])
                //fmt.Println(v.words)
            }
        } else { // иначе создаем такую запись
            var s = []string{w[i]}
            tmpRes[tmp] = &record{w[i],s}
        }
    }

    for _,v := range tmpRes {
        // в ответ не должны попасть множества из одного элемента
        // поэтому проверяем длинну слайса слов
        if len(v.words) > 1 {
            // если слов больше 1, добавляем это множество в ответ
            sort.Sort(v.words)
            Res[v.firstWord] = v.words
        }
    }
    return Res
}

// для нахождения общего вида слова
func toForm(s string) string {
    var tmpS = []rune(s)
    var res form
    for i := range tmpS {
        if p,ok := inRange(rune(tmpS[i]),&res); ok {
            res[p].count++
        } else { 
            res = append(res, symbolCount{tmpS[i],1})
        }
    }
    sort.Sort(res)
    return res.toString()
}

func inRange(r rune, f *form) (int,bool) {
    var res int = 0
    for res < f.Len() {
        if (*f)[res].r == r {
            return res,true
        }
        res++
    }
    return 0, false
}

/****************
 *   Methods    *
 ****************/

// Words
func(w Words) Len() int {
    return len(w)
}

func(w Words) Less(i,j int) bool {
    return w[i] < w[j]
}

func(w Words) Swap(i,j int) {
    w[i],w[j] = w[j],w[i]
}

// form
func(f form) Len() int {
    return len(f)
}

func(f form) Less(i,j int) bool {
    return f[i].r < f[j].r
}

func(f form) Swap(i,j int) {
    f[i],f[j] = f[j],f[i]
}

func(f form) toString() string {
    var res []rune
    for i := range f {
        if f[i].r == 0 || f[i].count == 0 {
            fmt.Printf("rune can't be 0, as well as it's count\n")
            os.Exit(1)
        }
        res = append(res,f[i].r,rune(f[i].count + '0'))
    }
    return string(res)
}
