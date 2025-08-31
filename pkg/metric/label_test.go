package metric

import "testing"

func TestLabEncode(t *testing.T) {
	result := marshalMetricName("cpu.busy", []Label{
		{Name: "cpuid", Value: "1"},
		{Name: "cpuid2", Value: "22"},
	})
	t.Logf("r = %s", result)
}
