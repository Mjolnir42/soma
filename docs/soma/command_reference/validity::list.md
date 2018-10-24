# DESCRIPTION

This command is used to list all viability definitions from SOMA.

System property names must not contain `/` characters.

# SYNOPSIS

```
soma validity list
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
property | string | Name of the system property | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | global | | no | yes
global | validity | list | yes | no

# EXAMPLES

```
soma validity list
```
