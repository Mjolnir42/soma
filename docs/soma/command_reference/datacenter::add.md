# DESCRIPTION

This command is used to add datacenter definitions to SOMA.

The UN/Locode may not contain / characters.

# SYNOPSIS

```
soma datacenter add ${locode}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
locode | string | UN/Locode of the datacenter | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | global | | no | yes
global | datacenter | add | yes | no

# EXAMPLES

```
soma datacenter add de.fra
soma datacenter add us.chi
soma datacenter add de.ber.foo
```
