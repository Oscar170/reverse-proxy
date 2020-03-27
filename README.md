# Component proxy
#### Proxy to use with some backend service to add custom render whit http calls.
#### This can be used in no node backends to use some fontend lib to build the ui with ssr.
#### We can use whatever framework/lib that can use ssr in the render api implementation
#### The render api should return this json sheme.
```
{
    html: string,
    css: string?,
    initState: Object,
}
```
