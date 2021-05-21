package benchmarks

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"testing"
	"time"

	"github.com/guregu/null"
)

type Message struct {
	Message string `json:"message"`
}

type World struct {
	ID      int32  `json:"id"`
	Message string `json:"message"`
}

// func (w World) MarshalJSON() ([]byte, error) {
// 	// out := []byte(nil)
// 	// out = append(out, []byte(`{"id":"`)...)
// 	// strconv.AppendInt(out, int64(w.ID), 10)
// 	// out = append(out, []byte(`","randomnumber":"`)...)
// 	// out = append(out, []byte(w.Message)...)
// 	// out = append(out, []byte(`"}`)...)

// 	out := []byte(`{"id":1,"message":"hello,world"}`)
// 	return out, nil
// }

/*
json 编码推荐使用NewEncoder.Enocde
*/
func BenchmarkJSONMarshaler(b *testing.B) {
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		out := &bytes.Buffer{}
		enc := json.NewEncoder(out)
		m := World{
			ID:      1,
			Message: "Hello World",
		}
		for pb.Next() {
			out.Reset()
			if err := enc.Encode(&m); err != nil {
				b.Fatal("Encode:", err)
			}
		}
		b.Logf("%s", out.String())
	})

}
func BenchmarkJSONMarshal(b *testing.B) {
	b.ReportAllocs()

	b.RunParallel(func(pb *testing.PB) {
		var out []byte
		var err error
		m := World{
			ID:      1,
			Message: "Hello World",
		}
		for pb.Next() {
			if out, err = json.Marshal(m); err != nil {
				b.Fatal("Marshal:", err)
			}
		}
		b.Logf("%s", out)
	})

}
func BenchmarkEncodeMarshaler(b *testing.B) {
	b.ReportAllocs()

	m := struct {
		A null.Int
		B time.Time
		C time.Time
		D null.String
	}{
		A: null.IntFrom(42),
		B: time.Now(),
		C: time.Now().Add(-time.Hour),
		D: null.StringFrom(`hello`),
	}

	b.RunParallel(func(pb *testing.PB) {
		enc := json.NewEncoder(ioutil.Discard)

		for pb.Next() {
			if err := enc.Encode(&m); err != nil {
				b.Fatal("Encode:", err)
			}
		}
	})
}
