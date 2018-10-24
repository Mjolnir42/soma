# DESCRIPTION

This command is used to add system property validity definitions to
SOMA. For system properties it is necessary to declare where in the
tree the property can be used.
The absence of a validity definition represents invalidity.

System property names must not contain `/` characters.

Boolean values must be parsable by strconv.ParseBool, which accepts
1, t, T, TRUE, true, True, 0, f, F, FALSE, false, False.

# SYNOPSIS

```
soma validity add ${property} on ${entity} direct ${directAllowed} inherited ${inheritedAllowed}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
property | string | Name of the system property | | no
entity | string | Name of the entity | | no
directAllowed | boolean | Boolean value | | no
inheritedAllowed | boolean | Boolean value | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | global | | no | yes
global | validity | add | yes | no

# EXAMPLES

```
soma validity add tag on node direct true inherited true
soma validity add cluster_address on cluster direct true inherited false
soma validity add cluster_address on node direct false inherited true
```
