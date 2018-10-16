# DESCRIPTION

This command lists all system properties defined in SOMA.

# SYNOPSIS

```
soma property-mgmt system list
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
global | property-mgmt | list | yes | no
global | property-system | list | yes | no

# EXAMPLES

```
soma property-mgmt system list
```
