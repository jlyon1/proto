# Proto

Schema Based Development Without the SaaS. Like [buf.build](https://buf.build) but without charging to host schemas. Instead `proto` uses version control as a backing in a sane way.

## Dependencies

* protoc
* protoc plugins (eg protoc-gen-go)

## Features

These are the current and planned features

### Version 0.0.1 Pre-Release

* `proto mod init (dns.name)/org/repo` - Creates a proto repo in a directory with a given remote
* `proto mod tidy` - clean a repo
* `proto get` - When run within a proto repo will fetch all dependencies
* `proto compile` - Compile proto repo as specified in  `proto.yaml`

### Future Planned Features

* `proto lint` - Lint a proto repo
* `proto push` - push a repo
* `proto curl [url]` - curl using the protos defined in a proto repo
* `proto validate [hash1] [hash2]` - Validate a repo is backwards compatible with another commit
* build in `protoc` binaries and extensions so they don't need to be manually installed