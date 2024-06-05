# Custom Patterns
Creating Custom Patterns is a key feature of Fabric. This section will cover how to create and use custom patterns in Fabric.



## Table of Contents
- [Custom Patterns](#custom-patterns)



## Custom Patterns
You can also use Custom Patterns with Fabric, meaning Patterns you keep locally and don't upload to Fabric.

One possible place to store them is `~/.config/custom-fabric-patterns`. 

Then when you want to use them, simply copy them into `~/.config/fabric/patterns`.

```bash
cp -a ~/.config/custom-fabric-patterns/* ~/.config/fabric/patterns/`
```

Now you can run them with:

```bash
pbpaste | fabric -p your_custom_pattern
```

