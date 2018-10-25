# DESCRIPTION

This command is used to search and list all check instances in a specific
rollout workflow state.

# SYNOPSIS

```
soma workflow search ${status}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
status | string | Name of the status to search for | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | operation | | no | yes
operation | workflow | search | yes | no

# EXAMPLES

```
soma workflow search rollout_failed
soma workflow search blocked
```
