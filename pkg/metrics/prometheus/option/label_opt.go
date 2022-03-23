package option

import "github.com/prometheus/client_golang/prometheus"

func WithBindLabels(labels ...BindLabel) bindLabelOption {
	return labels
}

func NewBaseBindLabels(name, module, zone string) bindLabelOption {
	return []BindLabel{
		BindLabel{key: LabelName, value: name},
		BindLabel{key: LabelZone, value: zone},
		BindLabel{key: LabelModule, value: module},
	}
}

type BindLabel struct {
	key   string
	value string
}

func NewBindLabel(key, val string) BindLabel {
	return BindLabel{key: key, value: val}
}

type bindLabelOption []BindLabel

func (this bindLabelOption) Type() string {
	return "___bind_labels"
}

func (this bindLabelOption) Append(labels []BindLabel) bindLabelOption {
	return append(this, labels...)
}

func (this bindLabelOption) Invoke(reg prometheus.Registerer) prometheus.Registerer {
	labels := prometheus.Labels{}
	for _, label := range this {
		labels[label.key] = label.value
	}
	return prometheus.WrapRegistererWith(labels, reg)
}
