// SPDX-License-Identifier: Apache-2.0
// Copyright 2017-2021 Authors of Cilium

// Code generated by client-gen. DO NOT EDIT.

package v2

import (
	"context"
	"time"

	v2 "github.com/cilium/cilium/pkg/k8s/apis/cilium.io/v2"
	scheme "github.com/cilium/cilium/pkg/k8s/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// CiliumEndpointsGetter has a method to return a CiliumEndpointInterface.
// A group's client should implement this interface.
type CiliumEndpointsGetter interface {
	CiliumEndpoints(namespace string) CiliumEndpointInterface
}

// CiliumEndpointInterface has methods to work with CiliumEndpoint resources.
type CiliumEndpointInterface interface {
	Create(ctx context.Context, ciliumEndpoint *v2.CiliumEndpoint, opts v1.CreateOptions) (*v2.CiliumEndpoint, error)
	Update(ctx context.Context, ciliumEndpoint *v2.CiliumEndpoint, opts v1.UpdateOptions) (*v2.CiliumEndpoint, error)
	UpdateStatus(ctx context.Context, ciliumEndpoint *v2.CiliumEndpoint, opts v1.UpdateOptions) (*v2.CiliumEndpoint, error)
	Delete(ctx context.Context, name string, opts v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error
	Get(ctx context.Context, name string, opts v1.GetOptions) (*v2.CiliumEndpoint, error)
	List(ctx context.Context, opts v1.ListOptions) (*v2.CiliumEndpointList, error)
	Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v2.CiliumEndpoint, err error)
	CiliumEndpointExpansion
}

// ciliumEndpoints implements CiliumEndpointInterface
type ciliumEndpoints struct {
	client rest.Interface
	ns     string
}

// newCiliumEndpoints returns a CiliumEndpoints
func newCiliumEndpoints(c *CiliumV2Client, namespace string) *ciliumEndpoints {
	return &ciliumEndpoints{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the ciliumEndpoint, and returns the corresponding ciliumEndpoint object, and an error if there is any.
func (c *ciliumEndpoints) Get(ctx context.Context, name string, options v1.GetOptions) (result *v2.CiliumEndpoint, err error) {
	result = &v2.CiliumEndpoint{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("ciliumendpoints").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of CiliumEndpoints that match those selectors.
func (c *ciliumEndpoints) List(ctx context.Context, opts v1.ListOptions) (result *v2.CiliumEndpointList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	result = &v2.CiliumEndpointList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("ciliumendpoints").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested ciliumEndpoints.
func (c *ciliumEndpoints) Watch(ctx context.Context, opts v1.ListOptions) (watch.Interface, error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("ciliumendpoints").
		VersionedParams(&opts, scheme.ParameterCodec).
		Timeout(timeout).
		Watch(ctx)
}

// Create takes the representation of a ciliumEndpoint and creates it.  Returns the server's representation of the ciliumEndpoint, and an error, if there is any.
func (c *ciliumEndpoints) Create(ctx context.Context, ciliumEndpoint *v2.CiliumEndpoint, opts v1.CreateOptions) (result *v2.CiliumEndpoint, err error) {
	result = &v2.CiliumEndpoint{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("ciliumendpoints").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(ciliumEndpoint).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a ciliumEndpoint and updates it. Returns the server's representation of the ciliumEndpoint, and an error, if there is any.
func (c *ciliumEndpoints) Update(ctx context.Context, ciliumEndpoint *v2.CiliumEndpoint, opts v1.UpdateOptions) (result *v2.CiliumEndpoint, err error) {
	result = &v2.CiliumEndpoint{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("ciliumendpoints").
		Name(ciliumEndpoint.Name).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(ciliumEndpoint).
		Do(ctx).
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
func (c *ciliumEndpoints) UpdateStatus(ctx context.Context, ciliumEndpoint *v2.CiliumEndpoint, opts v1.UpdateOptions) (result *v2.CiliumEndpoint, err error) {
	result = &v2.CiliumEndpoint{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("ciliumendpoints").
		Name(ciliumEndpoint.Name).
		SubResource("status").
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(ciliumEndpoint).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the ciliumEndpoint and deletes it. Returns an error if one occurs.
func (c *ciliumEndpoints) Delete(ctx context.Context, name string, opts v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("ciliumendpoints").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *ciliumEndpoints) DeleteCollection(ctx context.Context, opts v1.DeleteOptions, listOpts v1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}
	return c.client.Delete().
		Namespace(c.ns).
		Resource("ciliumendpoints").
		VersionedParams(&listOpts, scheme.ParameterCodec).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched ciliumEndpoint.
func (c *ciliumEndpoints) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts v1.PatchOptions, subresources ...string) (result *v2.CiliumEndpoint, err error) {
	result = &v2.CiliumEndpoint{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("ciliumendpoints").
		Name(name).
		SubResource(subresources...).
		VersionedParams(&opts, scheme.ParameterCodec).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
