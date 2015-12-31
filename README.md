# Goshify

_hashify.me without JavaScript and optional hosting_

## What

Goshify takes Base 64-encoded Markdown in a URL and turns it into HTML. It's similar to [hashify.me](http://hashify.me), but without the editor and without the requirement for JavaScript on the client side.

Additionally, it can save your stuff for easier recalling, without giant URLs involved.

## Why

Hashify can be used to quickly put up a short document written in Markdown. As an added benefit, the document is self-contained in its URL. However, for people to see that page, their browsers must have JavaScript enabled, and they must download quite a bit of it, plus some assets for the Hashify editor.

Goshify removes the JavaScript requirement by doing the conversion to HTML server-side. All the browser needs to do is render some HTML and CSS. As for the editor, you can keep using the Hashify one, which is pretty good, or any other of your liking.

## How
In the following examples,

    QSBzYW1wbGUgZG9jdW1lbnQ=

is the document encoded in Base-64 as Hashify does it.
    
To **present a document**, as Hashify does:

    http://goshify.tny.im/d/QSBzYW1wbGUgZG9jdW1lbnQ=

To **store this document** for easier recalling:

    http://goshify.tny.im/s/QSBzYW1wbGUgZG9jdW1lbnQ=
You'll be given an ID to recall this document later.

You can also make a POST request with the raw Markdown (Base 64 encoded or not) to

    http://goshify.tny.im/s

To **recall a document**:

    http://goshify.tny.im/l/352c706a-bffc-4da0-a75e-e5b429a9c5ae
    
Where 352c706a-bffc-4da0-a75e-e5b429a9c5ae is the ID that was returned in the storing step. If a document with the specified ID doesn't exist, a 404 is returned.

To recall the Base 64-encoded document, as it was saved (for example, for editing it in Hashify):

    http://goshify.tny.im/r/352c706a-bffc-4da0-a75e-e5b429a9c5ae

_**But this is not a REST API!?**_ Nobody said it was :) It is designed to be easy to use by a human fiddling with the address bar of any web browser.

## How, how

Goshify is just a simple, 150 line Go program... [Here's the GitHub](https://github.com/tnyim/goshify). It uses:

  - https://github.com/boltdb/bolt
  - https://github.com/gorilla/mux
  - https://github.com/russross/blackfriday
  
## Who
![Segvault](http://s.lowendshare.com/11/1451510832.605.segvault-24.png), part of [tny. internet media](http://i.tny.im).
