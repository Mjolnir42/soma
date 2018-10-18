# DESCRIPTION

This command removes a global service template property from SOMA.

# SYNOPSIS

```
soma property-mgmt template remove ${property}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
property | string | Name of the property | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | global | | no | yes
global | property-mgmt | remove | yes | no
global | property-template | remove | yes | no

# EXAMPLES

```
soma property-mgmt template remove SSH
```
