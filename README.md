# netconf

This library is a simple NETCONF client based on [RFC6241](http://tools.ietf.org/html/rfc6241) and [RFC6242](http://tools.ietf.org/html/rfc6242) (although not fully compliant yet).

# Features
* Support for SSH transport using go.crypto/ssh. (Other transports are planned).
* Built in RPC support (in progress).
* Support for custom RPCs.
* Independent of XML library.  Free to choose encoding/xml or another third party library to parse the results.
* Client struct in the vain of UpGuard's sftp fork (might need to refactor this slightly)

# Install
    go get github.com/cloudhousetech/go-netconf/netconf

# Documentation

You can view full API documentation at GoDoc:

http://godoc.org/github.com/Juniper/go-netconf/netconf

# License

(BSD 2)

Copyright © 2013, Juniper Networks

All rights reserved.

Redistribution and use in source and binary forms, with or without modification, are permitted provided that the following conditions are met:

(1) Redistributions of source code must retain the above copyright notice, this list of conditions and the following disclaimer.

(2) Redistributions in binary form must reproduce the above copyright notice, this list of conditions and the following disclaimer in the documentation and/or other materials provided with the distribution.

THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS “AS IS” AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

The views and conclusions contained in the software and documentation are those of the authors and should not be interpreted as representing official policies, either expressed or implied, of Juniper Networks.

# Authors and Contributors

* [Brandon Bennett](https://github.com/nemith), Facebook
* [Charl Matthee](https://github.com/charl)
* [Nick Wood](https://github.com/nickswift) UpGuard