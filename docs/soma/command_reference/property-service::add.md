# DESCRIPTION

This command is used to add per-team service properties to SOMA, which
are key/array-of-value dictionary objects.

Possible attribute values that the service can contain depends on the
SOMA attributes command. Attributes with a cardinality of `multi` can be
specified more than once.

# SYNOPSIS

```
soma property-mgmt service add ${service} team ${team} [${attribute} ${attrValue}, ...]
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
service | string | Name of the service | | no
team | string | Name of the team | | no
attribute | string | Name of the attribute | | yes
attrValue | string | Value of the attribute | | yes

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient, all system or all required permissions.

Category | Section | Action | Required | System | Sufficient
 ------- | ------- | ------ | -------- | ------ | ----------
omnipotence | | | no | no | yes
system | global | | no | yes | no
system | team | | no | yes | no
global | property-mgmt | add | yes | no | no
team | property-service | add | yes | no | no

# EXAMPLES

```
soma property-mgmt service add PowerDNS \
  team ExampleTeam \
  transport_protocol udp \
  transport_protocol tcp \
  application_protocol dns \
  port 53
```
