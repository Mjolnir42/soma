# DESCRIPTION

This command is used to remove threshold predicates from SOMA.

Predicates must not contain `/` characters.

# SYNOPSIS

```
soma predicate remove ${pred}
```

# ARGUMENT TYPES

Name | Type |     Description   | Default | Optional
 --- |  --- | ----------------- | ------- | --------
pred | string | Symbol of the predicate | | no

# PERMISSIONS

The request is authorized if the user either has at least one
sufficient or all required permissions.

Category | Section | Action | Required | Sufficient
 ------- | ------- | ------ | -------- | ----------
omnipotence | | | no | yes
system | global | | no | yes
global | predicate | remove | yes | no

# EXAMPLES

```
soma predicate remove '=='
```
