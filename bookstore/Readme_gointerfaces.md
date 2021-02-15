https://www.alexedwards.net/blog/interfaces-explained

To help reduce duplication or boilerplate code.
To make it easier to use mocks instead of real objects in unit tests.
As an architectural tool, to help enforce decoupling between parts of your codebase.

- E.G Stringer interface, pass to 
fun WriteLog(s fmt.Stringer){
    log.Println(s.String())
}

But if you're new to Go, this still begs a couple of questions: How do you know that the io.Writer interface even exists? And how do you know in advance that bytes.Buffer and os.File both satisfy it?

