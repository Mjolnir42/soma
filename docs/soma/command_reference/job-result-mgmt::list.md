# DESCRIPTION

This command list all job result defined in SOMA by id.

# SYNOPSIS

```
soma job result-mgmt list
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
global | job-result-mgmt | list | yes | no

# EXAMPLES

```
soma job result-mgmt list
```
