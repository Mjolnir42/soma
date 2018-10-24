# DESCRIPTION

This command is used to remove all viability definitions for a specific
system property from SOMA.

Viability records can not be removed if the property is in use within a
configuration tree.

System property names must not contain `/` characters.

# SYNOPSIS

```
soma validity remove ${property}
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
global | validity | remove | yes | no

# EXAMPLES

```
soma validity remove tag
soma validity remove disable_all_monitoring
```
