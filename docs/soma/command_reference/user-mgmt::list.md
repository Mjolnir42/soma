# DESCRIPTION

This command lists all users in SOMA that are not flagged as removed.

# SYNOPSIS

```
soma user-mgmt list
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
identity | user-mgmt | list | yes | no

# EXAMPLES

```
soma user-mgmt list
```
