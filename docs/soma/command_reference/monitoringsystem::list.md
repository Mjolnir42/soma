# DESCRIPTION

This command lists all the monitoring systems defined in SOMA for which
the user has the appropriate permissions.

# SYNOPSIS

```
soma monitoringsystem-mgmt list
```

# ARGUMENT TYPES

This command takes no arguments.

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions. Permissions in category
`monitoring` must be granted on the specific monitoring systems.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | global | no | yes
global | monitoringsystem-mgmt | all | no | yes
system | monitoring | no | yes
monitoring | monitoringsystem | list | yes | no

# EXAMPLES

```
soma monitoringsystem-mgmt list
```
