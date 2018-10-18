# DESCRIPTION

This command shows details for a global service template defined in
SOMA.

# SYNOPSIS

```
soma property-mgmt template show ${property}
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
global | property-template | show | yes | no

# EXAMPLES

```
soma property-mgmt template show SSH
```
