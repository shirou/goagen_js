/*
Package goagen_js provides a goa generator for a javascript client module.
The module exposes functions for calling the API actions. It relies on the
fetch API (or using fetch polyfill like https://github.com/github/fetch) to perform the actual HTTP requests.

Output js is es2015 format. So you may need to build using babel https://babeljs.io/.

The generator also produces an validator of API parameters.  Since this validator is separeted on each fields, you can use as is, for example form validator.

The controller simply serves all the files under the "js" directory.
*/
package goagen_js
