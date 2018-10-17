# DESCRIPTION

This command shows details for a native introspection property defined in SOMA.

The native property name may not contain the `/` character.

# SYNOPSIS

```
soma property-mgmt native show ${property}
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
global | property-mgmt | show | yes | no
global | property-native | show | yes | no

# EXAMPLES

```
soma property-mgmt native show environment
```
