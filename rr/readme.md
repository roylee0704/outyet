# ResponseRecorder (Fakes)

It is actually a good example of `fake` which implements `ResponseWriter`.


## Good Practices
You should write functions that takes broad interfaces. Don't write functions that takes concrete type: `os.File`.


Instead you want to write functions that takes a `io.Reader`, like just an interface type that specifies the method you actually want from the file.
