# DESCRIPTION

These commands are used to destroy properties of various types attached
to a repository.

# SYNOPSIS

```
soma repository property destroy system ${system} on ${repository} view ${view}
soma repository property destroy custom ${custom} on ${repository} view ${view}
soma repository property destroy service ${service} on ${repository} view ${view}
soma repository property destroy oncall ${oncall} on ${repository} view ${view}
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

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | repository | | no | yes
repository | repository-config | property-destroy | yes | no

# EXAMPLES

```
soma repository property destroy service 'OpenSSH - Admin Access' on example view internal
soma repository property destroy system dns_zone on example view any
```
