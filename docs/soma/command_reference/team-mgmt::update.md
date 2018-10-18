# DESCRIPTION

This command is used to update teams in SOMA. The team is not patched
but replaced with the new information, so the new record has to be fully
specified.

This command is intended for use by automated LDAP synchronization
tooling.

# SYNOPSIS

```
soma team-mgmt add ${team} name ${name} ldap ${ldapID} [system ${system}]
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
team | string | Current name of the team | | no
name | string | New name of the team | | no
ldapID | integer | Numeric LDAP ID of the team | | no
system | boolean | Boolean flag if this is a system team | false | yes

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | identity | | no | yes
identity | team-mgmt | update | yes | no

# EXAMPLES

```
soma team-mgmt update wheel name wheel ldap 0 system true
```
