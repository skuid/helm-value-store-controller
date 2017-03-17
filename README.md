# helm-value-store-controller

[![Build Status](https://travis-ci.org/skuid/helm-value-store-controller.svg?branch=master)](https://travis-ci.org/skuid/helm-value-store-controller)
[![https://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square](https://img.shields.io/badge/godoc-reference-5272B4.svg?style=flat-square)](http://godoc.org/github.com/skuid/helm-value-store-controller/)


Helm-value-store-controller syncs with a value-store to ensure certain helm
releases are installed.

See [helm-value-store](https://github.com/skuid/helm-value-store) for more
details.

## Example

TODO

### AWS Prerequisite

You must have the ability to create, read, and write to a DynamoDB table.

Set the proper access key environment variables, or use the
`$HOME/.aws/{config/credentials}` and set the appropriate
`AWS_DEFAULT_PROFILE` environment variable.

## Usage

```
$ helm-value-store-controller -h
Usage of ./helm-value-store-controller:
  -b, --blacklist stringArray   A list of release names to not update or install
  -i, --interval string         The sync interval to check the value store (default "300s")
  -l, --labels string           The labels to search the value store for.
  -p, --port int                The port to listen on (default 3000)
  -t, --table string            The DynamoDB table to read from (default "helm-charts")
```

## License

MIT License (see [LICENSE](/LICENSE))
