/*
Copyright 2017 The Kubernetes Authors.

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

package v1

import (
	"github.com/kubernetes-incubator/kube-arbitrator/pkg/apis/v1"

	"k8s.io/apimachinery/pkg/runtime"
	serializer "k8s.io/apimachinery/pkg/runtime/serializer"
	rest "k8s.io/client-go/rest"
)

type CrdV1Interface interface {
	RESTClient() rest.Interface
	QueueGetter
	TasksetGetter
}

// CrdV1Client is used to interact with features provided by the  group.
type CrdV1Client struct {
	restClient rest.Interface
}

func (c *CrdV1Client) Queues(namespace string) QueueInterface {
	return newQueues(c, namespace)
}

func (c *CrdV1Client) Tasksets(namespace string) TasksetInterface {
	return newTasksets(c, namespace)
}

// NewForConfig creates a new CoreV1Client for the given config.
func NewForConfig(c *rest.Config) (*CrdV1Client, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}
	return &CrdV1Client{client}, nil
}

// NewForConfigOrDie creates a new CrdV1Client for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *CrdV1Client {
	client, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return client
}

// New creates a new CrdV1Client for the given RESTClient.
func New(c rest.Interface) *CrdV1Client {
	return &CrdV1Client{c}
}

func setConfigDefaults(config *rest.Config) error {
	scheme := runtime.NewScheme()
	if err := v1.AddToScheme(scheme); err != nil {
		return err
	}

	config.GroupVersion = &v1.SchemeGroupVersion
	config.APIPath = "/apis"
	config.ContentType = runtime.ContentTypeJSON
	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: serializer.NewCodecFactory(scheme)}

	return nil
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *CrdV1Client) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.restClient
}
