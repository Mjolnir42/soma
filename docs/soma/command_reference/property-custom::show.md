# DESCRIPTION

This command is used to show a custom property.

# SYNOPSIS

```
soma property-mgmt custom show ${property} in ${repository}
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
global | property-mgmt | show | yes | no | no
repository | property-custom | show | yes | no | no

# EXAMPLES

```
soma property-mgmt custom show foobar in testing
```
