package mysql

import (
	"reflect"
	"testing"

	"github.com/yangyuqian/genus"
	"github.com/yangyuqian/genus/types"
)

func TestGetPlanItem(t *testing.T) {
	type args struct {
		opts *planItemOpts
	}
	tests := []struct {
		name     string
		args     args
		wantItem *genus.PlanItem
	}{
		{"OK", args{&planItemOpts{
			PlanType: types.SINGLETON,
			TmplDir:  "./",
			Suffix:   ".tpl",
		}}, &genus.PlanItem{
			PlanType:    "SINGLETON",
			TemplateDir: "./",
			Suffix:      ".tpl",
		}},
	}
	for _, tt := range tests {
		if gotItem := GetPlanItem(tt.args.opts); !reflect.DeepEqual(gotItem, tt.wantItem) {
			t.Errorf("%q. GetPlanItem() = %+v, want %+v", tt.name, gotItem, tt.wantItem)
		}
	}
}
