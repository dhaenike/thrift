/*
 * Licensed to the Apache Software Foundation (ASF) under one
 * or more contributor license agreements. See the NOTICE file
 * distributed with this work for additional information
 * regarding copyright ownership. The ASF licenses this file
 * to you under the Apache License, Version 2.0 (the
 * "License"); you may not use this file except in compliance
 * with the License. You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing,
 * software distributed under the License is distributed on an
 * "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
 * KIND, either express or implied. See the License for the
 * specific language governing permissions and limitations
 * under the License.
 */

package thrift

import (
	"testing"

	"golang.org/x/net/context"
)

func TestHttpClient(t *testing.T) {
	l, addr := HttpClientSetupForTest(t)
	if l != nil {
		defer l.Close()
	}
	trans, err := NewTHttpPostClient("http://" + addr.String())
	if err != nil {
		l.Close()
		t.Fatalf("Unable to connect to %s: %s", addr.String(), err)
	}
	TransportTest(t, trans, trans)
}

func TestHttpClientHeaders(t *testing.T) {
	l, addr := HttpClientSetupForTest(t)
	if l != nil {
		defer l.Close()
	}
	trans, err := NewTHttpPostClient("http://" + addr.String())
	if err != nil {
		l.Close()
		t.Fatalf("Unable to connect to %s: %s", addr.String(), err)
	}
	TransportHeaderTest(t, trans, trans)
}

func HttpCancelTest(t *testing.T, writeTrans TTransport, readTrans TTransport, canceler context.CancelFunc) {
	bdata := []byte{1, 2, 3, 4, 5}

	if !writeTrans.IsOpen() {
		t.Fatalf("Transport %T not open: %s", writeTrans, writeTrans)
	}
	if !readTrans.IsOpen() {
		t.Fatalf("Transport %T not open: %s", readTrans, readTrans)
	}

	// this write should succeed
	_, err := writeTrans.Write(bdata)
	if err != nil {
		t.Fatalf("Transport %T cannot write binary data of length %d: %s", writeTrans, len(bdata), err)
	}

	// and this flush also
	err = writeTrans.Flush()
	if err != nil {
		t.Fatalf("Transport %T cannot flush write binary data 2: %s", writeTrans, err)
	}

	// now canceling transport so flush should fail after
	canceler()

	// this write should not fail as it doesn't use transport yet
	_, err = writeTrans.Write(bdata)
	if err != nil {
		t.Fatalf("Transport %T cannot write binary data after canceling %s", writeTrans, err)
	}

	// Flush should fail as we have canceled the operation
	err = writeTrans.Flush()
	if err == nil {
		t.Fatalf("Flush operation to transport %T could not be canceled", writeTrans)
	}
}

func TestHttpClientWithCtx(t *testing.T) {
	l, addr := HttpClientSetupForTest(t)
	if l != nil {
		defer l.Close()
	}
	ctx := context.Background()
	childCtx, canceler := context.WithCancel(ctx)

	trans, err := NewTHttpPostClientWithCtx("http://"+addr.String(), childCtx)
	if err != nil {
		l.Close()
		t.Fatalf("Unable to connect to %s: %s", addr.String(), err)
	}
	HttpCancelTest(t, trans, trans, canceler)
}
