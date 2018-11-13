# DESCRIPTION

This command list all job status defined in SOMA by id.

# SYNOPSIS

```
soma job status-mgmt list
```

# ARGUMENT TYPES

This command takes no arguments.

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | global | | no | yes
global | job-status-mgmt | list | yes | no

# EXAMPLES

```
soma job status-mgmt list
```
