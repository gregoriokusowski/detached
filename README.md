# detached

Next-gen thin clients.

### The idea

Worried about your backups?
Tired of burning your fingers when running the test suite?
Having different configs between your private and work machines?
Preparing yourself to buy a more powerful (and heavier) laptop next year?

The idea of detached is to group practices and make it easy to run your environment outside of your laptop.
Why would you need a workhorse if you can have a shiny and light poney?

### The concept

Basically, throw the heavy lifting to another machine. So you will need a server and a client.

##### Server

You can use VPSs (AWS EC2, Google Container/Compute engines, etc), some old desktop, or a server at your workplace.

##### Client

Any machine, depending on your OS preferences and mobility needs.
You can use your phone, tablet, netbook, laptop... Running Android, iOS, Linux, Windows, OSX, as long as you have a SSH (or mosh) client available.

### Ok, so how will detached help me?

At first, detached aims to provide information about the requirements (like getting started topics for the tooling you need to know).

Second step will be to provide a skeleton/framework where you can attach your basic data (dotfiles, credentials, etc) and it will bootstrap your environment.

Having that in mind, we can still provide more information about how to use Docker, X-Server, etc. within your setup and how those tools can help you to improve your workflow.

### Usage

Create your default config with:

```bash
detached init
```

Spin up your machine and start working:

```bash
detached attach
```
