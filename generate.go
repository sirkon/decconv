package decconv

//go:generate genny -in=conv32.go -out=conv64.go gen "Decode32=Decode64 int32=int64 pow32=pow64 Encode32=Encode64"
