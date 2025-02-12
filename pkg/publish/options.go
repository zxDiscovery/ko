// Copyright 2018 Google LLC All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package publish

import (
	"log"
	"net/http"

	"github.com/google/go-containerregistry/pkg/authn"
	"github.com/google/go-containerregistry/pkg/name"
)

// WithTransport is a functional option for overriding the default transport
// on a default publisher.
func WithTransport(t http.RoundTripper) Option {
	return func(i *defaultOpener) error {
		i.t = t
		return nil
	}
}

// WithAuth is a functional option for overriding the default authenticator
// on a default publisher.
func WithAuth(auth authn.Authenticator) Option {
	return func(i *defaultOpener) error {
		i.auth = auth
		return nil
	}
}

// WithAuthFromKeychain is a functional option for overriding the default
// authenticator on a default publisher using an authn.Keychain
func WithAuthFromKeychain(keys authn.Keychain) Option {
	return func(i *defaultOpener) error {
		// We parse this lazily because it is a repository prefix, which
		// means that docker.io/mattmoor actually gets interpreted as
		// docker.io/library/mattmoor, which gets tricky when we start
		// appending things to it in the publisher.
		repo, err := name.NewRepository(i.base)
		if err != nil {
			return err
		}
		auth, err := keys.Resolve(repo.Registry)
		if err != nil {
			return err
		}
		if auth == authn.Anonymous {
			log.Println("No matching credentials were found, falling back on anonymous")
		}
		i.auth = auth
		return nil
	}
}

// WithNamer is a functional option for overriding the image naming behavior
// in our default publisher.
func WithNamer(n Namer) Option {
	return func(i *defaultOpener) error {
		i.namer = n
		return nil
	}
}

// WithTags is a functional option for overriding the image tags
func WithTags(tags []string) Option {
	return func(i *defaultOpener) error {
		i.tags = tags
		return nil
	}
}

func Insecure(b bool) Option {
	return func(i *defaultOpener) error {
		i.insecure = b
		return nil
	}
}