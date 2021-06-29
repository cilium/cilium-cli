package install

import (
	"gopkg.in/check.v1"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (b *InstallSuite) TestTail(c *check.C) {
	cases := []struct {
		given  []*corev1.Pod
		expect []*corev1.Pod
	}{
		{
			given: []*corev1.Pod{
				{ObjectMeta: v1.ObjectMeta{Name: "foo"}}},
			expect: nil,
		},
		{
			given: []*corev1.Pod{
				{ObjectMeta: v1.ObjectMeta{Name: "foo"}},
				{ObjectMeta: v1.ObjectMeta{Name: "bar"}}},
			expect: []*corev1.Pod{{ObjectMeta: v1.ObjectMeta{Name: "bar"}}},
		},
		{
			given:  []*corev1.Pod{},
			expect: nil,
		},
		{
			given: []*corev1.Pod{
				{ObjectMeta: v1.ObjectMeta{Name: "foo"}},
				{ObjectMeta: v1.ObjectMeta{Name: "bar"}},
				{ObjectMeta: v1.ObjectMeta{Name: "qux"}}},
			expect: []*corev1.Pod{{ObjectMeta: v1.ObjectMeta{Name: "bar"}}, {ObjectMeta: v1.ObjectMeta{Name: "qux"}}},
		},
	}
	for _, cs := range cases {
		got := tail(cs.given)
		c.Assert(got, check.DeepEquals, cs.expect, check.Commentf("Got %v, expected %v", got, cs.expect))
	}
}

func (b *InstallSuite) TestGroupPodsByNamespace(c *check.C) {
	t := corev1.PodList{
		Items: []corev1.Pod{
			{ObjectMeta: v1.ObjectMeta{Name: "foo", Namespace: "berlin"}},
			{ObjectMeta: v1.ObjectMeta{Name: "bar", Namespace: "berlin"}},
			{ObjectMeta: v1.ObjectMeta{Name: "qux", Namespace: "caracas"}},
		},
	}

	nsToPods := groupPodsByNamespace(&t)
	c.Assert(nsToPods["berlin"], check.DeepEquals, []*corev1.Pod{{ObjectMeta: v1.ObjectMeta{Name: "foo", Namespace: "berlin"}},
		{ObjectMeta: v1.ObjectMeta{Name: "bar", Namespace: "berlin"}}})
	c.Assert(nsToPods["caracas"], check.DeepEquals, []*corev1.Pod{{ObjectMeta: v1.ObjectMeta{Name: "qux", Namespace: "caracas"}}})
}
