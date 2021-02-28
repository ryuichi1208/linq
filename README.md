## Sample

``` go
resp, err := http.Get("https://example.com/api")
if err != nil {
    return nil, err
}
defer func() {
    defer resp.Body.Close()
    io.Copy(ioutil.Discard, resp.Body)
}
if resp.StatusCode < 200 || 299 < resp.StatusCode {
    return nil, errors.New("something error message...")
}
var t T
if err := json.NewDecoder(resp.Body).Decode(&t); err != nil {
    return nil, err
}
return &t, nil
```
