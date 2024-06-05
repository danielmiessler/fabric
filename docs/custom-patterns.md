# Custom Patterns
Creating Custom Patterns is a key feature of Fabric. This section will cover how to create and use custom patterns in Fabric.



## Table of Contents
- [Storing Custom Patterns](#storing-custom-patterns)



## Storing Custom Patterns
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

**!!! Note:** If you store you patterns directly in `~/.config/fabric/patterns`, they will be overwritten when you update Fabric. So it's better to store them in a separate directory.



## Writing Custom Patterns
To write a custom pattern, you need to create a new folder with the name of your pattern or where you store your patterns. Inside this folder, you need to create a `system.md` file. This file will the system prompt you are giving to Fabric.


### System Prompt
- **Definition:** A system prompt is an initial instruction or context given to the AI by its developers or operators. It sets the tone, behavior, and boundaries for the AI's responses throughout a session.
- **Purpose:** It guides the AI on how to behave, what kind of responses to provide, and how to interpret user inputs. The system prompt can be used to define the AI's personality, scope of knowledge, and conversation style.


You can start by modifying a existing pattern which you can find in the [/patterns](https://github.com/danielmiessler/fabric/tree/main/patterns) directory. [Template for new patterns](https://github.com/danielmiessler/fabric/blob/main/patterns/official_pattern_template/)