# dcos-checks

### add a new check
1. `go get github.com/spf13/cobra/cobra`
2. `cobra add -t github.com/dcos/dcos-checks/cmd/checks/<subcommand> <subcommand>`
3. edit `cmd/checks/<subcommand>.go`
4. rename `init()` to `func Add(root *cobra.Command)` and change it accordingly:

    ```
    func Add(root *cobra.Command) {
      root.AddCommand(<subcommand>)
      // other flag and configuration goes here
    }
    ```
5. Modify `cmd/subcommands.go` to add your check.


see https://github.com/spf13/cobra for more details
