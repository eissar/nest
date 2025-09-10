package format

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"reflect"
	"slices"
	"strconv"
	"strings"

	"github.com/go-logfmt/logfmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// StructToURLValues converts a struct to a url.Values map based on its json tags.
// This allows you to easily serialize a struct into URL query parameters.
func StructToURLValues(data interface{}) (url.Values, error) {
	// The url.Values type is a map[string][]string, which is what http.Request.URL.Query() returns.
	// It's the standard way to represent query parameters in Go.
	values := url.Values{}

	// Use reflection to inspect the struct.
	// We expect data to be a struct, so we get its value.
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		// If it's a pointer, dereference it to get the struct.
		v = v.Elem()
	}

	// Ensure we are working with a struct.
	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("StructToURLValues only accepts structs; got %T", data)
	}

	// Get the type of the struct to access its fields and tags.
	t := v.Type()

	// Iterate over all the fields of the struct.
	for i := 0; i < v.NumField(); i++ {
		fieldValue := v.Field(i)
		fieldType := t.Field(i)

		// Get the json tag for the current field.
		jsonTag := fieldType.Tag.Get("json")

		// Skip this field if the json tag is "-"
		if jsonTag == "-" {
			continue
		}

		// Parse the tag to get the parameter name and options like "omitempty".
		tagParts := strings.Split(jsonTag, ",")
		paramName := tagParts[0]

		// If the paramName is empty, it means the field is unexported or has no tag.
		// We use the field name as a fallback, but this is often not desired.
		// A better practice is to ensure all exported fields have tags.
		if paramName == "" {
			// Skip unexported fields.
			if !fieldType.IsExported() {
				continue
			}
			paramName = fieldType.Name
		}

		// Check for the "omitempty" option.
		hasOmitempty := false
		if len(tagParts) > 1 {
			for _, part := range tagParts[1:] {
				if part == "omitempty" {
					hasOmitempty = true
					break
				}
			}
		}

		// If "omitempty" is present and the field has its zero value, skip it.
		if hasOmitempty && fieldValue.IsZero() {
			continue
		}

		// Convert the field's value to a string.
		var paramValue string
		switch fieldValue.Kind() {
		case reflect.String:
			paramValue = fieldValue.String()
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			paramValue = strconv.FormatInt(fieldValue.Int(), 10)
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			paramValue = strconv.FormatUint(fieldValue.Uint(), 10)
		case reflect.Float32, reflect.Float64:
			paramValue = strconv.FormatFloat(fieldValue.Float(), 'f', -1, 64)
		case reflect.Bool:
			paramValue = strconv.FormatBool(fieldValue.Bool())
		case reflect.Slice:
			// Handle slices by joining elements with a comma.
			// Another common approach is to add multiple parameters with the same name.
			// e.g., ?tags=go&tags=web
			// We'll demonstrate the comma-separated approach here.
			sliceVal := reflect.ValueOf(fieldValue.Interface())
			var elements []string
			for j := 0; j < sliceVal.Len(); j++ {
				elements = append(elements, fmt.Sprint(sliceVal.Index(j).Interface()))
			}
			paramValue = strings.Join(elements, ",")
		default:
			// For other types, you might need more complex logic.
			// For this example, we'll just skip them.
			continue
		}

		// Add the key-value pair to our url.Values map.
		values.Add(paramName, paramValue)
	}

	return values, nil
}

// TODO: move this somewhere else
// misses flag set
type CobraPFlagParams struct {
	// P is the pointer to the string variable that will hold the flag's value.
	P *string
	// Name is the long name of the flag.
	Name string
	// Shorthand is the single‑character abbreviation.
	Shorthand string
	// Value is the default value for the flag.
	Value string
	// Usage describes the purpose of the flag.
	Usage string
}

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

// TODO: try to unwrap
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

func WriteLogFmt(v any, w io.Writer) {
	if w == nil {
		w = os.Stdout
	}
	enc := logfmt.NewEncoder(w)

	// Handle different types of values
	switch val := v.(type) {
	case string:
		enc.EncodeKeyval("msg", val)
	case error:
		enc.EncodeKeyval("error", val.Error())
	case map[string]any:
		for k, v := range val {
			enc.EncodeKeyval(k, v)
		}
	case map[string]string:
		for k, v := range val {
			enc.EncodeKeyval(k, v)
		}
	default:
		// For any other type, use reflection to get field names and values
		enc.EncodeKeyval("value", fmt.Sprint(val))
	}

	enc.EndRecord()
}

// FormatType defines the allowed formats for the Format function.
// it implements pflag.VarP
type FormatType string

const (
	FormatJSON   FormatType = "json"
	FormatLogFmt FormatType = "logfmt"
	// FormatLog    FormatType = "log"
	// case  "yaml", "yml":
)

// Set implements pflag.Value. It validates and sets the format.
func (f *FormatType) Set(val string) error {
	switch FormatType(val) {
	case FormatJSON: //, FormatLog, FormatLogFmt:
		*f = FormatType(val)
		return nil
	default:
		return fmt.Errorf("unsupported format %q (supported: %s) ", val, FormatJSON) //, FormatLog, FormatLogFmt)
	}
}

// String implements pflag.Value.
func (f *FormatType) String() string {
	if f == nil {
		return ""
	}
	return string(*f)
}

// Type implements pflag.Value (required by cobra/pflag for help output).
func (f *FormatType) Type() string {
	return "format"
}

// allowedFormats := []f.FormatType{f.FormatJSON}
// format := apiCmd.Flag("format")
//
//	if format == nil {
//		panic("FORMAT")
//	}
//
// format.Usage = "output format. One of: "+f.HelpFmt(&allowedFormats)
func OutputFmtFlagUsage(allowedFormats []FormatType) string {
	return "output format. One of: " + fmt.Sprint(allowedFormats)
	// HelpFmt(&allowedFormats)
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
		WriteJSON(o, os.Stdout)
	case FormatLogFmt:
		WriteLogFmt(o, os.Stdout)
	default:
		WriteJSON(o, os.Stdout)
	}
}

// ValidateFlags returns a PersistentPreRunE that validates flag values
// against pre-defined sets of allowed values.  It uses generics so the
// allowed slice can be any comparable type (string, int, custom enums...).
func ValidateFlags[T comparable](rules map[string][]T) func(*cobra.Command, []string) error {
	fmt.Println("DEBUG")
	return func(cmd *cobra.Command, args []string) error {
		for name, allowed := range rules {
			fl := cmd.Flags().Lookup(name)
			if fl == nil || !fl.Changed {
				continue
			}

			// convert the flag value to the target type
			var val T
			switch any(val).(type) {
			case string:
				v, ok := any(fl.Value.String()).(T)
				if !ok {
					return fmt.Errorf("flag --%s: cannot convert %q to %T", name, fl.Value.String(), val)
				}
				val = v
			default:
				if err := json.Unmarshal([]byte(fl.Value.String()), &val); err != nil {
					return fmt.Errorf("flag --%s: cannot unmarshal %q into %T: %w", name, fl.Value.String(), val, err)
				}
			}

			if !slices.Contains(allowed, val) {
				return fmt.Errorf("flag --%s: invalid value %v (must be one of %v)", name, val, allowed)
			}
		}
		return nil
	}
}

// TODO: some kind of way to filter by fields

// fmtFields := func(props string) []string {
// 	// strings.ReplaceAll(props, " ", "") // remove whitespace
// 	fields := strings.Split(props, ",")
//
// 	return fields
// }
// switch format {
// case "json":
// 	var targetProperties []string
// 	if properties != "" {
// 		targetProperties = fmtFields(properties)
// 	} else {
// 		targetProperties = defaultFields
// 	}
//
// 	allFields := structToKeys(&api.ListItem{}) // no struct keys should have any whitespace
// 	// find fields which are in allfields but not in inputFilterFields
// 	exclFields := filterFieldsByReference(targetProperties, allFields)
// 	err = jsonFmtStdOut(cmd, data, exclFields)
// 	if err != nil {
// 		fmt.Printf("jsonFmtStdOut: %v\n", err)
// 	}
// case "logfmt":
// 	logFmtStdOut(data, strings.Split(properties, ","))
// }

/*
BindStructFlags automatically binds struct fields to Cobra command-line flags.

The function uses reflection to iterate over exported fields of opts, creating
a flag for each one. Fields must be exported (start with capital letter) to be
considered for flag binding.

Flag naming follows these rules:
- Use the struct tag `flagname:"custom-name"` to override the default name (WIP?)
- Default flag name converts Go field name to kebab-case ("MyField" → "my-field")

Flag usage is determined by:
- Use the struct tag `flag:"description"` to set help text
- Falls back to generated usage based on field type and name if not provided

Supported field types: int, string (bool, float64, etc. can be added to switch)

Parameters:
- cmd: The Cobra command to add flags to
- opts: Pointer to struct whose fields become flags. Must be pointer to struct.

Returns error if:
- opts is not a pointer to struct
- Unsupported field type is encountered

Example:

	```go
	type Options struct {
			Debug       bool   `flag:"Enable debug logging"`
			MaxItems    int    `flagname:"max-items" flag:"Maximum items to process"`
			OutputFile  string
	}
	var opts Options
	BindStructFlags(cmd, &opts)
	```
*/
func BindStructFlags[T any](cmd *cobra.Command, opts *T) error {
	val := reflect.ValueOf(opts)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("opts must be a pointer to struct")
	}
	val = val.Elem()
	typ := val.Type()

	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)

		if f.PkgPath != "" { // unexported
			continue
		}
		flagName := f.Tag.Get("flagname")
		if flagName == "" {
			flagName = CamelCase(f.Name)
		}
		usageTag := f.Tag.Get("flag") // e.g., flag:"max items to fetch"
		if usageTag == "" {
			usageTag = defaultUsageForField(f) // fallback
		}

		addr := val.Field(i).Addr().Interface()

		// TODO: MORE TYPES
		switch x := addr.(type) {
		case *int:
			cmd.Flags().IntVarP(x, flagName, "", *x, usageTag)
		case *string:
			cmd.Flags().StringVarP(x, flagName, "", *x, usageTag)
		default:
			return fmt.Errorf("unsupported field %s of type %s", f.Name, f.Type)
		}
	}
	return nil
}

func defaultUsageForField(f reflect.StructField) string {
	return fmt.Sprintf("sets the %s parameter", f.Name)
}

// CamelCase converts a string to camelCase, handling various cases.
// It first checks if the string is entirely uppercase and, if so, converts it to lowercase.
// Then, it processes the string to create camelCase by capitalizing letters
// that follow a space, underscore, or hyphen, and lowercasing the very first letter.
func CamelCase(s string) string {
	// 1) Detect if there is at least one ASCII letter, and if every ASCII letter is uppercase.
	hasLetter := false
	allUpper := true
	for i := 0; i < len(s); i++ {
		c := s[i]
		switch {
		case 'a' <= c && c <= 'z':
			hasLetter = true
			allUpper = false
		case 'A' <= c && c <= 'Z':
			hasLetter = true
		}
	}

	// If all ASCII letters are uppercase, lowercase them all and return.
	if hasLetter && allUpper {
		var b strings.Builder
		b.Grow(len(s))
		for i := 0; i < len(s); i++ {
			c := s[i]
			if 'A' <= c && c <= 'Z' {
				b.WriteByte(c + ('a' - 'A'))
			} else {
				b.WriteByte(c)
			}
		}
		return b.String()
	}

	// 2) Otherwise build camelCase:
	//    - First letter → lowercase (if ASCII letter)
	//    - Any letter after ' ', '_' or '-' → uppercase
	//    - All other letters remain as-is.
	var b strings.Builder
	b.Grow(len(s))

	upperNext := false
	first := true

	for i := 0; i < len(s); i++ {
		c := s[i]

		// delimiters trigger next letter uppercase, and are skipped
		if c == ' ' || c == '_' || c == '-' {
			upperNext = true
			continue
		}

		if first {
			// lowercase first letter if needed
			if 'A' <= c && c <= 'Z' {
				b.WriteByte(c + ('a' - 'A'))
			} else {
				b.WriteByte(c)
			}
			first = false
		} else if upperNext {
			// uppercase this letter if it’s ASCII lowercase
			if 'a' <= c && c <= 'z' {
				b.WriteByte(c - ('a' - 'A'))
			} else {
				b.WriteByte(c)
			}
			upperNext = false
		} else {
			b.WriteByte(c)
		}
	}

	return b.String()
}
