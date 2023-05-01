# Version checker
> a tool to check if there's new release of a GitHub repo

## usage
```go
versionchecker.Set(0, 0, 1)
versionchecker.Owner = "valyala"
versionchecker.Repo = "fasthttp"
latest, hasNew, err := versionchecker.CheckUpgrade()
if err != nil {
    return
}
if !hasNew {
    fmt.Printf("latest:%s, current is new", latest)
} else {
    fmt.Printf("latest:%s, current is not new", latest)
}


output:
    latest:1.47.0, current is not new
```
