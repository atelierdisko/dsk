# DSK
## Design System Kit

## Abstract

tbd

## Quickstart

[Download the binary](https://github.com/atelierdisko/hoi/releases) from GitHub and make it exectuable.

```
curl -L https://github.com/atelierdisko/dsk/releases/download/v0.5.0/dsk-darwin-amd64 -o dsk
chmod +x dsk
```

Now run the `dsk` command pointing it to the directory that contains your desing system definitions.
```
./dsk example
```

Finally go to [open the Web Application in your browser](http://localhost:8080).

## Use it with [Create React App](https://github.com/facebookincubator/create-react-app) 

After creating the react app, Modify the entry point in `src/index.js`. When
loaded in DSK it must provide a `renderComponent()` function. This global
function is the only little glue code required. 

```javascript
window.renderComponent = function(root, name, props) {
  return import('./' + name).then(c => {
	ReactDOM.render(React.createElement(c.default, props), root);
  });
};
```

Using the `INSIDE_DSK` constant, it's possible to conditionally activate certain
default behavior while developing and deactivate it when embedded in DSK.

```javascript
if (typeof INSIDE_DSK === "undefined") {
  ReactDOM.render(<App />, document.getElementById('root'));
  registerServiceWorker();
}
```

When ready, create the bundle and copy it into our DSK tree.

```
$ npm run build 
$ cp build/static/*/main*.{css,js} example/DataEntry/Password/
```

## Development

[Go](https://golang.org/) version 1.9 or later is needed for developing and
testing the application. The `make dev` command assumes your test design system
definitions are below `_test`.

```
$ go get github.com/atelierdisko/dsk
$ cd $GOPATH/github.com/atelierdisko/dsk
$ make dev
```

## Requirements

A filesystem and a browser that supports ES6.

## Copyright & License

DSK is Copyright (c) 2017 Atelier Disko if not otherwise
stated. Use of the source code is governed by a BSD-style
license that can be found in the LICENSE file.
