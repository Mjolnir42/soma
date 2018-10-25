# DESCRIPTION

This command is used to print out the number of check instances in every
workflow state.

# SYNOPSIS

```
soma workflow summary
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
operation | workflow | summary | yes | no

# EXAMPLES

```
soma workflow summary
```
