package metric

import (
	"encoding/binary"
	"sort"
)

const (
	// The maximum length of label name.
	//
	// Longer names are truncated.
	maxLabelNameLen = 256

	// The maximum length of label value.
	//
	// Longer values are truncated.
	maxLabelValueLen = 16 * 1024
)

type (
	Label struct {
		Name  string
		Value string
	}
)

// MarshalUint16 appends marshaled v to dst and returns the result.
func MarshalUint16(dst []byte, u uint16) []byte {
	return append(dst, byte(u>>8), byte(u))
}

// UnmarshalUint16 returns unmarshaled uint32 from src.
func UnmarshalUint16(src []byte) uint16 {
	// This is faster than the manual conversion.
	return binary.BigEndian.Uint16(src)
}

func marshalMetricName(metric string, labels []Label) string {
	if len(labels) == 0 {
		return metric
	}
	invalid := func(name, value string) bool {
		return name == "" || value == ""
	}

	// Determine the bytes size in advance.
	size := len(metric) + 2
	sort.Slice(labels, func(i, j int) bool {
		return labels[i].Name < labels[j].Name
	})
	for i := range labels {
		label := &labels[i]
		if invalid(label.Name, label.Value) {
			continue
		}
		if len(label.Name) > maxLabelNameLen {
			label.Name = label.Name[:maxLabelNameLen]
		}
		if len(label.Value) > maxLabelValueLen {
			label.Value = label.Value[:maxLabelValueLen]
		}
		size += len(label.Name)
		size += len(label.Value)
		size += 4
	}

	// Start building the bytes.
	out := make([]byte, 0, size)
	out = MarshalUint16(out, uint16(len(metric)))
	out = append(out, metric...)
	for i := range labels {
		label := &labels[i]
		if invalid(label.Name, label.Value) {
			continue
		}
		out = MarshalUint16(out, uint16(len(label.Name)))
		out = append(out, label.Name...)
		out = MarshalUint16(out, uint16(len(label.Value)))
		out = append(out, label.Value...)
	}
	return string(out)
}
