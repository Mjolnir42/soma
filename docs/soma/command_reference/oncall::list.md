# DESCRIPTION

This command lists all oncall duty teams in SOMA.

# SYNOPSIS

```
soma oncall list
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
global | oncall | list | yes | no

# EXAMPLES

```
soma oncall show list
```
