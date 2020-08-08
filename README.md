# metatag-gen

CLI tool of meta tag generator.



## Install

```bash
go get github.com/shinshin86/metatag-gen
```



## Usage

This tool can be used by giving the required information as parameters.

```bash
metatag-gen {parameters}

# If no parameter is specified, the output is empty. (Only template strings.)
```



### Example

```bash
go run main.go -u https://github.com/shinshin86/metatag-gen -t shinshin86/metatag-gen -d "CLI tool of meta tag generator." -k "meta tag, generator, cli" -i "/examples/images"
```

output(HTML)

```html
<meta charset="utf-8">
<title>shinshin86/metatag-gen</title>
<meta name="viewport" content="width=device-width,initial-scale=1">
<meta name="title" content="shinshin86/metatag-gen">
<meta name="description" content="CLI tool of meta tag generator.">
<meta name="keywords" content="meta tag, generator, cli">

<!-- Open Graph Meta Tags -->
<meta property="og:type" content="website">
<meta property="og:url" content="https://github.com/shinshin86/metatag-gen">
<meta property="og:title" content="shinshin86/metatag-gen">
<meta property="og:description" content="CLI tool of meta tag generator.">
<meta property="og:image" content="/examples/images">

<!-- Twitter -->
<meta property="twitter:card" content="summary_large_image">
<meta property="twitter:url" content="https://github.com/shinshin86/metatag-gen">
<meta property="twitter:title" content="shinshin86/metatag-gen">
<meta property="twitter:description" content="CLI tool of meta tag generator.">
<meta property="twitter:image" content="/examples/images">
```



### parameters

```bash
Usage of metatag-gen
  -d string
    	Description
  -i string
    	OGP Image
  -k string
    	Keywords (If you want to specify more than one, please separate them with a comma.)
  -t string
    	Title
  -tmpl string
    	Use template (default "html")
  -u string
    	URL
```



You can also check with this command.

```bash
metatag-gen --help
```



## TODO

I will also support these templates.

- [x] pug  (Node.js)
- [ ] HAML (Ruby) 
- [ ] Slim (Ruby)

