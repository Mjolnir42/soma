# DESCRIPTION

This command lists all teams defined in SOMA.

# SYNOPSIS

```
soma team-mgmt list
```

# ARGUMENT TYPES

This command takes no arguments.

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | identity | | no | yes
identity | team-mgmt | list | yes | no

# EXAMPLES

```
soma team-mgmt list
```
