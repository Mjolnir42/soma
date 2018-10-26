# DESCRIPTION

This command removes a monitoring system from SOMA. It must be unused
for this action to succeed.

# SYNOPSIS

```
soma monitoringsystem-mgmt remove ${name}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
name | string | Name of the monitoring system | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | global | | no | yes
global | monitoringsystem-mgmt | remove | yes | no

# EXAMPLES

```
soma monitoringsystem-mgmt remove ExampleMonitoring
```
