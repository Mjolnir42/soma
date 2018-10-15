# DESCRIPTION

This command lists all datacenters defined in SOMA.

# SYNOPSIS

```
soma datacenter list
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
global | datacenter | list | yes | no

# EXAMPLES

```
soma datacenter list
```
