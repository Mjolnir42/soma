# DESCRIPTION

This command lists all states defined in SOMA.

# SYNOPSIS

```
soma state list
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
global | state | list | yes | no

# EXAMPLES

```
soma state list
```
