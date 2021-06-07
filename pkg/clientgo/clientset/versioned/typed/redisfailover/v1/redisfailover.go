/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Code generated by client-gen. DO NOT EDIT.

package v1

import (
	"context"

	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"

	v1 "github.com/erda-project/erda/pkg/clientgo/apis/redisfailover/v1"
	scheme "github.com/erda-project/erda/pkg/clientgo/clientset/versioned/scheme"
)

// RedisFailoversGetter has a method to return a RedisFailoverInterface.
// A group's client should implement this interface.
type RedisFailoversGetter interface {
	RedisFailovers(namespace string) RedisFailoverInterface
}

// RedisFailoverInterface has methods to work with RedisFailover resources.
type RedisFailoverInterface interface {
	Create(context.Context, *v1.RedisFailover) (*v1.RedisFailover, error)
	Update(context.Context, *v1.RedisFailover) (*v1.RedisFailover, error)
	Delete(ctx context.Context, name string, options *meta_v1.DeleteOptions) error
	DeleteCollection(ctx context.Context, options *meta_v1.DeleteOptions, listOptions meta_v1.ListOptions) error
	Get(ctx context.Context, name string, options meta_v1.GetOptions) (*v1.RedisFailover, error)
	List(ctx context.Context, opts meta_v1.ListOptions) (*v1.RedisFailoverList, error)
	Watch(ctx context.Context, opts meta_v1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.RedisFailover, err error)
	RedisFailoverExpansion
}

// redisFailovers implements RedisFailoverInterface
type redisFailovers struct {
	client rest.Interface
	ns     string
}

// newRedisFailovers returns a RedisFailovers
func newRedisFailovers(c *DatabasesV1Client, namespace string) *redisFailovers {
	return &redisFailovers{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the redisFailover, and returns the corresponding redisFailover object, and an error if there is any.
func (c *redisFailovers) Get(ctx context.Context, name string, options meta_v1.GetOptions) (result *v1.RedisFailover, err error) {
	result = &v1.RedisFailover{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("redisfailovers").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of RedisFailovers that match those selectors.
func (c *redisFailovers) List(ctx context.Context, opts meta_v1.ListOptions) (result *v1.RedisFailoverList, err error) {
	result = &v1.RedisFailoverList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("redisfailovers").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do(ctx).
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested redisFailovers.
func (c *redisFailovers) Watch(ctx context.Context, opts meta_v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("redisfailovers").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch(ctx)
}

// Create takes the representation of a redisFailover and creates it.  Returns the server's representation of the redisFailover, and an error, if there is any.
func (c *redisFailovers) Create(ctx context.Context, redisFailover *v1.RedisFailover) (result *v1.RedisFailover, err error) {
	result = &v1.RedisFailover{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("redisfailovers").
		Body(redisFailover).
		Do(ctx).
		Into(result)
	return
}

// Update takes the representation of a redisFailover and updates it. Returns the server's representation of the redisFailover, and an error, if there is any.
func (c *redisFailovers) Update(ctx context.Context, redisFailover *v1.RedisFailover) (result *v1.RedisFailover, err error) {
	result = &v1.RedisFailover{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("redisfailovers").
		Name(redisFailover.Name).
		Body(redisFailover).
		Do(ctx).
		Into(result)
	return
}

// Delete takes name of the redisFailover and deletes it. Returns an error if one occurs.
func (c *redisFailovers) Delete(ctx context.Context, name string, options *meta_v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("redisfailovers").
		Name(name).
		Body(options).
		Do(ctx).
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *redisFailovers) DeleteCollection(ctx context.Context, options *meta_v1.DeleteOptions, listOptions meta_v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("redisfailovers").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do(ctx).
		Error()
}

// Patch applies the patch and returns the patched redisFailover.
func (c *redisFailovers) Patch(ctx context.Context, name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.RedisFailover, err error) {
	result = &v1.RedisFailover{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("redisfailovers").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do(ctx).
		Into(result)
	return
}
