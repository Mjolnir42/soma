# DESCRIPTION

This command is used to list all custom properties of a repository.

# SYNOPSIS

```
soma property-mgmt custom list in ${repository}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
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
global | property-mgmt | list | yes | no | no
repository | property-custom | list | yes | no | no

# EXAMPLES

```
soma property-mgmt custom list in testing
```
