# DESCRIPTION

This command is used to rename a state definitions in SOMA. Outside of
typofixing during system setup this command is probably unused.

# SYNOPSIS

```
soma state rename ${state} to ${new-state}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
state | string | Old name of the state | | no
new-state | string | New name of the state | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | global | | no | yes
global | state | rename | yes | no

# EXAMPLES

```
soma state rename stadnalone to standalone
```
