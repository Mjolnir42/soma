# DESCRIPTION

This command is used to add teams to SOMA. In addition to its name, a
team also has a numeric ldapID and a system flag. System teams are
internal and not synchronized from LDAP.

# SYNOPSIS

```
soma team-mgmt add ${team} ldap ${ldapID} [system ${system}]
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
team | string | Name of the team | | no
ldapID | integer | Numeric LDAP ID of the team | | no
system | boolean | Boolean flag if this is a system team | false | yes

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | identity | | no | yes
identity | team-mgmt | add | yes | no

# EXAMPLES

```
soma team-mgmt add wheel ldap 0 system true
```
