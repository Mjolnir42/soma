# DESCRIPTION

This command is used to print all check instances with their rollout
workflow states.

# SYNOPSIS

```
soma workflow list
```

# ARGUMENT TYPES

This command takes no arguments.

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | operation | | no | yes
operation | workflow | list | yes | no

# EXAMPLES

```
soma workflow list
```
