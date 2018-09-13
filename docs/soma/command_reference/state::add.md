# DESCRIPTION

This command is used to add state definitions to SOMA.

# SYNOPSIS

```
soma state add ${state}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
state | string | Name of the state | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | global | | no | yes
global | state | add | yes | no

# EXAMPLES

```
soma state add standalone
soma state add clustered
```
