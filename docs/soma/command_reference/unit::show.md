# DESCRIPTION

This command shows details for a measurement unit defined in SOMA.

# SYNOPSIS

```
soma unit show ${unit}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
unit | string | Glyph of the unit | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | global | | no | yes
global | unit | show | yes | no

# EXAMPLES

```
soma unit show B
```
