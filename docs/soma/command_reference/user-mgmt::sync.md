# DESCRIPTION

This command lists users in SOMA in a format suitable for external
sync tools.

In comparison to `list`, the `sync` command does not list users that are
flagged as system users. It does however list deleted users. It also
contains all details that can be updated using the `update` command:
userID, username, firstname, lastname, employeenr, mailaddr, team and
deleted.

# SYNOPSIS

```
soma user-mgmt sync
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
identity | user-mgmt | sync | yes | no

# EXAMPLES

```
soma user-mgmt sync
```
