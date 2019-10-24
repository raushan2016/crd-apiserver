/*
Copyright 2018 The Kubernetes Authors.

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

package options

import (
	"fmt"
	"io"
	"net"

	"github.com/raushan2016/crd-apiserver/pkg/apiserver"
	"k8s.io/apimachinery/pkg/runtime/schema"
	utilerrors "k8s.io/apimachinery/pkg/util/errors"
	genericapiserver "k8s.io/apiserver/pkg/server"
	genericoptions "k8s.io/apiserver/pkg/server/options"
)

const defaultEtcdPathPrefix = "/registry/instance.raushank.io"

// CustomResourceDefinitionsServerOptions describes the runtime options of an apiextensions-apiserver.
type CustomResourceDefinitionsServerOptions struct {
	RecommendedOptions *genericoptions.RecommendedOptions
	Admission          *genericoptions.AdmissionOptions
	StdOut             io.Writer
	StdErr             io.Writer
}

const GroupName = "instance.raushank.io"

// SchemeGroupVersion is group version used to register these objects
var SchemeGroupVersion = schema.GroupVersion{Group: GroupName, Version: "v1beta1"}

// NewCustomResourceDefinitionsServerOptions creates default options of an apiextensions-apiserver.
func NewCustomResourceDefinitionsServerOptions(out, errOut io.Writer) *CustomResourceDefinitionsServerOptions {
	o := &CustomResourceDefinitionsServerOptions{
		RecommendedOptions: genericoptions.NewRecommendedOptions(
			defaultEtcdPathPrefix,
			apiserver.Codecs.LegacyCodec(SchemeGroupVersion),
			genericoptions.NewProcessInfo("raushankapiserver", "mir-apiserver-system"),
		),
		Admission: genericoptions.NewAdmissionOptions(),

		StdOut: out,
		StdErr: errOut,
	}

	return o
}

// Validate validates the apiextensions-apiserver options.
func (o CustomResourceDefinitionsServerOptions) Validate() error {
	errors := []error{}
	errors = append(errors, o.RecommendedOptions.Validate()...)
	errors = append(errors, o.Admission.Validate()...)
	return utilerrors.NewAggregate(errors)
}

// Complete fills in missing options.
func (o *CustomResourceDefinitionsServerOptions) Complete() error {
	return nil
}

// Config returns an apiextensions-apiserver configuration.
func (o CustomResourceDefinitionsServerOptions) Config() (*apiserver.Config, error) {
	// TODO have a "real" external address
	if err := o.RecommendedOptions.SecureServing.MaybeDefaultWithSelfSignedCerts("localhost", nil, []net.IP{net.ParseIP("127.0.0.1")}); err != nil {
		return nil, fmt.Errorf("error creating self-signed certificates: %v", err)
	}

	serverConfig := genericapiserver.NewRecommendedConfig(apiserver.Codecs)
	if err := o.RecommendedOptions.ApplyTo(serverConfig); err != nil {
		return nil, err
	}

	config := &apiserver.Config{
		GenericConfig: serverConfig,
	}
	return config, nil
}
