# DESCRIPTION

This command is used to add global service template properties to SOMA,
which are key/array-of-value dictionary objects like services, but not
scoped to a specific team. They can not be used within configurations,
but need to be imported for the team that wants to use them first.

Possible attribute values that the service can contain depends on the
SOMA attributes command. Attributes with a cardinality of `multi` can be
specified more than once.

# SYNOPSIS

```
soma property-mgmt template add ${property} [${attribute} ${attrValue}, ...]
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
property | string | Name of the property | | no
attribute | string | Name of the attribute | | yes
attrValue | string | Value of the attribute | | yes

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | global | | no | yes
global | property-mgmt | add | yes | no
global | property-template | add | yes | no

# EXAMPLES

```
soma property-mgmt template add SSH \
  software_provider OpenSSH \
  transport_protocol tcp \
  application_protocol SSHv2 \
  port 22 \
  process_name sshd \
  uid 0
```
