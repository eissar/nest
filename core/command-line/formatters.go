package commandline

import (
	"encoding/json"
	"os"

	"github.com/eissar/nest/api"
	"github.com/go-logfmt/logfmt"
	"github.com/spf13/cobra"
)

func logFmtStdOut(data []*api.ListItem, props []string) {
	outp := logfmt.NewEncoder(os.Stdout)

	outp.EncodeKeyvals("test", nil)
	outp.EndRecord()
	// for _, item := range data {
	//	// outp.EncodeKeyval("url", item.URL)
	//	// outp.EncodeKeyval("folderIds", strings.Join(item.Folders, ", "))
	//	outp.EndRecord()
	// }

}

// exclude properties to exclude
func jsonFmtStdOut(cmd *cobra.Command, data []*api.ListItem, exclude []string) error {
	// json.NewEncoder(os.Stdout)

	b, err := json.Marshal(data)
	if err != nil {
		return err
	}

	var m []map[string]interface{}
	if err := json.Unmarshal(b, &m); err != nil {
		return err
	}

	for i, _ := range m {
		for _, key := range exclude {
			delete(m[i], key)
		}
	}

	st := json.NewEncoder(os.Stdout)

	st.Encode(m)

	return nil
}
