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

	"github.com/go-logfmt/logfmt"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

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

// bindStructToFlags populates every flag in cmd with a field from opts.
func BindStructToFlags(cmd *cobra.Command, opts any) error {

	val := reflect.ValueOf(opts)
	if val.Kind() != reflect.Ptr || val.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("opts must be a pointer to struct")
	}
	val = val.Elem()
	typ := val.Type()

	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)

		// #region try

		// #endregion end try

		if f.PkgPath != "" { // unexported
			continue
		}
		// "camelCase" field name works as long as the struct follows that pattern.
		flagName := toKebabCase(f.Name) // turn "Limit" into "limit"
		usageTag := f.Tag.Get("flag")   // e.g. flag:"max items to fetch"
		if usageTag == "" {
			usageTag = defaultUsageForField(f) // fallback
		}

		addr := val.Field(i).Addr().Interface()

		switch x := addr.(type) {
		case *int:
			cmd.Flags().IntVarP(x, flagName, "", *x, usageTag)
		case *string:
			cmd.Flags().StringVarP(x, flagName, "", *x, usageTag)
		// add other simple kinds here (bool, float64, ...) if needed
		default:
			return fmt.Errorf("unsupported field %s of type %s", f.Name, f.Type)
		}
	}
	return nil
}

func defaultUsageForField(f reflect.StructField) string {
	return fmt.Sprintf("sets the %s parameter", f.Name)
}

// --- fairly trivial helpers --------------------------------------------------
func toKebabCase(s string) string {
	// naive camel -> kebab, e.g. "OrderBy" -> "order-by"
	// actual implementation could use the flect package or github.com/iancoleman/strcase
	return s // simplified stub
}
