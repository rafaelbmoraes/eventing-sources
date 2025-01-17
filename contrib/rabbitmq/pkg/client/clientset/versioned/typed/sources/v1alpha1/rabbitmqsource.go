/*
Copyright 2019 The Knative Authors

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

package v1alpha1

import (
	v1alpha1 "github.com/knative/eventing-contrib/contrib/rabbitmq/pkg/apis/sources/v1alpha1"
	scheme "github.com/knative/eventing-contrib/contrib/rabbitmq/pkg/client/clientset/versioned/scheme"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// RabbitmqSourcesGetter has a method to return a RabbitmqSourceInterface.
// A group's client should implement this interface.
type RabbitmqSourcesGetter interface {
	RabbitmqSources(namespace string) RabbitmqSourceInterface
}

// RabbitmqSourceInterface has methods to work with RabbitmqSource resources.
type RabbitmqSourceInterface interface {
	Create(*v1alpha1.RabbitmqSource) (*v1alpha1.RabbitmqSource, error)
	Update(*v1alpha1.RabbitmqSource) (*v1alpha1.RabbitmqSource, error)
	UpdateStatus(*v1alpha1.RabbitmqSource) (*v1alpha1.RabbitmqSource, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.RabbitmqSource, error)
	List(opts v1.ListOptions) (*v1alpha1.RabbitmqSourceList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.RabbitmqSource, err error)
	RabbitmqSourceExpansion
}

// rabbitmqSources implements RabbitmqSourceInterface
type rabbitmqSources struct {
	client rest.Interface
	ns     string
}

// newRabbitmqSources returns a RabbitmqSources
func newRabbitmqSources(c *SourcesV1alpha1Client, namespace string) *rabbitmqSources {
	return &rabbitmqSources{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the rabbitmqSource, and returns the corresponding rabbitmqSource object, and an error if there is any.
func (c *rabbitmqSources) Get(name string, options v1.GetOptions) (result *v1alpha1.RabbitmqSource, err error) {
	result = &v1alpha1.RabbitmqSource{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("rabbitmqsources").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of RabbitmqSources that match those selectors.
func (c *rabbitmqSources) List(opts v1.ListOptions) (result *v1alpha1.RabbitmqSourceList, err error) {
	result = &v1alpha1.RabbitmqSourceList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("rabbitmqsources").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested rabbitmqSources.
func (c *rabbitmqSources) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("rabbitmqsources").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a rabbitmqSource and creates it.  Returns the server's representation of the rabbitmqSource, and an error, if there is any.
func (c *rabbitmqSources) Create(rabbitmqSource *v1alpha1.RabbitmqSource) (result *v1alpha1.RabbitmqSource, err error) {
	result = &v1alpha1.RabbitmqSource{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("rabbitmqsources").
		Body(rabbitmqSource).
		Do().
		Into(result)
	return
}

// Update takes the representation of a rabbitmqSource and updates it. Returns the server's representation of the rabbitmqSource, and an error, if there is any.
func (c *rabbitmqSources) Update(rabbitmqSource *v1alpha1.RabbitmqSource) (result *v1alpha1.RabbitmqSource, err error) {
	result = &v1alpha1.RabbitmqSource{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("rabbitmqsources").
		Name(rabbitmqSource.Name).
		Body(rabbitmqSource).
		Do().
		Into(result)
	return
}

// UpdateStatus was generated because the type contains a Status member.
// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().

func (c *rabbitmqSources) UpdateStatus(rabbitmqSource *v1alpha1.RabbitmqSource) (result *v1alpha1.RabbitmqSource, err error) {
	result = &v1alpha1.RabbitmqSource{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("rabbitmqsources").
		Name(rabbitmqSource.Name).
		SubResource("status").
		Body(rabbitmqSource).
		Do().
		Into(result)
	return
}

// Delete takes name of the rabbitmqSource and deletes it. Returns an error if one occurs.
func (c *rabbitmqSources) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("rabbitmqsources").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *rabbitmqSources) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("rabbitmqsources").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched rabbitmqSource.
func (c *rabbitmqSources) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.RabbitmqSource, err error) {
	result = &v1alpha1.RabbitmqSource{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("rabbitmqsources").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
