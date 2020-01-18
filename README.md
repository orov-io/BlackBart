# BlackBeart

> A template to work easily with gin as a GO & JSON microservice framework.
> *Based on a [micro framework template](https://github.com/orovium/micro)*

## Motivation

This library wraps the [gin](https://github.com/gin-gonic/gin) http router in order to provide an easy way to enable some tools for develop API servers. In the future, gin will be removed to provide a framework agnostic wrapper.

Also, a great GCLOUD integration is provided.

__Freak alert__ this library is called [BlackBart](https://en.wikipedia.org/wiki/Bartholomew_Roberts) by the famous pirate, in order to increase the pirate popularity to [fight against globar warming](https://pastafarians.org.au/pastafarianism/pirates-and-global-warming/)

## Importing

You mus import this module:

```Go
import "github.com/orov-io/BlackBart/server"
```

## Functionalities

### The router

You can:

* Run the gin router as is.
* Run the gin router on an app engine environment.

So you can use this flexibility for initializing your server based on a env variable:

```Go
service, err := server.StartDefaultService()
if err != nil {
    log.WithError(err).Panic("Can't initialize the service ...")
}

environment := os.Getenv(envKey)

if environment == local {
    err = service.Run(server.GetEnvPort(portKey))
} else {
    err = nil
    service.RunAppEngine()
}

if err != nil {
    log.WithError(err).Panic("Can't start the server")
}
```

The above code starts a service with the default [options](#Using-options-to-configure-the-service) attached

### Responses functions
  
You can user the convenience functions on [response.go](./server/response.go) to standardize your http responses.

### The logger

You can configure the [Logrus](https://github.com/sirupsen/logrus) logger as you want or use the defaults logger options.
To set a logger with a global scope on your package, you can use this snipped:

```Go
var log = server.GetLogger()
```

As soon as you import this module, a basic logrus logger is initialized, so you always will
fetch a logger with this instruction.

On default server configuration (non-local envs) logger is configured to log the provided service name and version.

### Firebase integration

This service can use firebase as a authentication service out the box. You can pass to it your firebase admin credentials (read only or readWrite credential) storing it on a bucket or simply with a file. Put either the bucket or the file in the firebaseOptions struct before initialize the service.

If a credential is provided, service wil try to give the firebase auth client and leave it ready to use on _server.GetAuthClient()_ and _server.GetFirebaseApp()_ functions.

### Database initialization

This service has a module to stablish connections to postgres databases given the connection parameters. It also leave you inject an initialized sql database. See the _WithInjectedDB()_ function on __[options.go](./server/options.go)__

### Pre-defined errors

You can find a set of ready to use gin http responses on __[response.go](./server/response.go)__.

This includes quick ways to send 404 and 500 errors, and 200, 201 and 204 responses.

## Configuring the service

You can configure resulting service in two complementary ways:

* Via env variables.
* Providing needed options on code.

Please, note that providing options on code will overwrite the env variables way.

### Using env variables

Each module expect to find some env variables in order to execute their default initialize process. If this variables are no set or not found, then the module will not be initialized.

Below, there are a list with env variables that each module is waiting for:

#### Database module
  
* DATABASE_MIGRATIONS_DIR: The folder that store migration files. Is relative from calling file.
* DATABASE_HOST
* DATABASE_PASSWORD
* DATABASE_USER
* DATABASE_SSL_MODE
* SERVICE_DATABASE_NAME

#### Redis module

* REDIS_ADDRESS
* REDIS_PASSWORD

#### Logger module

  Always a logger is provided by the server, based on your environment declaration. If no environment is provided, the service will be initialized with a JSON print format. To declare your environment, you can use ENV env variable

* ENV: one of [local, dev, pre, prod]
  * ENV == local: Log with pretty formatter and logrus.TaceLevel as the lowest level to log.
  * ENV == dev, pre: Log with JSON formatter and logrus.DebugLevel as the lowest level to log.
  * ENV == prod: Log with JSON formatter and logrus.InfoLevel as the lowest level to log.
  * default: Log with JSON formatter and logrus.DebugLevel as the lowest level to log.

#### Response module

  As logger, response module is always available. Anyway, if you are in production (ENV == prod) errors are hidden and only a trace UUID is written to the JSON response. A log with the Fatal level will be written with an error and trace_id attributes.

#### Firebase module

* FIREBASE_BUCKET and FIREBASE_BUCKET_FILE_NAME: Use this two variables to provide to the service a GCLOUD bucket to search for the firebase json credential. By default, service assumes that file is called firebase.json. If you are not in an app engine context in the same gcloud project, you will need to set the GOOGLE_APPLICATION_CREDENTIALS variable to the json file that stores your google IAM with permissions to read in the provided bucket.
* FIREBASE_CONFIG_PATH: You also can provide the firebase json credential as a local file. Use this variable to says to the service where the file is allocated.

#### Profiler

By default, _BlackBart_ try to start the GCLOUD profiler, but it not panic if it was not possible.

#### Local key/value database

_BlackBart_ provides a local key/value database. It uses [badger](https://github.com/dgraph-io/badger). To initialize it, you have 2 ways:

* Set ENABLE_BADGER env variable to true. _BlackBart_ will provide you with a in-memory badger database with badger default options.
* Set a new InternalDBOptions Option with your badger options. See the InternalDBOptions in _[server options](./server/options.go)_.

If badger is enabled, you can get it with the `server.GetInternalDB()`function.

### Using options to configure the service

The _server.StartDefaultService()_ function starts for you a fresh service based on your env variables. It uses internally the _WithDefaultOptions()_ method of the _Options_ struct.

Otherwise, you can start a new service with your own options (or call the _WithDefaultOptions()_ method and after updating the options that you need).

As an example:

```Go
// initialize your custom sql db...

options := server.NewOptions().WithDefaultOptions()

dbOptions := server.NewDBOptions().WithInjectedDB(db)

server.Init(options)

// After that, you can run your service with your injected db

```
