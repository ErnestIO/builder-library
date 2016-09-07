# Builder Library

Builder-Library is a library to create ernest builders, it works as a bridge between the workflow manager and the corresponding adapters, and receiving a bunch of components to process will generate as much 'component.create' messages as components we have.

Once all of them are fully processed it will emit back a message 'component.process.done' as confirmation or 'component.process.error' in case something is broken

# Using it

```
	s.Setup(os.Getenv("NATS_URI"))
	s.ProcessRequest("components.process", "component.process")
	s.ProcessSuccessResponse("component.process.done", "component.process", "components.process.done")
	s.ProcessFailedResponse("component.process.error", "components.process.error")
```


## Build status

* master: [![CircleCI](https://circleci.com/gh/ernestio/builder-library/tree/master.svg?style=svg)](https://circleci.com/gh/ernestio/builder-library/tree/master)

## Installation

```
make deps
make install
```

## Running Tests

```
make test
```

## Contributing

Please read through our
[contributing guidelines](CONTRIBUTING.md).
Included are directions for opening issues, coding standards, and notes on
development.

Moreover, if your pull request contains patches or features, you must include
relevant unit tests.

## Versioning

For transparency into our release cycle and in striving to maintain backward
compatibility, this project is maintained under [the Semantic Versioning guidelines](http://semver.org/).

## Copyright and License

Code and documentation copyright since 2015 r3labs.io authors.

Code released under
[the Mozilla Public License Version 2.0](LICENSE).

