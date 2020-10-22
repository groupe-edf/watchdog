package util

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"
	"unicode/utf8"
)

// ElapsedTime log task elapsed time
func ElapsedTime(ctx context.Context, action string) {
	elapsed := time.Since(GetStartTime(ctx))
	fmt.Printf("%s took %s", action, elapsed)
	fmt.Println()
}

// PrintMessage print message in the stdout
func PrintMessage(message string) {
	fmt.Println()
	fmt.Println("###" + strings.Repeat("#", utf8.RuneCountInString(message)) + "###")
	fmt.Println("#  " + strings.Repeat(" ", utf8.RuneCountInString(message)) + "  #")
	fmt.Println("#  " + message + "  #")
	fmt.Println("#  " + strings.Repeat(" ", utf8.RuneCountInString(message)) + "  #")
	fmt.Println("###" + strings.Repeat("#", utf8.RuneCountInString(message)) + "###")
	fmt.Println()
}

// ItemExists check if an element exists in array
func ItemExists(arrayType interface{}, item interface{}) bool {
	arr := reflect.ValueOf(arrayType)
	if arr.Kind() != reflect.Array {
		panic("Invalid data-type")
	}
	for i := 0; i < arr.Len(); i++ {
		if arr.Index(i).Interface() == item {
			return true
		}
	}
	return false
}
