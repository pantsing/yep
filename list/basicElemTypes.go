package list

// Element

type Elem interface {
	// Array function MUST ruturn a slice with capbility == n.
	// The slice capability < n would lead to an out of range panic.
	Array(n int) Elems
}

type StringElem string

func (e StringElem) Array(n int) Elems {
	return StringElems(make([]StringElem, n, n))
}

type IntElem int

func (e IntElem) Array(n int) Elems {
	return IntElems(make([]IntElem, n, n))
}

type Int8Elem int8

func (e Int8Elem) Array(n int) Elems {
	return Int8Elems(make([]Int8Elem, n, n))
}

type Int16Elem int16

func (e Int16Elem) Array(n int) Elems {
	return Int16Elems(make([]Int16Elem, n, n))
}

type Int32Elem int32

func (e Int32Elem) Array(n int) Elems {
	return Int32Elems(make([]Int32Elem, n, n))
}

type Int64Elem int64

func (e Int64Elem) Array(n int) Elems {
	return Int64Elems(make([]Int64Elem, n, n))
}

type UIntElem uint

func (e UIntElem) Array(n int) Elems {
	return UIntElems(make([]UIntElem, n, n))
}

type UInt8Elem uint8

func (e UInt8Elem) Array(n int) Elems {
	return UInt8Elems(make([]UInt8Elem, n, n))
}

type UInt16Elem uint16

func (e UInt16Elem) Array(n int) Elems {
	return UInt16Elems(make([]UInt16Elem, n, n))
}

type UInt32Elem uint32

func (e UInt32Elem) Array(n int) Elems {
	return UInt32Elems(make([]UInt32Elem, n, n))
}

type UInt64Elem uint64

func (e UInt64Elem) Array(n int) Elems {
	return UInt64Elems(make([]UInt64Elem, n, n))
}

type Float32Elem float32

func (e Float32Elem) Array(n int) Elems {
	return Float32Elems(make([]Float32Elem, n, n))
}

type Float64Elem float64

func (e Float64Elem) Array(n int) Elems {
	return Float64Elems(make([]Float64Elem, n, n))
}

type Complex64Elem complex64

func (e Complex64Elem) Array(n int) Elems {
	return Complex64Elems(make([]Complex64Elem, n, n))
}

type Complex128Elem complex128

func (e Complex128Elem) Array(n int) Elems {
	return Complex128Elems(make([]Complex128Elem, n, n))
}

type UintptrElem uintptr

func (e UintptrElem) Array(n int) Elems {
	return UintptrElems(make([]UintptrElem, n, n))
}

type BoolElem bool

func (e BoolElem) Array(n int) Elems {
	return BoolElems(make([]BoolElem, n, n))
}

type ByteElem byte

func (e ByteElem) Array(n int) Elems {
	return ByteElems(make([]ByteElem, n, n))
}

type ByteSliceElem []byte

func (e ByteSliceElem) Array(n int) Elems {
	return ByteSliceElems(make([]ByteSliceElem, n, n))
}

// Elements
type Elems interface {
	Get(i int) Elem
	Set(i int, e Elem)
}

type StringElems []StringElem

func (es StringElems) Get(i int) Elem {
	return es[i]
}
func (es StringElems) Set(i int, e Elem) {
	es[i] = e.(StringElem)
}

type IntElems []IntElem

func (es IntElems) Get(i int) Elem {
	return es[i]
}
func (es IntElems) Set(i int, e Elem) {
	es[i] = e.(IntElem)
}

type Int8Elems []Int8Elem

func (es Int8Elems) Get(i int) Elem {
	return es[i]
}
func (es Int8Elems) Set(i int, e Elem) {
	es[i] = e.(Int8Elem)
}

type Int16Elems []Int16Elem

func (es Int16Elems) Get(i int) Elem {
	return es[i]
}
func (es Int16Elems) Set(i int, e Elem) {
	es[i] = e.(Int16Elem)
}

type Int32Elems []Int32Elem

func (es Int32Elems) Get(i int) Elem {
	return es[i]
}
func (es Int32Elems) Set(i int, e Elem) {
	es[i] = e.(Int32Elem)
}

type Int64Elems []Int64Elem

func (es Int64Elems) Get(i int) Elem {
	return es[i]
}
func (es Int64Elems) Set(i int, e Elem) {
	es[i] = e.(Int64Elem)
}

type UIntElems []UIntElem

func (es UIntElems) Get(i int) Elem {
	return es[i]
}
func (es UIntElems) Set(i int, e Elem) {
	es[i] = e.(UIntElem)
}

type UInt8Elems []UInt8Elem

func (es UInt8Elems) Get(i int) Elem {
	return es[i]
}
func (es UInt8Elems) Set(i int, e Elem) {
	es[i] = e.(UInt8Elem)
}

type UInt16Elems []UInt16Elem

func (es UInt16Elems) Get(i int) Elem {
	return es[i]
}
func (es UInt16Elems) Set(i int, e Elem) {
	es[i] = e.(UInt16Elem)
}

type UInt32Elems []UInt32Elem

func (es UInt32Elems) Get(i int) Elem {
	return es[i]
}
func (es UInt32Elems) Set(i int, e Elem) {
	es[i] = e.(UInt32Elem)
}

type UInt64Elems []UInt64Elem

func (es UInt64Elems) Get(i int) Elem {
	return es[i]
}
func (es UInt64Elems) Set(i int, e Elem) {
	es[i] = e.(UInt64Elem)
}

type Float32Elems []Float32Elem

func (es Float32Elems) Get(i int) Elem {
	return es[i]
}
func (es Float32Elems) Set(i int, e Elem) {
	es[i] = e.(Float32Elem)
}

type Float64Elems []Float64Elem

func (es Float64Elems) Get(i int) Elem {
	return es[i]
}
func (es Float64Elems) Set(i int, e Elem) {
	es[i] = e.(Float64Elem)
}

type Complex64Elems []Complex64Elem

func (es Complex64Elems) Get(i int) Elem {
	return es[i]
}
func (es Complex64Elems) Set(i int, e Elem) {
	es[i] = e.(Complex64Elem)
}

type Complex128Elems []Complex128Elem

func (es Complex128Elems) Get(i int) Elem {
	return es[i]
}
func (es Complex128Elems) Set(i int, e Elem) {
	es[i] = e.(Complex128Elem)
}

type UintptrElems []UintptrElem

func (es UintptrElems) Get(i int) Elem {
	return es[i]
}
func (es UintptrElems) Set(i int, e Elem) {
	es[i] = e.(UintptrElem)
}

type BoolElems []BoolElem

func (es BoolElems) Get(i int) Elem {
	return es[i]
}
func (es BoolElems) Set(i int, e Elem) {
	es[i] = e.(BoolElem)
}

type ByteElems []ByteElem

func (es ByteElems) Get(i int) Elem {
	return es[i]
}
func (es ByteElems) Set(i int, e Elem) {
	es[i] = e.(ByteElem)
}

type ByteSliceElems []ByteSliceElem

func (es ByteSliceElems) Get(i int) Elem {
	return es[i]
}
func (es ByteSliceElems) Set(i int, e Elem) {
	es[i] = e.(ByteSliceElem)
}
