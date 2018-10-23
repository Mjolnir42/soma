# DESCRIPTION

This command is used to look up a notification level in SOMA when it
is unclear if the search term is the long or short name.

Level names must not contain `/` characters.

# SYNOPSIS

```
soma level search ${lvl}|${abbrev}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
lvl | string | Name of the level | | no
abbrev | string | Abbreviation of the level name | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | global | | no | yes
global | level | search | yes | no

# EXAMPLES

```
soma level search debug
soma level search err
```
