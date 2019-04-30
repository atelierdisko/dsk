# Building your own Frontend 

The following sections describe everything you need to know when building your own frontend
and bundle it with the `dsk` binary. By default DSK uses a built-in minimal frontend. The default frontend
is a good starting point when developing your own. It can be found in the `frontend` directory of
this project.

## Available API Endpoints

Please read our [API document](Architecture/API) for in depth information.

## Designing the URL Schema

Your frontend and its subdirectories will be mounted directly at the root path
`/`. Requests to anything under `/api` are routed to the backend, anything else
is routed into your application in `index.html`. 

Relative asset source paths inside documents will be made absolute to allow you
displaying a document's content wherever you like to.

## Starting DSK with the custom frontend

You "tell" DSK, to use your frontend by invoking it with the `frontend` flag,
providing an absolute path or a path relative to the current working directory,
that contains the frontend.

```
./dsk -frontend=/my/frontend example
```

## Replacing the built-in frontend

To persistently _bake_ your frontend into DSK, effectively replacing the
built-in one, install the the development tools as described in the _Development_
section first. After doing so, you create your custom DSK build by running the
following command.

```
$ FRONTEND=/my/frontend make dist
```