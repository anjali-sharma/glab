## glab mr list

List merge requests

### Synopsis

List merge requests

```
glab mr list [flags]
```

### Options

```
  -a, --all                Get all merge requests
  -c, --closed             Get only closed merge requests
  -l, --label string       Filter merge request by label <name>
  -m, --merged             Get only merged merge requests
      --milestone string   Filter merge request by milestone <id>
      --mine               Get only merge requests assigned to me
  -o, --opened             Get only opened merge requests
  -p, --page int           Page number (default 1)
  -P, --per-page int       Number of items to list per page (default 20)
```

### Options inherited from parent commands

```
      --help          Show help for command
  -R, --repo string   Select another repository using the OWNER/REPO format or the project ID. Supports group namespaces
```

### SEE ALSO

* [glab mr](glab_mr.md)	 - Create, view and manage merge requests

###### Auto generated by spf13/cobra on 1-Oct-2020
