# tmux

When connecting to remote servers, usually you don't need to keep a long-lasting session.
Considering that detached aims to give you a remote environment, it should keep its state even if you get disconnected.

So here we have [tmux](http://tmux.github.io/). In the terminal multiplexer category, tmux seems to be the most common choice nowadays.

When stabilishing connection to a remote machine, tmux will provide ways to manage different sessions, that are not directly dependent to your connection.

### Getting started

In order to get started with terminal multiplexing, you just need to know the basics:

* Create a new tmux session
* Detach from the current session
* Reattach to the session you created
* Windows/Panes (Create and navigate through them)

First of all, this [simple tutorial](https://www.sitepoint.com/tmux-a-simple-start/) will give you a simple overview of installation and concepts.

After this, you might consider to learn some commands in this [cheatsheet](https://gist.github.com/MohamedAlaa/2961058).

You can customize your setup or you can use some defaults, like [byobu](http://byobu.co/) or .tmux.

### Oh My Tmux!

