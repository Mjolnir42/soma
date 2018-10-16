# DESCRIPTION

This command is used to add system properties to SOMA, which are global
key/value properties available to all repositories.

The system property name may not contain the `/` character.

# SYNOPSIS

```
soma property-mgmt system add ${property}
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
global | property-system | add | yes | no

# EXAMPLES

```
soma property-mgmt system add dns_zone
soma property-mgmt system add fqdn
```
