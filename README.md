[![Build Status](https://travis-ci.org/jonnybazookatone/governor.svg?branch=master)](https://travis-ci.org/jonnybazookatone/governor)
[![Coverage Status](https://coveralls.io/repos/jonnybazookatone/governor/badge.svg?branch=master&service=github)](https://coveralls.io/github/jonnybazookatone/governor?branch=master)

# Governor
A simple go-lang binary that collects required file contents from Consul and writes them to disk.

## Setup
You will need to define the IP address of the Consul service. Currently, this is carried out using environment variables. The current variables are:
  - **CONSUL_HOST**: the host name of the Consul service
  - **CONSUL_PORT**: the port of the Consul service (defaults to 8500)

## Usage

```
governor -c govern.conf
```

The contents of your govern.conf file should be in a JSON format, that contains a set of key-values, where, each key is the name of the key in the Consul store, and each value is the intended output file on the disk of the machine running governor. For example:

```
{
  "NGINX_CONFIGURATION": "/etc/nginx/nginx.conf"
  "VARNISH_CONFIGURATION": /etc/services/varnish.d/varnish.conf"
}
```

In this example, it would look up the value for `NGINX_CONFIGURATION` within the Consul service, and write the received content to `/etc/nginx/nginx.conf`, and so on for each key/value item.
