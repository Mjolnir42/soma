# DESCRIPTION

This command is used to add a native introspection property to SOMA,
which are fixed special properties. Support for them is coded into the
middleware.

By not adding an existing native introspection property to property
system, it is possible to effectively disable use of that property
as a check constraint.

The native introspection property name may not contain the `/` character.

# SYNOPSIS

```
soma property-mgmt native add ${property}
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
global | property-mgmt | add | yes | no
global | property-native | add | yes | no

# EXAMPLES

```
soma property-mgmt native add environment
soma property-mgmt native add entity
```
