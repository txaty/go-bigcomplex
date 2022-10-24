package complex

import (
	"math/big"
	"reflect"
	"testing"
)

func TestHurwitzInt_Prod(t *testing.T) {
	type fields struct {
		dblR *big.Int
		dblI *big.Int
		dblJ *big.Int
		dblK *big.Int
	}
	type args struct {
		a *HurwitzInt
		b *HurwitzInt
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *HurwitzInt
	}{
		{
			name: "test_(1+i+j+k) * (1+i+j+k)",
			fields: fields{
				dblR: nil,
				dblI: nil,
				dblJ: nil,
				dblK: nil,
			},
			args: args{
				a: NewHurwitzInt(big.NewInt(1), big.NewInt(1), big.NewInt(1), big.NewInt(1), false),
				b: NewHurwitzInt(big.NewInt(1), big.NewInt(1), big.NewInt(1), big.NewInt(1), false),
			},
			want: NewHurwitzInt(big.NewInt(-2), big.NewInt(2), big.NewInt(2), big.NewInt(2), false),
		},
		{
			name: "test_2i * 1",
			fields: fields{
				dblR: nil,
				dblI: nil,
				dblJ: nil,
				dblK: nil,
			},
			args: args{
				a: NewHurwitzInt(big.NewInt(0), big.NewInt(2), big.NewInt(0), big.NewInt(0), false),
				b: NewHurwitzInt(big.NewInt(1), big.NewInt(0), big.NewInt(0), big.NewInt(0), false),
			},
			want: NewHurwitzInt(big.NewInt(0), big.NewInt(2), big.NewInt(0), big.NewInt(0), false),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HurwitzInt{
				dblR: tt.fields.dblR,
				dblI: tt.fields.dblI,
				dblJ: tt.fields.dblJ,
				dblK: tt.fields.dblK,
			}
			if got := h.Prod(tt.args.a, tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Prod() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHurwitzInt_String(t *testing.T) {
	type fields struct {
		dblR *big.Int
		dblI *big.Int
		dblJ *big.Int
		dblK *big.Int
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "test_1+i+j+k",
			fields: fields{
				dblR: big.NewInt(2),
				dblI: big.NewInt(2),
				dblJ: big.NewInt(2),
				dblK: big.NewInt(2),
			},
			want: "1+i+j+k",
		},
		{
			name: "test_0",
			fields: fields{
				dblR: big.NewInt(0),
				dblI: big.NewInt(0),
				dblJ: big.NewInt(0),
				dblK: big.NewInt(0),
			},
			want: "0",
		},
		{
			name: "test_1.5+1.5i+1.5j+1.5k",
			fields: fields{
				dblR: big.NewInt(3),
				dblI: big.NewInt(3),
				dblJ: big.NewInt(3),
				dblK: big.NewInt(3),
			},
			want: "1.5+1.5i+1.5j+1.5k",
		},
		{
			name: "test_0.5+0.5i+0.5j+0.5k",
			fields: fields{
				dblR: big.NewInt(1),
				dblI: big.NewInt(1),
				dblJ: big.NewInt(1),
				dblK: big.NewInt(1),
			},
			want: "0.5+0.5i+0.5j+0.5k",
		},
		{
			name: "test_-0.5i-0.5j+0.5k",
			fields: fields{
				dblR: big.NewInt(0),
				dblI: big.NewInt(-1),
				dblJ: big.NewInt(-1),
				dblK: big.NewInt(1),
			},
			want: "-0.5i-0.5j+0.5k",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HurwitzInt{
				dblR: tt.fields.dblR,
				dblI: tt.fields.dblI,
				dblJ: tt.fields.dblJ,
				dblK: tt.fields.dblK,
			}
			if got := h.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHurwitzInt_ValInt(t *testing.T) {
	type fields struct {
		dblR *big.Int
		dblI *big.Int
		dblJ *big.Int
		dblK *big.Int
	}
	tests := []struct {
		name   string
		fields fields
		wantR  *big.Int
		wantI  *big.Int
		wantJ  *big.Int
		wantK  *big.Int
	}{
		{
			name: "test_1+i+j+k",
			fields: fields{
				dblR: big.NewInt(2),
				dblI: big.NewInt(2),
				dblJ: big.NewInt(2),
				dblK: big.NewInt(2),
			},
			wantR: big.NewInt(1),
			wantI: big.NewInt(1),
			wantJ: big.NewInt(1),
			wantK: big.NewInt(1),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &HurwitzInt{
				dblR: tt.fields.dblR,
				dblI: tt.fields.dblI,
				dblJ: tt.fields.dblJ,
				dblK: tt.fields.dblK,
			}
			gotR, gotI, gotJ, gotK := h.ValInt()
			if !reflect.DeepEqual(gotR, tt.wantR) {
				t.Errorf("ValInt() gotR = %v, want %v", gotR, tt.wantR)
			}
			if !reflect.DeepEqual(gotI, tt.wantI) {
				t.Errorf("ValInt() gotI = %v, want %v", gotI, tt.wantI)
			}
			if !reflect.DeepEqual(gotJ, tt.wantJ) {
				t.Errorf("ValInt() gotJ = %v, want %v", gotJ, tt.wantJ)
			}
			if !reflect.DeepEqual(gotK, tt.wantK) {
				t.Errorf("ValInt() gotK = %v, want %v", gotK, tt.wantK)
			}
		})
	}
}
