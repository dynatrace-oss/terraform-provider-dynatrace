package download

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"
)

func PrintPreview(dlConfig DownloadConfig) error {
	var count int
	fmt.Printf("Processing preview...\n\n")
	w := tabwriter.NewWriter(os.Stdout, 0, 1, 1, ' ', 0)
	fmt.Fprintf(w, "%v\t%v\t\n", "RESOURCE NAME", "COUNT")
	for resName, resStruct := range ResourceInfoMap {
		trimRes := strings.TrimPrefix(resName, "dynatrace_")
		if len(dlConfig.ResourceNames) != 0 {
			if (!dlConfig.Exclude && !dlConfig.MatchResource(resName)) || (dlConfig.Exclude && dlConfig.MatchResource(resName)) {
				continue
			}
		}
		if ResourceInfoMap[resName].NoListClient != nil {
			fmt.Fprintf(w, "%v\t%v\t\n", trimRes, 1)
		} else {
			clients := resStruct.RESTClient(
				dlConfig.EnvironmentURL,
				dlConfig.APIToken,
			)
			count = 0
			for _, client := range clients {
				ids, err := client.LIST()
				if err != nil {
					return err
				}
				count += len(ids)
			}
			fmt.Fprintf(w, "%v\t%v\t\n", trimRes, count)
		}
	}
	fmt.Println("-------- Export Preview --------")
	w.Flush()
	return nil
}
