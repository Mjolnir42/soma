# DESCRIPTION

This command is used to add a new oncall duty team to SOMA.

Oncall duty names must not be formatted as UUIDs.
Phone extensions must be numbers of 4 or 5 digits length.

# SYNOPSIS

```
soma oncall add ${name} phone ${extension}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
name | string | Name of the oncall duty | | no
extension | integer | Numeric phone extension of this oncall duty | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | global | | no | yes
global | oncall | add | yes | no

# EXAMPLES

```
soma oncall add "Emergency Phone" phone 9111
```
