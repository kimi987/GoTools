package main
import (
	"strings"
	"fmt"
	"io/ioutil"
	"os"
)

func main() {
	data, err := ioutil.ReadFile("files.txt")

	if err != nil {
		fmt.Println("[Error] err = ", err)
		return
	}

	fileContents := strings.Split(string(data), "\n")
	result := ""
	for _,v := range fileContents {
		contents := strings.Split(v, "|")

		if len(contents) < 2 {
			continue
		}

		if(strings.Contains(contents[0], ".manifest")) {
			continue
		}

		val := strings.Replace(contents[0], ".ab", ".*", 1) 
		vals := strings.Split(val, "/")
		
		if len(vals) < 2 {
			continue
		}

		if vals[0] == "Avatar" {
			result = fmt.Sprintf("%s%s,%s\n", result, vals[1], "Assets/Resources/High/Avatar")
			result = fmt.Sprintf("%s%s,%s\n", result, vals[1], "Assets/Resources/Low/Avatar")
		} else if vals[0] == "Scenes" {
			result = fmt.Sprintf("%s%s,%s\n", result, vals[1], "Assets/Scenes")
		} else {
			path := ""
			if len(vals) > 2 {
				for k1, v1 := range vals {
					if k1 == len(vals) - 1 {
						continue
					}
					path = fmt.Sprintf("%s/%s", path, v1)
				}
			} else {
				path = fmt.Sprintf("%s/%s", path, vals[0])
			}

			result = fmt.Sprintf("%s%s,%s%s\n", result, vals[len(vals) - 1], "Assets/Resources", path)
		}
	}

	ioutil.WriteFile("AssetBundleIgnore.csv", []byte(result), os.ModePerm)
}

