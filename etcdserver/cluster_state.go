/*
   Copyright 2014 CoreOS, Inc.

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

package etcdserver

import (
	"errors"
)

const (
	ClusterStateValueNew      = "new"
	ClusterStateValueExisting = "existing"
)

var (
	ClusterStateValues = []string{
		ClusterStateValueNew,
		ClusterStateValueExisting,
	}
)

// ClusterState implements the flag.Value interface.
type ClusterState string

// Set verifies the argument to be a valid member of ClusterStateFlagValues
// before setting the underlying flag value.
func (cs *ClusterState) Set(s string) error {
	for _, v := range ClusterStateValues {
		if s == v {
			*cs = ClusterState(s)
			return nil
		}
	}

	return errors.New("invalid value")
}

func (cs *ClusterState) String() string {
	return string(*cs)
}
