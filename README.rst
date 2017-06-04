goagen js
==========

This package provides a `goa <https://goa.design/>`_ generator for a modern javascript client module.
This module enabels calling the API actions and varidate these parameters using goa design definition.

Generated JS code are using

- ES 2015
- fetch API (using `whatwg-fetch` library)
- Promise API

You need to *compile* using babel, webpack, or other tool. This compilation is a standard way on the modern JavaScript era.


Status
------------------

Very alpha. you can not use in the production.

- handling JS and goa kind difference.
- array is not work
- required is not work
- re-thinking about `validate()` calling interface. Is this easy to use or not?
- support Flowtype and TypeScript


How to generate JS from your design
---------------------------------------------


At first, you have to do `go get`

::

  % go get github.com/shirou/goagen_js

Then, you can `goagen gen` with your design.

::

  % goagen gen --pkg-path=github.com/shirou/goagen_js -d github.com/some/your/great/design


These two files are generated under `js` directory.

- js/api_request.js
- js/api_validator.js


How to use API
------------------------------------

Invoking API
````````````````````

In the `api_request.js`, each API function is generated and you can use just normal invokation with Promise.

::

  import * as api from "./api_request.js";

  // This API is invoked with query parameter.
  // Query parameter should be passed as object.
  api.FooBar(payload).then((response) => {
    if (response.status !== 200){
      console.log("error", response.status);
    }
    return response.json();
  });


  // This API is invoked with Path Param and POST payload
  // To use this API, passing PathParam and JSON payload as object.
  // This API request to /v2/:fooID/:barType
  api.FooParam(fooID, barType, payload).then((response) => {
    if (response.status !== 200){
      console.log("error", response.status);
    }
    return response.json();
  });


Validation
````````````````

In the `api_validation.js`, there are many `rules` for each APIs and a `validate` function.

::

  import * as v from "./api_validator.js";

  // validate payload
  console.log(v.validate(v.FooBarGet.payload, {
    id: 1,
    too_large_int: 99999,
  }));

  // validate Path Params
  console.log(v.validate(v.FooParamPost.fooID, 99999));

This validate function is used in `api_request.js`. Also you can use to validate before request, for example Form Validation.


Type
``````````




LICENSE
---------------------

MIT License
