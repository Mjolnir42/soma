# DESCRIPTION

This command is used to rename datacenter definitions in SOMA.

# SYNOPSIS

```
soma datacenter rename ${old} to ${new}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
old | string | old UN/Locode of the datacenter | | no
new | string | new UN/Locode of the datacenter | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | global | | no | yes
global | datacenter | rename | yes | no

# EXAMPLES

```
soma datacenter rename ln.sam to nl.ams
```
