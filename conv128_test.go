package decconv

import (
	"testing"
)

func TestDecode128(t *testing.T) {
	type args struct {
		precision int
		scale     int
		input     string
	}
	tests := []struct {
		name    string
		args    args
		wantLo  uint64
		wantHi  uint64
		conv    string
		wantErr bool
	}{
		{
			name: "only-integral-positive",
			args: args{
				precision: 9,
				scale:     0,
				input:     "12",
			},
			wantLo:  12,
			conv:    "12",
			wantErr: false,
		},
		{
			name: "only-integral-negative",
			args: args{
				precision: 9,
				scale:     0,
				input:     "-123",
			},
			wantLo:  ^uint64(123) + 1,
			wantHi:  ^uint64(0),
			conv:    "-123",
			wantErr: false,
		},
		{
			name: "generic-empty-fraction",
			args: args{
				precision: 9,
				scale:     1,
				input:     "12.0",
			},
			wantLo:  120,
			conv:    "12",
			wantErr: false,
		},
		{
			name: "generic-empty-fraction",
			args: args{
				precision: 9,
				scale:     2,
				input:     "12.02",
			},
			wantLo:  1202,
			conv:    "12.02",
			wantErr: false,
		},
		{
			name: "check-input-positive",
			args: args{
				precision: 9,
				scale:     5,
				input:     "3015.07654",
			},
			wantLo:  301507654,
			conv:    "3015.07654",
			wantErr: false,
		},
		{
			name: "check-input-negative",
			args: args{
				precision: 9,
				scale:     5,
				input:     "-3015.07654",
			},
			wantLo:  ^uint64(301507654) + 1,
			wantHi:  ^uint64(0),
			conv:    "-3015.07654",
			wantErr: false,
		},
		{
			name: "real-128-bit-stuff",
			args: args{
				precision: 38,
				scale:     2,
				input:     "295147905179352825855.25",
			},
			wantLo:  18446744073709551541,
			wantHi:  1599,
			conv:    "295147905179352825855.25",
			wantErr: false,
		},
		{
			name: "passing-leading-zeroes",
			args: args{
				precision: 9,
				scale:     2,
				input:     "0000123.25",
			},
			wantLo:  12325,
			conv:    "123.25",
			wantErr: false,
		},
		{
			name: "passing-trailing-zeroes",
			args: args{
				precision: 9,
				scale:     3,
				input:     "123.02500000000000",
			},
			wantLo:  123025,
			conv:    "123.025",
			wantErr: false,
		},
		{
			name: "passing-both-sides-zeroes",
			args: args{
				precision: 9,
				scale:     4,
				input:     "-0000000123.12300000000",
			},
			wantLo:  ^uint64(1231230) + 1,
			wantHi:  ^uint64(0),
			conv:    "-123.123",
			wantErr: false,
		},
		{
			name: "error-empty-input",
			args: args{
				precision: 9,
				scale:     4,
				input:     "",
			},
			wantLo:  0,
			wantErr: true,
		},
		{
			name: "error-invalid-input-integral-part",
			args: args{
				precision: 9,
				scale:     4,
				input:     "1A2.3",
			},
			wantLo:  0,
			wantErr: true,
		},
		{
			name: "error-invalid-input-fraction-part",
			args: args{
				precision: 9,
				scale:     4,
				input:     "12.b3",
			},
			wantLo:  0,
			wantErr: true,
		},
		{
			name: "error-overflow-integral",
			args: args{
				precision: 9,
				scale:     6,
				input:     "1234.12",
			},
			wantLo:  0,
			wantErr: true,
		},
		{
			name: "error-overflow-fraction",
			args: args{
				precision: 9,
				scale:     3,
				input:     "1234.1234",
			},
			wantLo:  0,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotLo, gotHi, err := Decode128(tt.args.precision, tt.args.scale, []byte(tt.args.input))
			if (err != nil) != tt.wantErr {
				t.Errorf("Decode32() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotLo != tt.wantLo || gotHi != tt.wantHi {
				t.Errorf("Decode32() = (%v, %v), want (%v, %v)", gotLo, gotHi, tt.wantLo, tt.wantHi)
			}
			if conv := Encode128(tt.args.scale, gotLo, gotHi); conv != tt.conv && err == nil {
				t.Errorf("Encode32() = %v, want %v", conv, tt.conv)
			}
		})
	}
}
