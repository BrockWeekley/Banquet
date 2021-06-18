# <span style="color:#8accee">Welcome to Banquet</span>

---

**A distributable Software as a Service tool that allows for branding and automatic 
deployment of custom-built JavaScript applications.**


## Available Commands:
###

---

```shell
> banquet course <option> <project_name> <repository_link?>
```
Provides the ability to add or remove a project from Banquet using a GitHub repository
link.
###

> <span style="color:#37ffa7">option:</span> 
> 
> ```add``` (Clones the provided repository link into the file structure and makes it
> available to the waiter for serving.)
> 
> ```remove``` (Removes a provided project from the file structure and removes it from
> the menu.)

> <span style="color:#37ffa7">project name:</span>
> 
> Any string that will become the name of the project to be cloned, which will be used
> as an identifier for that project.

> <span style="color:#37ffa7">repository link:</span>
> 
> An HTTP link to a public GitHub repository, only required for ```banquet add```.
###

---

```shell
> banquet -serve(?) reserve <port>
```
Builds a custom React web portal that can be used to customize and manage courses 
that have been added to the menu.
###

> <span style="color:#37ffa7">port:</span>
>
> Specifies the port that the waiter web portal should serve on.
###

---

```shell
> banquet help
```
Prints a list of available commands and information about the project.
###

---
###

## Waiter:
**The waiter is the name of the web portal that you can access to serve various
applications from the menu across various ports or the Firebase console.**
###

---
> <span style="color:#ffa611">Firebase setup</span>
> 
> In order to begin deploying applications from the menu with Firebase hosting,
> you must first set up a Firebase account and initial project.