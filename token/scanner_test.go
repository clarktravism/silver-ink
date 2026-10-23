package token

import (
	"bytes"
	"encoding/json"
	"os"
	"silver-ink/data"
	"strings"
	"testing"
)

type MockLabeler struct {
	Labels map[data.Name][]int `json:"ExpectedLabels"`
}

func (ml *MockLabeler) AddLabel(name data.Name, pos int) {
	if ml.Labels == nil {
		ml.Labels = make(map[data.Name][]int, 1)
	}
	ml.Labels[name] = append(ml.Labels[name], pos)
}

func (ml *MockLabeler) Equal(ml2 *MockLabeler) bool {
	for k, v := range ml.Labels {
		v2 := ml2.Labels[k]
		if len(v) != len(v2) {
			return false
		}
		for i := range v {
			found := false
			for j := range v2 {
				if v[i] == v2[j] {
					found = true
					break
				}
			}
			if !found {
				return false
			}
		}
	}
	return true
}

type MockNameRepo struct {
	i int
}

func (mn *MockNameRepo) GetName(name data.Name) int {
	mn.i++
	return mn.i
}

func (mn *MockNameRepo) IsElement(data.Name) bool { return false }
func (mn *MockNameRepo) IsPrefab(data.Name) bool  { return false }
func (mn *MockNameRepo) IsBoard(data.Name) bool   { return false }

type ByteCodeScannerTest struct {
	Name          string `json:"Name"`
	Text          string `json:"Text"`
	ExpectedBytes []byte `json:"ExpectedBytes"`
	ExpectedError string `json:"ExpectedError"`
}

func TestByteCodeLine(t *testing.T) {
	line := "#IF ALIGNED PLAYER #TRY SEEK PLAYER #THUD"
	
	labler := MockLabeler{}
	namer := MockNameRepo{}
	scanner := NewByteCodeScanner(nil, &labler, &namer)
	out, _ := scanner.byteCodeLine([]byte(line))

	t.Log(out)
}

func TestByteCodeScanner(t *testing.T) {

	// load tests from json file
	b, err := os.ReadFile("./testdata/bytecodescanner.json")
	if err != nil {
		t.Error(err)
		return
	}
	tests := []ByteCodeScannerTest{}
	if err := json.Unmarshal(b, &tests); err != nil {
		t.Error(err)
		return
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			r := strings.NewReader(test.Text)
			labler := MockLabeler{}
			namer := MockNameRepo{}
			p := 0
			scanner := NewByteCodeScanner(r, &labler, &namer)
			for scanner.Scan() {
				if err := scanner.Err(); err != nil {
					if test.ExpectedError == "" {
						t.Fail()
						t.Logf("unexpected error: %v", err)
					} else if test.ExpectedError != err.Error() {
						t.Fail()
						t.Logf("unexpected error: %v", err)
					}
				} else {
					b := scanner.Bytes()
					if len(test.ExpectedBytes) < p+len(b) {
						t.Fail()
						t.Log("expected byte slice not matching")
						return
					}
					if !bytes.Equal(b, test.ExpectedBytes[p:p+len(b)]) {
						t.Fail()
						t.Log("expected byte slice not matching")
						return
					}
					p += len(b)
				}
			}
		})
	}
}
