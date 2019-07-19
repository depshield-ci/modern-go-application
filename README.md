# Modern Go Application

![depshield](https://ci.dev.depshield.sonatype.org/badges/depshield-ci/modern-go-application/depshield.svg)

**Go application boilerplate and example applying modern practices**

This repository tries to collect the best practices of application development using Go language.
In addition to the language specific details, it also implements various language independent practices.

Some of the areas Modern Go Application touches:

- architecture
- package structure
- building the application
- testing
- configuration
- running the application (eg. in Docker)
- developer environment/experience
- instrumentation

To help adopting these practices, this repository also serves as a boilerplate for new applications.


## Features

- configuration (using [spf13/viper](https://github.com/spf13/viper))
- logging (using [goph/logur](https://github.com/goph/logur) and [sirupsen/logrus](https://github.com/sirupsen/logrus))
- error handling (using [goph/emperror](https://github.com/goph/emperror))
- metrics and tracing using [Prometheus](https://prometheus.io/) and [Jaeger](https://www.jaegertracing.io/) (via [OpenCensus](https://opencensus.io/))
- health checks (using [InVisionApp/go-health](https://github.com/InVisionApp/go-health))
- graceful restart (using [cloudflare/tableflip](https://github.com/cloudflare/tableflip)) and shutdown
- support for multiple server/daemon instances (using [oklog/run](https://github.com/oklog/run))
- messaging (using [ThreeDotsLabs/watermill](https://github.com/ThreeDotsLabs/watermill))
- MySQL database connection (using [go-sql-driver/mysql](https://github.com/go-sql-driver/mysql))
- Redis connection (using [gomodule/redigo](https://github.com/gomodule/redigo))


## First steps

To create a new application from the boilerplate clone this repository (if you haven't done already) into your GOPATH
then execute the following:

```bash
chmod +x init.sh && ./init.sh
? Package name (github.com/sagikazarmark/modern-go-application)
? Project name (modern-go-application)
? Binary name (modern-go-application)
? Service name (modern-go-application)
? Friendly service name (Modern Go Application)
? Update README (Y/n)
? Remove init script (y/N) y
```

It updates every import path and name in the repository to your project's values.
**Review** and commit the changes.


### Load generation

To test or demonstrate the application it comes with a simple load generation tool.
You can use it to test the example endpoints and generate some load (for example in order to fill dashboards with data).

Follow the instructions in [etc/loadgen](etc/loadgen).


## Inspiration

See [INSPIRATION.md](INSPIRATION.md) for links to articles, projects, code examples that somehow inspired
me while working on this project.


## License

The MIT License (MIT). Please see [License File](LICENSE) for more information.
