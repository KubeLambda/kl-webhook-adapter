## About

API Server application provides very basic functionality to work with Address Book contacts.
You can create new contact with phone numbers, update existing contact, fetch contact by id,
list all contacts, and delete existing contact.


## How to build application

Open the directory with newly created project and run:

```sh
go build -o bin/webhook-adapter
```
it will result in building executable file "webhook-adapter" (feel free to name it differently).


## How to run application

**IMPORTANT:** depending on your configuration, e.g. if you have added database support
etc, starting of your code may fail because you need to complete configuration settings (e.g.
your database URL and credentials). So in this case keep reading README past this section.

To perform an initial launch an application, run this in shell:

```sh
make run --deployment=local
```

We launch API server by specifying `run` command. `--deployment=local` tells our code to
perform a local deployment. Local deployment settings can contain options such as URL of
your local database or other local specific environment. `--deployment` flag loads
file `local.yaml` from `config/` directory that resides in the same directory where
your executable file is. You can create a copy of this file and name it,
for example `prod.yaml` where you can add production-specific settings, then running

```sh
make run --deployment=prod
```

will load this production settings for your API server.


## How to override configuration values

Sometimes editing configuration file to add values is not the best strategy. As an example,
if you have database settings in your `prod.yaml` file, having URL of database specified
there is not a bad idea, but storing a password there - is not good. The better approach
would be to pass sensitive settings via environment variables. And because we use
viper library to load yaml configuration file, it allows us to override values specified
in it with something different. The typical syntax of environment variable is:
`[EnvPrefix]_[YamlConfigKey] = value`. `EnvPrefix` is something you have previously entered when
generated project with GoQuick.

**internal/core/app/configload.go:**
```
const envVarPrefix = "WEBHOOK_ADAPTER"
```

Let's give it a try. Let's say you want to change a listening port for your API server.
If you open `local.yaml` you can find something like this:

```yaml
server:
  port: 8080
```

`server/port` translates to SERVER_PORT and combined with environment prefix WEBHOOK_ADAPTER
you can override it as:

```shell
export WEBHOOK_ADAPTER_SERVER_PORT=9090
```
```sh
./WEBHOOK_ADAPTER run --deployment=local
```

now API server code will be listing port 9090 instead of 8080.

## Access REST API

Generated application uses REST protocol to store and fetch address book records.
Once you have the application launched, you can perform HTTP calls to test REST APIs
exposed by API server.

Please note that each HTTP response contains **X-Request-Id** header with value that
is displayed with application logs (as **requestId** field). It helps you to troubleshoot
application, because logger provided with generated code prints request id with
every log line.

### Examples of REST requests

#### Get service version 

Request:

```sh
curl --location 'http://localhost:8080/api/version'
```
Response (could be slightly different):

```json
{
  "service": "rest-net/http",
  "version": "0.1.0",
  "build": "1"
}
```

#### Send request

Request:

```sh
curl --data '{"name":"bob"}' --header 'Content-Type: application/json' http://0.0.0.0:3001/api/adapter
'
```

Response:

```json
```

### Logging

Each HTTP request returns `X-Request-Id` header as part of response. This `X-Request-Id`
is always unique, unless you specify it explicitly as part of request. What makes it useful
is that each application log line contains `{requestId="...."}` tag, and it matches
`X-Request-Id` value. It makes debugging code much easier because you can filter logs
scoped to specific request.
