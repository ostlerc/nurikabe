package stats

import (
	"bufio"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"strconv"
)

type Records struct {
	Stats map[string]*LevelRecord `json:"stats"`
}

type LevelRecord struct {
	Steps   int `json:"steps,omitempty"`
	Seconds int `json:"seconds,omitempty"`
}

func New() *Records {
	return &Records{Stats: make(map[string]*LevelRecord, 30)}
}

func Load(file string) (*Records, error) {
	recs := &Records{}
	reader, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	r := bufio.NewReader(reader)
	dat, err := r.ReadBytes('\n')
	if err != nil && err != io.EOF {
		return nil, err
	}
	err = json.Unmarshal(dat, recs)
	if err != nil {
		return nil, err
	}
	return recs, nil
}

func (r *Records) String(file string) string {
	if rec, ok := r.Stats[file]; !ok {
		return ""
	} else {
		return "record: " + strconv.Itoa(rec.Steps) + " steps, " + strconv.Itoa(rec.Seconds) + " seconds  "
	}
}

func (r *Records) Save(file string) error {
	dat, err := json.Marshal(r)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(file, dat, 0600)
}

//Returns true if new stats were better than previous
func (r *Records) Log(name string, steps, seconds int) bool {
	rec, ok := r.Stats[name]
	defer func() {
		r.Stats[name] = rec
	}()
	if !ok {
		rec = &LevelRecord{Steps: steps, Seconds: seconds}
		return true
	}
	if steps < rec.Steps {
		rec.Steps = steps
		rec.Seconds = seconds
		return true
	} else if steps == rec.Steps {
		if seconds < rec.Seconds {
			rec.Seconds = seconds
			return true
		}
	}
	return false
}
