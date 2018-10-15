# DESCRIPTION

This command dumps all datacenter information from SOMA suitable for
processing by an external sync command.

# SYNOPSIS

```
soma datacenter sync
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
global | datacenter | sync | yes | no

# EXAMPLES

```
soma datacenter sync
```
