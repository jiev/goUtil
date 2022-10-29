package fileUtil


import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
)

func Readln(r *bufio.Reader) (string, error) {
	var (
		isPrefix bool  = true
		err      error = nil
		line, ln []byte
	)
	for isPrefix && err == nil {
		line, isPrefix, err = r.ReadLine()
		ln = append(ln, line...)
	}
	return string(ln), err
}

func ReadFileLines(fileName string) (ret []string, err error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("can not find the file %v", file))
	}
	defer file.Close()
	br := bufio.NewReader(file)

	for {
		line, err := Readln(br)
		if err != nil {
			if err == io.EOF {
				err = nil
			}
			break
		}
		ret = append(ret, string(line))
	}
	return ret, err
}





