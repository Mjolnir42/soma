# DESCRIPTION

This command shows details for a system property defined in SOMA.

The system property name may not contain the `/` character.

# SYNOPSIS

```
soma property-mgmt system show ${property}
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
global | property-system | show | yes | no

# EXAMPLES

```
soma property-mgmt system show dns_zone
soma property-mgmt system show fqdn
```
