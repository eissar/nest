package cmd

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/eissar/nest/api"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// What is this??
func rebuildCmdText(cmd *cobra.Command, args []string) string {
	var commandBuilder strings.Builder

	commandBuilder.WriteString(cmd.CommandPath())

	cmd.Flags().Visit(func(flag *pflag.Flag) {
		commandBuilder.WriteString(fmt.Sprintf(" --%s=%s", flag.Name, flag.Value))
	})
	if len(args) > 0 {
		commandBuilder.WriteString(" ")
		commandBuilder.WriteString(strings.Join(args, " "))
	}

	return commandBuilder.String()
}

// validateIsEagleServerRunning checks if the nest server is running at the specified URL.
func isServerRunning(url string) bool {
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return false //, fmt.Errorf("error creating request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return false //, fmt.Errorf("error making request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return false //, fmt.Errorf("received code other than 200: %v", resp.StatusCode)
	}

	return true //, nil
}

// TODO: add to command line (all below)
func removeItem(baseUrl string, itemIds []string) error {
	err := api.ItemMoveToTrash(baseUrl, itemIds)
	if err != nil {
		return fmt.Errorf("failed to move item to trash: %s", err)
	}

	fmt.Println("Item moved to Trash successfully")
	return nil
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

// filterStructFields takes a struct and two lists of allowed field names.
// It returns a map containing only the fields that are present in BOTH lists.
func filterStructFields(s any, list1, list2 []string) map[string]any {
	// Use maps for efficient O(1) lookups.
	allowList1 := make(map[string]struct{})
	for _, fieldName := range list1 {
		allowList1[fieldName] = struct{}{}
	}

	allowList2 := make(map[string]struct{})
	for _, fieldName := range list2 {
		allowList2[fieldName] = struct{}{}
	}

	// The map to store the filtered result.
	result := make(map[string]interface{})

	// Use reflection to inspect the struct.
	val := reflect.ValueOf(s)
	// If it's a pointer, get the element it points to.
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldName := field.Name

		// Check if the field name exists in both allow-lists.
		_, inList1 := allowList1[fieldName]
		_, inList2 := allowList2[fieldName]

		if inList1 && inList2 {
			// If it matches, add the field name and value to our result map.
			result[fieldName] = val.Field(i).Interface()
		}
	}

	return result
}

// extractFields takes a struct and a list of desired field names.
// It returns a map containing only the fields from the struct that are present in the list.
func extractFields(s any, fieldsToKeep []string) map[string]any {
	allowList := make(map[string]struct{}, len(fieldsToKeep))
	for _, fieldName := range fieldsToKeep {
		allowList[fieldName] = struct{}{}
	}

	// The map to store the filtered result.
	result := make(map[string]any)

	val := reflect.ValueOf(s)
	val = reflect.Indirect(val)

	// Ensure we are working with a struct.
	if val.Kind() != reflect.Struct {
		return result // Or return an error
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := typ.Field(i)
		fieldName := field.Name

		// Check if the field name exists in the allow-list.
		if _, allowed := allowList[fieldName]; allowed {
			// If it's allowed, add the field name and value to our result map.
			result[fieldName] = val.Field(i).Interface()
		}
	}

	return result
}

//	func extractFields1(a any, fields []string) (filteredFields []string) {
//		referenceKeys := structToKeys(&a)
//		for _, field := range fields {
//			for i, _ := range referenceKeys {
//				if referenceKeys[i] == field {
//					filteredFields = append(filteredFields, field)
//				}
//			}
//
//		}
//		return filteredFields
//	}
func filterFieldsByReference(referenceFields []string, fields []string) []string {
	referenceSet := make(map[string]struct{})
	for _, field := range referenceFields {
		referenceSet[field] = struct{}{}
	}

	var output []string
	for _, field := range fields {
		if _, found := referenceSet[field]; !found {
			output = append(output, field)
		}
	}

	return output
}
