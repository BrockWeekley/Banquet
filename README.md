#<span style="color:#8accee">Welcome to Banquet</span>

---

#### A distributable Software as a Service tool that allows for branding and automatic
#### deployment of custom-built JavaScript applications


## Available Commands:
###

---

```shell
> banquet course <option> <repository_link?> <project_name>
```
Provides the ability to add or remove a project from Banquet using a GitHub repository
link.
###

> <span style="color:#37ffa7">option:</span> 
> 
> ```add``` (Clones the provided repository link into the file structure and makes it
> available to the waiter for serving)
> 
> ```remove``` (Removes a provided project from the file structure and removes it from
> the menu)

> <span style="color:#37ffa7">repository link:</span>
> 
> An HTTP link to a public GitHub repository

> <span style="color:#37ffa7">project name:</span>
> 
> Any string that does not already exist within the menu folder that will be used to
> identify the cloned project
###

---

```shell
> banquet reserve
```
Builds a custom React web portal that can be used to customize and manage courses 
that have been added to the menu.
###

---

```shell
> banquet help
```
Prints a list of available commands and information about the project.