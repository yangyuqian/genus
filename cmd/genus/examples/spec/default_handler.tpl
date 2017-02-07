var DefaultHandler = http.HandlerFunc(func(w http.ResponseWriter,r *http.Request) {
    w.Write([]byte(`This is a generated command-line interface for general web application.

Go to [Genus](https://github.com/yangyuqian/genus) for more details.`))
})
