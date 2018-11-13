# DESCRIPTION

This command is used to add job workflow status definitions to SOMA.

# SYNOPSIS

```
soma job status-mgmt add ${status}
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
global | job-status-mgmt | add | yes | no

# EXAMPLES

```
soma job status-mgmt add queued
soma job status-mgmt add in_progress
soma job status-mgmt add processed
```
