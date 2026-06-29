package leetcode

import "testing"

func TestAllOne(t *testing.T) {
	obj := AllOneConstructor()

	if obj.GetMaxKey() != "" || obj.GetMinKey() != "" {
		t.Error("empty structure should return empty strings")
	}

	obj.Inc("a")
	obj.Inc("b")
	obj.Inc("b")
	obj.Inc("c")
	obj.Inc("c")
	obj.Inc("c")

	if obj.GetMaxKey() != "c" {
		t.Errorf("GetMaxKey = %q, want %q", obj.GetMaxKey(), "c")
	}
	if obj.GetMinKey() != "a" {
		t.Errorf("GetMinKey = %q, want %q", obj.GetMinKey(), "a")
	}

	obj.Dec("c")
	obj.Dec("c")
	obj.Dec("c")

	if obj.GetMaxKey() != "b" {
		t.Errorf("after dec c, GetMaxKey = %q, want %q", obj.GetMaxKey(), "b")
	}
}
