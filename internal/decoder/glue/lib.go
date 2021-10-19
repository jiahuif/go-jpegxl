// Package glue generates Go bindings to libjxl.
// Please run "go generate" every time you modify decoder.i

//go:generate swig -c++ -go -cgo -intgosize 64 decoder.i
package glue

/*
#cgo pkg-config: libjxl
*/
import "C"
