# DESCRIPTION

This command is used to remove unit definitions from SOMA.

# SYNOPSIS

```
soma unit remove ${unit}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
unit | string | Unit of the unit | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | global | | no | yes
global | unit | remove | yes | no

# EXAMPLES

```
soma unit remove B
soma unit remove s
```
