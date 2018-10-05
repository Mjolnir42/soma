# DESCRIPTION

This command lists all workflow status defined in SOMA.

# SYNOPSIS

```
soma status list
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
global | status | list | yes | no

# EXAMPLES

```
soma status list
```
