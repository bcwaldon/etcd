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

package client

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"time"

	"github.com/coreos/etcd/Godeps/_workspace/src/code.google.com/p/go.net/context"
)

var (
	ErrTimeout            = context.DeadlineExceeded
	DefaultRequestTimeout = 5 * time.Second
)

type roundTripResponse struct {
	resp *http.Response
	err  error
}

type httpClient struct {
	transport CancelableTransport
	endpoint  url.URL
	timeout   time.Duration
}

func (c *httpClient) Do(ctx context.Context, act HTTPAction) (*http.Response, []byte, error) {
	req := act.HTTPRequest(c.endpoint)

	rtchan := make(chan roundTripResponse, 1)
	go func() {
		resp, err := c.transport.RoundTrip(req)
		rtchan <- roundTripResponse{resp: resp, err: err}
		close(rtchan)
	}()

	var resp *http.Response
	var err error

	select {
	case rtresp := <-rtchan:
		resp, err = rtresp.resp, rtresp.err
	case <-ctx.Done():
		c.transport.CancelRequest(req)
		// wait for request to actually exit before continuing
		<-rtchan
		err = ctx.Err()
	}

	// always check for resp nil-ness to deal with possible
	// race conditions between channels above
	defer func() {
		if resp != nil {
			resp.Body.Close()
		}
	}()

	if err != nil {
		return nil, nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	return resp, body, err
}
