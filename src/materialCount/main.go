package main
import (
	"fmt"
	"os"
	"path/filepath"
	"io/ioutil"
	"strings"
)

var RepeatMap map[string]int

func main() {
	RepeatMap = make(map[string]int)

	err := filepath.Walk("Material", func(path string, info os.FileInfo, err error) error {
		if info != nil {
			name := info.Name()

			name = strings.Split(name, "#")[0]

			RepeatMap[name] = RepeatMap[name] + 1
		}
		return nil
	})

	if err != nil {
		fmt.Println("err = ", err)
	}
	result := ""
	for k, v := range RepeatMap {
		if v > 2 {
			name := strings.ReplaceAll(k, " ","")
			result = fmt.Sprintf("%s\nfileList.Add(\"Assets/Arts/Particles/Texture/Materials/%s.mat\");",result, name)
			result = fmt.Sprintf("%s\nfileList.Add(\"Assets/Arts/Particles/Materials/%s.mat\");",result, name)
		}
	}
	ioutil.WriteFile("materials.txt", []byte(result), os.ModePerm)
}

