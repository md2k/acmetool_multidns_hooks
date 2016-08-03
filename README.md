## acmetool_multidns_hooks
DNS Hook for https://github.com/hlandau/acme with support of Multi-DNS providers
Project partly based on work of https://github.com/xenolf/lego (DNS Providers workflow and DNS tools)

### Configuration
Hook looking for its configuration file in acmetool `ACME_STATE_DIR`/conf/ location
Name of config file equal to name of executable with `.yaml` extension (Example: dns.hook.yaml)
Configuration example below:
```yaml
---
providers:
  cloudflare:
    accounts:
      myAccount1:
        authdata:
          email: myAccount1@example.com
          api_key: XXXXXXXXXXXXXXXXXXXXXXXXXXXX
        domains:
          - example.com
          - example2.com
      myAccount2:
        authdata:
          email: myAccount2@example.com
          api_key: XXXXXXXXXXXXXXXXXXXXXXXXXXXX
        domains:
          - domain2.com
          - domain3.com
  route53:
    accounts:
      myAccount1:
        authdata:
          key_id: YYYYYYYYYYYYYY
          access_key: XXXXXXXXXXXXXXXXXXXXXXXXXXXX
        domains:
          - moreexamle.com
      myAccount2:
        authdata:
          key_id: YYYYYYYYYYYYYY
          access_key: XXXXXXXXXXXXXXXXXXXXXXXXXXXX
        domains:
          - domain10.com

```

Where:
***myAccount1/2/n...***  - this is used more for separation of different accounts if you manage several AWS/Cloudflare accounts.
***domains*** - list of domains which can be used with `dns.hook` to Update DNS records with ACME DNS-01 Challenge.

When this hook used with Acmetool, any certificate requested for single domain or multiple domains will be checked by `dns.hook` with Yaml configuration
and required DNS-01 challenges will be updated on correct DNS providers depends to which provider-account this domain belong to.


### Shell wrapper
`acmecli` - small shell wrapper to simplify acmetool usage with custom ACME_STATE_DIR directory.


### Test Hook
`dns.hook` can be easily tested by execution of it with next parameters:

```
ACME_STATE_DIR="/PATH/TO/ACME/DIR" ./dns.hook challenge-dns-start test.example.com test.example.com-SOMEDUMMYSTRING OTHER_DUMMY_STRING
ACME_STATE_DIR="/PATH/TO/ACME/DIR" ./dns.hook challenge-dns-stop test.example.com test.example.com-SOMEDUMMYSTRING OTHER_DUMMY_STRING
```

More information about Hooks integration with Acmetool is [here](https://github.com/hlandau/acme/blob/master/_doc/SCHEMA.md#hooks)

`MIT License.`
