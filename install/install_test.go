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
				{ObjectMeta: v1.ObjectMeta{Name: "Foo"}}},
			expect: nil,
		},
		{
			given: []*corev1.Pod{
				{ObjectMeta: v1.ObjectMeta{Name: "Foo"}},
				{ObjectMeta: v1.ObjectMeta{Name: "Bar"}}},
			expect: []*corev1.Pod{{ObjectMeta: v1.ObjectMeta{Name: "Bar"}}},
		},
		{
			given:  []*corev1.Pod{},
			expect: nil,
		},
		{
			given: []*corev1.Pod{
				{ObjectMeta: v1.ObjectMeta{Name: "Foo"}},
				{ObjectMeta: v1.ObjectMeta{Name: "Bar"}},
				{ObjectMeta: v1.ObjectMeta{Name: "Qux"}}},
			expect: []*corev1.Pod{{ObjectMeta: v1.ObjectMeta{Name: "Bar"}}, {ObjectMeta: v1.ObjectMeta{Name: "Qux"}}},
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
			{ObjectMeta: v1.ObjectMeta{Name: "Foo", Namespace: "Berlin"}},
			{ObjectMeta: v1.ObjectMeta{Name: "Bar", Namespace: "Berlin"}},
			{ObjectMeta: v1.ObjectMeta{Name: "Qux", Namespace: "Caracas"}},
		},
	}

	nsToPods := groupPodsByNamespace(&t)
	c.Assert(nsToPods["Berlin"], check.DeepEquals, []*corev1.Pod{{ObjectMeta: v1.ObjectMeta{Name: "Foo", Namespace: "Berlin"}},
		{ObjectMeta: v1.ObjectMeta{Name: "Bar", Namespace: "Berlin"}}})
	c.Assert(nsToPods["Caracas"], check.DeepEquals, []*corev1.Pod{{ObjectMeta: v1.ObjectMeta{Name: "Qux", Namespace: "Caracas"}}})
}

func (b *InstallSuite) TestRunEvictions(c *check.C) {
	// it should not evict cep pods
	// it should not evict HostNetwork pods
	// t := corev1.PodList{
	// 	Items: []corev1.Pod{
	// 		{ObjectMeta: v1.ObjectMeta{Name: "Foo", Namespace: "Berlin"}, Spec: corev1.PodSpec{HostNetwork: true}},
	// 		{ObjectMeta: v1.ObjectMeta{Name: "ciliumEP", Namespace: "Berlin"}},
	// 		{ObjectMeta: v1.ObjectMeta{Name: "Qux", Namespace: "Caracas"}},
	// 	},
	// }
	// cepMap := map[string]struct{}{"ciliumEP": struct{}{}}
	// k := K8sInstaller{}
	// k.runEvictions(context.Background(), cepMap, &t)
}
