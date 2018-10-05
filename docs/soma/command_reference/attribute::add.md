# DESCRIPTION

This command is used to add a new attribute to the system that can then
be used when defining a service.

Attributes have a cardinality how many values a service is allowed to
have for said attribute. The valid cardinalities are `once` if the
attribute is only allowed to have a single value or `multi` if multiple
values are allowed (not required).

# SYNOPSIS

```
soma attribute add ${attribute} cardinality ${cardinality}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
attribute | string | Name of the attribute | | no
cardinality | string | Cardinality keyword for this attribute | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | global | | no | yes
global | attribute | add | yes | no

# EXAMPLES

```
soma attribute add credential_password cardinality once
soma attribute add credential_user cardinality once
soma attribute add port cardinality multi
```
