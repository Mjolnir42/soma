# DESCRIPTION

This command looks up the ID of a monitoring system defined in SOMA by
its name.

# SYNOPSIS

```
soma monitoringsystem-mgmt search ${name}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
name | string | Name of the monitoring system | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions. Permissions in category
`monitoring` must be granted on the specific monitoring system being
shown.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | global | no | yes
global | monitoringsystem-mgmt | search-all | no | yes
system | monitoring | no | yes
monitoring | monitoringsystem | search | yes | no

# EXAMPLES

```
soma monitoringsystem-mgmt search ExampleMonitoring
```
