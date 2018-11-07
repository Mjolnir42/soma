# DESCRIPTION

This command is used to list monitoring system capabilities within SOMA.

# SYNOPSIS

```
soma capability list
```

# ARGUMENT TYPES

This command takes no arguments.

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Permissions in category monitoring must be granted on the specific
monitoring systems.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | monitoring | | no | yes
monitoring | capability | list | yes | no

# EXAMPLES

```
soma capability list
```
