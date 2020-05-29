# DnsInverse
This is a tool to do inverse dns lookups based on dns queries that happened on your dns server.
Note that this isn't reverse dns as pointer records, especially with cloud, don't tell you much nowadays.
The data is collected using [dnstap](https://dnstap.info/).
## Dns configuration
I only tested this with [coredns](https://coredns.io/), but any dns server that supports sending dnstap over tcp should work. I used the following configuration:
```
.:53 {
	forward . 1.1.1.1
	dnstap tcp://127.0.0.1:1058 full
}
```
## Api
use /?ip=y.x.y.z to get a json array of domains connected to that IP.
## Disclaimer
The tool doesn't do any pruning and will fill up your memory, but luckily a lot slower than google chrome.
I only need v4 for now, support for v6 will probably be added later.