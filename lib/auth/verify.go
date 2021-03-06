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

package auth

import (
	"crypto/hmac"
	"net"
)

// Verify checks a user supplied username and token pair
func Verify(name, addr string, token, key, seed, expires, salt []byte) bool {
	bname := []byte(name)
	bip := []byte(net.ParseIP(extractAddress(addr)).String())

	// whiteout unstable subsecond timestamp part with "random" value
	copy(expires[9:], []byte{0xde, 0xad, 0xca, 0xfe})

	calculated := computeToken(
		bname,
		key,
		seed,
		expires,
		salt,
		bip,
	)
	return hmac.Equal(token, calculated)
}

// VerifyExtracted checks a user supplied username and token pair
func VerifyExtracted(name, addr string, token, key, seed, expires, salt []byte) bool {
	bname := []byte(name)
	bip := []byte(net.ParseIP(addr).String())

	// whiteout unstable subsecond timestamp part with "random" value
	copy(expires[9:], []byte{0xde, 0xad, 0xca, 0xfe})

	calculated := computeToken(
		bname,
		key,
		seed,
		expires,
		salt,
		bip,
	)
	return hmac.Equal(token, calculated)
}

// vim: ts=4 sw=4 sts=4 noet fenc=utf-8 ffs=unix
