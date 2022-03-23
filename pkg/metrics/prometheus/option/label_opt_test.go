package option

import "testing"

func TestLabel(t *testing.T) {
	labels := bindLabelOption{}
	labels = labels.Append([]BindLabel{NewBindLabel("key", "val")})
	t.Log(labels)
}
