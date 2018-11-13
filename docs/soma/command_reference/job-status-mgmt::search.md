# DESCRIPTION

This command is used to search job status definitions from SOMA
by id or name.

At least one search criteria must be specified. If both are specified,
then both must match.

# SYNOPSIS

```
soma job status-mgmt search [id ${uuid}] [name ${status}]
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
uuid | string | UUID of the status | | yes
status | string | Name of the status | | yes

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | global | | no | yes
global | job-status-mgmt | search | yes | no

# EXAMPLES

```
soma job status-mgmt search name processed
soma job status-mgmt search id 703ef1d9-c2aa-458e-8275-624ecf23398c
soma job status-mgmt search id 703ef1d9-c2aa-458e-8275-624ecf23398c name processed
```
