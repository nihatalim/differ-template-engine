package domain

import (
	"reflect"
	"testing"
)

func TestUrl_GetParameters(t *testing.T) {
	tests := []struct {
		name string
		u    Url
		want []string
	}{
		{
			name: "",
			u:    Url("http://localhost:8080/productcontents/${contentId}?storefrontId=${storefrontId}&culture=${culture}"),
			want: []string{"contentId", "storefrontId", "culture"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parameters := tt.u.GetParameters()
			if got := parameters; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetParameters() = %v, want %v", got, tt.want)
			}
		})
	}
}
