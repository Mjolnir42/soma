# DESCRIPTION

This command shows details for a datacenter defined in SOMA.

# SYNOPSIS

```
soma datacenter show ${locode}
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
global | datacenter | show | yes | no

# EXAMPLES

```
soma datacenter show de.fra
```
