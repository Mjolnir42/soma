/*-
Copyright (c) 2016, Jörg Pernfuß <code.jpe@gmail.com>
All rights reserved.

Redistribution and use in source and binary forms, with or without
modification, are permitted provided that the following conditions are met:

* Redistributions of source code must retain the above copyright notice, this
  list of conditions and the following disclaimer.

* Redistributions in binary form must reproduce the above copyright notice,
  this list of conditions and the following disclaimer in the documentation
  and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

package somaauth

import (
	"crypto/hmac"
	"errors"

	"github.com/dchest/blake2b"
)

// Format string for millisecond precision RFC3339 timestamps
const rfc3339Milli string = "2006-01-02T15:04:05.000Z07:00"

// TokenExpirySeconds can be set to regulate the lifetime of newly
// issued authentication tokens. The default value is 43200, or 12
// hours.
var TokenExpirySeconds uint64 = 43200

// ErrAuth indicates an authentication failure
var ErrAuth = errors.New("Authentication failed")

// ErrInput is returned if tokens can not be generated due to
// misconfiguration
var ErrInput = errors.New("Invalid input")

// computeToken does what it says on the label and computes the HMAC
// token. As input it takes the username, hmac key, token seed, token
// expiry time, token salt and client ip address
func computeToken(name, key, seed, expires, salt, ip []byte) []byte {
	mac := hmac.New(blake2b.New256, key)
	mac.Write(seed)
	mac.Write(name)
	mac.Write(ip)
	mac.Write(expires)
	mac.Write(salt)
	return mac.Sum(nil)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
