# DESCRIPTION

This command removes a per-team service property from SOMA.

# SYNOPSIS

```
soma property-mgmt service remove ${service} team ${team}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
service | string | Name of the service | | no
team | string | Name of the team | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient, all system or all required permissions.

Category | Section | Action | Required | System | Sufficient
 ------- | ------- | ------ | -------- | ------ | ----------
omnipotence | | | no | no | yes
system | global | | no | yes | no
system | team | | no | yes | no
global | property-mgmt | remove | yes | no | no
team | property-service | remove | yes | no | no

# EXAMPLES

```
soma property-mgmt service remove PowerDNS team ExampleTeam
```
