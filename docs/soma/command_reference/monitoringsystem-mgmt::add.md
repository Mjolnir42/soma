# DESCRIPTION

This command adds a monitoring system to SOMA. It's name must not
contain any of the following characters: `/`, `:`, `.`. It has a mode, a
primary responsible contact user, an owning team and optionally a
callback URI to which deployment requests will be signaled.

# SYNOPSIS

```
soma monitoringsystem-mgmt add ${name} mode ${mode} contact ${user} team ${team} [ callback ${callback} ]
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
name | string | Name of the monitoring system | | no
mode | string | Mode of the monitoring system | | no
contact | string | Name of the primary contact user of this monitoring system | | no
team | string | Name of the team owning this monitoring system | | no
callback | string | Callback request URI for this monitoring system | | yes

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | global | | no | yes
global | monitoringsystem-mgmt | add | yes | no

# EXAMPLES

```
soma monitoringsystem-mgmt add ExampleMonitoring \
  mode private \
  contact root \
  team wheel \
  callback 'https://[::1]:666/poke'
```
