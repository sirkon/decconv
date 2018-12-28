package decdec

//go:generate genny -in=decoding32.go -out=decoding64.go gen "Decode32=Decode64 uint32=uint64 pow32=pow64"
