// Copyright 2020 The go-eth Authors
// This file is part of the go-eth library.
//
// The go-eth library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-eth library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-eth library. If not, see <http://www.gnu.org/licenses/>.

package node

// Lifecycle encompasses the behavior of services that can be started and stopped
// on the node. Lifecycle management is delegated to the node, but it is the
// responsibility of the server-specific package to configure and register the
// server on the node using the `RegisterLifecycle` method.
type Lifecycle interface {
	// Start is called after all services have been constructed and the networking
	// layer was also initialized to spawn any goroutines required by the server.
	Start() error

	// Stop terminates all goroutines belonging to the server, blocking until they
	// are all terminated.
	Stop() error
}
