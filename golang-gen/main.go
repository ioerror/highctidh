package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"text/template"

	"github.com/joncalhoun/pipe"
)

type data struct {
	Type string
	Name string
	Bits int
}

func main() {
	var d data
	flag.IntVar(&d.Bits, "bits", 1024, "number of bits")
	flag.StringVar(&d.Type, "type", "", "The subtype used for the queue being generated")
	flag.StringVar(&d.Name, "name", "", "The name used for the queue being generated. This should start with a capital letter so that it is exported.")
	flag.Parse()
	err := os.WriteFile(fmt.Sprintf("%s.go", strings.ToLower(d.Name)), doTemplate(BindingCode, d), 0644)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(fmt.Sprintf("%s_test.go", strings.ToLower(d.Name)), doTemplate(BindingTests, d), 0644)
	if err != nil {
		panic(err)
	}
	err = os.WriteFile(fmt.Sprintf("%s_bench_test.go", strings.ToLower(d.Name)), doTemplate(BenchmarkTests, d), 0644)
	if err != nil {
		panic(err)
	}
}

func doTemplate(templateString string, d data) []byte {
	t := template.Must(template.New("").Parse(templateString))
	rc, wc, errCh := pipe.Commands(
		exec.Command("gofmt"),
		exec.Command("goimports"),
	)
	go func() {
		select {
		case err, ok := <-errCh:
			if ok && err != nil {
				panic(err)
			}
		}
	}()
	t.Execute(wc, d)
	wc.Close()
	var b bytes.Buffer
	io.Copy(&b, rc)
	return b.Bytes()
}
