// Задача: Реализовать собственную shell-утилиту(терминал). 
// Утилита должна поддерживать следующие команды:
// cd,pwd,echo,kill,ps,exit,quit, а также запускать процесс.
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"
)

func main() {
    fmt.Printf(">")
    scan := bufio.NewScanner(os.Stdin)
    for scan.Scan() {
        childPid := 0
        var scannedString = scan.Text()
        scannedWords := strings.Split(scannedString," ")
        cmd := scannedWords[0]
        switch cmd {
            // команда cd без аргументов совершает переход в домашнюю директорию
            case "cd" : {
                // cd без аргументов
                if len(scannedWords) == 1 {
                    fmt.Println("no args")
                    if ok := os.Chdir(os.Getenv("HOME")); ok == nil {
                        os.Setenv("PWD",os.Getenv("HOME"))
                        fmt.Println(os.Getenv("PWD"))
                    } else {
                        fmt.Println(ok.Error())
                    }
                // читаем только первый аргумент после cd
                } else {
                    dir := scannedWords[1]
                    switch dir {
                        // "выходим" из текущей директории
                        case ".." : {
                            dirWords := strings.Split(os.Getenv("PWD"),"/")
                            // отбрасываем последнюю директорию
                            dirWords = dirWords[:len(dirWords)-1]
                            dstDir := "/" + strings.Join(dirWords,"/")
                            if ok := os.Chdir(dstDir); ok == nil {
                                os.Setenv("PWD",dstDir)
                                fmt.Println(os.Getenv("PWD"))
                            } else {
                                fmt.Println(ok.Error())
                            }
                        }
                        // нужно выделить этот случай как отдельный,
                        // но при этом делать ничего не нужно
                        case "." : {
                        }
                        // пробуем перейти по указанному пути
                        default : {
                            // "путь назначения"
                            var dstDir string
                            // если путь абсолютный, он и будет конечным путем
    //fmt.Printf(">")
                            if scannedWords[1][0] == '/' {
                                dstDir = scannedWords[1]
                            } else {
                                // иначе относительно текущей директории
                                tmpDir := os.Getenv("PWD")
                                // чтобы избежать двойного слэша
                                // при переходи из корня
                                if tmpDir == "/" {
                                    dstDir = "/" + scannedWords[1]
                                } else {
                                    dstDir = tmpDir + "/" + scannedWords[1]
                                }
                            }
                            // перед переходом нужно убедиться что dstDir
                            // это путь а не файл
                            if dirInfo,ok := os.Stat(dstDir); ok == nil {
                                if dirInfo.IsDir() {
                                    os.Setenv("PWD",dstDir)
                                    os.Chdir(dstDir)
                                    fmt.Println(os.Getenv("PWD"))
                                }
                                // если dstDir это файл, то ничего не делаем
                            } else {
                                fmt.Println(ok.Error())
                            }
                        }
                    }
                }
            }
            case "pwd" : {
                fmt.Println(os.Getenv("PWD"))
            }
            case "echo" : {
                for i := 1; i < len(scannedWords); i++ {
                    if scannedWords[i][0] == '$' {
                        fmt.Printf("%v ",os.Getenv(scannedWords[i][1:]))
                    } else {
                        fmt.Printf("%v ", scannedWords[i])
                    }
                }
                fmt.Printf("\n")
            }
            // "убийство" процесса происходит по его PID
            case "kill" : {
                pid,_ := strconv.Atoi(scannedWords[1])
                proc,ok := os.FindProcess(pid)
                if ok != nil {
                    fmt.Println(ok.Error())
                } else {
                    proc.Kill()
                }
            }
            // просто выводит содержимое /proc как это вывел бы ls
            // другим вариантом реализации было-бы игнорирование данного
            // case'а т.к. в этом случае вызвался бы ps линукса
            case "ps" : {
                attr := syscall.ProcAttr {
                    Dir : "/proc",
                    Files : []uintptr{0, 1, 2},
                    Env : os.Environ(),
                }
                var ok error
                childPid, ok = syscall.ForkExec("/usr/bin/ls", []string{}, &attr)
                if ok != nil {
                    //fmt.Println(os.Args[0])
                    fmt.Println(ok.Error())
                }
            }
            case "exit","quit" : {
                os.Exit(0)
            }
            // пытаемся запустить процесс
            default : {
                // все пути перечисленные в переменной окружения PATH
                var pathsToFile []string = strings.Split(os.Getenv("PATH"), ":")
                var words []string = strings.Split(scannedString," ")
                //args := strings.Join(words[1:]," ")
                for i := 0; i < len(pathsToFile); i++ {
                    if fileExists(pathsToFile[i] + "/" + cmd) {
                        attr := syscall.ProcAttr {
                            Dir : os.Getenv("PWD"),
                            Files : []uintptr{0, 1, 2},
                            Env : os.Environ(),
                        }
                        _,ok := syscall.ForkExec(pathsToFile[i] + "/" +
                                cmd, words[0:], &attr)
                        if ok != nil {
                            fmt.Println(ok.Error())
                        }
                    }
                }
            }
                
        }
        var ws syscall.WaitStatus
        // нужно подождать окончания процесса перед выводом '>'
        syscall.Wait4(childPid, &ws, 0, nil)
        fmt.Printf(">")
    }
}

// данная функция проверяет существует ли файл
func fileExists(name string) bool {
    if _, ok := os.Stat(name); ok != nil {
        return os.IsExist(ok)
    } else { return true }
}
