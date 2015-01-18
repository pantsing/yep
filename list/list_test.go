package list

import (
	"bytes"
	"fmt"
	"strconv"
	"testing"
)

func TestList(t *testing.T) {
	l, _ := New(FIFO, 1024)
	fmt.Println(l.Len())

	l.Push(ByteSliceElem([]byte("1234567890")))
	e, err := l.Pop()
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(e.(ByteSliceElem)))
}

func Benchmark_FIFO_Int(b *testing.B) {

	dl, _ := New(FIFO, 1024)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for i := 0; i < 1000000; i++ {
			dl.Push(IntElem(i))
		}

		for i := 0; i < 1000000; i++ {
			j, err := dl.Pop()
			k := int(j.(IntElem))
			if err != nil || i != k {
				b.Error(j, err)
				return
			}
		}

	}
}

func Benchmark_LIFO_Int(b *testing.B) {

	dl, _ := New(LIFO, 1024)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for i := 0; i < 1000000; i++ {
			dl.Push(IntElem(i))
		}

		for i := 999999; i >= 0; i-- {
			j, err := dl.Pop()
			k := int(j.(IntElem))
			if err != nil || i != k {
				b.Error(i, k, err)
				return
			}
		}

	}
}

func Benchmark_RAND_Int(b *testing.B) {

	dl, _ := New(RAND, 1024)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for i := 0; i < 1000000; i++ {
			dl.Push(IntElem(i))
		}

		for i := 0; i < 1000000; i++ {
			j, err := dl.Pop()
			k := int(j.(IntElem))
			if err != nil {
				b.Error(i, k, err)
				return
			}
		}

	}
}

func Benchmark_Str(b *testing.B) {
	S := make([]string, 1000000, 1000000)
	for i := 0; i < 1000000; i++ {
		S[i] = strconv.Itoa(i)
	}
	dl, _ := New(FIFO, 1024)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for i := 0; i < 1000000; i++ {
			dl.Push(StringElem(S[i]))
		}

		for i := 0; i < 1000000; i++ {
			j, err := dl.Pop()
			k := string(j.(StringElem))
			if err != nil || S[i] != k {
				b.Error(j, err)
				return
			}
		}

	}
}

func Benchmark_Byte(b *testing.B) {
	S := make([][]byte, 1000000, 1000000)
	for i := 0; i < 1000000; i++ {
		S[i] = []byte(strconv.Itoa(i))
	}
	dl, _ := New(FIFO, 1024)
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		for i := 0; i < 1000000; i++ {
			dl.Push(ByteSliceElem(S[i]))
		}

		for i := 0; i < 1000000; i++ {
			j, err := dl.Pop()
			k := []byte(j.(ByteSliceElem))
			if err != nil || !bytes.Equal(S[i], k) {
				b.Error(j, err)
				return
			}
		}

	}
}
