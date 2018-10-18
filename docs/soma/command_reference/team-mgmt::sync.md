# DESCRIPTION

This command dumps all teams defined in SOMA in a format suitable for
processing by sync tools.

# SYNOPSIS

```
soma team-mgmt sync
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
identity | team-mgmt | sync | yes | no

# EXAMPLES

```
soma team-mgmt sync
```
