package stats

import (
	"bufio"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
)

type Records struct {
	Stats  map[string]*LevelRecord `json:"stats"`
	sorter map[string]int
}

func (r *Records) Level(difficulty string, lvl int) (*LevelRecord, bool) {
	v, ok := r.Stats[strconv.Itoa(lvl)+difficulty]
	return v, ok
}

func (r *Records) Length() int {
	return len(r.Stats)
}

func (r *Records) All() []*LevelRecord {
	recs := make([]*LevelRecord, len(r.Stats), len(r.Stats))
	i := 0
	for _, rec := range r.Stats {
		recs[i] = rec
		i++
	}
	ret := &recordList{recs: recs}
	sort.Sort(ret)
	return ret.recs
}

type recordList struct {
	recs []*LevelRecord
}

// Len is part of sort.Interface.
func (r *recordList) Len() int {
	return len(r.recs)
}

// Swap is part of sort.Interface.
func (r *recordList) Swap(i, j int) {
	r.recs[i], r.recs[j] = r.recs[j], r.recs[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (r *recordList) Less(i, j int) bool {
	if r.recs[i].Difficulty == r.recs[j].Difficulty {
		return r.recs[i].Lvl < r.recs[j].Lvl
	}
	return r.recs[i].Difficulty < r.recs[j].Difficulty
}

type LevelRecord struct {
	Difficulty string
	Lvl        int
	Steps      int `json:"steps,omitempty"`
	Seconds    int `json:"seconds,omitempty"`
}

func New(sortMap map[string]int) *Records {
	return &Records{Stats: make(map[string]*LevelRecord, 30), sorter: sortMap}
}

func Load(file string, sorter map[string]int) (*Records, error) {
	recs := &Records{sorter: sorter}
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

func (r *Records) String(difficulty string, lvl int) string {
	if rec, ok := r.Stats[strconv.Itoa(lvl)+difficulty]; !ok {
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
func (r *Records) Log(difficulty string, lvl, steps, seconds int) bool {
	key := strconv.Itoa(lvl) + difficulty
	rec, ok := r.Stats[key]
	defer func() {
		r.Stats[key] = rec
	}()
	if !ok {
		rec = &LevelRecord{Lvl: lvl, Difficulty: difficulty, Steps: steps, Seconds: seconds}
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
