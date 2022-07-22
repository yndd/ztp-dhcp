package structs

import (
	"encoding/json"
	"fmt"
	"testing"

	"bou.ke/monkey"
)

func TestClientIdentifier_String(t *testing.T) {
	type fields struct {
		CIType CITypeEnum
		Value  string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		preRun  func()
		postRun func()
	}{
		{
			name:    "Json Marshall Fails",
			preRun:  patchJsonMarshal,
			postRun: func() { monkey.Unpatch(json.Marshal) },
			fields: fields{
				CIType: String,
				Value:  "FooBar",
			},
			want: "ERROR Marshalling",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cir := &ClientIdentifier{
				CIType: tt.fields.CIType,
				Value:  tt.fields.Value,
			}
			tt.preRun()
			if got := cir.String(); got != tt.want {
				t.Errorf("ClientIdentifier.String() = %v, want %v", got, tt.want)
			}
			tt.postRun()
		})
	}
}

func patchJsonMarshal() {
	monkey.Patch(json.Marshal, func(v any) ([]byte, error) {
		return nil, fmt.Errorf("Fake Error")
	})
}
