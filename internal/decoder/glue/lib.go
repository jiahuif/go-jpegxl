//go:generate swig -c++ -go -cgo -intgosize 64 decoder.i
package glue

/*
#cgo pkg-config: libjxl
 */
import "C"