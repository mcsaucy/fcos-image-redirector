# Fedora CoreOS Image Redirector

Need a way to find the most up-to-date FCOS artifact? Don't use this. Use [the
official page](https://getfedora.org/coreos?stream=stable) instead.
[FCOS #625](https://github.com/coreos/fedora-coreos-tracker/issues/625) has
details explaining why you should use the official approach.

Trying to find a stable URL you can hit to pull down the most up-to-date image
(for example via iPXE)? This may be what you want.

## How it works

The Fedora CoreOS project exposes
[streams](https://builds.coreos.fedoraproject.org/streams/stable.json) files
to track the most up-to-date versions of things for a given stream (stable,
testing, or next). When we get a request, we pull down the streams JSON blob,
parse it, cache it, and then redirect you to the FCOS URL.

At the end of the day, this just parses a JSON blob, uses the request path to
walk the object tree and then redirects to what it finds.

## What's supported

We should be able to redirect to any location within the
"architectures/\*/artifacts" structure of the [streams JSON
blob](https://builds.coreos.fedoraproject.org/streams/stable.json).

### Basic lookups
Example: `/artifacts/x86_64/metal/pxe/kernel`

This just redirects to the `location` for that resource in the streams JSON.

### Peeking
Example: `/artifacts/x86_64/metal/pxe/kernel?peek`

This doesn't redirect. Instead of redirecting to a URL, we write it to the
response body.

### Signature fetching
Example: `/artifacts/x86_64/metal/pxe/kernel?sig`

Redirects to the `.sig` file for the resource. This can also be used with
`?peek`.

### Digest fetching
Example: `/artifacts/x86_64/metal/pxe/kernel?sha256`

Writes the SHA256 digest for the resource (as stored in the streams blob) in
the response body. There's no redirection here.

## How do I know this isn't tampering with images?
Because it's not serving them. The requester is merely being redirected to the
a resource hosted by the Fedora CoreOS project.

Now, we could totally redirect you to some other source that hosts compromised
or non-FCOS images and pass them off as legitimate, but if you're worried about
that you should probably not be using this (or at least audit all 200 lines of
code and host it yourself).

At the end of the day, the best thing you can do is use the `coreos-installer`
to grab your images, as it's the Official way and it performs signature
verification.
