// Copyright 2014 Brady Catherman
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//   http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// This module was written with the sole intention of making writing unit tests
// in Go use far less boiler plate checks. As such most of the boiler plate
// functionality is reduced into helpful functions that require little to no
// additional checking or setup.
//
// Full Stack Traces
//
// Sometimes it is helpful to write a wrapper function to do a bunch of the
// boiler plate testing for you, however this causes the testing library
// to output a useless line number when reporting an error. As such this
// library implements a Fatal/Fatalf that will output the full stack
// trace where the error was encountered.
//
// Easier Error Checking
//
// Often error checking around initialization or setup functionality can become
// extensive which means that often error returns will get ignored since
// developers don't want to keep typing up boiler plate error checking. Since
// the errors rarely ever happen the pain isn't really uncovered until later.
// This library makes testing for expected errors or success super trivial.
// See the Examples section for simple examples explaining how to do this.
//
// Temporary Files are Not
//
// When writing temporary files people often assume that something will
// clean them up for them but in Go this is far from the truth. Often an
// external process cal clean them after some period but it is not rare to find
// /tmp completely full of random temporary cruft that explodes when running
// a large suite of unit tests.
//
// This library implements functionality to ensure that the temporary files
// that get created will be removed. It does this via a child process launched
// that will wait for the parent to die before removing all temporary files
// created. Using a new processes helps ensure that a panic or hard crash
// won't prevent files from being cleaned.
//
// Simple Equality
//
// It can be quite tedious to validate that two objects are equal. Especially
// if they are maps of complex data types. Often this leads to dozens of
// helper functions or code explosion. The Equal() and NotEqual() functions
// are designed to eliminate this completely. Using reflection they walk
// through objects verifying that they are in fact equal, regardless of types.
// Unlike most other implementations this will also follow pointers to ensure
// that referenced data is also equal.
package testlib
