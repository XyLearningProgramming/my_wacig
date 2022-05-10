package my_engine

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	exampleCodeDir string
	fileCodes      []string
)

func TestMain(m *testing.M) {
	// prepare examples
	// _, filename, _, _ := runtime.Caller(0)
	wd, err := os.Getwd()
	if err != nil {
		panic(fmt.Errorf("getting working dir failed: %w", err))
	}
	exampleCodeDir = path.Join(path.Dir(wd), "examples")
	files, err := ioutil.ReadDir(exampleCodeDir)
	if err != nil {
		panic(fmt.Errorf("error opening from %s: %w", exampleCodeDir, err))
	}

	for _, f := range files {
		fn := f.Name()
		if strings.HasSuffix(fn, ".monkey") {
			code, err := ioutil.ReadFile(path.Join(exampleCodeDir, fn))
			if err != nil {
				panic(fmt.Errorf("error reading from %s: %w", fn, err))
			}
			fileCodes = append(fileCodes, string(code))
		}
	}

	os.Exit(m.Run())
}

func TestEngineWithExamples(t *testing.T) {
	enginesToTest := []Engine{
		NewEvalEngine(),
		NewVMEngine(),
	}
	for _, eg := range enginesToTest {
		for _, code := range fileCodes {
			res, err := eg.Evaluate(code)
			assert.NoError(
				t,
				err,
				"error evaluating: engine: %v: code: %s: res: %s",
				eg, code, res,
			)
		}
	}
}
