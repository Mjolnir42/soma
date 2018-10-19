# DESCRIPTION

This command is used to create a new custom property. Custom properties
are key/value pairs with fixed keys similar to system properties, except
that they are per-repository properties instead of global like system
properties.

# SYNOPSIS

```
soma property-mgmt custom add ${property} to ${repository}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
property | string | Name of the custom property | | no
repository | string | Name or UUID of the repository | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient, all system or all required permissions. Repository scoped
permissions must be on the repository the property should be created in.

Category | Section | Action | Required | System | Sufficient
 ------- | ------- | ------ | -------- | ------ | ----------
omnipotence | | | no | no | yes
system | global | | no | yes | no
system | repository | | no | yes | no
global | property-mgmt | add | yes | no | no
repository | property-custom | add | yes | no | no

# EXAMPLES

```
soma property-mgmt custom add foobar to testing
```
