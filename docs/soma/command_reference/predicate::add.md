# DESCRIPTION

This command is used to add threshold predicates to SOMA.

Predicates must not contain `/` characters.

# SYNOPSIS

```
soma predicate add ${pred}
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
global | predicate | add | yes | no

# EXAMPLES

```
soma predicate add '!='
```
