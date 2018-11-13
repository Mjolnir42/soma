# DESCRIPTION

This command is used to remove job status definitions from SOMA.

# SYNOPSIS

```
soma job status-mgmt remove ${status}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
status | string | Name of the status | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | global | | no | yes
global | job-status-mgmt | remove | yes | no

# EXAMPLES

```
soma job status-mgmt remove queued
```
