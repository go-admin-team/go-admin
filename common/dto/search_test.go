package dto

import (
	"reflect"
	"testing"
)

// GetIds 曾在 else 分支中重复追加 Id：仅传 Id 时返回 [5 5]，
// 导致删除接口对同一条记录执行两次。
func TestGeneralDelDtoGetIds(t *testing.T) {
	cases := []struct {
		name string
		dto  GeneralDelDto
		want []int
	}{
		{"仅 Id", GeneralDelDto{Id: 5}, []int{5}},
		{"仅 Ids", GeneralDelDto{Ids: []int{1, 2}}, []int{1, 2}},
		{"Id 与 Ids 并存", GeneralDelDto{Id: 5, Ids: []int{1, 2}}, []int{5, 1, 2}},
		{"Ids 含非正数被过滤", GeneralDelDto{Ids: []int{0, -1, 3}}, []int{3}},
		{"Id 为 0 视为未传", GeneralDelDto{Id: 0, Ids: []int{7}}, []int{7}},
		{"全部为空时回退到 0（全量删除约定）", GeneralDelDto{}, []int{0}},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := c.dto.GetIds()
			if !reflect.DeepEqual(got, c.want) {
				t.Errorf("GetIds() = %v, want %v", got, c.want)
			}
		})
	}
}
