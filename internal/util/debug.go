package util

import (
	"context"
	"fmt"
	"os"
	"reflect"
	"strings"
	"text/template"
	"time"
	"unicode/utf8"

	"github.com/groupe-edf/watchdog/internal/config"
	"github.com/groupe-edf/watchdog/internal/version"
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

// PrintBanner Print watchdog banner
func PrintBanner(options *config.Options) error {
	t, err := template.New("watchdog").Parse(config.Banner)
	if err != nil {
		return err
	}
	data := map[string]interface{}{
		"Options":   options,
		"BuildInfo": version.GetBuildInfo(),
	}
	return t.Execute(os.Stdout, data)
}
