# DESCRIPTION

This command is used to remove datacenter definitions from SOMA.

# SYNOPSIS

```
soma datacenter remove ${locode}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
locode | string | UN/Locode of the datacenter | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | global | | no | yes
global | datacenter | remove | yes | no

# EXAMPLES

```
soma datacenter remove us.chi
```
