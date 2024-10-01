// SPDX-License-Identifier: Apache-2.0
// Copyright Authors of Cilium

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	"context"

	v2alpha1 "github.com/cilium/cilium/pkg/k8s/apis/cilium.io/v2alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeCiliumLoadBalancerIPPools implements CiliumLoadBalancerIPPoolInterface
type FakeCiliumLoadBalancerIPPools struct {
	Fake *FakeCiliumV2alpha1
}

var ciliumloadbalancerippoolsResource = v2alpha1.SchemeGroupVersion.WithResource("ciliumloadbalancerippools")

var ciliumloadbalancerippoolsKind = v2alpha1.SchemeGroupVersion.WithKind("CiliumLoadBalancerIPPool")

// Get takes name of the ciliumLoadBalancerIPPool, and returns the corresponding ciliumLoadBalancerIPPool object, and an error if there is any.
func (c *FakeCiliumLoadBalancerIPPools) Get(ctx context.Context, name string, options v1.GetOptions) (result *v2alpha1.CiliumLoadBalancerIPPool, err error) {
	emptyResult := &v2alpha1.CiliumLoadBalancerIPPool{}
	obj, err := c.Fake.
		Invokes(testing.NewRootGetActionWithOptions(ciliumloadbalancerippoolsResource, name, options), emptyResult)
	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v2alpha1.CiliumLoadBalancerIPPool), err
}

// List takes label and field selectors, and returns the list of CiliumLoadBalancerIPPools that match those selectors.
func (c *FakeCiliumLoadBalancerIPPools) List(ctx context.Context, opts v1.ListOptions) (result *v2alpha1.CiliumLoadBalancerIPPoolList, err error) {
	emptyResult := &v2alpha1.CiliumLoadBalancerIPPoolList{}
	obj, err := c.Fake.
		Invokes(testing.NewRootListActionWithOptions(ciliumloadbalancerippoolsResource, ciliumloadbalancerippoolsKind, opts), emptyResult)
	if obj == nil {
		return emptyResult, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v2alpha1.CiliumLoadBalancerIPPoolList{ListMeta: obj.(*v2alpha1.CiliumLoadBalancerIPPoolList).ListMeta}
	for _, item := range obj.(*v2alpha1.CiliumLoadBalancerIPPoolList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested ciliumLoadBalancerIPPools.
func (c *FakeCiliumLoadBalancerIPPools) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchActionWithOptions(ciliumloadbalancerippoolsResource, opts))
}

// Create takes the representation of a ciliumLoadBalancerIPPool and creates it.  Returns the server's representation of the ciliumLoadBalancerIPPool, and an error, if there is any.
func (c *FakeCiliumLoadBalancerIPPools) Create(ctx context.Context, ciliumLoadBalancerIPPool *v2alpha1.CiliumLoadBalancerIPPool, opts v1.CreateOptions) (result *v2alpha1.CiliumLoadBalancerIPPool, err error) {
	emptyResult := &v2alpha1.CiliumLoadBalancerIPPool{}
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateActionWithOptions(ciliumloadbalancerippoolsResource, ciliumLoadBalancerIPPool, opts), emptyResult)
	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v2alpha1.CiliumLoadBalancerIPPool), err
}

// Update takes the representation of a ciliumLoadBalancerIPPool and updates it. Returns the server's representation of the ciliumLoadBalancerIPPool, and an error, if there is any.
func (c *FakeCiliumLoadBalancerIPPools) Update(ctx context.Context, ciliumLoadBalancerIPPool *v2alpha1.CiliumLoadBalancerIPPool, opts v1.UpdateOptions) (result *v2alpha1.CiliumLoadBalancerIPPool, err error) {
	emptyResult := &v2alpha1.CiliumLoadBalancerIPPool{}
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateActionWithOptions(ciliumloadbalancerippoolsResource, ciliumLoadBalancerIPPool, opts), emptyResult)
	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v2alpha1.CiliumLoadBalancerIPPool), err
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *FakeCiliumLoadBalancerIPPools) UpdateStatus(ctx context.Context, ciliumLoadBalancerIPPool *v2alpha1.CiliumLoadBalancerIPPool, opts v1.UpdateOptions) (result *v2alpha1.CiliumLoadBalancerIPPool, err error) {
	emptyResult := &v2alpha1.CiliumLoadBalancerIPPool{}
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateSubresourceActionWithOptions(ciliumloadbalancerippoolsResource, "status", ciliumLoadBalancerIPPool, opts), emptyResult)
	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v2alpha1.CiliumLoadBalancerIPPool), err
}

// Delete takes name of the ciliumLoadBalancerIPPool and deletes it. Returns an error if one occurs.
func (c *FakeCiliumLoadBalancerIPPools) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteActionWithOptions(ciliumloadbalancerippoolsResource, name, opts), &v2alpha1.CiliumLoadBalancerIPPool{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeCiliumLoadBalancerIPPools) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionActionWithOptions(ciliumloadbalancerippoolsResource, opts, listOpts)

	_, err := c.Fake.Invokes(action, &v2alpha1.CiliumLoadBalancerIPPoolList{})
	return err
}

// Patch applies the patch and returns the patched ciliumLoadBalancerIPPool.
func (c *FakeCiliumLoadBalancerIPPools) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v2alpha1.CiliumLoadBalancerIPPool, err error) {
	emptyResult := &v2alpha1.CiliumLoadBalancerIPPool{}
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceActionWithOptions(ciliumloadbalancerippoolsResource, name, pt, data, opts, subresources...), emptyResult)
	if obj == nil {
		return emptyResult, err
	}
	return obj.(*v2alpha1.CiliumLoadBalancerIPPool), err
}
