package format

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"slices"
	"strings"

	"github.com/spf13/pflag"
)

// don't use in big loop, uses reflection
func structToKeys[T any](a *T) []string {
	// Get the reflect.Value of the struct
	v := reflect.ValueOf(*a)

	// Get the reflect.Type of the struct
	t := v.Type()

	var arr []string
	// Iterate over the fields
	for i := 0; i < v.NumField(); i++ {
		// Get the field's Type (contains name, type, tags)
		fieldType := t.Field(i)
		// Get the field's Value
		// fieldValue := v.Field(i)

		jsonName := strings.Split(fieldType.Tag.Get("json"), ",")[0]
		if jsonName != "" {
			arr = append(arr, jsonName)
		} else {
			arr = append(arr, fieldType.Name)
		}

		// fmt.Printf("Field Name: %s, Field Type: %s, Field Value: %v\n",
		//	fieldType.Name,
		//	fieldType.Type,
		//	fieldValue.Interface(), // Use .Interface() to get the actual value
		// )
	}

	return arr
}

func HelpFmt[T any](a *T) string {
	val := reflect.Indirect(reflect.ValueOf(a))

	switch val.Kind() {
	case reflect.Slice:
		return strings.Join(val.Interface().([]string), ", ")

	case reflect.Struct:
		// Assuming structToKeys correctly inspects the struct and returns its key names.
		return strings.Join(structToKeys(a), ", ")

	default:
		// Provide a fallback for any other types.
		return fmt.Sprint(val.Interface())
	}
}

func WriteJSON(v any, w io.Writer) {
	if w == nil {
		w = os.Stdout
	}
	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(v); err != nil {
		log.Fatalf("WriteJSON: %v", err)
	}
}

// FormatType defines the allowed formats for the Format function.

type FormatType string

const (
	FormatJSON   FormatType = "json"
	FormatLog    FormatType = "log"
	FormatLogFmt FormatType = "logfmt"
)

// allowedFormats := []f.FormatType{f.FormatJSON}
// format := apiCmd.Flag("format")
//
//	if format == nil {
//		panic("FORMAT")
//	}
//
// format.Usage = "output format. One of: "+f.HelpFmt(&allowedFormats)
func OutputFmtFlagUsage(allowedFormats []FormatType) string {
	return "output format. One of: " + HelpFmt(&allowedFormats)
}

// o outputFormat
// updates o.Usage and validates..
func UpdateUsageAndAssertContains(o *pflag.Flag, allowedFormats []FormatType) error {
	o.Usage = OutputFmtFlagUsage(allowedFormats)
	val := o.Value.String()
	if slices.Contains(allowedFormats, FormatType(val)) {
		return nil
	}
	return fmt.Errorf("invalid value:%s\n    %s.", val, o.Usage)
}

// Format formats the output according to the given format.
func Format(f FormatType, o ...any) {
	switch f {
	case FormatJSON:
		{
			WriteJSON(o, os.Stdout)
		}
	}
}
