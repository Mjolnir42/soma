# DESCRIPTION

These commands are used to attach properties to a repository.

Properties with inheritance enabled are passed down the tree to all
children of the repository.
Properties with the childrenonly flag are not active on the tree object
they have been defined on during check constraint evaluation.

Properties with `inheritance false childrenonly true` are essentially
inert.

# SYNOPSIS

```
soma repository property create system ${system} on ${repository} view ${view} value ${value} [inheritance ${inherit}] [childrenonly ${child}]
soma repository property create custom ${custom} on ${repository} view ${view} value ${value} [inheritance ${inherit}] [childrenonly ${child}]
soma repository property create service ${service} on ${repository} view ${view} [inheritance ${inherit}] [childrenonly ${child}]
soma repository property create oncall ${oncall} on ${repository} view ${view} [inheritance ${inherit}] [childrenonly ${child}]
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
repository | string | Name of the repository | | no
view | string | Name of the view the property is attached in | | no
system | string | Name of the system property | | no
custom | string | Name of the custom property | | no
service | string | Name of the service property | | no
oncall | string | Name of the oncall property | | no
value | string | Value of the property | | no
inherit | boolean | Flag to enable/disable inheritance | true | yes
child | boolean | Flag to enable/disable childrenonly | false | yes

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | repository | | no | yes
repository | repository-config | property-create | yes | no

# EXAMPLES

```
soma repository property create system dns_zone on example view any value example.org childrenonly true
soma repository property create oncall 24/7-Support on example view external
soma repository property create service 'OpenSSH - Admin Access' on example view internal
soma repository property create custom foobar on example view local value snafu inheritance false
```
