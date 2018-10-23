# DESCRIPTION

This command is used to add notification levels to SOMA.

Level names must not contain `/` characters.

# SYNOPSIS

```
soma level add ${lvl} shortname ${abbrev} numeric ${num}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
lvl | string | Name of the level | | no
abbrev | string | Abbreviation of the level name | | no
num | uint16 | Numeric value of the level | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | global | | no | yes
global | level | add | yes | no

# EXAMPLES

```
soma level add ok shortname ok numeric 0
soma level add informational shortname info numeric 2
soma level add error shortname err numeric 5
```
