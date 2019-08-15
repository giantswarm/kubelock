/*
Copyright 2019 Giant Swarm GmbH.

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

package fake

import (
	v1alpha1 "github.com/giantswarm/apiextensions/pkg/apis/application/v1alpha1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeAppCatalogs implements AppCatalogInterface
type FakeAppCatalogs struct {
	Fake *FakeApplicationV1alpha1
}

var appcatalogsResource = schema.GroupVersionResource{Group: "application.giantswarm.io", Version: "v1alpha1", Resource: "appcatalogs"}

var appcatalogsKind = schema.GroupVersionKind{Group: "application.giantswarm.io", Version: "v1alpha1", Kind: "AppCatalog"}

// Get takes name of the appCatalog, and returns the corresponding appCatalog object, and an error if there is any.
func (c *FakeAppCatalogs) Get(name string, options v1.GetOptions) (result *v1alpha1.AppCatalog, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootGetAction(appcatalogsResource, name), &v1alpha1.AppCatalog{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.AppCatalog), err
}

// List takes label and field selectors, and returns the list of AppCatalogs that match those selectors.
func (c *FakeAppCatalogs) List(opts v1.ListOptions) (result *v1alpha1.AppCatalogList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootListAction(appcatalogsResource, appcatalogsKind, opts), &v1alpha1.AppCatalogList{})
	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &v1alpha1.AppCatalogList{ListMeta: obj.(*v1alpha1.AppCatalogList).ListMeta}
	for _, item := range obj.(*v1alpha1.AppCatalogList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested appCatalogs.
func (c *FakeAppCatalogs) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewRootWatchAction(appcatalogsResource, opts))
}

// Create takes the representation of a appCatalog and creates it.  Returns the server's representation of the appCatalog, and an error, if there is any.
func (c *FakeAppCatalogs) Create(appCatalog *v1alpha1.AppCatalog) (result *v1alpha1.AppCatalog, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootCreateAction(appcatalogsResource, appCatalog), &v1alpha1.AppCatalog{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.AppCatalog), err
}

// Update takes the representation of a appCatalog and updates it. Returns the server's representation of the appCatalog, and an error, if there is any.
func (c *FakeAppCatalogs) Update(appCatalog *v1alpha1.AppCatalog) (result *v1alpha1.AppCatalog, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootUpdateAction(appcatalogsResource, appCatalog), &v1alpha1.AppCatalog{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.AppCatalog), err
}

// Delete takes name of the appCatalog and deletes it. Returns an error if one occurs.
func (c *FakeAppCatalogs) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewRootDeleteAction(appcatalogsResource, name), &v1alpha1.AppCatalog{})
	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeAppCatalogs) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewRootDeleteCollectionAction(appcatalogsResource, listOptions)

	_, err := c.Fake.Invokes(action, &v1alpha1.AppCatalogList{})
	return err
}

// Patch applies the patch and returns the patched appCatalog.
func (c *FakeAppCatalogs) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.AppCatalog, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewRootPatchSubresourceAction(appcatalogsResource, name, pt, data, subresources...), &v1alpha1.AppCatalog{})
	if obj == nil {
		return nil, err
	}
	return obj.(*v1alpha1.AppCatalog), err
}
