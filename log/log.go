package log

import "fmt"

func Info(msg string,opt ...interface{}){
	fmt.Printf(msg,opt)
	fmt.Println()
}

func Error(msg string,opt ...interface{}){
	fmt.Printf(msg,opt)
	fmt.Println()
}

func Debug(msg string,opt ...interface{}){
	fmt.Printf(msg,opt)
	fmt.Println()
}