package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	n "net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"sync"
)

type ImgData struct {
	MessageType string      `json:"type"`
	Data        [][][]uint8 `json:"data"`
}

func visit(files *[]string) filepath.WalkFunc {
	return func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return nil
		}
		if !info.IsDir() {
			*files = append(*files, path)
		}
		return nil
	}
}

func findLowestNumberForFile(files []string) int {
	numbers := make([]int, len(files))
	for i := 0; i < len(files); i++ {
		fileName := filepath.Base(files[i])
		numbers[i], _ = strconv.Atoi(strings.Split(fileName, ".")[0])
	}
	if len(numbers) == 0 {
		return 0
	}
	sort.Ints(numbers)
	return numbers[len(numbers)-1]
}

func RunServer(port string, outputPath string, startNumber int) {

	var files []string

	if len(outputPath) == 0 {
		outputPath = filepath.Join(os.TempDir(), "Data")
	}

	info, err := os.Stat(outputPath)
	if os.IsNotExist(err) {
		err = os.Mkdir(outputPath, 0755)
		if err != nil { // 0755 is Unix permission code in Windows it's ignored
			log.Fatal(err)
		}

	} else if !info.IsDir() {
		fmt.Println("There entry with same path, but it's not directory. Please choose another name")
		return
	}
	err = filepath.Walk(outputPath, visit(&files))
	if err != nil {
		fmt.Println(err)
		return
	}
	var smallestNumber int
	if startNumber == -1 {
		smallestNumber = findLowestNumberForFile(files)
	} else {
		smallestNumber = startNumber
	}

	mutex := sync.Mutex{}

	handler := func(w n.ResponseWriter, r *n.Request) {
		println("Message received")
		mutex.Lock()
		localNumber := smallestNumber
		smallestNumber++
		mutex.Unlock()
		builder := strings.Builder{}
		bufData := make([]byte, 1024)
		err := error(nil)
		bytesRead := 0
		for err == nil {
			bytesRead, err = r.Body.Read(bufData)
			if bytesRead == 0 {
				break
			}
			builder.Write(bufData[:bytesRead])
		}
		var data ImgData
		err = json.Unmarshal([]byte(builder.String()), &data)
		if err != nil {
			println(err.Error())
			return
		}

		img := image.NewRGBA(image.Rect(0, 0, len(data.Data[0]), len(data.Data)))
		for x := 0; x < len(data.Data[0]); x++ {
			for y := 0; y < len(data.Data); y++ {
				img.Set(x, y, color.RGBA{
					R: data.Data[y][x][0],
					G: data.Data[y][x][1],
					B: data.Data[y][x][2],
					A: 255,
				})
			}
		}

		file, err := os.Create(filepath.Join(outputPath, strconv.Itoa(localNumber)+".png"))
		if err != nil {
			println(err.Error())
			return
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				println(err.Error())
			}
		}(file)
		err = png.Encode(file, img)
		if err != nil {
			println(err.Error())
			return
		}
		println("Done!")
	}
	println("Server started")
	err = n.ListenAndServe(port, n.HandlerFunc(handler))
	if err != nil {
		panic(err)
	}
}

func main() {
	var port string
	var outputPath string
	var startNumber int
	flag.StringVar(&port, "port", ":8080", "Port to be listen by server, default 8080")
	flag.StringVar(&port, "p", ":8080", "Port to be listen by server, default 8080")
	flag.StringVar(&outputPath, "output_path", "", "Path to folder to save data in, default os.TempDir()/Data")
	flag.StringVar(&outputPath, "o", "", "Path to folder to save data in, default os.TempDir()/Data")
	flag.IntVar(&startNumber, "start_number", -1, "Number to start file naming with. If not set, the "+
		"largest number in output dir will be used")
	flag.IntVar(&startNumber, "s", -1, "Number to start file naming with. If not set, the "+
		"largest number in output dir will be used")
	flag.Parse()
	RunServer(port, outputPath, startNumber)
}
